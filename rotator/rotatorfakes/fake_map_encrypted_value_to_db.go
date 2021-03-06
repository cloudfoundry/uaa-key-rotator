// Code generated by counterfeiter. DO NOT EDIT.
package rotatorfakes

import (
	"sync"

	"github.com/cloudfoundry/uaa-key-rotator/crypto"
	"github.com/cloudfoundry/uaa-key-rotator/rotator"
)

type FakeMapEncryptedValueToDB struct {
	MapStub        func(value crypto.EncryptedValue) ([]byte, error)
	mapMutex       sync.RWMutex
	mapArgsForCall []struct {
		value crypto.EncryptedValue
	}
	mapReturns struct {
		result1 []byte
		result2 error
	}
	mapReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	MapBase64ToCipherValueStub        func(value string) ([]byte, error)
	mapBase64ToCipherValueMutex       sync.RWMutex
	mapBase64ToCipherValueArgsForCall []struct {
		value string
	}
	mapBase64ToCipherValueReturns struct {
		result1 []byte
		result2 error
	}
	mapBase64ToCipherValueReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeMapEncryptedValueToDB) Map(value crypto.EncryptedValue) ([]byte, error) {
	fake.mapMutex.Lock()
	ret, specificReturn := fake.mapReturnsOnCall[len(fake.mapArgsForCall)]
	fake.mapArgsForCall = append(fake.mapArgsForCall, struct {
		value crypto.EncryptedValue
	}{value})
	fake.recordInvocation("Map", []interface{}{value})
	fake.mapMutex.Unlock()
	if fake.MapStub != nil {
		return fake.MapStub(value)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.mapReturns.result1, fake.mapReturns.result2
}

func (fake *FakeMapEncryptedValueToDB) MapCallCount() int {
	fake.mapMutex.RLock()
	defer fake.mapMutex.RUnlock()
	return len(fake.mapArgsForCall)
}

func (fake *FakeMapEncryptedValueToDB) MapArgsForCall(i int) crypto.EncryptedValue {
	fake.mapMutex.RLock()
	defer fake.mapMutex.RUnlock()
	return fake.mapArgsForCall[i].value
}

func (fake *FakeMapEncryptedValueToDB) MapReturns(result1 []byte, result2 error) {
	fake.MapStub = nil
	fake.mapReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeMapEncryptedValueToDB) MapReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.MapStub = nil
	if fake.mapReturnsOnCall == nil {
		fake.mapReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.mapReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeMapEncryptedValueToDB) MapBase64ToCipherValue(value string) ([]byte, error) {
	fake.mapBase64ToCipherValueMutex.Lock()
	ret, specificReturn := fake.mapBase64ToCipherValueReturnsOnCall[len(fake.mapBase64ToCipherValueArgsForCall)]
	fake.mapBase64ToCipherValueArgsForCall = append(fake.mapBase64ToCipherValueArgsForCall, struct {
		value string
	}{value})
	fake.recordInvocation("MapBase64ToCipherValue", []interface{}{value})
	fake.mapBase64ToCipherValueMutex.Unlock()
	if fake.MapBase64ToCipherValueStub != nil {
		return fake.MapBase64ToCipherValueStub(value)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.mapBase64ToCipherValueReturns.result1, fake.mapBase64ToCipherValueReturns.result2
}

func (fake *FakeMapEncryptedValueToDB) MapBase64ToCipherValueCallCount() int {
	fake.mapBase64ToCipherValueMutex.RLock()
	defer fake.mapBase64ToCipherValueMutex.RUnlock()
	return len(fake.mapBase64ToCipherValueArgsForCall)
}

func (fake *FakeMapEncryptedValueToDB) MapBase64ToCipherValueArgsForCall(i int) string {
	fake.mapBase64ToCipherValueMutex.RLock()
	defer fake.mapBase64ToCipherValueMutex.RUnlock()
	return fake.mapBase64ToCipherValueArgsForCall[i].value
}

func (fake *FakeMapEncryptedValueToDB) MapBase64ToCipherValueReturns(result1 []byte, result2 error) {
	fake.MapBase64ToCipherValueStub = nil
	fake.mapBase64ToCipherValueReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeMapEncryptedValueToDB) MapBase64ToCipherValueReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.MapBase64ToCipherValueStub = nil
	if fake.mapBase64ToCipherValueReturnsOnCall == nil {
		fake.mapBase64ToCipherValueReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.mapBase64ToCipherValueReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeMapEncryptedValueToDB) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.mapMutex.RLock()
	defer fake.mapMutex.RUnlock()
	fake.mapBase64ToCipherValueMutex.RLock()
	defer fake.mapBase64ToCipherValueMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeMapEncryptedValueToDB) recordInvocation(key string, args []interface{}) {
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

var _ rotator.MapEncryptedValueToDB = new(FakeMapEncryptedValueToDB)
