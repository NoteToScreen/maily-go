package maily

import "errors"

// ErrInvalidConfig is reported when a required parameter was not set in the Context.
var ErrInvalidConfig = errors.New("maily: required parameter not set in Context")

// ErrTemplateMissing is reported when the requested template could not be found.
var ErrTemplateMissing = errors.New("maily: couldn't find requested template")

// ErrTemplateMissingFile is reported when the requested template was missing a required file.
var ErrTemplateMissingFile = errors.New("maily: requested template is missing files")
