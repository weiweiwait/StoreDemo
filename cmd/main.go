package main

import (
	"Go-SecKill/config"
	"Go-SecKill/routes"
)

func main() {
	config.Init()
	r := routes.NewRouter()
	_ = r.Run(config.HttpPort)
}
