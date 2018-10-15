package slb

import (
	"fmt"
	"math"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

func getAllSLBs(c *Client) []string {
	var SLBs []string
	var err error
	request := slb.CreateDescribeLoadBalancersRequest()
	response := slb.CreateDescribeLoadBalancersResponse()

	// Set the request.PageSize
	request.PageSize = requests.NewInteger(pageSize)
	request.Domain = "slb.aliyuncs.com"
	fmt.Printf("DescribeLoadBalancers request domain is set to %s\n", request.Domain)
	totalPages := 1
	for i := 1; i <= totalPages; i++ {
		request.PageNumber = requests.NewInteger(int(i))
		response, err = c.client.DescribeLoadBalancers(request)
		// fmt.Printf("Requested for page number %d with %d data \n", i, response.TotalCount)
		if err != nil {
			// Handle exceptions
			fmt.Printf("could not send request 'DescribeLoadBalancers' to alibaba: %s", err)
			os.Exit(1)
		}
		totalSLBs := response.TotalCount
		totalPages = int(math.Ceil(float64(totalSLBs) / float64(pageSize)))
		for _, item := range response.LoadBalancers.LoadBalancer {
			SLBs = append(SLBs, item.LoadBalancerId)
		}
	}

	return SLBs
}
