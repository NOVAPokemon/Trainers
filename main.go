package main

import (
	"github.com/NOVAPokemon/utils"
	trainerdb "github.com/NOVAPokemon/utils/database/trainer"
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
		commsManager = utils.CreateDefaultDelayedManager(false)
	}

	trainerdb.InitTrainersDBClient(*flags.ArchimedesEnabled)
	utils.StartServer(serviceName, host, port, routes, commsManager)
}
