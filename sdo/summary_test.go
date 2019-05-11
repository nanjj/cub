package sdo_test

import (
	"reflect"
	"testing"

	"github.com/nanjj/cub/sdo"
)

func TestSummaryClone(t *testing.T) {
	tcs := []struct {
		name    string
		summary sdo.Summary
	}{
		{"empty", sdo.Summary{}},
		{"ToEmpty", sdo.Summary{
			To: []string{},
		}},
		{"ToOne", sdo.Summary{
			To: []string{"one"},
		}},
		{"ToTwo", sdo.Summary{
			To: []string{"one", "two"},
		}},
		{"LensEmpty", sdo.Summary{
			Lens: []int{},
		}},
		{"LensOne", sdo.Summary{
			Lens: []int{1},
		}},
		{"LensTwo", sdo.Summary{
			Lens: []int{1},
		}},
		{"CarrierEmpty", sdo.Summary{
			Carrier: map[string]string{},
		}},
		{"CarrierOne", sdo.Summary{
			Carrier: map[string]string{"uber-id": "one"},
		}},
		{"CarrierTwo", sdo.Summary{
			Carrier: map[string]string{"uber-id": "one", "key": "value"},
		}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			want := tc.summary.Clone()
			if !reflect.DeepEqual(want, tc.summary) {
				t.Fatalf("%v\n%v", want, tc.summary)
			}
			modified := false
			if len(want.To) > 0 {
				want.To[0] = "nowhere"
				modified = true
			}
			if len(want.Lens) > 0 {
				want.Lens[0] = -1
				modified = true
			}
			if len(want.Carrier) > 0 {
				for k := range want.Carrier {
					want.Carrier[k] = "fake"
				}
				modified = true
			}
			if modified {
				if reflect.DeepEqual(want, tc.summary) {
					t.Fatalf("%v\n%v", want, tc.summary)
				}
			}
		})
	}
}
