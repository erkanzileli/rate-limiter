package configs

type algorithm string

const (
	FixedWindowCount algorithm = "fixed-window-count"
)

type algorithmConfig struct {
	// Name of the algorithm. Default is fixed-window-count.
	Name algorithm

	// Options represents the specified algorithm options.
	Options algorithmOptionsConfig
}

type algorithmOptionsConfig struct {
	// WindowLengthInSeconds is for FixedWindowCount. Default is 60.
	WindowLengthInSeconds int64
}

// IsValid is for validating given algorithm name.
func (a algorithm) IsValid() bool {
	switch a {
	case FixedWindowCount:
		return true
	}
	return false
}
