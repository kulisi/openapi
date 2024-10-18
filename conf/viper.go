package conf

type Viper struct {
	ConfigName  string   `json:"name"`
	ConfigType  string   `json:"type"`
	ConfigPaths []string `json:"paths"`
}
