package tcp_server

import (
	"radio_simulator/src/simulator_context"
	"radio_simulator/src/simulator_ngap"
	"regexp"
	"strconv"
)

var stringFormat = regexp.MustCompile(`\S+`)

// Parse cli Command from Tcp Server
// Usage: [ options ]
//
// options:
//   show [all|${PduSessionId}]
// 		show ue current state
//
//   reg {ip}
// 		register to CoreNetwork
// 		ip is specific Ran Ip to Connect(default={firstRanIp})
//
// 	 dereg
// 		deregister to CoreNetwork
//
//   sess i [add|del]
// 		pduSessionId i add or delete
//
// Output Format
//   show  [all|${PduSessionId}]
//		first line in all case:
//			"[SHOW] REGISTERED\n" or "[SHOW] REGISTERING\n" or "[SHOW] DEREGISTERED\n"
//		sessInfo:
//			"[SHOW] ID=%d,ULIP=%s,ULTEID=%d,DLIP=%s,DLTEID=%d\n"
// 		all means show all Pdu Session Id
//
//   reg {ip}
// 		"[REG] SUCCESS\n" or
// 		"[REG] FAIL\n" or
//
// 	 dereg
// 		"[DEREG] SUCCESS\n" or
// 		"[DEREG] FAIL\n" or
//
//   sess i [add|del]
//		"[SESSION] ID=%d,ULIP=%s,ULTEID=%d,DLIP=%s,DLTEID=%d\n" for add case, "[SESSION] DEL %d\n" for del case or
// 		"[SESSION] ADD/DEL %d FAIL\n"
//
func parseCmd(ue *simulator_context.UeContext, cmd string) {
	params := stringFormat.FindAllString(cmd, -1)
	cnt := len(params)
	if cnt == 0 {
		return
	}
	var msg string
	switch params[0] {
	case "show":
		if cnt == 1 {
			msg = "show missing action[all/{id}]"
		} else {
			switch params[1] {
			case "all":
				msg = "[SHOW] " + ue.RegisterState + "\n"
				for _, sess := range ue.PduSession {
					sessInfo := sess.GetTunnelMsg()
					if sessInfo == "" {
						continue
					}
					msg = msg + "[SHOW] " + sessInfo + "\n"
				}
			default:
				id, err := strconv.Atoi(params[1])
				if err != nil {
					msg = "sess id is not digit"
					break
				}
				sess := ue.PduSession[id]
				if sess == nil {
					msg = "sess " + params[1] + " has not established yet"
					break
				}
				sessInfo := sess.GetTunnelMsg()
				if sessInfo == "" {
					msg = "sess " + params[1] + " is still establishing"
					break
				}
				msg = "[SHOW] " + sessInfo + "\n"
			}
		}
	case "reg":
		switch ue.RegisterState {
		case simulator_context.RegisterStateRegitered:
			ue.SendSuccessRegister()
			return
		case simulator_context.RegisterStateRegitering:
			return
		}
		ran := self.RanPool[self.DefaultRanUri]
		if cnt > 1 {
			ran = self.RanPool[params[1]]
			if ran == nil {
				msg = "ranIp " + params[1] + " does not exist"
				break
			}
			// Use Default RanUri
		}
		ue.AttachRan(ran)
		ue.RegisterState = simulator_context.RegisterStateRegitering
		simulator_ngap.SendInitailUeMessage_RegistraionRequest(ran, ue)
	case "dereg":
		if ue.RegisterState == simulator_context.RegisterStateDeregitered {
			ue.SendSuccessDeregister()
			return
		} else {
			// TODO: Send Degister Request
		}
	case "sess":
		if cnt <= 2 {
			msg = "sess need id and action[add/del]"
			break
		}
		id, err := strconv.Atoi(params[1])
		if err != nil {
			msg = "sess id is not digit"
			break
		}
		sess := ue.PduSession[id]
		switch params[2] {
		case "add":
			if sess == nil {
				// TODO: Send Pdu Session Estblishment
				break
			}
			sessInfo := sess.GetTunnelMsg()
			if sessInfo == "" {
				msg = "sess " + params[1] + " is still establishing"
				break
			}
			msg = "[SESSION] " + sessInfo + "\n"
		case "del":
			if sess == nil {
				msg = "[SESSION] DEL " + params[1] + "\n"
				break
			} else {
				// TODO: Send Pdu Session Release
			}
		default:
			msg = "sess action is not [add/del]"
		}
	}

	if msg != "" {
		if msg[0] != '[' {
			msg = "[ERROR] " + msg + "\n"
		}
		ue.TcpConn.Write([]byte(msg))
	}

}
