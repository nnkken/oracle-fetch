// Code generated by mockery v2.22.1. DO NOT EDIT.

package chainlink

import (
	big "math/big"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	mock "github.com/stretchr/testify/mock"
)

// MockChainLinkContract is an autogenerated mock type for the ChainLinkContract type
type MockChainLinkContract struct {
	mock.Mock
}

type MockChainLinkContract_Expecter struct {
	mock *mock.Mock
}

func (_m *MockChainLinkContract) EXPECT() *MockChainLinkContract_Expecter {
	return &MockChainLinkContract_Expecter{mock: &_m.Mock}
}

// LatestRoundData provides a mock function with given fields: opts
func (_m *MockChainLinkContract) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	ret := _m.Called(opts)

	var r0 struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	}
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	}, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	}); ok {
		r0 = rf(opts)
	} else {
		r0 = ret.Get(0).(struct {
			RoundId         *big.Int
			Answer          *big.Int
			StartedAt       *big.Int
			UpdatedAt       *big.Int
			AnsweredInRound *big.Int
		})
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockChainLinkContract_LatestRoundData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LatestRoundData'
type MockChainLinkContract_LatestRoundData_Call struct {
	*mock.Call
}

// LatestRoundData is a helper method to define mock.On call
//   - opts *bind.CallOpts
func (_e *MockChainLinkContract_Expecter) LatestRoundData(opts interface{}) *MockChainLinkContract_LatestRoundData_Call {
	return &MockChainLinkContract_LatestRoundData_Call{Call: _e.mock.On("LatestRoundData", opts)}
}

func (_c *MockChainLinkContract_LatestRoundData_Call) Run(run func(opts *bind.CallOpts)) *MockChainLinkContract_LatestRoundData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.CallOpts))
	})
	return _c
}

func (_c *MockChainLinkContract_LatestRoundData_Call) Return(_a0 struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, _a1 error) *MockChainLinkContract_LatestRoundData_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockChainLinkContract_LatestRoundData_Call) RunAndReturn(run func(*bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error)) *MockChainLinkContract_LatestRoundData_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockChainLinkContract interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockChainLinkContract creates a new instance of MockChainLinkContract. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockChainLinkContract(t mockConstructorTestingTNewMockChainLinkContract) *MockChainLinkContract {
	mock := &MockChainLinkContract{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
