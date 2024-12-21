package impl

import (
	"context"
	"eth_mpc/config"
	"fmt"
	"github.com/getamis/alice/example/node"
	"github.com/getamis/sirius/log"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

type PeerManager struct {
	id       string
	host     host.Host
	protocol protocol.ID
	peers    map[string]string
	logger   log.Logger
}

var peerManagers = make(map[protocol.ID]*PeerManager)

func NewPeerManager(host host.Host, protocol protocol.ID) *PeerManager {
	if peerManagers[protocol] != nil {
		return peerManagers[protocol]
	}
	peerManagers[protocol] = &PeerManager{
		id:       host.ID().String(),
		host:     host,
		protocol: protocol,
		peers:    make(map[string]string),
		logger:   log.New(),
	}
	for _, p := range config.Cfg.Peers {
		peerManagers[protocol].AddPeer(p.Id, node.GetPeerAddr(p.Port, p.Id))
	}
	return peerManagers[protocol]
}

func (p *PeerManager) NumPeers() uint32 {
	// #nosec: G115: integer overflow conversion int -> uint32
	return uint32(len(p.peers))
}

func (p *PeerManager) SelfID() string {
	return p.id
}

func (p *PeerManager) PeerIDs() []string {
	ids := make([]string, len(p.peers))
	i := 0
	for id := range p.peers {
		ids[i] = id
		i++
	}
	return ids
}

func (p *PeerManager) MustSend(peerId string, message interface{}) {
	msg, ok := message.(proto.Message)
	if !ok {
		p.logger.Warn("invalid proto message")
		return
	}

	target := p.peers[peerId]
	maddr, err := multiaddr.NewMultiaddr(target)
	if err != nil {
		p.logger.Warn("Cannot parse the target address", "target", target, "err", err)
		return
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		p.logger.Warn("Cannot parse addr", "addr", maddr, "err", err)
		return
	}

	s, err := p.host.NewStream(context.Background(), info.ID, p.protocol)
	if err != nil {
		p.logger.Warn("Cannot create a new stream", "from", p.host.ID(), "to", info.ID, "protocol", p.protocol, "err", err)
		return
	}

	bs, err := proto.Marshal(msg)
	if err != nil {
		p.logger.Warn("Cannot marshal message", "err", err)
		return
	}

	if _, err = s.Write(bs); err != nil {
		p.logger.Warn("Cannot write message to IO", "err", err)
		return
	}

	if err = s.Close(); err != nil {
		p.logger.Warn("Cannot close the stream", "err", err)
		return
	}
	log.Debug("Sent message", "peer", target)
}

// EnsureAllConnected connects the host to specified peer and sends the message to it.
func (p *PeerManager) EnsureAllConnected() {
	var wg sync.WaitGroup

	connect := func(ctx context.Context, host host.Host, target string) error {
		maddr, err := multiaddr.NewMultiaddr(target)
		if err != nil {
			p.logger.Warn("Cannot parse the target address", "target", target, "err", err)
			return err
		}

		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			p.logger.Error("Cannot parse addr", "addr", maddr, "err", err)
			return err
		}

		err = host.Connect(ctx, *info)
		if err != nil {
			// log.Warn("Failed to connect to peer", "err", err)
			return err
		}

		for {
			protocols, err := host.Peerstore().GetProtocols(info.ID)
			if err != nil {
				log.Warn("Failed to get peer protocols", "err", err)
			}
			for _, protocol := range protocols {
				fmt.Printf("protocols:%v\r\n", protocols)
				fmt.Printf("p.protocol:%v\r\n", p.protocol)
				if protocol == p.protocol {
					return nil
				}
			}
			log.Debug("Waiting for peers")
			time.Sleep(50 * time.Millisecond)
		}
	}

	for _, peerAddr := range p.peers {
		wg.Add(1)
		addr := peerAddr
		go func() {
			defer wg.Done()
			logger := log.New("peer", addr)
			for {
				err := connect(context.Background(), p.host, addr)
				if err != nil {
					time.Sleep(50 * time.Millisecond)
					continue
				}
				logger.Info("Successfully connect to peer")
				return
			}
		}()
	}
	wg.Wait()
}

// AddPeer adds a peer to the peer list.
func (p *PeerManager) AddPeer(peerId string, peerAddr string) {
	p.peers[peerId] = peerAddr
}

// GetPeers get the peer list.
func (p *PeerManager) GetPeers() map[string]string {
	return p.peers
}
