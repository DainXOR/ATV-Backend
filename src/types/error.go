package types

type ComposedError struct {
	Errors []error `json:"errors"`
}

func (m *ComposedError) Error() string {
	msg := "Errors: "
	for _, err := range m.Errors {
		msg += err.Error() + " > "
	}

	return msg
}

func (m *ComposedError) Add(err error) {
	m.Errors = append(m.Errors, err)
}
