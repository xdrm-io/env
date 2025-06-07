package env

// Err is an error that can be returned by the env package
type Err string

func (e Err) Error() string { return string(e) }

// avail env errors
const (
	ErrNotPtr       Err = "not a pointer"
	ErrNotStructPtr Err = "not a pointer to struct"

	ErrFieldUnexported  Err = "field is unexported"
	ErrFieldNoEnvTag    Err = "no env tag"
	ErrFieldDecode      Err = "field decode"
	ErrFieldUnsupported Err = "unsupported field type"
	ErrFieldRequired    Err = "field is required"
)
