// Code generated by counterfeiter. DO NOT EDIT.
package cryptofakes

import (
	"sync"

	"github.com/cloudfoundry/uaa-key-rotator/crypto"
)

type FakeCipherAccessor struct {
	GetCipherStub        func([]byte) ([]byte, error)
	getCipherMutex       sync.RWMutex
	getCipherArgsForCall []struct {
		arg1 []byte
	}
	getCipherReturns struct {
		result1 []byte
		result2 error
	}
	getCipherReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCipherAccessor) GetCipher(arg1 []byte) ([]byte, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.getCipherMutex.Lock()
	ret, specificReturn := fake.getCipherReturnsOnCall[len(fake.getCipherArgsForCall)]
	fake.getCipherArgsForCall = append(fake.getCipherArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	fake.recordInvocation("GetCipher", []interface{}{arg1Copy})
	fake.getCipherMutex.Unlock()
	if fake.GetCipherStub != nil {
		return fake.GetCipherStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getCipherReturns.result1, fake.getCipherReturns.result2
}

func (fake *FakeCipherAccessor) GetCipherCallCount() int {
	fake.getCipherMutex.RLock()
	defer fake.getCipherMutex.RUnlock()
	return len(fake.getCipherArgsForCall)
}

func (fake *FakeCipherAccessor) GetCipherArgsForCall(i int) []byte {
	fake.getCipherMutex.RLock()
	defer fake.getCipherMutex.RUnlock()
	return fake.getCipherArgsForCall[i].arg1
}

func (fake *FakeCipherAccessor) GetCipherReturns(result1 []byte, result2 error) {
	fake.GetCipherStub = nil
	fake.getCipherReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeCipherAccessor) GetCipherReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetCipherStub = nil
	if fake.getCipherReturnsOnCall == nil {
		fake.getCipherReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getCipherReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeCipherAccessor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getCipherMutex.RLock()
	defer fake.getCipherMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCipherAccessor) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ crypto.CipherAccessor = new(FakeCipherAccessor)
