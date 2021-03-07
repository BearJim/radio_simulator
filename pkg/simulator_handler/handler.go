package simulator_handler

import (
	"time"

	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_handler/simulator_message"
	"github.com/jay16213/radio_simulator/pkg/simulator_ngap/ngap_handler"
	"github.com/sirupsen/logrus"
)

var HandlerLog *logrus.Entry
var NgapLog *logrus.Entry

func init() {
	HandlerLog = logger.HandlerLog
	NgapLog = logger.NgapLog
}

func Handle(laddr string) {
	for {
		select {
		case msg, ok := <-simulator_message.SimChannel[laddr]:
			if ok {
				ngap_handler.Dispatch(laddr, msg.Value)
			} else {
				HandlerLog.Errorln("Channel closed!")
			}

		case <-time.After(time.Second * 1):

		}
	}
}
