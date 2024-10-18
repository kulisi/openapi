package conf

type Config struct {
	Viper   Viper   `json:"viper" mapstructure:"viper" yaml:"viper"`
	Gin     Gin     `json:"gin" mapstructure:"gin" yaml:"gin"`
	Zap     Zap     `json:"zap" mapstructure:"zap" yaml:"zap"`
	Gorm    Gorm    `json:"gorm" mapstructure:"gorm" yaml:"gorm"`
	Service Service `json:"service" mapstructure:"service" yaml:"service"`
}
