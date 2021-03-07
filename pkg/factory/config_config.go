package factory

import (
	"github.com/jay16213/radio_simulator/pkg/simulator_context"

	"github.com/free5gc/openapi/models"
)

var TestAmDataTable = make(map[string]models.AccessAndMobilitySubscriptionData)
var TestSmfSelDataTable = make(map[string]models.SmfSelectionSubscriptionData)
var TestAmPolicyDataTable = make(map[string]models.AmPolicyData)

type Config struct {
	DBName     string       `yaml:"dbName"`
	DBUrl      string       `yaml:"dbUrl"`
	RanInfo    []RanContext `yaml:"ranInfo"`
	TcpUri     string       `yaml:"tcpUri"`
	UeInfoFile []string     `yaml:"ueInfoFile"`
	TunnelInfo TunnelInfo   `yaml:"gtp5gTunnelInfo"`
	// ListenIp   string       `yaml:"listenIp"`
	Logger Logger `yaml:"logger"`
}

type TunnelInfo struct {
	TunDev    string `yaml:"tunDev"`
	Gtp5gPath string `yaml:"path"`
}

type RanContext struct {
	AmfUri        string                       `yaml:"amfUri"`
	RanSctpUri    string                       `yaml:"ranSctpUri"`
	RanGtpUri     simulator_context.AddrInfo   `yaml:"ranGtpUri"`
	UpfUriList    []simulator_context.AddrInfo `yaml:"upfUriList"`
	RanName       string                       `yaml:"ranName"`
	GnbId         GnbId                        `yaml:"gnbId"`
	SupportTAList []SupportTAItem              `yaml:"taiList"`
}

type GnbId struct {
	PlmnId    models.PlmnId `yaml:"plmnId"`
	BitLength int           `yaml:"length"`
	Value     string        `yaml:"value"`
}

type SupportTAItem struct {
	Tac      string            `yaml:"tac"`
	Plmnlist []PlmnSupportItem `yaml:"plmnList,omitempty"`
}

type PlmnSupportItem struct {
	PlmnId     models.PlmnId   `yaml:"plmnId"`
	SNssaiList []models.Snssai `yaml:"snssaiList,omitempty"`
}
