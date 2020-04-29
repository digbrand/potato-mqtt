package broker

import (
	"fmt"
	"github.com/digbrand/potato-mqtt/topic"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"go.uber.org/zap"
	"net"
	"time"
)

type Client struct {
	Id        string
	conn      net.Conn
	broker    *Broker
	keepAlive uint16
	timeout   time.Time
	done      chan bool
}

func newClient(id string, conn net.Conn) *Client {
	return &Client{
		Id:   id,
		conn: conn,
		done: make(chan bool),
	}
}

func (c *Client) Publish(packet *packets.PublishPacket) {
	mp, err := c.broker.topicManager.GetSubscribers(packet.TopicName)
	if err != nil {
		zap.L().Error("publish failure", zap.Error(err))
		return
	}
	switch packet.Qos {
	case 0:
	case 1:
		if err := c.publishAck(packet); err != nil {
			zap.L().Error("publish ack write failure", zap.Error(err))
			return
		}
	case 2:
		if err := c.publishAck(packet); err != nil {
			zap.L().Error("publish ack write failure", zap.Error(err))
			return
		}
	default:
		zap.L().Error("unknow publish qos", zap.String("packet", packet.String()))
		return
	}

	for _, v := range mp {
		sub := v.Client.(*Client)
		if err := packet.Write(sub.conn); err != nil {
			zap.L().Error("publish to subscribers failure", zap.Error(err))
			continue
		}
	}
}

func (c *Client) publishAck(packet *packets.PublishPacket) error {
	puback := packets.NewControlPacket(packets.Puback).(*packets.PubackPacket)
	puback.MessageID = packet.MessageID
	if err := puback.Write(c.conn); err != nil {
		return err
	}
	return nil
}

const (
	SubAckQos0    = 0x00
	SubAckQos1    = 0x01
	SubAckQos2    = 0x02
	SubAckFailure = 0x80
)

func (c *Client) Ping(packet *packets.PingreqPacket) {
	msg := packets.NewControlPacket(packets.Pingresp).(*packets.PingrespPacket)
	if err := msg.Write(c.conn); err != nil {
		zap.L().Error("write ping resp packet error", zap.Error(err))
	}
}

func (c *Client) Subscribe(packet *packets.SubscribePacket) {
	suback := packets.NewControlPacket(packets.Suback).(*packets.SubackPacket)
	suback.MessageID = packet.MessageID
	ret := make([]byte, 0)

	for _, tp := range packet.Topics {
		if err := c.broker.topicManager.Subscribe(tp, topic.NewSubscribeThing(c.Id, c)); err != nil {
			ret = append(ret, SubAckFailure)
			continue
		}
		ret = append(ret, suback.Qos)
	}
	suback.ReturnCodes = ret
	if err := suback.Write(c.conn); err != nil {
		zap.L().Error("sub ack write error", zap.Error(err))
	}
}

func (c *Client) processPacket(pk packets.ControlPacket) {
	switch pk.(type) {
	case *packets.ConnackPacket:
	case *packets.ConnectPacket:
	case *packets.PublishPacket:
		packet, ok := pk.(*packets.PublishPacket)
		if !ok {
			zap.L().Error("invalid publish packet", zap.String("packet", packet.String()))
			return
		}
		c.Publish(packet)
	case *packets.PubackPacket:
	case *packets.PubrecPacket:
	case *packets.PubrelPacket:
	case *packets.PubcompPacket:
	case *packets.SubscribePacket:
		packet, ok := pk.(*packets.SubscribePacket)
		if !ok {
			zap.L().Error("invalid subscribe packet", zap.String("packet", packet.String()))
			return
		}
		c.Subscribe(packet)
	case *packets.SubackPacket:
	case *packets.UnsubscribePacket:
		//packet := pk.(*packets.UnsubscribePacket)
		//c.ProcessUnSubscribe(packet)
	case *packets.UnsubackPacket:
	case *packets.PingreqPacket:
		packet, ok := pk.(*packets.PingreqPacket)
		if !ok {
			zap.L().Error("invalid ping packet", zap.String("packet", packet.String()))
			return
		}
		fmt.Println("get ping packet")
		c.Ping(packet)
	case *packets.PingrespPacket:
	case *packets.DisconnectPacket:
		c.conn.Close()
		zap.L().Info("receive disconnect packet,connection closed")
		return
	}
}

func (c *Client) readLoop() {
	if c.keepAlive > 0 {
		go func() {
			for {
				if time.Now().After(c.timeout) {
					c.done <- true
					return
				}
				time.Sleep(time.Second)
			}
		}()
	}

	for {
		select {
		case <-c.done:
			zap.L().Info("client time out,send disconnect packet to local client")
			c.processPacket(packets.NewControlPacket(packets.Disconnect))
			close(c.done)
			return
		default:
			pk, err := packets.ReadPacket(c.conn)
			if err != nil {
				continue
			}
			if c.keepAlive > 0 {
				c.setTimeout()
			}
			c.processPacket(pk)
		}
	}
}

func (c *Client) setTimeout() {
	c.timeout = time.Now().Add(time.Second * time.Duration(c.keepAlive))
}
