package errors

import "fmt"

// WrappedError wraps the given error with the given message and arguments.
func WrappedError(err error, msg string, args ...any) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %s", err, fmt.Sprintf(msg, args...))
}
