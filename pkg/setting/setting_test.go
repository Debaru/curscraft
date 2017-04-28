package setting_test

import (
	"testing"

	"github.com/Debaru/curscraft/pkg/setting"
)

var testCase = []struct {
	path      string
	wantError bool
}{
	{"test_data/Blizzard/World Of Warcraft/", false},
	{"test_data/Blizzard/World Of Warcraft", false},
	{"test_data/Bli", true},
	{"test_data/Blizzard/World Of Warcraft/Interface", true},
	{"test_data/Blizzard_T/World Of Warcraft", true},
	{"", true},
}

func TestSetPath(t *testing.T) {
	for _, tc := range testCase {
		t.Run(tc.path, func(t *testing.T) {
			var s setting.Settings
			err := s.SetPath(tc.path)

			if (err == nil && tc.wantError) || (err != nil && !tc.wantError) {
				t.Errorf("got %s; want %t", err, tc.wantError)
			}

		})
	}
}
