package config

type Server struct {
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Mysql  Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Pgsql  Pgsql `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
}
