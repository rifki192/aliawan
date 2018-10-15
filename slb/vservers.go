package slb

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

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

	// Create look-up hash table to check if a vServerName exists
	vServerNames := make(map[string]bool)
	vServerNameList := strings.Split(vServerName, ",")
	for _, name := range vServerNameList {
		vServerNames[name] = true
	}

	for _, slbId := range AllSLBs {
		request.LoadBalancerId = slbId
		response, err = c.client.DescribeVServerGroups(request)
		if err != nil {
			// Handle exceptions
			fmt.Printf("could not send request DescribeVServerGroups to alibaba: %s", err)
			os.Exit(1)
		}
		for _, vs := range response.VServerGroups.VServerGroup {
			if vServerNames[vs.VServerGroupName] {
				Vservers = append(Vservers, vs.VServerGroupId)
			}

		}
	}

	fmt.Printf("Found %d vservergroup with name %s\n", len(Vservers), vServerName)
	fmt.Println(Vservers)

	return Vservers
}

func (c *Client) AddInstanceToVServerGroup(vServerName string, instanceID string) error {
	var err error

	request := slb.CreateAddVServerGroupBackendServersRequest()
	response := slb.CreateAddVServerGroupBackendServersResponse()

	vServerGroups := getVServerGroupsIdByVServerName(c, vServerName)

	for _, vg := range vServerGroups {
		var backendServer BackendServer
		var backendServers []BackendServer
		request.VServerGroupId = vg
		backendServer.ServerId = instanceID
		backendServer.Port = Port
		backendServer.Weight = Weight
		backendServers = append(backendServers, backendServer)
		sJSON, _ := json.Marshal(backendServers)
		request.BackendServers = string(sJSON)
		fmt.Printf("Adding this backend server: %v\n", request.BackendServers)

		retryCount := 0
		for response, err = c.client.AddVServerGroupBackendServers(request); err != nil && retryCount < 10; {
			retryCount++
			fmt.Printf("Retry count: %v\n", retryCount)
			time.Sleep(10 * time.Second)

			response, err = c.client.AddVServerGroupBackendServers(request)
		}
		if err != nil {
			return err
		}
		if response.IsSuccess() {
			fmt.Printf("Instance [%s] successfully added to %s vgroup (%s).\n", instanceID, vServerName, vg)
		}
	}

	return nil
}
