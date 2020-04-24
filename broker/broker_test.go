package broker

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"testing"
	"time"
)

func TestClusterListen(t *testing.T) {
	go NewBroker(&Config{
		TcpPort: 8089,
		Routes: []*Route{{
			IP:   "",
			Port: 8082,
		}},
	}).StartTCPListen()
	time.Sleep(time.Second)
	go NewBroker(&Config{
		TcpPort: 8082,
		Routes: []*Route{{
			IP:   "",
			Port: 8089,
		}},
	}).StartTCPListen()

	beginClient("client1", 8089)
	beginClient("client2", 8082)

	client3 := beginClient("client3", 8089)
	time.Sleep(100 * time.Millisecond)

	if token := client3.Publish("ab", 0, false, "hello world"); token.Wait() && token.Error() != nil {
		t.Fatal(token.Error())
	}
	time.Sleep(time.Second)
}

func beginClient(clientId string, port int) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:" + strconv.Itoa(port)).SetUsername("test").SetClientID(clientId)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe("ab", 0, func(c mqtt.Client, message mqtt.Message) {
		fmt.Println("clientId:", clientId, " sub:", string(message.Payload()))
	}); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

func TestClientListen(t *testing.T) {
	go NewBroker(&Config{
		TcpPort: 8089,
	}).StartTCPListen()
	time.Sleep(time.Second)

	client1 := beginClient("client1", 8089)

	if token := client1.Publish("ab", 0, false, "hello world"); token.Wait() && token.Error() != nil {
		t.Fatal(token.Error())
	}
	time.Sleep(time.Second)
}
