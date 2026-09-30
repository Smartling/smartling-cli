package files

import (
	"testing"
	"time"
)

func TestLocalTimeZoneName_IsLoadableZone(t *testing.T) {
	name := localTimeZoneName()
	if name == "" || name == "Local" {
		t.Fatalf("localTimeZoneName() = %q, want an IANA zone name", name)
	}
	if _, err := time.LoadLocation(name); err != nil {
		t.Errorf("localTimeZoneName() = %q, not a loadable zone: %v", name, err)
	}
}
