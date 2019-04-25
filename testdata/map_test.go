package drilling

import (
	"math/rand"
	"sync"
	"testing"
)

var (
	testStringKey   = "abcdefghijlkmn"
	testStringKeyNo = "abcdefghijlkmnno"
	testMapString   = func() map[string]interface{} {
		m := map[string]interface{}{}
		for i := 0; i < 100000; i++ {
			m[RandomString(16)] = true
		}
		m[testStringKey] = true
		return m
	}()
	testIntKey   = 999999
	testIntKeyNo = -1
	testMapInt   = func() map[int]interface{} {
		m := map[int]interface{}{}
		for i := 0; i < 100000; i++ {
			m[rand.Intn(999999)] = true
		}
		m[testIntKey] = true
		return m
	}()
	testSyncMapString = func() sync.Map {
		m := sync.Map{}
		for i := 0; i < 100000; i++ {
			m.Store(RandomString(16), true)
		}
		m.Store(testStringKey, true)
		return m
	}()
)

func BenchmarkMap(b *testing.B) {
	b.Run("map[string]", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, ok := testMapString[testStringKey]; !ok {
				b.Fatal()
			}
		}
	})
	b.Run("map[int]", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, ok := testMapInt[testIntKey]; !ok {
				b.Fatal()
			}
		}
	})
	b.Run("map[string]No", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, ok := testMapString[testStringKeyNo]; ok {
				b.Fatal()
			}
		}
	})
	b.Run("map[int]No", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, ok := testMapInt[testIntKeyNo]; ok {
				b.Fatal()
			}
		}
	})
	b.Run("syncMapString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, ok := testSyncMapString.Load(testStringKey); !ok {
				b.Fatal()
			}
		}
	})
	b.Run("syncMapStringNo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, ok := testSyncMapString.Load(testStringKeyNo); ok {
				b.Fatal()
			}
		}
	})

}

func TestRandomString(t *testing.T) {
	t.Log(RandomString(10))
}
