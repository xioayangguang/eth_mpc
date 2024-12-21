package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"eth_mpc/config"
	"eth_mpc/impl"
	"fmt"
	"github.com/getamis/sirius/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
)

const notifyProtocol = "/notify/1.0.0"

var ServerCmd = &cobra.Command{
	Use:  "server",
	Long: `server for using the secret shares to generate a signature.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := gin.Default()
		r.POST("/hello/", func(c *gin.Context) {
			var json = config.MonitorNotify{
				ProtocolId: uuid.New().String(),
			}
			if err := c.BindJSON(&json); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			result, err := Send(json)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}
			c.JSON(200, result)
		})
		return r.Run(fmt.Sprintf(":%v", config.Cfg.HttpPort))
	},
}

func MonitorP2p() error {
	host, err := impl.GetHostByCfg()
	if err != nil {
		return err
	}
	pm := impl.NewPeerManager(host, notifyProtocol)
	host.SetStreamHandler(notifyProtocol, func(s network.Stream) {
		buf, err := io.ReadAll(s)
		if err != nil {
			log.Warn("Cannot read data from stream", "err", err)
			return
		}
		_ = s.Close()
		var message config.MonitorNotify
		if err = json.Unmarshal(buf, &message); err != nil {
			log.Error("Cannot unmarshal data", "err", err)
			return
		} else {
			if result, err := call(message); err != nil {
				log.Error("Cannot unmarshal data", "err", err)
			} else {
				rawResult, _ := yaml.Marshal(result)
				fmt.Println(string(rawResult))
			}
		}
	})
	pm.EnsureAllConnected()
	return nil
}

func Send(message config.MonitorNotify) (interface{}, error) {
	host, err := impl.GetHostByCfg()
	if err != nil {
		return nil, err
	}
	pm := impl.NewPeerManager(host, notifyProtocol)
	pm.EnsureAllConnected()
	peers := pm.GetPeers()
	for _, peerAddr := range peers {
		maddr, err := multiaddr.NewMultiaddr(peerAddr)
		if err != nil {
			return nil, err
		}
		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			return nil, err
		}
		s, err := host.NewStream(context.Background(), info.ID, notifyProtocol)
		if err != nil {
			return nil, err
		}
		bs, err := json.Marshal(message)
		if err != nil {
			return nil, err
		}
		if _, err = s.Write(bs); err != nil {
			return nil, err
		}
		if err = s.Close(); err != nil {
			return nil, err
		}
	}
	return call(message)
}

func call(message config.MonitorNotify) (interface{}, error) {
	p2pProtocol := protocol.ID(fmt.Sprintf("/%v/%v", message.Type, message.ProtocolId))
	if message.Type == "dkg" {
		return Dkg(p2pProtocol)
	} else if message.Type == "reshare" {
		return Reshare(p2pProtocol)
	} else if message.Type == "signer" {
		msg, err := base64.StdEncoding.DecodeString(message.Msg)
		if err != nil {
			return nil, err
		}
		return Signer(p2pProtocol, msg)
	} else {
		return nil, errors.New("type that does not exist")
	}
}
