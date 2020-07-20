package main

import (
	"github.com/NOVAPokemon/utils"
)

const (
	host        = utils.ServeHost
	port        = utils.TrainersPort
	serviceName = "TRAINERS"
)

func main() {
	flags := utils.ParseFlags(serverName)

	if !*flags.LogToStdout {
		utils.SetLogFile(serverName)
	}

	if !*flags.DelayedComms {
		commsManager = utils.CreateDefaultCommunicationManager()
	} else {
		locationTag := utils.GetLocationTag(utils.DefaultLocationTagsFilename, serverName)
		commsManager = utils.CreateDelayedCommunicationManager(utils.DefaultDelayConfigFilename, locationTag)
	}

	utils.StartServer(serviceName, host, port, routes, commsManager)
}
