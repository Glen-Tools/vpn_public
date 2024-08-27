package model

type Inbound struct {
	Proxy []Proxy `mapstructure:"proxy"`
}

type Proxy struct {
	Protocol string `mapstructure:"protocol"`
	Listen   string `mapstructure:"listen"`
	Port     string `mapstructure:"port"`
}
