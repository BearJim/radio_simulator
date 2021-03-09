package simulator_handler

import (
	"time"

	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
	"github.com/jay16213/radio_simulator/pkg/simulator_ngap/ngap_handler"
	"github.com/sirupsen/logrus"
)

var HandlerLog *logrus.Entry
var NgapLog *logrus.Entry

func init() {
	HandlerLog = logger.HandlerLog
	NgapLog = logger.NgapLog
}

func Handle(ran *simulator_context.RanContext, msgChan chan []byte) {
	for {
		select {
		case msg, ok := <-msgChan:
			if ok {
				ngap_handler.Dispatch(ran, msg)
			} else {
				HandlerLog.Errorln("Channel closed!")
			}

		case <-time.After(time.Second * 1):

		}
	}
}
