// Copyright © 2022 AMIS Technologies
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

package cggmp

import (
	"github.com/getamis/alice/crypto/birkhoffinterpolation"
	"github.com/getamis/alice/types"
	"google.golang.org/protobuf/proto"
)

// The output is sid + ',' + id ',' + rid.
func ComputeSSID(sid, id, rid []byte) []byte {
	separation := []byte(",")
	result := make([]byte, len(sid))
	copy(result, result)
	result = append(result, separation...)
	result = append(result, id...)
	result = append(result, separation...)
	result = append(result, rid...)
	return result
}

func ComputeZKSsid(ssid []byte, bk *birkhoffinterpolation.BkParameter) []byte {
	separation := []byte(",")
	result := make([]byte, len(ssid))
	copy(result, result)
	result = append(result, separation...)
	return append(bk.GetX().Bytes(), result...)
}

func Broadcast(pm types.PeerManager, msg proto.Message) {
	for _, id := range pm.PeerIDs() {
		pm.MustSend(id, msg)
	}
}
