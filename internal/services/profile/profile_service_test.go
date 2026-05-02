package profile_test

import (
	"testing"

	"github.com/guarzo/canifly/internal/services/profile"
	"github.com/guarzo/canifly/internal/testutil"
)

func TestNewService_Constructs(t *testing.T) {
	s := profile.NewService(nil, nil, nil, nil, &testutil.MockLogger{})
	if s == nil {
		t.Fatal("NewService returned nil")
	}
}
