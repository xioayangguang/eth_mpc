package impl

import (
	"crypto/ecdsa"
	"errors"
	"eth_mpc/config"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/getamis/alice/crypto/birkhoffinterpolation"
	"github.com/getamis/alice/crypto/ecpointgrouplaw"
	"github.com/getamis/alice/crypto/elliptic"
	"github.com/getamis/alice/crypto/tss/dkg"
	"github.com/getamis/sirius/log"
	"math/big"
)

var (
	ErrConversion = errors.New("conversion error")
)

func GetCurve() elliptic.Curve {
	return elliptic.Secp256k1()
}

func ConvertDKGResult(cfgPubkey config.Pubkey, cfgShare string, cfgBKs map[string]config.BK) (*dkg.Result, error) {
	x, ok := new(big.Int).SetString(cfgPubkey.X, 10)
	if !ok {
		log.Error("Cannot convert string to big int", "x", cfgPubkey.X)
		return nil, ErrConversion
	}
	y, ok := new(big.Int).SetString(cfgPubkey.Y, 10)
	if !ok {
		log.Error("Cannot convert string to big int", "y", cfgPubkey.Y)
		return nil, ErrConversion
	}
	pubkey, err := ecpointgrouplaw.NewECPoint(GetCurve(), x, y)
	if err != nil {
		log.Error("Cannot get public key", "err", err)
		return nil, err
	}
	share, ok := new(big.Int).SetString(cfgShare, 10)
	if !ok {
		log.Error("Cannot convert string to big int", "share", share)
		return nil, ErrConversion
	}

	dkgResult := &dkg.Result{
		PublicKey: pubkey,
		Share:     share,
		Bks:       make(map[string]*birkhoffinterpolation.BkParameter),
	}

	for peerID, bk := range cfgBKs {
		x, ok := new(big.Int).SetString(bk.X, 10)
		if !ok {
			log.Error("Cannot convert string to big int", "x", bk.X)
			return nil, ErrConversion
		}
		dkgResult.Bks[peerID] = birkhoffinterpolation.NewBkParameter(x, bk.Rank)
	}

	return dkgResult, nil
}

func DecodeLondonSignature(sig []byte) (r, s, v *big.Int) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64]})
	return r, s, v
}

func EncodeLondonSignature(r, s, v *big.Int) []byte {
	sig := make([]byte, crypto.SignatureLength)
	rBytes := r.Bytes()
	copy(sig[:32-len(rBytes)], make([]byte, 32-len(rBytes))) // 填充前导零（如果需要）
	copy(sig[:32], rBytes)
	sBytes := s.Bytes()
	copy(sig[32:64-len(sBytes)], make([]byte, 32-len(sBytes))) // 填充前导零（如果需要）
	copy(sig[32:64], sBytes)
	vByte := make([]byte, 1)
	vByte[0] = byte(v.Int64())
	copy(sig[64:], vByte)
	return sig
}

func AddressFromXY(x, y string) (common.Address, error) {
	xPoint, ok := new(big.Int).SetString(x, 10)
	if !ok {
		return common.Address{}, errors.New("cannot convert string to big int")
	}
	yPoint, ok := new(big.Int).SetString(y, 10)
	if !ok {
		return common.Address{}, errors.New("cannot convert string to big int")
	}
	pub := ecdsa.PublicKey{Curve: crypto.S256(), X: xPoint, Y: yPoint}
	return crypto.PubkeyToAddress(pub), nil
}
