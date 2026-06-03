package jobs

import (
	"os"
	"testing"

	"github.com/Smartling/smartling-cli/services/helpers/rlog"
)

// TestMain initializes the global logger; RunList/RunView/RunFiles call
// rlog.Debugf, which panics on a nil logger without Init.
func TestMain(m *testing.M) {
	rlog.Init()
	os.Exit(m.Run())
}
