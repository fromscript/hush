package models

import "sync"

type Room struct {
	ID      string
	Members sync.Map // map[string]*Client
}
