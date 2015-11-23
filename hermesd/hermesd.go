package main

import (
	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git/packets"
	"github.com/shafreeck/hermes/hermes"
	"io"
	"log"
	"net"
)

var (
	WAITING     = 0
	CONNECT     = 1
	CONNACK     = 2
	ESTABLISHED = 3
	DISCONNECT  = 4
)

type client struct {
	state int
	conn  net.Conn
	id    string
	hms   *hermes.Hermes
}

func (c *client) processConnect(cp packets.ControlPacket) {
	// load all the retained sessions
	// setup the cursor goroution
}

func (c *client) process(in, out chan packets.ControlPacket) {
	conn := c.conn
	for {
		cp, err := packets.ReadPacket(conn)
		if err != nil {
			log.Printf("%s\n", err.Error())
			if err == io.EOF {
				return
			}
		}
		switch cp.(type) {
		case *packets.PublishPacket:
			log.Printf("%s\n", cp.String())
			p := cp.(*packets.PublishPacket)
			log.Printf("%s\n", p.TopicName)

			c.hms.Publish(p.TopicName, p.Payload)
			if p.Qos == 1 {
				pa := packets.NewControlPacket(packets.Puback).(*packets.PubackPacket)
				pa.MessageID = p.MessageID
				pa.Write(conn)
			}

		case *packets.ConnectPacket:
			log.Printf("%s\n", cp.String())
			p := cp.(*packets.ConnectPacket)
			log.Printf("%s\n", p.ProtocolName)
			c.id = p.ClientIdentifier
			ca := packets.NewControlPacket(packets.Connack).(*packets.ConnackPacket)
			ca.Write(conn)

		case *packets.SubscribePacket:
			p := cp.(*packets.SubscribePacket)
			c.hms.Subscribe(p.Topics[0], c.id)
			if p.Qos > 0 {
				pa := packets.NewControlPacket(packets.Suback).(*packets.SubackPacket)
				pa.MessageID = p.MessageID
				pa.Write(conn)
			}
		}
	}
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:1883")
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	hms := hermes.NewHermes()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("%s\n", err.Error())
		}
		c := &client{conn: conn, hms: hms}
		c.state = WAITING
		go c.process()
	}
}
