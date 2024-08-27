package model

type Routing struct {
	Host []Host `mapstructure:"host"`
}
type Host struct {
	Name string `mapstructure:"name"`
	Url  string `mapstructure:"url"`
}
