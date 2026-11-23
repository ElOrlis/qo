package qo

import "time"

type Store interface{}

type Cache struct {
	store     Store
	enable    bool
	globalTTL time.Duration
	keyFunc   CacheKeyFunc
}
