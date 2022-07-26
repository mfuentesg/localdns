package main

import (
	"github.com/mfuentesg/localdns/cli/command"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})
}

func main() {
	command.Execute()
}
