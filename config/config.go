package config

import (
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"os"
)

var Cfg = Config{}

func Initialization() {
	yamlFile, err := os.ReadFile(viper.GetString("config"))
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(yamlFile, &Cfg); err != nil {
		panic(err)
	}
}

type Config struct {
	PeerConfig `yaml:",omitempty,inline"`
	DKGResult  `yaml:",omitempty,inline"`
	Rank       uint32 `yaml:"rank"`
	Threshold  uint32 `yaml:"threshold"`
}
type PeerConfig struct {
	Port     int64  `yaml:"port"`
	HttpPort int64  `yaml:"httpport"`
	Identity string `yaml:"identity"`
	Peers    []struct {
		Id   string `yaml:"id"`
		Port int64  `yaml:"port"`
	} `yaml:"peers"`
}

type DKGResult struct {
	Share  string        `yaml:"share" json:"share"`
	Pubkey Pubkey        `yaml:"pubkey" json:"pubkey"`
	BKs    map[string]BK `yaml:"bks" json:"BKs"`
}

type ReshareResult struct {
	Share string `yaml:"share" json:"share"`
}

type SignerResult struct {
	Sign string `yaml:"sign" json:"sign"`
}

type BK struct {
	X    string `yaml:"x" json:"x"`
	Rank uint32 `yaml:"rank" json:"rank"`
}
type Pubkey struct {
	X string `yaml:"x" json:"x"`
	Y string `yaml:"y" json:"y"`
}

type MonitorNotify struct {
	Type       string `json:"type"`
	ProtocolId string `json:"protocolId"`
	Msg        string `json:"msg"`
}
