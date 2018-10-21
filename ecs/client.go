package ecs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/rifki192/alicloud-image-overwriter/config"
)

type Client struct {
	client *ecs.Client
}

func New(cfg *config.Config) *Client {
	var c Client
	config := sdk.NewConfig()

	credential := &credentials.BaseCredential{
		AccessKeyId:     cfg.AccesKey,
		AccessKeySecret: cfg.SecretKey,
	}

	ecsClient, err := ecs.NewClientWithOptions(cfg.Region, config, credential)
	if err != nil {
		// Handle exceptions
		panic(err)
	}
	c.client = ecsClient

	return &c
}
