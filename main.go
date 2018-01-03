package epiclogger

import (
	"os"

	logrus_stack "github.com/andela/logrus-stack"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/grpclog"
)

func replaceGrpcLogger() {
	grpclog.SetLogger(baseLogger)
}

func init() {
	// do something here to set environment depending on an environment variable
	// or command-line flag
	Environment := os.Getenv("GO_ENV")
	if Environment == "production" || Environment == "staging" {
		SetFormatter(&EpicFormatter{})
		SetLevel(log.InfoLevel)
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		SetFormatter(&TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "15:04:05",
		})
		SetLevel(log.DebugLevel)
	}
	AddHook(logrus_stack.StandardHook())
	// if Environment == "production" {
	// 	bugsnag.Configure(bugsnag.Configuration{
	// 		APIKey: os.Getenv("BUGSNAG_API_KEY"),
	// 	})
	// 	hook, err := logrus_bugsnag.NewBugsnagHook()
	// 	if err != nil {
	// 		AddHook(hook)
	// 	}
	// }
	replaceGrpcLogger()
}
