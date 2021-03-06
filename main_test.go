package epiclogger

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestBasicError(t *testing.T) {
	hook := test.NewLocal(baseLogger.Logger)
	Error(errors.New("Helloerror"))
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "Helloerror", hook.LastEntry().Message)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestBasicInfo(t *testing.T) {
	// logger := WithField("sse", []string{"adebayo", "thuita", "james"})
	hook := test.NewLocal(baseLogger.Logger)
	WithField("sse", []string{"adebayo", "thuita", "james"}).Info("I am a simple info")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "I am a simple info", hook.LastEntry().Message)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestServiceHook(t *testing.T) {
	hook := test.NewLocal(baseLogger.Logger)
	WithField("sse", map[string]string{"ais": "adebayo", "pulse": "thuita", "lv": "james"}).Error("I am a simple error")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "I am a simple error", hook.LastEntry().Message)
	assert.Equal(t, "golang-service", hook.LastEntry().Data["service"])
	assert.Equal(t, "123", hook.LastEntry().Data["version"])

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestServiceHooks(t *testing.T) {
	WithField("name", "ikem").Info("I am a simple error")
}
