package simulator_context

import (
	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/free5gc/openapi/models"
)

type AMFContext struct {
	Name            string // AMF Name
	ServedGUAMIList []ServedGUAMI
	Addr            *sctp.SCTPAddr
}

type ServedGUAMI struct {
	Guami         models.Guami
	BackupAMFName string
}
