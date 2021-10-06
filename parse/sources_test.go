package parse

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/tylerchambers/aptlib/repo"
)

func logAndFail(t *testing.T, got, want []*repo.SourceEntry) {
	t.Logf("SourcesList() = ")
	for _, v := range got {
		t.Logf("%v,\n", *v)
	}
	t.Logf("\nwant ")
	for _, v := range want {
		t.Logf("%v,\n", *v)
	}
	t.Fail()
}
func TestSourcesList(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []*repo.SourceEntry
		wantErr bool
	}{
		{
			name: "test valid sources.list",
			args: args{
				r: strings.NewReader("# deb-src http://security.ubuntu.com/ubuntu focal-security universe\n" +
					"deb http://security.ubuntu.com/ubuntu focal-security multiverse\n" +
					"deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu   focal stable test"),
			},
			want: []*repo.SourceEntry{
				{
					RepoType:     "deb",
					Options:      map[string]string{},
					URI:          "http://security.ubuntu.com/ubuntu",
					Distribution: "focal-security",
					Components:   []string{"multiverse"},
				},
				{
					RepoType: "deb",
					Options: map[string]string{
						"arch":      "amd64",
						"signed-by": "/usr/share/keyrings/docker-archive-keyring.gpg",
					},
					URI:          "https://download.docker.com/linux/ubuntu",
					Distribution: "focal",
					Components:   []string{"stable", "test"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SourcesList(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("SourcesList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				logAndFail(t, got, tt.want)
			}
			for i := range got {
				if !reflect.DeepEqual(got[i].Components, tt.want[i].Components) {
					logAndFail(t, got, tt.want)
				}
				if !reflect.DeepEqual(got[i].Distribution, tt.want[i].Distribution) {
					logAndFail(t, got, tt.want)
				}
				if len(got[i].Options) != len(tt.want[i].Options) {
					t.Fatalf("options mismatch: got %v want %v", got[i].Options, tt.want[i].Options)
				}
				for k, v := range got[i].Options {
					if v != tt.want[i].Options[k] {
						t.Fatalf("options mismatch: got %v want %v", got[i].Options, tt.want[i].Options)
					}
				}
				if !reflect.DeepEqual(got[i].RepoType, tt.want[i].RepoType) {
					logAndFail(t, got, tt.want)
				}
				if !reflect.DeepEqual(got[i].URI, tt.want[i].URI) {
					logAndFail(t, got, tt.want)
				}
			}
		})
	}
}
