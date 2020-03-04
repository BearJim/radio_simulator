package ue_factory

import ()
import "radio_simulator/lib/openapi/models"

type Config struct {
	Supi          string                              `yaml:"supi"`
	Gpsis         []string                            `yaml:"gpsis"`
	Nssai         models.Nssai                        `yaml:"nssai"`
	UeAmbr        UeAmbr                              `yaml:"ueAmbr"`
	SmfSelData    models.SmfSelectionSubscriptionData `yaml:"smfSelData"`
	AuthData      AuthData                            `yaml:"auths"`
	SubscCats     []string                            `json:"subscCats,omitempty"`
	ServingPlmnId string                              `yaml:"servingPlmn"`
}

type UeAmbr struct {
	Upink    string `yaml:"uplink"`
	DownLink string `yaml:"downlink"`
}

type AuthData struct {
	AuthMethod string `yaml:"authMethod"`
	K          string `yaml:"K"`
	Opc        string `yaml:"Opc,omitempty"`
	Op         string `yaml:"Op,omitempty"`
	AMF        string `yaml:"AMF"`
	SQN        string `yaml:"SQN"`
}
