package exceptions

import "errors"

// ErrRequirementNotFound indicates that no service requirement exists for the
// given ID.
var ErrRequirementNotFound = errors.New("service requirement not found")
