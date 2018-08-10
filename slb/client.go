package slb

import (
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
	config := sdk.NewConfig()

	credential := &credentials.BaseCredential{
		AccessKeyId:     cfg.AccesKey,
		AccessKeySecret: cfg.SecretKey,
	}

	slbClient, err := slb.NewClientWithOptions(cfg.Region, config, credential)
	if err != nil {
		// Handle exceptions
		panic(err)
	}
	c.client = slbClient

	return &c
}
