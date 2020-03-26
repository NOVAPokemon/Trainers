package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

const host = "localhost"
const Port = 8004

var addr = fmt.Sprintf("%s:%d", host, Port)

func main() {
	rand.Seed(time.Now().Unix())

	r := utils.NewRouter(routes)

	log.Infof("Starting TRADES server in port %d...\n", Port)
	log.Fatal(http.ListenAndServe(addr, r))
}
