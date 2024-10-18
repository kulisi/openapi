package conf

type Gorm struct {
	Use   string `mapstructure:"use" json:"use" yaml:"use"`
	Mssql Mssql  `json:"mssql" mapstructure:"mssql" yaml:"mssql"`
	Mysql Mysql  `json:"mysql" mapstructure:"mysql" yaml:"mysql"`
}
