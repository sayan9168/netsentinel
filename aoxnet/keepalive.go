package aoxnet

import "time"

type KeepAliveConfig struct {
	Interval  time.Duration
	Timeout   time.Duration
	MaxMisses int
}

func (k KeepAliveConfig) Normalize() KeepAliveConfig {
	if k.Interval <= 0 {
		k.Interval = 30 * time.Second
	}
	if k.Timeout <= 0 {
		k.Timeout = 10 * time.Second
	}
	if k.MaxMisses <= 0 {
		k.MaxMisses = 3
	}
	return k
}
