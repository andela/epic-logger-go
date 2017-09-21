package main

import "github.com/andela/epic-logger-go"

func main() {
	epiclogger.Info("I am an info")
	epiclogger.Error("I am an error")
	epiclogger.Debug("I am done")
	epiclogger.Fatal("I am fatal")
}
