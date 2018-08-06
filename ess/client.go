package ess

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/rifki192/alicloud-image-overwriter/config"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
)

type Client struct {
	client *ess.Client
}

func New(cfg *config.Config) *Client {
	var c Client
	config := sdk.NewConfig()

	credential := &credentials.BaseCredential{
		AccessKeyId:     cfg.AccesKey,
		AccessKeySecret: cfg.SecretKey,
	}

	ecsClient, err := ess.NewClientWithOptions(cfg.Region, config, credential)
	if err != nil {
		// Handle exceptions
		panic(err)
	}
	c.client = ecsClient

	return &c
}
