package impl

import (
	"encoding/base64"
	"eth_mpc/config"
	"fmt"
	"github.com/getamis/sirius/log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
)

var basicHost host.Host

// MakeBasicHost creates a LibP2P host.
func MakeBasicHost(port int64, priv crypto.PrivKey) (host.Host, error) {
	if basicHost != nil {
		return basicHost, nil
	}
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port))
	if err != nil {
		return nil, err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(priv),
	}
	basicHost, err = libp2p.New(opts...)
	if err != nil {
		return nil, err
	}
	return basicHost, nil
}

func GetHostByCfg() (host.Host, error) {
	if basicHost != nil {
		return basicHost, nil
	}
	rawIdentity, err := base64.StdEncoding.DecodeString(config.Cfg.Identity)
	priv, err := crypto.UnmarshalPrivateKey(rawIdentity)
	if err != nil {
		log.Crit("Failed to unmarshal", "err", err)
	}
	basicHost, err = MakeBasicHost(config.Cfg.Port, priv)
	if err != nil {
		log.Crit("Failed to create a basic host", "err", err)
	}
	selfId := basicHost.ID().String()
	log.Debug("my ID", "id", selfId, "addr", basicHost.Addrs())
	return basicHost, nil
}

func GetPeerAddr(port int64, peerId string) string {
	return fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/p2p/%s", port, peerId)
}
