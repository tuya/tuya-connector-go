package message

import (
	"context"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"sync/atomic"
	"time"

	"github.com/tuya/pulsar-client-go/core/manage"
	"github.com/tuya/pulsar-client-go/core/msg"
)

type consumer struct {
	topic            string
	flowPeriodSecond int
	flowPermit       uint32
	csm              *manage.ManagedConsumer
	cancelFunc       context.CancelFunc
	stopFlag         uint32
	stopped          chan int
}

func (c *consumer) receiveAsync(ctx context.Context, queue chan msg.Message) {
	go func() {
		err := c.csm.ReceiveAsync(ctx, queue)
		if err != nil {
			logger.Log.Debug("consumer stopped, topic=%s", c.topic)
		}
	}()
}

func (c *consumer) receive(ctx context.Context, revHandler messageFunc) {
	queue := make(chan msg.Message, 228)
	ctx, cancel := context.WithCancel(ctx)
	c.cancelFunc = cancel
	go c.cronFlow()
	go c.receiveAsync(ctx, queue)

	for {
		select {
		case <-ctx.Done():
			close(c.stopped)
			return
		case m := <-queue:
			if atomic.LoadUint32(&c.stopFlag) == 1 {
				close(c.stopped)
				return
			}
			logger.Log.Debugf("consumer receive message, topic=%s", c.topic)
			c.handler(context.Background(), &m, revHandler)
		}
	}
}

// redelivered message
func (c *consumer) cronFlow() {
	if c.flowPeriodSecond == 0 {
		return
	}
	if c.flowPermit == 0 {
		return
	}
	tk := time.NewTicker(time.Duration(c.flowPeriodSecond) * time.Second)
	for {
		select {
		case <-c.stopped:
			tk.Stop()
			logger.Log.Infof("stop CronFlow, topic=%s", c.topic)
			return
		case <-tk.C:
			mc := c.csm.Consumer(context.Background())
			if mc == nil {
				continue
			}
			if len(mc.Overflow) > 0 {
				logger.Log.Infof("RedeliverOverflow, topic=%s, num=%d", mc.Topic, len(mc.Overflow))
				_, err := mc.RedeliverOverflow(context.Background())
				if err != nil {
					logger.Log.Warnf("RedeliverOverflow failed, topic=%s, err=%s", mc.Topic, err.Error())
				}
			}

			if mc.Unactive || len(mc.Queue) > 0 {
				continue
			}
			err := mc.Flow(c.flowPermit)
			if err != nil {
				logger.Log.Errorf("flow failed, topic=%s, err=%s", mc.Topic, err.Error())
			}
		}
	}
}

func (c *consumer) handler(ctx context.Context, m *msg.Message, revHandler messageFunc) {
	fields := make([]interface{}, 0, 10)
	fields = append(fields, "Handler trace info, msgID="+m.Msg.GetMessageId().String())
	fields = append(fields, ", topic="+m.Topic)
	defer func(start time.Time) {
		spend := time.Since(start)
		fields = append(fields, ", total spend="+spend.String())
		logger.Log.Debug(fields...)
	}(time.Now())

	now := time.Now()
	var list []*msg.SingleMessage
	var err error
	num := m.Meta.GetNumMessagesInBatch()
	if num > 0 && m.Meta.NumMessagesInBatch != nil {
		list, err = msg.DecodeBatchMessage(m)
		if err != nil {
			logger.Log.Error("DecodeBatchMessage failed", err.Error())
			return
		}
	}
	spend := time.Since(now)
	fields = append(fields, ", decode spend="+spend.String())

	now = time.Now()
	if c.csm.Unactive() {
		logger.Log.Warnf("unused msg because of consumer is unactivated, payload=%s", string(m.Payload))
		return
	}
	spend = time.Since(now)
	fields = append(fields, ", Unactive spend="+spend.String())

	idCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	now = time.Now()
	if c.csm.ConsumerID(idCtx) != m.ConsumerID {
		logger.Log.Warnf("unused msg because of different ConsumerID, payload=%s", string(m.Payload))
		return
	}
	spend = time.Since(now)
	fields = append(fields, ", ConsumerID spend="+spend.String())
	cancel()

	now = time.Now()
	if len(list) == 0 {
		go revHandler(m.Payload)
	} else {
		for i := 0; i < len(list); i++ {
			go revHandler(list[i].SinglePayload)
		}
	}
	spend = time.Since(now)
	fields = append(fields, ", HandlePayload spend="+spend.String())
	if err != nil {
		logger.Log.Errorf("handle message failed, topic=%s, err=%+v", m.Topic, err)
	}

	now = time.Now()
	ackCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err = c.csm.Ack(ackCtx, *m)
	cancel()
	spend = time.Since(now)
	fields = append(fields, ", Ack spend="+spend.String())
	if err != nil {
		logger.Log.Error("ack failed", err.Error())
	}

}

func (c *consumer) stop() {
	atomic.AddUint32(&c.stopFlag, 1)
	c.cancelFunc()
	<-c.stopped
}
