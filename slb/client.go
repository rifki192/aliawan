package slb

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/rifki192/alicloud-image-overwriter/config"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
)

const pageSize = 10

type Client struct {
	client *slb.Client
}

func New(cfg *config.Config) *Client {
	var c Client
	var err error
	var slbClient *slb.Client

	conf := sdk.NewConfig()

	if cfg.AccesKey == "" {
		fmt.Println("Load apps with Sts Credential")
		stsCred := config.GetStsCredentials()
		slbClient, err = slb.NewClientWithStsToken(stsCred.Region, stsCred.AccessKey, stsCred.SecretKey, stsCred.StsToken)
	} else {
		fmt.Println("Load apps with Env Credential")
		credential := &credentials.BaseCredential{
			AccessKeyId:     cfg.AccesKey,
			AccessKeySecret: cfg.SecretKey,
		}
		slbClient, err = slb.NewClientWithOptions(cfg.Region, conf, credential)
	}

	if err != nil {
		// Handle exceptions
		panic(err)
	}
	c.client = slbClient

	return &c
}
