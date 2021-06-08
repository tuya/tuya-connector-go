package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLog(t *testing.T) {
	tl := NewDefaultLogger("logtest", true)
	assert.False(t, tl == nil, "log handler is nil")
	for i := 0; i < 9; i++ {
		tl.Infof("logtesting, i=%d", i)
	}
}
