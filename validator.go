package go_cqrs

// Validator is an optional interface that request types can implement.
// When the [Validation] decorator is active, it calls Validate() before
// the request reaches the handler. If Validate returns an error, the
// pipeline short-circuits and returns that error immediately.
//
//	type CreateUserCmd struct {
//	    Name  string
//	    Email string
//	}
//
//	func (c CreateUserCmd) Validate() error {
//	    if c.Name == "" {
//	        return errors.New("name is required")
//	    }
//	    return nil
//	}
type Validator interface {
	Validate() error
}
