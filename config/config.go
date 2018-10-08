package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/caarlos0/env"
)

// Config is defined all the configuration needed by the app
type Config struct {
	Region    string `env:"ALICLOUD_REGION" envDefault:"ap-southeast-1"`
	AccesKey  string `env:"ALICLOUD_ACCESS_KEY"`
	SecretKey string `env:"ALICLOUD_SECRET_KEY"`
}

type StsCred struct {
	Region    string `json:"RegionId,omitempty"`
	AccessKey string `json:"AccessKeyId"`
	SecretKey string `json:"AccessKeySecret"`
	StsToken  string `json:"SecurityToken"`
}

const (
	Credentials_URL = "http://100.100.100.200/latest/meta-data/ram/security-credentials/AutoScaleSLB"
	RegionId_URL    = "http://100.100.100.200/latest/meta-data/region-id"
)

// LoadConfig for get all the configuration from Env Variables
func LoadConfig() *Config {
	//Init Config from ENV Variable
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg
}

func GetStsCredentials() *StsCred {
	var sts *StsCred

	respSts := GetHttpCall(Credentials_URL)
	json.Unmarshal([]byte(respSts), &sts)
	fmt.Println("Got Credentials with AccessKey: ", sts.AccessKey)
	// fmt.Println("Got Credentials with SecretKey: ", sts.SecretKey)
	// fmt.Println("Got Credentials with StsToken: ", sts.StsToken)

	respRegion := GetHttpCall(RegionId_URL)
	sts.Region = respRegion
	fmt.Println("Got Credentials with Region: ", sts.Region)

	return sts
}

func GetHttpCall(uri string) string {
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Erro when get http call, ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Println("Erro when read response body, ", err)
		}
		bodyString := string(bodyBytes)

		return bodyString
	}

	return ""
}
