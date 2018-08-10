package slb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

type BackendServer struct {
	Port     int    `json:"Port" xml:"Port"`
	ServerId string `json:"ServerId" xml:"ServerId"`
	Weight   int    `json:"Weight" xml:"Weight"`
}

const (
	Port   = 80
	Weight = 100
)

func getVServerGroupsIdByVServerName(c *Client, vServerName string) []string {
	var Vservers []string
	var err error

	request := slb.CreateDescribeVServerGroupsRequest()
	response := slb.CreateDescribeVServerGroupsResponse()

	AllSLBs := getAllSLBs(c)

	for _, slbId := range AllSLBs {
		request.LoadBalancerId = slbId
		response, err = c.client.DescribeVServerGroups(request)
		if err != nil {
			// Handle exceptions
			log.Printf("could not send request DescribeLoadBalancers to alibaba: %s", err)
			os.Exit(1)
		}
		for _, vs := range response.VServerGroups.VServerGroup {
			if vs.VServerGroupName == vServerName {
				Vservers = append(Vservers, vs.VServerGroupId)
			}

		}
	}

	return Vservers
}

func (c *Client) AddInstanceToVServerGroup(vServerName string, instanceID string) error {
	var err error
	var backendServer BackendServer
	var backendServers []BackendServer

	request := slb.CreateAddVServerGroupBackendServersRequest()
	response := slb.CreateAddVServerGroupBackendServersResponse()

	vServerGroups := getVServerGroupsIdByVServerName(c, vServerName)

	for _, vg := range vServerGroups {
		request.VServerGroupId = vg
		backendServer.ServerId = instanceID
		backendServer.Port = Port
		backendServer.Weight = Weight
		backendServers = append(backendServers, backendServer)
		sJson, _ := json.Marshal(backendServers)
		request.BackendServers = string(sJson)
		response, err = c.client.AddVServerGroupBackendServers(request)
		if err != nil {
			return err
		}
		if response.IsSuccess() {
			fmt.Printf("Instance %s success added to %s with name %s.\n", instanceID, vg, vServerName)
		}
	}

	return nil
}
