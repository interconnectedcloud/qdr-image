package tcpadaptor

import (
	"os"
	"testing"

	"github.com/skupperproject/skupper/test/utils/base"
)

// TestMain helps parsing the common test flags and running package level tests
func TestMain(m *testing.M) {
	base.ParseFlags()
	os.Exit(m.Run())
}
