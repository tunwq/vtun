package main

import (
	"flag"

	"github.com/net-byte/vtun/common/config"
	"github.com/net-byte/vtun/tcp"
	"github.com/net-byte/vtun/udp"
	"github.com/net-byte/vtun/ws"
)

func main() {
	config := config.Config{}
	flag.StringVar(&config.DeviceName, "dn", "tun0", "device name")
	flag.StringVar(&config.CIDR, "c", "172.16.0.10/24", "tun interface cidr")
	flag.StringVar(&config.IPv6CIDR, "6c", "fced:9999::9999/64", "tun interface ipv6 cidr")
	flag.IntVar(&config.MTU, "mtu", 1500, "tun mtu")
	flag.StringVar(&config.LocalAddr, "l", ":3000", "local address")
	flag.StringVar(&config.ServerAddr, "s", ":3001", "server address")
	flag.StringVar(&config.Key, "k", "freedom@2022", "key")
	flag.StringVar(&config.Protocol, "p", "wss", "protocol tcp/udp/ws/wss")
	flag.StringVar(&config.DNS, "d", "8.8.8.8:53", "dns address")
	flag.StringVar(&config.WebSocketPath, "path", "/freedom", "websocket path")
	flag.BoolVar(&config.ServerMode, "S", false, "server mode")
	flag.BoolVar(&config.GlobalMode, "g", false, "client global mode")
	flag.BoolVar(&config.Obfs, "obfs", false, "enable data obfuscation")
	flag.IntVar(&config.Timeout, "t", 30, "dial timeout in seconds")
	flag.Parse()
	config.Init()
	switch config.Protocol {
	case "udp":
		if config.ServerMode {
			udp.StartServer(config)
		} else {
			udp.StartClient(config)
		}
	case "tcp":
		if config.ServerMode {
			tcp.StartServer(config)
		} else {
			tcp.StartClient(config)
		}
	case "ws":
		if config.ServerMode {
			ws.StartServer(config)
		} else {
			ws.StartClient(config)
		}
	default:
		if config.ServerMode {
			ws.StartServer(config)
		} else {
			ws.StartClient(config)
		}
	}
}
