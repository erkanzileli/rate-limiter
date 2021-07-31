package configs

import (
	"fmt"
	"strings"
	"time"
)

const (
	colon      = ":"
	comma      = ","
	addrFormat = "http://%s:%s"
)

type appConfig struct {
	// Port is the app port
	Port string

	// Hosts can be IP or IP:PORT list
	Hosts []string

	// Timeout is client timeout when redirecting the requests
	Timeout time.Duration
}

func (a appConfig) GetAddresses() (hosts string) {
	for i, host := range a.Hosts {
		if strings.ContainsAny(host, colon) {
			hosts += host
		} else {
			hosts += fmt.Sprint(addrFormat, host, a.Port)
		}

		if i+1 < len(a.Hosts) {
			hosts += comma
		}
	}
	return hosts
}
