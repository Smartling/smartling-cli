package helpers

import (
	"os"
	"testing"

	"github.com/Smartling/smartling-cli/services/helpers/rlog"
)

func TestMain(m *testing.M) {
	rlog.Init()
	os.Exit(m.Run())
}
