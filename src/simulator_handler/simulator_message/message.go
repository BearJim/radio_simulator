package simulator_message

import (
	"sync"
)

var SimChannel chan NGAPMessage
var mtx sync.Mutex

type NGAPMessage struct {
	NgapAddr string // NGAP Connection Addr
	Value    []byte // input/request value
}

const (
	MaxChannel int = 100000
)

func init() {
	// init Pool
	SimChannel = make(chan NGAPMessage, MaxChannel)
}

func SendMessage(msg NGAPMessage) {
	mtx.Lock()
	SimChannel <- msg
	mtx.Unlock()
}
