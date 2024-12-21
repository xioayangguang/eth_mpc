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

package refresh

import (
	"github.com/getamis/alice/types"
	"github.com/minio/blake2b-simd"
	"google.golang.org/protobuf/proto"
)

func (m *Message) IsValid() bool {
	switch m.Type {
	case Type_Round1:
		return m.GetRound1() != nil
	case Type_Round2:
		return m.GetRound2() != nil
	case Type_Round3:
		return m.GetRound3() != nil
	}
	return false
}

func (m *Message) GetMessageType() types.MessageType {
	return types.MessageType(m.Type)
}

func (m *Message) EchoHash() ([]byte, error) {
	echoMsg := m.GetEchoMessage()
	if echoMsg == nil {
		return nil, nil
	}
	// NOTE: there's an issue if there's a map field in the message
	// https://developers.google.com/protocol-buffers/docs/encoding#implications
	// Deterministic serialization only guarantees the same byte output for a particular binary.
	bs, err := proto.Marshal(echoMsg.(*Message))
	if err != nil {
		return nil, err
	}
	got := blake2b.Sum256(bs)
	return got[:], nil
}

func (m *Message) GetEchoMessage() types.Message {
	mm := &Message{
		Type: m.Type,
		Id:   m.Id,
	}
	switch m.Type {
	case Type_Round1:
		mm.Body = &Message_Round1{
			Round1: &Round1Msg{
				Commitment: m.GetRound1().GetCommitment(),
			},
		}
	}
	return nil
}
