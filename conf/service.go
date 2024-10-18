package conf

type Service struct {
	Name        string `json:"name" mapstructure:"name" yaml:"name"`
	DisplayName string `json:"display-name" mapstructure:"display-name" yaml:"display-name"`
	Description string `json:"description" mapstructure:"description" yaml:"description"`
}
