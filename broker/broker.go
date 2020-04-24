package broker

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"sync"
)

type Broker struct {
	clients      *sync.Map
	topicManager *topicManager
	config       *Config
	cls          *cluster
}

func NewBroker(config *Config) *Broker {
	b := &Broker{
		clients:      &sync.Map{},
		config:       config,
		topicManager: newTopicManager(),
	}
	b.cls = newCluster(b)
	return b
}

func (b *Broker) Start() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go b.StartTCPListen()
	go b.StartHttpListen()
	go b.StartWsListener()
	go newCluster(b).start()
	wg.Wait()
}

func (b *Broker) StartHttpListen() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:8089").SetUsername("test").SetClientID("server")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	g := gin.New()
	g.POST("/publish", func(context *gin.Context) {

	})
	g.POST("/subscribe", func(context *gin.Context) {

	})
	g.Run(":8099")
}

func (b *Broker) StartWsListener() {
}

func (b *Broker) StartTCPListen() {
	addr := fmt.Sprintf(":%d", b.config.TcpPort)
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ls.Accept()
		if err != nil {
			log.Println("err:", err)
			continue
		}
		go b.handleConnection(conn)
	}
}

func (b *Broker) handleConnection(conn net.Conn) {
	packet, err := packets.ReadPacket(conn)
	if err != nil {
		panic(err)
	}
	cp, ok := packet.(*packets.ConnectPacket)
	if !ok {
		panic(err)
	}

	connack := packets.NewControlPacket(packets.Connack).(*packets.ConnackPacket)
	connack.SessionPresent = cp.CleanSession
	connack.ReturnCode = cp.Validate()

	if cp.Username != "test" {
		connack.ReturnCode = packets.ErrRefusedBadUsernameOrPassword
		err = connack.Write(conn)
		if err != nil {
			panic(err)
		}
	}

	err = connack.Write(conn)
	if err != nil {
		panic(err)
	}

	client := newClient(cp.ClientIdentifier, conn)
	client.broker = b
	client.cls = b.cls
	b.addClient(client)
}

func (b *Broker) addClient(c *Client) {
	b.clients.Store(c.Id, c)
	c.readLoop()
}
