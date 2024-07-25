package ierr

// IErr interface that should be implemented by all custom errors.
type IErr interface {
	// Type returns the type of the error
	Type() string
}
