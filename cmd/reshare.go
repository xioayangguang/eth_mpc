package cmd

import (
	"eth_mpc/config"
	"eth_mpc/impl"
	"fmt"
	"github.com/getamis/alice/crypto/tss/ecdsa/gg18/reshare"
	"github.com/getamis/sirius/log"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const reshareProtocol = "/reshare/1.0.0"

var ReshareCmd = &cobra.Command{
	Use:  "reshare",
	Long: `Refresh the secret shares without changing the public key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		signerResult, err := Reshare(reshareProtocol)
		rawResult, _ := yaml.Marshal(signerResult)
		fmt.Println(string(rawResult))
		return err
	},
}

func Reshare(dkgProtocol protocol.ID) (*config.ReshareResult, error) {
	result, err := common(dkgProtocol, func(pm *impl.PeerManager, l impl.Listener) (*impl.Node[*reshare.Message, *reshare.Result], error) {
		dkgResult, err := impl.ConvertDKGResult(config.Cfg.Pubkey, config.Cfg.Share, config.Cfg.BKs)
		if err != nil {
			log.Warn("Cannot get DKG result", "err", err)
			return nil, err
		}
		reshareCore, err := reshare.NewReshare(pm, config.Cfg.Threshold, dkgResult.PublicKey, dkgResult.Share, dkgResult.Bks, l)
		if err != nil {
			log.Warn("Cannot create a new reshare core", "err", err)
			return nil, err
		}
		node := impl.New[*reshare.Message, *reshare.Result](reshareCore, l)
		if err != nil {
			log.Crit("Failed to new service", "err", err)
		}
		return node, nil
	})
	if err != nil {
		return nil, err
	}
	return &config.ReshareResult{
		Share: (*result).Share.String(),
	}, nil
}
