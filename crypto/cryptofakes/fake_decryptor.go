// Code generated by counterfeiter. DO NOT EDIT.
package cryptofakes

import (
	"sync"

	"github.com/cloudfoundry/uaa-key-rotator/crypto"
)

type FakeDecryptor struct {
	DecryptStub        func(encryptedValue crypto.EncryptedValue) (string, error)
	decryptMutex       sync.RWMutex
	decryptArgsForCall []struct {
		encryptedValue crypto.EncryptedValue
	}
	decryptReturns struct {
		result1 string
		result2 error
	}
	decryptReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDecryptor) Decrypt(encryptedValue crypto.EncryptedValue) (string, error) {
	fake.decryptMutex.Lock()
	ret, specificReturn := fake.decryptReturnsOnCall[len(fake.decryptArgsForCall)]
	fake.decryptArgsForCall = append(fake.decryptArgsForCall, struct {
		encryptedValue crypto.EncryptedValue
	}{encryptedValue})
	fake.recordInvocation("Decrypt", []interface{}{encryptedValue})
	fake.decryptMutex.Unlock()
	if fake.DecryptStub != nil {
		return fake.DecryptStub(encryptedValue)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.decryptReturns.result1, fake.decryptReturns.result2
}

func (fake *FakeDecryptor) DecryptCallCount() int {
	fake.decryptMutex.RLock()
	defer fake.decryptMutex.RUnlock()
	return len(fake.decryptArgsForCall)
}

func (fake *FakeDecryptor) DecryptArgsForCall(i int) crypto.EncryptedValue {
	fake.decryptMutex.RLock()
	defer fake.decryptMutex.RUnlock()
	return fake.decryptArgsForCall[i].encryptedValue
}

func (fake *FakeDecryptor) DecryptReturns(result1 string, result2 error) {
	fake.DecryptStub = nil
	fake.decryptReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeDecryptor) DecryptReturnsOnCall(i int, result1 string, result2 error) {
	fake.DecryptStub = nil
	if fake.decryptReturnsOnCall == nil {
		fake.decryptReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.decryptReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeDecryptor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.decryptMutex.RLock()
	defer fake.decryptMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeDecryptor) recordInvocation(key string, args []interface{}) {
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

var _ crypto.Decryptor = new(FakeDecryptor)
