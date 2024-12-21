package impl

import (
	"fmt"
	"github.com/getamis/alice/types"
	"github.com/getamis/sirius/log"
	"github.com/libp2p/go-libp2p/core/network"
	"google.golang.org/protobuf/proto"
	"io"
	"reflect"
)

type Message interface {
	types.Message
	proto.Message
}

type Backend[M Message, R any] interface {
	AddMessage(senderId string, msg types.Message) error
	Start()
	Stop()
	GetResult() (R, error)
}

type Node[M Message, R any] struct {
	backend  Backend[M, R]
	listener Listener
}

func New[M Message, R any](backend Backend[M, R], l Listener) *Node[M, R] {
	return &Node[M, R]{
		backend:  backend,
		listener: l,
	}
}

func (n *Node[M, R]) Handle(s network.Stream) {
	var data M
	buf, err := io.ReadAll(s)
	if err != nil {
		log.Warn("Cannot read data from stream", "err", err)
		return
	}
	s.Close()

	msgType := reflect.TypeOf(data).Elem()
	data = reflect.New(msgType).Interface().(M)
	err = proto.Unmarshal(buf, data)
	if err != nil {
		log.Error("Cannot unmarshal data", "err", err)
		return
	}
	err = n.backend.AddMessage(data.GetId(), data)
	if err != nil {
		log.Warn("Cannot add message to DKG", "err", err)
		return
	}
}

func (n *Node[M, R]) Process() (r R, _ error) {
	n.backend.Start()
	defer n.backend.Stop()
	if err := <-n.listener.Done(); err != nil {
		return r, err
	}
	fmt.Println("=================本次完成======================")
	return n.backend.GetResult()
}
