package resolve

import (
	"testing"

	"github.com/spf13/pflag"
)

// newStringFlag returns a *pflag.Flag for a string flag with the given default
// value. If userValue is non-empty, the flag is treated as explicitly set
// (Changed=true) with that value.
func newStringFlag(t *testing.T, name, defaultValue, userValue string) *pflag.Flag {
	t.Helper()
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.String(name, defaultValue, "")
	if userValue != "" {
		if err := fs.Set(name, userValue); err != nil {
			t.Fatalf("flag.Set: %v", err)
		}
	}
	return fs.Lookup(name)
}

func strPtr(s string) *string { return &s }

func TestFallbackString(t *testing.T) {
	tests := []struct {
		name   string
		flag   *pflag.Flag
		param  StringParam
		envVar string // when non-empty, sets SMARTLING_CLI_<UPPER(FLAG)>=envVar for the test
		want   string
	}{
		{
			name:  "user-set flag wins over env, config, and default",
			flag:  newStringFlag(t, "threads", "20", "5"),
			param: StringParam{FlagName: "threads", Config: strPtr("99")},
			// SMARTLING_CLI_THREADS would be set, but flag.Changed=true short-circuits first.
			envVar: "42",
			want:   "5",
		},
		{
			name:   "env wins over config and default when flag not changed",
			flag:   newStringFlag(t, "threads", "20", ""),
			param:  StringParam{FlagName: "threads", Config: strPtr("99")},
			envVar: "42",
			want:   "42",
		},
		{
			name:  "config wins over default when flag not changed and env unset",
			flag:  newStringFlag(t, "threads", "20", ""),
			param: StringParam{FlagName: "threads", Config: strPtr("99")},
			want:  "99",
		},
		{
			name:  "flag default used when flag not changed, env unset, config nil",
			flag:  newStringFlag(t, "threads", "20", ""),
			param: StringParam{FlagName: "threads"},
			want:  "20",
		},
		{
			name:  "config used when flag is nil",
			flag:  nil,
			param: StringParam{FlagName: "threads", Config: strPtr("99")},
			want:  "99",
		},
		{
			name:  "empty string when flag is nil and config is nil",
			flag:  nil,
			param: StringParam{FlagName: "threads"},
			want:  "",
		},
		{
			name:   "env wins when flag is nil and env is set",
			flag:   nil,
			param:  StringParam{FlagName: "threads"},
			envVar: "42",
			want:   "42",
		},
		{
			name:  "config pointing to empty string still wins over default",
			flag:  newStringFlag(t, "threads", "20", ""),
			param: StringParam{FlagName: "threads", Config: strPtr("")},
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVar != "" {
				t.Setenv("SMARTLING_CLI_"+upper(tt.param.FlagName), tt.envVar)
			}
			if got := FallbackString(tt.flag, tt.param); got != tt.want {
				t.Errorf("FallbackString() = %q, want %q", got, tt.want)
			}
		})
	}
}

// upper is an ASCII-only uppercaser to keep this test file dependency-light.
func upper(s string) string {
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out[i] = c
	}
	return string(out)
}
