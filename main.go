package epiclogger

import (
	"os"

	"github.com/Gurpartap/logrus-stack"
	"github.com/Shopify/logrus-bugsnag"
	bugsnag "github.com/bugsnag/bugsnag-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/grpclog"
)

func replaceGrpcLogger() {
	grpclog.SetLogger(std)
}

func init() {
	// do something here to set environment depending on an environment variable
	// or command-line flag
	Environment := os.Getenv("GO_ENV")
	if Environment == "production" {
		log.SetFormatter(&EpicFormatter{})
		log.SetLevel(log.InfoLevel)
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "15:04:05",
		})
		log.SetLevel(log.DebugLevel)
	}
	log.AddHook(NewServiceHook())
	log.AddHook(logrus_stack.StandardHook())
	if Environment == "production" {
		bugsnag.Configure(bugsnag.Configuration{
			APIKey: os.Getenv("BUGSNAG_API_KEY"),
		})
		hook, err := logrus_bugsnag.NewBugsnagHook()
		if err != nil {
			log.AddHook(hook)
		}
	}
	replaceGrpcLogger()
}
