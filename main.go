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
	utils.CheckLogFlag(serviceName)
	utils.StartServer(serviceName, host, port, routes)
}
