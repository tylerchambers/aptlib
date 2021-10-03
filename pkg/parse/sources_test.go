package parse

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/tylerchambers/goapt/pkg/repo"
)

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
			name: ,
			args: args{},
			want: &[]*repo.SourceEntry{"deb-src"},
			wantErr: false,
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SourcesList(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("SourcesList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SourcesList() = %v, want %v", got, tt.want)
			}
		})
	}
}
