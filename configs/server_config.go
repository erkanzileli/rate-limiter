package configs

type serverConfig struct {
	// Addr represents listen address of the server
	Addr string

	// ReadTimeout is wait time for reading request and its body in milliseconds format
	ReadTimeout int

	// WriteTimeout is wait time for writing response in milliseconds format
	WriteTimeout int
}
