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

package commitment

import (
	"math/big"

	"github.com/getamis/alice/crypto/elliptic"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	pt "github.com/getamis/alice/crypto/ecpointgrouplaw"
)

var _ = Describe("message test", func() {
	DescribeTable("EcPoints()", func(curveType pt.EcPointMessage_Curve, curve elliptic.Curve) {
		var err error
		points := make([]*pt.EcPointMessage, 3)
		expected := make([]*pt.ECPoint, 3)
		for i := 0; i < 3; i++ {
			expected[i] = pt.ScalarBaseMult(curve, big.NewInt(int64(i)))
			Expect(err).Should(Succeed())
			points[i], err = expected[i].ToEcPointMessage()
			Expect(err).Should(Succeed())
		}
		commitment := PointCommitmentMessage{Points: points}

		got, err := commitment.EcPoints()
		Expect(err).Should(Succeed())
		for i, p := range got {
			Expect(p).Should(Equal(expected[i]))
		}
	},
		Entry("S256", pt.EcPointMessage_S256, elliptic.Secp256k1()),
	)

	It("invalid point", func() {
		commitment := PointCommitmentMessage{Points: []*pt.EcPointMessage{nil}}
		got, err := commitment.EcPoints()
		Expect(err).Should(Equal(pt.ErrInvalidPoint))
		Expect(got).Should(BeNil())
	})
})
