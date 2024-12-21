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

package paillier

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	config   = NewS256()
	S256N    = new(big.Int).Set(config.Curve.Params().N)
	p0, _    = new(big.Int).SetString("104975615121222854384410219330480259027041155688835759631647658735069527864919393410352284436544267374160206678331198777612866309766581999589789442827625308608614590850591998897357449886061863686453412019330757447743487422636807387508460941025550338019105820406950462187693188000168607236389735877001362796259", 10)
	q0, _    = new(big.Int).SetString("102755306389915984635356782597494195047102560555160692696207839728487252530690043689166546890155633162017964085393843240989395317546293846694693801865924045225783240995686020308553449158438908412088178393717793204697268707791329981413862246773904710409946848630083569401668855899757371993960961231481357354607", 10)
	n0       = new(big.Int).Mul(p0, q0)
	n0Square = new(big.Int).Exp(n0, big2, nil)
	p1, _    = new(big.Int).SetString("153358525493066047718272004328038648786913482509822520103254406755798143949604410622032791793048759895858718956427015341713113638900522567684085899778510408923255203608960957421132295490505636485549184519743624924400108443180927830685150312863987140912921698352623279965264898212539356851473952484018512691291", 10)
	q1, _    = new(big.Int).SetString("151738703367302097684435199030265294883567293365147694453436673424842637513382760820514891602796438794348681340963592185873039318978378212962222259077236055424192721882122757948608565796721667337602583758461213724029369212314671893322455516989072919317378126757038418828801663436111527900407845693045582540363", 10)
	n1       = new(big.Int).Mul(p1, q1)
	n1Square = new(big.Int).Exp(n1, big2, nil)
	ssIDInfo = []byte("Mark HaHa")
	pedp, _  = new(big.Int).SetString("172321190316317406041983369591732729491350806968006943303929709788136215251460267633420533682689046013587054841341976463526601587002102302546652907431187846060997247514915888514444763709031278321293105031395914163838109362462240334430371455027991864100292721059079328191363601847674802011142994248364894749407", 10)
	pedq, _  = new(big.Int).SetString("133775161118873760646458598449594229708046435932335011961444226591456542241216521727451860331718305184791260558214309464515443345834395848652314690639803964821534655704923535199917670451716761498957904445631495169583566095296670783502280310288116580525460451464561679063318393570545894032154226243881186182059", 10)
	pedN     = new(big.Int).Mul(pedp, pedq)
	pedt     = big.NewInt(9)
	peds     = big.NewInt(729)
	ped      = NewPedersenOpenParameter(pedN, peds, pedt)
)

