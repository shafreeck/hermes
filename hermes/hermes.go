package hermes

import (
	"bytes"
	"encoding/binary"
	"github.com/shafreeck/hermes/store"
	"log"
)

/* Store key format
 * topic and clientid are all null-terminated
 * Topic:
 *  type + topic + seq
 * Cursor:
 *  type + clientid + topic
 */

var (
	META   byte = 0
	TOPIC  byte = 1
	CURSOR byte = 2
)

type PubSub interface {
	Publish(topic string, message []byte) error
	Subscribe(topic string) (chan []byte, error)
	UnSubscribe(topic string) error
}

type Hermes struct {
	s      store.Store
	writer *Cursor
}
type Cursor struct {
	topic string
	seq   uint64
}

func NewHermes() *Hermes {
	var err error
	h := &Hermes{}
	h.s, err = store.OpenLevelDB("db.data")
	if err != nil {
		log.Printf("NewHermes %s\n", err)
	}
	return h
}
func (h *Hermes) Publish(topic string, message []byte) error {
	keyType := TOPIC

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, keyType)
	buf.Write([]byte(topic))
	binary.Write(buf, binary.LittleEndian, byte(0))
	binary.Write(buf, binary.LittleEndian, uint64(0))
	key := buf.Bytes()
	return h.s.Set(key, message)
}

func (h *Hermes) Subscribe(topic string, clientid string) error {
	keyType := CURSOR
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, keyType)
	buf.Write([]byte(clientid))
	binary.Write(buf, binary.LittleEndian, byte(0))
	buf.Write([]byte(topic))
	binary.Write(buf, binary.LittleEndian, byte(0))
	key := buf.Bytes()
	value := make([]byte, 10)
	return h.s.Set(key, value)
}
