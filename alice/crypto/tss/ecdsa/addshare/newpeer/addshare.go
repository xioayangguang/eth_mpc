// Copyright © 2020 AMIS Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package newpeer

import (
	"math/big"

	"github.com/getamis/alice/crypto/birkhoffinterpolation"
	"github.com/getamis/alice/crypto/ecpointgrouplaw"
	"github.com/getamis/alice/crypto/tss"
	"github.com/getamis/alice/crypto/tss/ecdsa/addshare"
	"github.com/getamis/alice/types"
	"github.com/getamis/alice/types/message"
	"github.com/getamis/sirius/log"
)

type AddShare struct {
	ph *peerHandler
	*message.MsgMain
}

type Result struct {
	PartialPublicKeys map[string]*ecpointgrouplaw.ECPoint
	PublicKey         *ecpointgrouplaw.ECPoint
	Share             *big.Int
	Bks               map[string]*birkhoffinterpolation.BkParameter
}

func NewAddShare(peerManager types.PeerManager, pubkey *ecpointgrouplaw.ECPoint, threshold, newPeerRank uint32, listener types.StateChangedListener) *AddShare {
	ph := newPeerHandler(peerManager, pubkey, threshold, newPeerRank)
	return &AddShare{
		ph:      ph,
		MsgMain: message.NewMsgMain(peerManager.SelfID(), peerManager.NumPeers(), listener, ph, types.MessageType(addshare.Type_OldPeer), types.MessageType(addshare.Type_Result)),
	}
}

// GetResult returns the final result: public key, share, bks (including self bk)
func (a *AddShare) GetResult() (*Result, error) {
	if a.GetState() != types.StateDone {
		return nil, tss.ErrNotReady
	}

	h := a.GetHandler()
	rh, ok := h.(*resultHandler)
	if !ok {
		log.Error("We cannot convert to result handler in done state")
		return nil, tss.ErrNotReady
	}

	// Total bks = peer bks + self bk
	bks := make(map[string]*birkhoffinterpolation.BkParameter, a.ph.peerManager.NumPeers()+1)
	pks := make(map[string]*ecpointgrouplaw.ECPoint, a.ph.peerManager.NumPeers()+1)

	bks[a.ph.peerManager.SelfID()] = rh.bk
	pks[a.ph.peerManager.SelfID()] = rh.siG

	for id, peer := range a.ph.peers {
		bks[id] = peer.peer.bk
		pks[id] = peer.peer.siG
	}

	return &Result{
		PartialPublicKeys: pks,
		PublicKey:         rh.pubkey,
		Share:             rh.share,
		Bks:               bks,
	}, nil
}