var _ = Describe("Operation commitment test", func() {
	var x, y, rhox, rhoy, rho, C, X, Y, D *big.Int
	x = big.NewInt(3)
	y = big.NewInt(5)
	rhox = big.NewInt(555)
	rhoy = big.NewInt(101)
	rho = big.NewInt(103)
	C = big.NewInt(108)
	X = new(big.Int).Mul(new(big.Int).Exp(new(big.Int).Add(big1, n1), x, n1Square), new(big.Int).Exp(rhox, n1, n1Square))
	Y = new(big.Int).Mul(new(big.Int).Exp(new(big.Int).Add(big1, n1), y, n1Square), new(big.Int).Exp(rhoy, n1, n1Square))
	Y.Mod(Y, n1Square)
	D = new(big.Int).Exp(C, x, n0Square)
	D.Mul(D, new(big.Int).Exp(new(big.Int).Add(big1, n0), y, n0Square))
	D.Mul(D, new(big.Int).Exp(rho, n0, n0Square))
	D.Mod(D, n0Square)
	n0 = new(big.Int).Mul(p0, q0)
	n1 = new(big.Int).Mul(p1, q1)
	pedN = new(big.Int).Mul(pedp, pedq)
	Context("NewPaillierOperationAndPaillierCommitment tests", func() {
		BeforeEach(func() {
			config = NewS256()
		})

		It("over Range, should be ok", func() {
			zkproof, err := NewPaillierOperationAndPaillierCommitment(config, ssIDInfo, x, y, rho, rhox, rhoy, n0, n1, X, Y, C, D, ped)
			Expect(err).Should(BeNil())
			err = zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).Should(BeNil())
		})

		It("wrong config", func() {
			config.TwoExpLAddepsilon = big.NewInt(-1)
			zkproof, err := NewPaillierOperationAndPaillierCommitment(config, ssIDInfo, x, y, rho, rhox, rhoy, n0, n1, X, Y, C, D, ped)
			Expect(err).ShouldNot(BeNil())
			Expect(zkproof).Should(BeNil())
		})

		It("wrong config", func() {
			config.TwoExpLpaiAddepsilon = big.NewInt(-1)
			zkproof, err := NewPaillierOperationAndPaillierCommitment(config, ssIDInfo, x, y, rho, rhox, rhoy, n0, n1, X, Y, C, D, ped)
			Expect(err).ShouldNot(BeNil())
			Expect(zkproof).Should(BeNil())
		})

		It("wrong config", func() {
			config.TwoExpL = big.NewInt(-1)
			zkproof, err := NewPaillierOperationAndPaillierCommitment(config, ssIDInfo, x, y, rho, rhox, rhoy, n0, n1, X, Y, C, D, ped)
			Expect(err).ShouldNot(BeNil())
			Expect(zkproof).Should(BeNil())
		})

		It("negative n0", func() {
			n0 = big.NewInt(-1)
			zkproof, err := NewPaillierOperationAndPaillierCommitment(config, ssIDInfo, x, y, rho, rhox, rhoy, n0, n1, X, Y, C, D, ped)
			Expect(err).ShouldNot(BeNil())
			Expect(zkproof).Should(BeNil())
		})

		It("negative n1", func() {
			n1 = big.NewInt(-1)
			zkproof, err := NewPaillierOperationAndPaillierCommitment(config, ssIDInfo, x, y, rho, rhox, rhoy, n0, n1, X, Y, C, D, ped)
			Expect(err).ShouldNot(BeNil())
			Expect(zkproof).Should(BeNil())
		})

		It("negative ped", func() {
			ped.n = big.NewInt(-1)
			zkproof, err := NewPaillierOperationAndPaillierCommitment(config, ssIDInfo, x, y, rho, rhox, rhoy, n0, n1, X, Y, C, D, ped)
			Expect(err).ShouldNot(BeNil())
			Expect(zkproof).Should(BeNil())
		})
	})

	Context("Verify tests", func() {
		var zkproof *PaillierOperationAndCommitmentMessage
		BeforeEach(func() {
			config = NewS256()
			config.Curve.Params().N = S256N
			var err error
			n0 = new(big.Int).Mul(p0, q0)
			n0Square = new(big.Int).Exp(n0, big2, nil)
			n1 = new(big.Int).Mul(p1, q1)
			n1Square = new(big.Int).Exp(n1, big2, nil)
			ped = NewPedersenOpenParameter(pedN, peds, pedt)
			zkproof, err = NewPaillierOperationAndPaillierCommitment(config, ssIDInfo, x, y, rho, rhox, rhoy, n0, n1, X, Y, C, D, ped)
			Expect(err).Should(BeNil())
		})

		It("wrong range", func() {
			zkproof.S = pedN.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("not coprime", func() {
			zkproof.S = pedp.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("wrong range", func() {
			zkproof.T = pedN.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("not coprime", func() {
			zkproof.T = pedp.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("wrong range", func() {
			zkproof.A = n0Square.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("not coprime", func() {
			zkproof.A = p0.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("wrong range", func() {
			zkproof.By = n1Square.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("not coprime", func() {
			zkproof.By = p1.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("wrong range", func() {
			zkproof.E = pedN.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("not coprime", func() {
			zkproof.E = pedp.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("wrong range", func() {
			zkproof.F = pedN.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("not coprime", func() {
			zkproof.F = pedq.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("wrong range", func() {
			zkproof.W = n0.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("not coprime", func() {
			zkproof.W = p0.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("wrong range", func() {
			zkproof.Wy = n1.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("not coprime", func() {
			zkproof.Wy = q1.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("wrong range", func() {
			zkproof.Bx = n1Square.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("not coprime", func() {
			zkproof.Bx = q1.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("wrong range", func() {
			zkproof.Wx = n1.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("wrong fieldOrder", func() {
			config.Curve.Params().N = big1
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("not coprime", func() {
			zkproof.Wx = q1.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("verify failure", func() {
			zkproof.A = big1.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("verify failure", func() {
			zkproof.Wx = big1.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("verify failure", func() {
			zkproof.Wy = big1.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("verify failure", func() {
			zkproof.Z1 = new(big.Int).Add(config.TwoExpLAddepsilon, config.TwoExpLAddepsilon).String()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("verify failure", func() {
			zkproof.Z2 = new(big.Int).Add(config.TwoExpLpaiAddepsilon, config.TwoExpLpaiAddepsilon).String()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("verify failure", func() {
			zkproof.Z3 = big0.String()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("verify failure", func() {
			zkproof.Z4 = big0.String()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})

		It("verify failure", func() {
			zkproof.F = big1.Bytes()
			err := zkproof.Verify(config, ssIDInfo, n0, n1, C, D, X, Y, ped)
			Expect(err).ShouldNot(BeNil())
		})
	})
})
