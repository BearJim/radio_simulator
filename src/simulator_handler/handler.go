package simulator_handler

import (
	"github.com/sirupsen/logrus"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_handler/simulator_message"
	"radio_simulator/src/simulator_ngap/ngap_handler"
	"time"
)

var HandlerLog *logrus.Entry
var NgapLog *logrus.Entry

func init() {
	HandlerLog = logger.HandlerLog
	NgapLog = logger.NgapLog
}

func Handle() {
	for {
		select {
		case msg, ok := <-simulator_message.SimChannel:
			if ok {
				ngap_handler.Dispatch(msg.NgapAddr, msg.Value)
			} else {
				HandlerLog.Errorln("Channel closed!")
			}

		case <-time.After(time.Second * 1):

		}
	}
}
