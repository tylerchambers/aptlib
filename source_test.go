package aptlib

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
			name: "new source from valid source file entry",
			args: args{
				sourceEntry: "deb http://us.archive.ubuntu.com/ubuntu/ focal main restricted",
			},
			want:    &SourceEntry{"deb", nil, "http://us.archive.ubuntu.com/ubuntu/", "focal", []string{"main", "restricted"}},
			wantErr: false,
		},
		{
			name: "new source from valid source file entry including options",
			args: args{
				sourceEntry: "deb [arch=amd64 hello=world] http://dl.google.com/linux/chrome/deb/ stable main",
			},
			want:    &SourceEntry{"deb", map[string]string{"arch": "amd64", "hello": "world"}, "http://dl.google.com/linux/chrome/deb/", "stable", []string{"main"}},
			wantErr: false,
		},
		{
			name: "invalid options (empty)",
			args: args{
				sourceEntry: "deb [] http://dl.google.com/linux/chrome/deb/ stable main",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid options (too short)",
			args: args{
				sourceEntry: "deb [aa] http://dl.google.com/linux/chrome/deb/ stable main",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid option (no key value seperated by equals)",
			args: args{
				sourceEntry: "deb [foo=bar baz] http://dl.google.com/linux/chrome/deb/ stable main",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid option (nested brackets)",
			args: args{
				sourceEntry: "deb [[]] http://dl.google.com/linux/chrome/deb/ stable main",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "new source from source file entry with invalid type",
			args: args{
				sourceEntry: "zzz http://us.archive.ubuntu.com/ubuntu/ focal main restricted",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "new source from source file entry with invalid URI",
			args: args{
				sourceEntry: "deb asdfasdfasdf focal main restricted",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "new source from source file entry with invalid URI",
			args: args{
				sourceEntry: "deb asdfasdfasdf focal main restricted",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "new source from source file entry with too few fields",
			args: args{
				sourceEntry: "deb asdfasdfasdf focal",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "new source from blank line",
			args: args{
				sourceEntry: "",
			},
			want:    nil,
			wantErr: true,
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
