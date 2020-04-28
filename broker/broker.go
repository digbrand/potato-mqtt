package broker

import (
	"fmt"
	"github.com/digbrand/potato-mqtt/topic"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map"
	"go.uber.org/zap"
	"log"
	"net"
	"sync"
)

type Broker struct {
	clients      cmap.ConcurrentMap
	topicManager topic.Manager
	config       *Config
}

func NewBroker(config *Config) *Broker {
	b := &Broker{
		clients:      cmap.New(),
		config:       config,
		topicManager: topic.NewManager(),
	}
	return b
}

func (b *Broker) Start() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go b.StartTCPListen()
	go b.StartHttpListen()
	go b.StartWsListener()
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
		zap.L().Error("connection error,invalid connect packet", zap.String("packet", cp.String()))
		return
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
		zap.L().Error("write to client error", zap.Error(err))
		return
	}

	client := newClient(cp.ClientIdentifier, conn)
	client.broker = b
	b.addClient(client)
}

func (b *Broker) addClient(c *Client) {
	b.clients.Set(c.Id, c)
	c.readLoop()
}
