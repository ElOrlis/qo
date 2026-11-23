package circuit

import (
	"sync"
	"time"
)

type Func func() error

type AlertFunc Func

type PingFunc Func

type HealthFunc Func

type Monitor struct {
	Alert  AlertFunc
	Ping   PingFunc
	Health HealthFunc
}

type Status uint8

const (
	CLOSED Status = iota
	HALF_OPENED
	OPENED
)

type Block struct {
	mu                sync.Mutex
	InvocationTimeout time.Time
	FailureTreshold   uint16
	status            Status
}
