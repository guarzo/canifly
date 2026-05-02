package skillplan_test

import (
	"testing"

	"github.com/guarzo/canifly/internal/services/skillplan"
	"github.com/guarzo/canifly/internal/testutil"
)

func TestNewService_Constructs(t *testing.T) {
	s := skillplan.NewService(&testutil.MockLogger{}, nil)
	if s == nil {
		t.Fatal("NewService returned nil")
	}
}
