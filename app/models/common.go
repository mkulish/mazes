package models

import "github.com/revel/revel"

// ValidationError represents an input validation error
// swagger:model ValidationError
type ValidationError struct {
	// Operation success flag
	// required: true
	// type: boolean
	OK bool `json:"ok"`

	// Validation errors
	// required: true
	// type: object
	Errors []*revel.ValidationError `json:"errors"`
}

// UnauthorizedError represents an unauthorized access error
// swagger:model UnauthorizedError
type UnauthorizedError struct {
	// Operation success flag
	// required: true
	// type: boolean
	OK bool `json:"ok"`

	// Unauthorized error
	// required: true
	// type: string
	Error string `json:"error"`
}

// InternalError represents an unexpected internal error
// swagger:model InternalError
type InternalError struct {
	// Operation success flag
	// required: true
	// type: boolean
	OK bool `json:"ok"`

	// Internal error
	// required: true
	// type: string
	Error string `json:"error"`
}
