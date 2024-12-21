package cmd

import (
	"eth_mpc/config"
	"eth_mpc/impl"
	"fmt"
	"github.com/getamis/alice/crypto/tss/dkg"
	"github.com/getamis/sirius/log"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const dkgProtocol = "/dkg/1.0.0"

var DkgCmd = &cobra.Command{
	Use:  "dkg",
	Long: `Distributed key generation for creating secret shares without any dealer.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dkgResult, err := Dkg(dkgProtocol)
		rawResult, _ := yaml.Marshal(dkgResult)
		fmt.Println(string(rawResult))
		return err
	},
}

func Dkg(dkgProtocol protocol.ID) (*config.DKGResult, error) {
	result, err := common(dkgProtocol, func(pm *impl.PeerManager, l impl.Listener) (*impl.Node[*dkg.Message, *dkg.Result], error) {
		dkgCore, err := dkg.NewDKG(impl.GetCurve(), pm, config.Cfg.Threshold, config.Cfg.Rank, l)
		if err != nil {
			log.Warn("Cannot create a new DKG", "config", config.Cfg, "err", err)
			return nil, err
		}
		n := impl.New[*dkg.Message, *dkg.Result](dkgCore, l)
		return n, nil
	})
	if err != nil {
		return nil, err
	}
	dkgResult := &config.DKGResult{
		Share: (*result).Share.String(),
		Pubkey: config.Pubkey{
			X: (*result).PublicKey.GetX().String(),
			Y: (*result).PublicKey.GetY().String(),
		},
		BKs: make(map[string]config.BK),
	}
	for peerID, bk := range (*result).Bks {
		dkgResult.BKs[peerID] = config.BK{
			X:    bk.GetX().String(),
			Rank: bk.GetRank(),
		}
	}
	return dkgResult, nil
}
