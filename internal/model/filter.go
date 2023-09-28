package model

import "sync"

type Filter struct {
	Mu   sync.RWMutex
	Word []string
}
