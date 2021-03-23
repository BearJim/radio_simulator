package factory

import (
	"net"

	"github.com/free5gc/openapi/models"
)

type Config struct {
	DBName          string          `yaml:"dbName"`
	DBUrl           string          `yaml:"dbUrl"`
	AmfSCTPEndpoint SCTPEndpoint    `yaml:"amfSctpEndpoint"`
	RanSctpEndpoint SCTPEndpoint    `yaml:"ranSctpEndpoint"`
	RanGtpUri       net.UDPAddr     `yaml:"ranGtpUri"`
	UpfUriList      []net.UDPAddr   `yaml:"upfUriList"`
	RanName         string          `yaml:"ranName"`
	GnbId           GnbId           `yaml:"gnbId"`
	SupportTAList   []SupportTAItem `yaml:"taiList"`
	ApiServerAddr   string          `yaml:"apiServerAddr"`
	UeInfoFile      []string        `yaml:"ueInfoFile"`
	TunnelInfo      TunnelInfo      `yaml:"gtp5gTunnelInfo"`
	Logger          Logger          `yaml:"logger"`
}

type TunnelInfo struct {
	TunDev    string `yaml:"tunDev"`
	Gtp5gPath string `yaml:"path"`
}

type SCTPEndpoint struct {
	IPs  []net.IP `yaml:"ips"`
	Port int      `yaml:"port"`
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
