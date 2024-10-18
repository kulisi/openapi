package conf

import "time"

type Gin struct {
	Use     bool          `mapstructure:"use" json:"use" yaml:"use"`
	Addr    string        `mapstructure:"addr" json:"addr" yaml:"addr"`
	WaitFor time.Duration `mapstructure:"waitfor" json:"waitfor" yaml:"waitfor"`
}
