package exceptions

import (
	"errors"
	"fmt"
	"testing"
)

func TestCatalogErrorsCanBeIdentified(t *testing.T) {
	for _, sentinel := range []error{
		ErrPublicServiceNotFound,
		ErrPublicServiceInactive,
		ErrInvalidRequirement,
	} {
		if !errors.Is(fmt.Errorf("catalog: %w", sentinel), sentinel) {
			t.Fatalf("wrapped error is not identifiable: %v", sentinel)
		}
	}
}
