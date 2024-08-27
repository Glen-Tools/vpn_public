package model

type Outbound struct {
	User []struct {
		Name string `mapstructure:"name"`
		Id   string `mapstructure:"id"`
	} `mapstructure:"user"`
}
