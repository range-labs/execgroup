package execgroup

import "strings"

// MultiError contains several errors.
type MultiError []error

// Err returns the error, if it is not empty, nil otherwise.
func (e MultiError) Err() error {
	if e.Empty() {
		return nil
	}
	return e
}

// Errors returns the individual errors.
func (e MultiError) Errors() []error {
	return e
}

// Empty returns true if there are no errors.
func (e MultiError) Empty() bool {
	return len(e) == 0
}

// Len returns the number of errors.
func (e MultiError) Len() int {
	return len(e)
}

// Append adds more errors.
func (e MultiError) Append(errs ...error) MultiError {
	for _, err := range errs {
		if me, ok := err.(MultiError); ok {
			e = e.Append(me...)
		} else if err != nil {
			e = append(e, err)
		}
	}
	return e
}

// Error concatenates the errors.
func (e MultiError) Error() string {
	if len(e) == 0 {
		return ""
	}
	if len(e) == 1 {
		return e[0].Error()
	}
	var errors = make([]string, len(e))
	for i, err := range e {
		errors[i] = err.Error()
	}
	return "multiple (" + strings.Join(errors, "; ") + ")"
}

// NewMultiError returns a multi error from the provided.
func NewMultiError(errs ...error) MultiError {
	var e MultiError
	return e.Append(errs...)
}
