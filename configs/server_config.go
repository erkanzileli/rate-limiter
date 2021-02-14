package configs

type serverConfig struct {
	// Addr represents listen address of the server
	Addr string `yaml:"addr"`

	// ReadTimeout is wait time for reading request and its body in milliseconds format
	ReadTimeout int `yaml:"readTimeout"`

	// WriteTimeout is wait time for writing response in milliseconds format
	WriteTimeout int `yaml:"writeTimeout"`
}
