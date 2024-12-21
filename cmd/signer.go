package cmd

import (
	"encoding/base64"
	"eth_mpc/config"
	"eth_mpc/impl"
	"fmt"
	"github.com/getamis/alice/crypto/homo/paillier"
	"github.com/getamis/alice/crypto/tss/ecdsa/gg18/signer"
	"github.com/getamis/sirius/log"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const signerProtocol = "/signer/1.0.0"

var SignerCmd = &cobra.Command{
	Use:  "signer",
	Long: `Signing for using the secret shares to generate a signature.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		signerResult, err := Signer(signerProtocol, []byte("hello"))
		rawResult, _ := yaml.Marshal(signerResult)
		fmt.Println(string(rawResult))
		return err
	},
}

func Signer(dkgProtocol protocol.ID, msg []byte) (*config.SignerResult, error) {
	result, err := common(dkgProtocol, func(pm *impl.PeerManager, l impl.Listener) (*impl.Node[*signer.Message, *signer.Result], error) {
		dkgResult, err := impl.ConvertDKGResult(config.Cfg.Pubkey, config.Cfg.Share, config.Cfg.BKs)
		if err != nil {
			log.Warn("Cannot get DKG result", "err", err)
			return nil, err
		}
		p, err := paillier.NewPaillier(2048)
		if err != nil {
			log.Warn("Cannot create a paillier function", "err", err)
			return nil, err
		}
		signerCore, err := signer.NewSigner(pm, dkgResult.PublicKey, p, dkgResult.Share, dkgResult.Bks, msg, l)
		if err != nil {
			log.Warn("Cannot create a new signer", "err", err)
			return nil, err
		}
		node := impl.New[*signer.Message, *signer.Result](signerCore, l)
		if err != nil {
			log.Crit("Failed to new service", "err", err)
		}
		return node, nil
	})
	if err != nil {
		return nil, err
	}
	return &config.SignerResult{
		Sign: base64.StdEncoding.EncodeToString((*result).EthSignature()),
	}, nil
}
