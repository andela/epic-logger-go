package main

import log "github.com/andela/epic-logger-go"

func main() {
	log.Info("I am an info")
	log.Error("I am an error")
	log.Debug("I am done")
	log.Fatal("I am fatal")
}
