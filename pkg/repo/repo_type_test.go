package repo

import "testing"

func TestRepoType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		rt   RepoType
		want bool
	}{
		{"binary repository type",
			"deb",
			true,
		},
		{"source repository type",
			"deb-src",
			true,
		}, {"invalid repository type",
			"zzz",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rt.IsValid(); got != tt.want {
				t.Errorf("RepoType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
