package broker

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
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

func TestClientPublish(t *testing.T) {
	go NewBroker(&Config{
		TcpPort: 8089,
	}).StartTCPListen()
	time.Sleep(time.Second)

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:8089").SetUsername("test").SetClientID("client1")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//if token := client.Subscribe("ab", 0, func(c mqtt.Client, message mqtt.Message) {
	//	fmt.Println("clientId:", clientId, " sub:", string(message.Payload()))
	//}); token.Wait() && token.Error() != nil {
	//	panic(token.Error())
	//}

	//if token := client1.Publish("ab", 0, false, "hello world"); token.Wait() && token.Error() != nil {
	//	t.Fatal(token.Error())
	//}
	time.Sleep(time.Second)
}

func TestClient(t *testing.T) {
	go NewBroker(&Config{
		TcpPort: 8089,
	}).StartTCPListen()
	time.Sleep(time.Second)

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:8089").SetUsername("test").SetClientID("client1")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe("/a/b", 0, func(c mqtt.Client, message mqtt.Message) {
		fmt.Println("client1:", string(message.Payload()))
	}); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	opts = mqtt.NewClientOptions().AddBroker("tcp://localhost:8089").SetUsername("test").SetClientID("client2")

	client2 := mqtt.NewClient(opts)
	if token := client2.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client2.Subscribe("/a/b", 0, func(c mqtt.Client, message mqtt.Message) {
		fmt.Println("client2:", string(message.Payload()))
	}); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	time.Sleep(100 * time.Millisecond)

	if token := client.Publish("/a/b", 0, false, "hello world"); token.Wait() && token.Error() != nil {
		t.Fatal(token.Error())
	}

	fmt.Println("publish finished")
	time.Sleep(time.Second)
}

func TestClientTimeout(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	go NewBroker(&Config{
		TcpPort: 8089,
	}).StartTCPListen()
	time.Sleep(time.Second)

	go func() {
		opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:8089").
			SetKeepAlive(2 * time.Second).SetUsername("test").SetClientID("client1")

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		time.Sleep(3 * time.Second)
	}()

	time.Sleep(5 * time.Second)
}

func TestTimeLoop(t *testing.T) {
	tm := time.Now().Add(2 * time.Second)
	go func() {
		for {
			if time.Now().After(tm) {
				fmt.Println("exit for loop")
				return
			}
		}
	}()

	time.Sleep(100 * time.Second)
}
