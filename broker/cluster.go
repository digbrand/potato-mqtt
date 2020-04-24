package broker

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"log"
	"net"
)

type cluster struct {
	broker    *Broker
	discovers []net.Conn
}

func newCluster(broker *Broker) *cluster {
	c := &cluster{
		broker: broker,
	}
	return c
}

func (c *cluster) pub(packet *packets.PublishPacket) {
	for _, v := range c.discovers {
		packet.Write(v)
	}
}

func (c *cluster) start() {
	addr := fmt.Sprintf(":%d", c.broker.config.TcpPort)
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, err := ls.Accept()
			if err != nil {
				log.Println("err:", err)
				continue
			}
			go func() {
				go c.broker.handleConnection(conn)
			}()
		}
	}()

	go func() {
		if c.broker.config.Routes == nil {
			return
		}
		for _, v := range c.broker.config.Routes {
			addr := fmt.Sprintf("%s:%d", v.IP, v.Port)
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				panic(err)
			}
			c.discovers = append(c.discovers, conn)

			ck := packets.NewControlPacket(packets.Connect)
			ck.Write(conn)
		}
	}()

	select {}
}
