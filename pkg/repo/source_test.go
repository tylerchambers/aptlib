package repo

import (
	"reflect"
	"testing"
)

func TestNewSource(t *testing.T) {
	type args struct {
		rt         RepoType
		options    map[string]string
		URI        string
		dist       string
		components []string
	}
	tests := []struct {
		name string
		args args
		want *SourceEntry
	}{
		{"new source test",
			args{"deb", map[string]string{"arch": "amd64"}, "http://us.archive.ubuntu.com/ubuntu/", "focal", []string{"contrib", "main"}},
			&SourceEntry{"deb", map[string]string{"arch": "amd64"}, "http://us.archive.ubuntu.com/ubuntu/", "focal", []string{"contrib", "main"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSource(tt.args.rt, tt.args.options, tt.args.URI, tt.args.dist, tt.args.components); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSourceFromString(t *testing.T) {
	type args struct {
		sourceEntry string
	}
	tests := []struct {
		name    string
		args    args
		want    *SourceEntry
		wantErr bool
	}{
		{
			"new source from valid source file entry",
			args{"deb http://us.archive.ubuntu.com/ubuntu/ focal main restricted"},
			&SourceEntry{"deb", nil, "http://us.archive.ubuntu.com/ubuntu/", "focal", []string{"main", "restricted"}},
			false,
		},
		{
			"new source from valid source file entry including options",
			args{"deb [arch=amd64 hello=world] http://dl.google.com/linux/chrome/deb/ stable main"},
			&SourceEntry{"deb", map[string]string{"arch": "amd64", "hello": "world"}, "http://dl.google.com/linux/chrome/deb/", "stable", []string{"main"}},
			false,
		},
		{
			"invalid options (empty)",
			args{"deb [] http://dl.google.com/linux/chrome/deb/ stable main"},
			nil,
			true,
		},
		{
			"invalid options (too short)",
			args{"deb [aa] http://dl.google.com/linux/chrome/deb/ stable main"},
			nil,
			true,
		},
		{
			"invalid option (no key value seperated by equals)",
			args{"deb [foo=bar baz] http://dl.google.com/linux/chrome/deb/ stable main"},
			nil,
			true,
		},
		{
			"invalid option (nested brackets)",
			args{"deb [[]] http://dl.google.com/linux/chrome/deb/ stable main"},
			nil,
			true,
		},
		{
			"new source from source file entry with invalid type",
			args{"zzz http://us.archive.ubuntu.com/ubuntu/ focal main restricted"},
			nil,
			true,
		},
		{
			"new source from source file entry with invalid URI",
			args{"deb asdfasdfasdf focal main restricted"},
			nil,
			true,
		},
		{
			"new source from source file entry with invalid URI",
			args{"deb asdfasdfasdf focal main restricted"},
			nil,
			true,
		},
		{
			"new source from source file entry with too few fields",
			args{"deb asdfasdfasdf focal"},
			nil,
			true,
		},
		{
			"new source from blank line",
			args{""},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SourceFromString(tt.args.sourceEntry)
			if (err != nil) != tt.wantErr {
				t.Errorf("SourceFromEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SourceFromEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}
