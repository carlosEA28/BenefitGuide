package exceptions

import "errors"

// ErrRequirementNotFound indicates that no service requirement exists for the
// given ID.
var ErrRequirementNotFound = errors.New("service requirement not found")

// ErrOfficialSourceNotFound indicates that no official source exists for the given ID.
var ErrOfficialSourceNotFound = errors.New("official source not found")

// ErrPublicServiceNotFound indicates that a public service ID is unknown.
var ErrPublicServiceNotFound = errors.New("public service not found")

// ErrPublicServiceInactive indicates that the public service cannot receive requirements.
var ErrPublicServiceInactive = errors.New("public service inactive")

// ErrInvalidRequirement indicates invalid requirement data.
var ErrInvalidRequirement = errors.New("invalid service requirement")
