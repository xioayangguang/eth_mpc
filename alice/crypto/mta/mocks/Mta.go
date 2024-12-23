// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import (
	big "math/big"

	ecpointgrouplaw "github.com/getamis/alice/crypto/ecpointgrouplaw"
	elliptic "github.com/getamis/alice/crypto/elliptic"

	homo "github.com/getamis/alice/crypto/homo"

	mock "github.com/stretchr/testify/mock"

	mta "github.com/getamis/alice/crypto/mta"

	zkproof "github.com/getamis/alice/crypto/zkproof"
)

// Mta is an autogenerated mock type for the Mta type
type Mta struct {
	mock.Mock
}

// Compute provides a mock function with given fields: publicKey, encMessage
func (_m *Mta) Compute(publicKey homo.Pubkey, encMessage []byte) (*big.Int, *big.Int, error) {
	ret := _m.Called(publicKey, encMessage)

	if len(ret) == 0 {
		panic("no return value specified for Compute")
	}

	var r0 *big.Int
	var r1 *big.Int
	var r2 error
	if rf, ok := ret.Get(0).(func(homo.Pubkey, []byte) (*big.Int, *big.Int, error)); ok {
		return rf(publicKey, encMessage)
	}
	if rf, ok := ret.Get(0).(func(homo.Pubkey, []byte) *big.Int); ok {
		r0 = rf(publicKey, encMessage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(homo.Pubkey, []byte) *big.Int); ok {
		r1 = rf(publicKey, encMessage)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*big.Int)
		}
	}

	if rf, ok := ret.Get(2).(func(homo.Pubkey, []byte) error); ok {
		r2 = rf(publicKey, encMessage)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Decrypt provides a mock function with given fields: c
func (_m *Mta) Decrypt(c *big.Int) (*big.Int, error) {
	ret := _m.Called(c)

	if len(ret) == 0 {
		panic("no return value specified for Decrypt")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*big.Int) (*big.Int, error)); ok {
		return rf(c)
	}
	if rf, ok := ret.Get(0).(func(*big.Int) *big.Int); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*big.Int) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAG provides a mock function with given fields: curve
func (_m *Mta) GetAG(curve elliptic.Curve) *ecpointgrouplaw.ECPoint {
	ret := _m.Called(curve)

	if len(ret) == 0 {
		panic("no return value specified for GetAG")
	}

	var r0 *ecpointgrouplaw.ECPoint
	if rf, ok := ret.Get(0).(func(elliptic.Curve) *ecpointgrouplaw.ECPoint); ok {
		r0 = rf(curve)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecpointgrouplaw.ECPoint)
		}
	}

	return r0
}

// GetAK provides a mock function with given fields:
func (_m *Mta) GetAK() *big.Int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAK")
	}

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// GetAProof provides a mock function with given fields: curve
func (_m *Mta) GetAProof(curve elliptic.Curve) (*zkproof.SchnorrProofMessage, error) {
	ret := _m.Called(curve)

	if len(ret) == 0 {
		panic("no return value specified for GetAProof")
	}

	var r0 *zkproof.SchnorrProofMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(elliptic.Curve) (*zkproof.SchnorrProofMessage, error)); ok {
		return rf(curve)
	}
	if rf, ok := ret.Get(0).(func(elliptic.Curve) *zkproof.SchnorrProofMessage); ok {
		r0 = rf(curve)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*zkproof.SchnorrProofMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(elliptic.Curve) error); ok {
		r1 = rf(curve)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEncK provides a mock function with given fields:
func (_m *Mta) GetEncK() []byte {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetEncK")
	}

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// GetProductWithK provides a mock function with given fields: v
func (_m *Mta) GetProductWithK(v *big.Int) *big.Int {
	ret := _m.Called(v)

	if len(ret) == 0 {
		panic("no return value specified for GetProductWithK")
	}

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(*big.Int) *big.Int); ok {
		r0 = rf(v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// GetProofWithCheck provides a mock function with given fields: curve, beta
func (_m *Mta) GetProofWithCheck(curve elliptic.Curve, beta *big.Int) ([]byte, error) {
	ret := _m.Called(curve, beta)

	if len(ret) == 0 {
		panic("no return value specified for GetProofWithCheck")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(elliptic.Curve, *big.Int) ([]byte, error)); ok {
		return rf(curve, beta)
	}
	if rf, ok := ret.Get(0).(func(elliptic.Curve, *big.Int) []byte); ok {
		r0 = rf(curve, beta)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(elliptic.Curve, *big.Int) error); ok {
		r1 = rf(curve, beta)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetResult provides a mock function with given fields: alphas, betas
func (_m *Mta) GetResult(alphas []*big.Int, betas []*big.Int) (*big.Int, error) {
	ret := _m.Called(alphas, betas)

	if len(ret) == 0 {
		panic("no return value specified for GetResult")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func([]*big.Int, []*big.Int) (*big.Int, error)); ok {
		return rf(alphas, betas)
	}
	if rf, ok := ret.Get(0).(func([]*big.Int, []*big.Int) *big.Int); ok {
		r0 = rf(alphas, betas)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func([]*big.Int, []*big.Int) error); ok {
		r1 = rf(alphas, betas)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OverrideA provides a mock function with given fields: newA
func (_m *Mta) OverrideA(newA *big.Int) (mta.Mta, error) {
	ret := _m.Called(newA)

	if len(ret) == 0 {
		panic("no return value specified for OverrideA")
	}

	var r0 mta.Mta
	var r1 error
	if rf, ok := ret.Get(0).(func(*big.Int) (mta.Mta, error)); ok {
		return rf(newA)
	}
	if rf, ok := ret.Get(0).(func(*big.Int) mta.Mta); ok {
		r0 = rf(newA)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mta.Mta)
		}
	}

	if rf, ok := ret.Get(1).(func(*big.Int) error); ok {
		r1 = rf(newA)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyProofWithCheck provides a mock function with given fields: proof, curve, alpha
func (_m *Mta) VerifyProofWithCheck(proof []byte, curve elliptic.Curve, alpha *big.Int) (*ecpointgrouplaw.ECPoint, error) {
	ret := _m.Called(proof, curve, alpha)

	if len(ret) == 0 {
		panic("no return value specified for VerifyProofWithCheck")
	}

	var r0 *ecpointgrouplaw.ECPoint
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, elliptic.Curve, *big.Int) (*ecpointgrouplaw.ECPoint, error)); ok {
		return rf(proof, curve, alpha)
	}
	if rf, ok := ret.Get(0).(func([]byte, elliptic.Curve, *big.Int) *ecpointgrouplaw.ECPoint); ok {
		r0 = rf(proof, curve, alpha)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecpointgrouplaw.ECPoint)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte, elliptic.Curve, *big.Int) error); ok {
		r1 = rf(proof, curve, alpha)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMta creates a new instance of Mta. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMta(t interface {
	mock.TestingT
	Cleanup(func())
}) *Mta {
	mock := &Mta{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
