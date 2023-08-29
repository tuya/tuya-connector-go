package message

import (
	"context"
	"crypto/tls"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/callmegema/tuya-connector-go/connector/constant"
	"github.com/callmegema/tuya-connector-go/connector/env"
	"github.com/callmegema/tuya-connector-go/connector/env/extension"
	"github.com/callmegema/tuya-connector-go/connector/logger"
	"github.com/callmegema/tuya-connector-go/connector/utils"
	"github.com/tuya/pulsar-client-go/core/manage"
)

const (
	DefaultFlowPeriodSecond = 30
	DefaultFlowPermit       = 10
)

func init() {
	extension.SetMessage(constant.TUYA_MESSAGE, newMessageInstance)
	fmt.Println("init message extension......")
}

func newMessageInstance() extension.IEventMessage {
	return NewEventMsgWrapper()
}

type eventMessage struct {
	done uint32
	c    *client
}

func NewEventMsgWrapper() extension.IEventMessage {
	return &eventMessage{}
}

func (e *eventMessage) SubEventMessage(f interface{}) {
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		logger.Log.Errorf("event message handler is not function")
		return
	}
	msgType := fv.Type().In(0).Elem().Name()
	e.c.mu.Lock()
	e.c.eventSubPool[msgType] = f
	e.c.mu.Unlock()
}

// stop event message receive
func (e *eventMessage) Stop() {
	for _, m := range e.c.csmPool {
		m.stop()
	}
}

// init event message client
func (e *eventMessage) InitMessageClient() {
	if atomic.LoadUint32(&e.done) == 1 {
		return
	}
	defer atomic.StoreUint32(&e.done, 1)
	ak := env.Config.GetAccessID()
	sk := env.Config.GetAccessKey()
	topic := fmt.Sprintf("persistent://%s/out/event", ak)
	authKey := utils.StrToMD5(ak + utils.StrToMD5(sk))[8:24]
	addrs := strings.Split(strings.TrimSuffix(env.Config.GetMsgHost(), "/"), "+ssl")

	// create client
	c := newClient(&clientConfig{
		accessID:   ak,
		accessKey:  sk,
		topic:      topic,
		authMethod: "auth1",
		authData:   fmt.Sprintf(`{"username":"%s","password":"%s"}`, ak, authKey),
		addr:       strings.Join(addrs, ""),
	})

	// sub
	errs := make(chan error, 10)
	go printAsyncErr(errs)
	c.subscribe(manage.ConsumerConfig{
		ClientConfig: manage.ClientConfig{
			Addr:       c.cfg.addr,
			AuthData:   []byte(c.cfg.authData),
			AuthMethod: "auth1",
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Errs: errs,
		},
		Topic:              topic,
		SubMode:            manage.SubscriptionModeFailover,
		Name:               subscriptionName(topic),
		NewConsumerTimeout: time.Minute,
	})
	e.c = c
}

func printAsyncErr(e chan error) {
	for err := range e {
		logger.Log.Error("async errors", err.Error())
	}
}

type clientConfig struct {
	accessID   string
	accessKey  string
	topic      string
	addr       string
	authMethod string
	authData   string
}

type client struct {
	pool         *manage.ClientPool
	cfg          *clientConfig
	csmPool      map[int]*consumer
	eventSubPool map[string]interface{}
	mu           *sync.RWMutex
}

func newClient(c *clientConfig) *client {
	return &client{
		pool:         manage.NewClientPool(),
		csmPool:      make(map[int]*consumer),
		cfg:          c,
		mu:           &sync.RWMutex{},
		eventSubPool: make(map[string]interface{}),
	}
}

func (c *client) subscribe(csmCfg manage.ConsumerConfig) {
	logger.Log.Infof("start creating consumer, pulsar=%s, topic=%s", c.cfg.addr, c.cfg.topic)
	size := 1
	isPartitioned := false
	prt, err := c.pool.Partitions(context.Background(), csmCfg.ClientConfig, c.cfg.topic)
	if err == nil {
		isPartitioned = true
		size = int(prt.GetPartitions())
	}
	originTopic := csmCfg.Topic
	for i := 0; i < size; i++ {
		if isPartitioned {
			csmCfg.Topic = fmt.Sprintf("%s-partition-%d", originTopic, i)
		}
		mc := manage.NewManagedConsumer(c.pool, csmCfg)
		csm := &consumer{
			csm:              mc,
			flowPeriodSecond: DefaultFlowPeriodSecond,
			flowPermit:       DefaultFlowPermit,
			topic:            csmCfg.Topic,
			stopped:          make(chan int),
		}
		go csm.receive(context.Background(), c.receiveMsg())
		c.csmPool[i] = csm
	}
	logger.Log.Infof("create consumer success, pulsar=%s, topic=%s", c.cfg.addr, c.cfg.topic)
}

func subscriptionName(topic string) string {
	topic = strings.TrimPrefix(topic, "persistent://")
	end := strings.Index(topic, "/")
	return topic[:end] + "-sub"
}
