package main

import (
	"github.com/digbrand/potato-mqtt/broker"
)

func main() {
	broker.NewBroker().Start()
}
