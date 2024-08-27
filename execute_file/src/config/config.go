package config

// import (
//     "fmt"
//     "runtime"
// )

const (
	// 目标文件在 config 目录下
	InboundYmlPath = "config/inbound.yml"
	RoutingYmlPath = "config/routing.yml"

	// 目标文件在 config 目录下
	ConfigJsonPath = "config/config.json"

	V2rayInboundSocksProtocol = "socks"

	RoutingDomainPrefix = "domain:"

	OutboundsUserIdPath = "outbounds[0].settings.vnext[0].users[0].id"
	RoutingHostPath     = "routing.rules[0].domain"
	LogAccessPath       = "log.access"
	LogErrorPath        = "log.error"
	LogAccessName       = "access.log"
	LogErrorName        = "error.log"
	InboundPath         = "inbounds"
	WindowsV2rayPath    = "v2ray_file/windows/v2ray.exe"
	MacV2rayPath        = "mac/v2ray"
	Debug               = false
)
