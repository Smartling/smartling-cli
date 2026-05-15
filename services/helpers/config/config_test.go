package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFromFile_Threads(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want uint32
	}{
		{
			name: "threads set as integer",
			yaml: "user_id: u\nsecret: s\nproject_id: p\nthreads: 42\n",
			want: 42,
		},
		{
			name: "threads absent defaults to zero",
			yaml: "user_id: u\nsecret: s\nproject_id: p\n",
			want: 0,
		},
		{
			name: "threads explicit zero",
			yaml: "user_id: u\nsecret: s\nthreads: 0\n",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "smartling.yml")
			if err := os.WriteFile(path, []byte(tt.yaml), 0o644); err != nil {
				t.Fatalf("setup: %v", err)
			}
			cfg, err := LoadConfigFromFile(path)
			if err != nil {
				t.Fatalf("LoadConfigFromFile: %v", err)
			}
			if cfg.Threads != tt.want {
				t.Errorf("Threads = %d, want %d", cfg.Threads, tt.want)
			}
		})
	}
}
