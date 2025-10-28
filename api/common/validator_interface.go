package common

type Validatable interface {
	Validate() error
}
