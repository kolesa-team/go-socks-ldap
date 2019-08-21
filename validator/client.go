package validator

type (
	// general validators interface
	Client interface {
		// get the field related with the validator
		GetFieldName() string
		// password validation
		Validate(hash string, password string) bool
	}
)
