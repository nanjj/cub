package sca

import (
	"reflect"
	"testing"
	"time"
)

func TestDataObjectString(t *testing.T) {
	tcs := []struct {
		s string
	}{
		{""},
		{"hello"},
	}
	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			d := DataObject{}
			if err := d.Encode(tc.s); err != nil {
				t.Fatal(err)
			}
			s := ""
			if err := d.Decode(&s); err != nil {
				t.Fatal(err)
			}
			if s != tc.s {
				t.Fatal(s, d)
			}
		})
	}
}

func BenchmarkDataObjectString(b *testing.B) {
	hello := "hello"
	b.Run("DataObject", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d := DataObject{}
			if err := d.Encode(hello); err != nil {
				b.Fatal(err)
			}
			s := ""
			if err := d.Decode(&s); err != nil {
				b.Fatal(err)
			}
			if s != hello {
				b.Fatal(s, d)
			}
		}
	})
	b.Run("[]byte", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d := []byte(hello)
			s := string(d)
			if s != hello {
				b.Fatal(s, d)
			}
		}
	})
}

func TestDataObjectBytes(t *testing.T) {
	tcs := []struct {
		b []byte
	}{
		{},
		{[]byte{}},
		{[]byte("hello")},
	}
	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			d := DataObject{}
			if err := d.Encode(tc.b); err != nil {
				t.Fatal(err)
			}
			b := []byte{}
			if err := d.Decode(&b); err != nil {
				t.Fatal(err)
			}
			if tc.b == nil {
				tc.b = []byte{}
			}
			if !reflect.DeepEqual(b, tc.b) {
				t.Log(b == nil)
				t.Log(tc.b == nil)
				t.Fatal(b, d)
			}
		})
	}
}
func TestDataObjectStruct(t *testing.T) {
	now := time.Now().Round(0)
	type TestCase struct {
		CreatedAt time.Time
		Name      string
		Id        int64
	}
	tcs := []TestCase{
		{now, "", 0},
	}
	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			d := DataObject{}
			if err := d.Encode(tc); err != nil {
				t.Fatal(err)
			}
			want := TestCase{}
			if err := d.Decode(&want); err != nil {
				t.Fatal(err)
			}
			if !want.CreatedAt.Equal(tc.CreatedAt) ||
				want.Name != tc.Name ||
				want.Id != tc.Id {
				t.Fatal(tc, want)
			}
		})
	}
}
