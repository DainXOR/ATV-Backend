package models

// Later will see how to use this interface to
// simplify the code on the db interface layer.
//
// To implement this interface, a model must
// define a method TableName that returns the name of the database table
// Yes, that's it, nothing else is required.
type DBModelInterface interface {
	TableName() string
}
