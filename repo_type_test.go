package aptlib

import "testing"

func TestRepoType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		rt   RepoType
		want bool
	}{
		{
			name: "binary repository type",
			rt:   "deb",
			want: true,
		},
		{
			name: "source repository type",
			rt:   "deb-src",
			want: true,
		},
		{
			name: "invalid repository type",
			rt:   "zzz",
			want: false,
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
