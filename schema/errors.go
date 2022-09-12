package schema

import (
	"fmt"
	"strings"
)

// ConstraintError indicates that the passed data violated one or more constraints defined in the schema.
// The message holds the exact path of the problematic field, as well as a message explaining the error.
// If this error is not easily understood, please open an issue on the Arcaflow plugin SDK.
type ConstraintError struct {
	Message string
	Path    []string
	Cause   error
}

// Error returns the error message.
func (c ConstraintError) Error() string {
	result := fmt.Sprintf("Validation failed for '%s': %s", strings.Join(c.Path, "' -> '"), c.Message)
	if c.Cause != nil {
		result += " (" + c.Cause.Error() + ")"
	}
	return result
}

// Unwrap returns the underlying error if any.
func (c ConstraintError) Unwrap() error {
	return c.Cause
}

// NoSuchStepError indicates that the given step is not supported by the plugin.
type NoSuchStepError struct {
	Step string
}

// Error returns the error message.
func (n NoSuchStepError) Error() string {
	return fmt.Sprintf("No such step: %s", n.Step)
}

// BadArgumentError indicates that an invalid configuration was passed to a schema component. The message will
// explain what exactly the problem is, but may not be able to locate the exact error as the schema may be manually
// built.
type BadArgumentError struct {
	Message string
	Cause   error
}

// Error returns the error message.
func (b BadArgumentError) Error() string {
	result := b.Message
	if b.Cause != nil {
		result += " (" + b.Cause.Error() + ")"
	}
	return result
}

// Unwrap returns the underlying error if any.
func (b BadArgumentError) Unwrap() error {
	return b.Cause
}

// UnitParseError indicates that it failed to parse a unit string.
type UnitParseError struct {
	Message string
	Cause   error
}

func (u UnitParseError) Error() string {
	result := u.Message
	if u.Cause != nil {
		result += " (" + u.Cause.Error() + ")"
	}
	return result
}

// Unwrap returns the underlying error if any.
func (u UnitParseError) Unwrap() error {
	return u.Cause
}
