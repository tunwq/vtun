package config

import (
	"encoding/json"
	"log"

	"github.com/net-byte/vtun/common/cipher"
)

type Config struct {
	DeviceName    string
	LocalAddr     string
	ServerAddr    string
	CIDR          string
	IPv6CIDR	  string
	Key           string
	Protocol      string
	DNS           string
	WebSocketPath string
	ServerMode    bool
	GlobalMode    bool
	Obfs          bool
	MTU           int
	Timeout       int
}

func (config *Config) Init() {
	cipher.GenerateKey(config.Key)
	json, _ := json.Marshal(config)
	log.Printf("init config:%s", string(json))
}
