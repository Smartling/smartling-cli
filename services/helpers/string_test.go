package helpers

import (
	"reflect"
	"testing"
)

func TestKeyValueToMap(t *testing.T) {
	tests := []struct {
		str     string
		want    map[string]string
		wantErr bool
	}{
		{
			str:     `file_charset="utf-8"`,
			want:    map[string]string{"file_charset": "utf-8"},
			wantErr: false,
		},
		{
			str:     `file_charset = "utf-8"`,
			want:    map[string]string{"file_charset": "utf-8"},
			wantErr: false,
		},
		{
			str:     `file_charset = utf-8`,
			want:    map[string]string{"file_charset": "utf-8"},
			wantErr: false,
		},
		{
			str:     `file_charset utf-8`,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := KeyValueToMap(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyValueToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyValueToMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}
