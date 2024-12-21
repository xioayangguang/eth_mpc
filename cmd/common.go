package cmd

import (
	"eth_mpc/impl"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func common[M impl.Message, R any](p2pProtocol protocol.ID, handler func(*impl.PeerManager, impl.Listener) (*impl.Node[M, R], error)) (*R, error) {
	host, err := impl.GetHostByCfg()
	pm := impl.NewPeerManager(host, p2pProtocol)
	l := impl.NewListener()
	n, err := handler(pm, l)
	if err != nil {
		return nil, err
	}
	host.SetStreamHandler(p2pProtocol, func(s network.Stream) {
		n.Handle(s)
	})
	pm.EnsureAllConnected()
	result, err := n.Process()
	return &result, err
}
