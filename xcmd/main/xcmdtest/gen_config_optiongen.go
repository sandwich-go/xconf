// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package xcmdtest

import "time"

// Config should use NewConfig to initialize it
type Config struct {
	HttpAddress string
	Timeouts    map[string]time.Duration
}

// NewConfig new Config
func NewConfig(opts ...ConfigOption) *Config {
	cc := newDefaultConfig()
	for _, opt := range opts {
		opt(cc)
	}
	if watchDogConfig != nil {
		watchDogConfig(cc)
	}
	return cc
}

// ApplyOption apply mutiple new option and return the old ones
// sample:
// old := cc.ApplyOption(WithTimeout(time.Second))
// defer cc.ApplyOption(old...)
func (cc *Config) ApplyOption(opts ...ConfigOption) []ConfigOption {
	var previous []ConfigOption
	for _, opt := range opts {
		previous = append(previous, opt(cc))
	}
	return previous
}

// ConfigOption option func
type ConfigOption func(cc *Config) ConfigOption

// WithHttpAddress option func for filed HttpAddress
func WithHttpAddress(v string) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.HttpAddress
		cc.HttpAddress = v
		return WithHttpAddress(previous)
	}
}

// WithTimeouts option func for filed Timeouts
func WithTimeouts(v map[string]time.Duration) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Timeouts
		cc.Timeouts = v
		return WithTimeouts(previous)
	}
}

// InstallConfigWatchDog the installed func will called when NewConfig  called
func InstallConfigWatchDog(dog func(cc *Config)) { watchDogConfig = dog }

// watchDogConfig global watch dog
var watchDogConfig func(cc *Config)

// newDefaultConfig new default Config
func newDefaultConfig() *Config {
	cc := &Config{}

	for _, opt := range [...]ConfigOption{
		WithHttpAddress(":3001"),
		WithTimeouts(map[string]time.Duration{
			"read":  time.Duration(10) * time.Second,
			"write": time.Duration(20) * time.Second,
		}),
	} {
		opt(cc)
	}

	return cc
}

// all getter func
func (cc *Config) GetHttpAddress() string                { return cc.HttpAddress }
func (cc *Config) GetTimeouts() map[string]time.Duration { return cc.Timeouts }

// ConfigVisitor visitor interface for Config
type ConfigVisitor interface {
	GetHttpAddress() string
	GetTimeouts() map[string]time.Duration
}

// ConfigInterface visitor + ApplyOption interface for Config
type ConfigInterface interface {
	ConfigVisitor
	ApplyOption(...ConfigOption) []ConfigOption
}