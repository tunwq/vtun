package tcp

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/net-byte/vtun/common/cache"
	"github.com/net-byte/vtun/common/cipher"
	"github.com/net-byte/vtun/common/config"
	"github.com/net-byte/vtun/common/counter"
	"github.com/net-byte/vtun/common/netutil"
	"github.com/net-byte/vtun/tun"
	"github.com/songgao/water"
)

// Start tcp server
func StartServer(config config.Config) {
	log.Printf("vtun tcp server started on %v", config.LocalAddr)
	iface := tun.CreateTun(config)
	// server -> client
	go toClient(config, iface)
	ln, err := net.Listen("tcp", config.LocalAddr)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		// client -> server
		go toServer(config, conn, iface)
	}

}

func toClient(config config.Config, iface *water.Interface) {
	packet := make([]byte, config.MTU)
	for {
		n, err := iface.Read(packet)
		if err != nil || err == io.EOF || n == 0 {
			continue
		}
		b := packet[:n]
		if key := netutil.GetDestinationKey(b); key != "" {
			if v, ok := cache.GetCache().Get(key); ok {
				if config.Obfs {
					b = cipher.XOR(b)
				}
				counter.IncrWriteByte(n)
				v.(net.Conn).Write(b)
			}
		}
	}
}

func toServer(config config.Config, tcpconn net.Conn, iface *water.Interface) {
	defer tcpconn.Close()
	packet := make([]byte, config.MTU)
	for {
		tcpconn.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
		n, err := tcpconn.Read(packet)
		if err != nil || err == io.EOF {
			break
		}
		b := packet[:n]
		if config.Obfs {
			b = cipher.XOR(b)
		}
		if key := netutil.GetSourceKey(b); key != "" {
			cache.GetCache().Set(key, tcpconn, 10*time.Minute)
			counter.IncrReadByte(len(b))
			iface.Write(b)
		}
	}
}
