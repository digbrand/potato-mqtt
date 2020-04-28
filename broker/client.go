package broker

import (
	"errors"
	"github.com/digbrand/potato-mqtt/topic"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"go.uber.org/zap"
	"net"
)

type Client struct {
	Id     string
	conn   net.Conn
	broker *Broker
}

func newClient(id string, conn net.Conn) *Client {
	return &Client{
		Id:   id,
		conn: conn,
	}
}

func (c *Client) Publish(packet *packets.PublishPacket) {
	//mp,err := c.broker.topicManager.GetSubscribers(packet.TopicName)
	//if err!=nil || len(mp)<1{
	//	packet:=packets.NewControlPacket(packets.Puback)
	//	aa:=packet.(*packets.PubackPacket)
	//	aa.
	//	//todo
	//}
}

func (c *Client) Subscribe(packet *packets.SubscribePacket) {
	suback := packets.NewControlPacket(packets.Suback).(*packets.SubackPacket)
	suback.MessageID = packet.MessageID
	ret := make([]byte, 0)
	suback.ReturnCodes = ret

	if packet.Topics == nil {
		zap.L().Error("topics is nil", zap.Error(errors.New("packet topics must not nil")))
		ret = append(ret, 0x80)
		suback.Write(c.conn)
		return
	}

	for _, tp := range packet.Topics {
		if err := c.broker.topicManager.Subscribe(tp, topic.NewSubscribeThing(c.Id, c)); err != nil {
			ret = append(ret, packets.ErrProtocolViolation)
			continue
		}
		ret = append(ret, packets.Suback)
	}

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
				c.Publish(packet)
			case *packets.PubackPacket:
			case *packets.PubrecPacket:
			case *packets.PubrelPacket:
			case *packets.PubcompPacket:
			case *packets.SubscribePacket:
				packet, ok := pk.(*packets.SubscribePacket)
				if !ok {
					zap.L().Error("invalid subscribe packet", zap.String("packet", packet.String()))
					continue
				}
				c.Subscribe(packet)
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
