package logging

// AppName represents the name of the application for the logging.
type AppName string

// String returns the string representation of the AppName.
func (n AppName) String() string {
	return string(n)
}
