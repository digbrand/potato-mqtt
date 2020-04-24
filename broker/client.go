package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"net"
)

type Client struct {
	Id     string
	conn   net.Conn
	broker *Broker
	cls    *cluster
}

func newClient(id string, conn net.Conn) *Client {
	return &Client{
		Id:   id,
		conn: conn,
	}
}

func (c *Client) processPub(packet *packets.PublishPacket) {
	list := c.broker.topicManager.getSubs(packet.TopicName)
	if list == nil {
		return
	}

	for _, v := range list {
		packet.Write(v.conn)
	}

	c.cls.pub(packet)
}

func (c *Client) processSub(packet *packets.SubscribePacket) {
	for _, t := range packet.Topics {
		c.broker.topicManager.sub(t, c)
	}
	suback := packets.NewControlPacket(packets.Suback).(*packets.SubackPacket)
	suback.MessageID = packet.MessageID
	suback.Write(c.conn)
}

func (c *Client) readLoop() {
	for {
		select {
		default:
			pk, err := packets.ReadPacket(c.conn)
			if err != nil {
				continue
			}

			switch pk.(type) {
			case *packets.ConnackPacket:
			case *packets.ConnectPacket:
			case *packets.PublishPacket:
				packet := pk.(*packets.PublishPacket)
				c.processPub(packet)
			case *packets.PubackPacket:
			case *packets.PubrecPacket:
			case *packets.PubrelPacket:
			case *packets.PubcompPacket:
			case *packets.SubscribePacket:
				packet := pk.(*packets.SubscribePacket)
				c.processSub(packet)
			case *packets.SubackPacket:
			case *packets.UnsubscribePacket:
				//packet := pk.(*packets.UnsubscribePacket)
				//c.ProcessUnSubscribe(packet)
			case *packets.UnsubackPacket:
			case *packets.PingreqPacket:
				//c.ProcessPing()
			case *packets.PingrespPacket:
			case *packets.DisconnectPacket:
				c.conn.Close()
			}
		}
	}
}
