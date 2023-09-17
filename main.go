package main

import (
	"WBL0/model"
	"WBL0/nats"
)

func main() {
	nats.Pub()
	nats.Sub()
	model.InitServer()
}
