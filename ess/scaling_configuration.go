package ess

import (
	"fmt"
	"math"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

const (
	pageSize = 10
)

func (c *Client) ReplaceScalingConfigurationsWithImageId(oldImageID string, newImageID string) error {
	scList := getScalingConfigurationsByImageID(c, oldImageID)

	if len(scList) == 0 {
		fmt.Printf("Can not find any resource using this image: %s\n", oldImageID)
		return nil
	}

	for _, item := range scList {
		err := changeScalingGroupImageID(c, item, newImageID)
		if err != nil {
			// Handle exceptions
			fmt.Printf("could not send request to alibaba: %s", err)
			os.Exit(1)
		}
	}
	return nil
}

func changeScalingGroupImageID(c *Client, scID string, imageID string) error {
	var err error

	request := ess.CreateModifyScalingConfigurationRequest()
	response := ess.CreateModifyScalingConfigurationResponse()

	request.ScalingConfigurationId = scID
	request.ImageId = imageID
	response, err = c.client.ModifyScalingConfiguration(request)
	if response.IsSuccess() {
		fmt.Printf("Scaling group configuration %s replaced with new image id %s\n", scID, imageID)
	}
	if err != nil {
		// Handle exceptions
		return fmt.Errorf("could not send request to alibaba: %s", err)
	}
	return nil
}

func getScalingConfigurationsByImageID(c *Client, imageId string) []string {
	var scalingConfList []string
	var err error

	request := ess.CreateDescribeScalingConfigurationsRequest()
	response := ess.CreateDescribeScalingConfigurationsResponse()

	// Set the request.PageSize
	request.PageSize = requests.NewInteger(pageSize)
	totalPages := 1
	for i := 1; i <= totalPages; i++ {
		request.PageNumber = requests.NewInteger(int(i))
		response, err = c.client.DescribeScalingConfigurations(request)
		if err != nil {
			// Handle exceptions
			fmt.Printf("could not send request to alibaba: %s", err)
			os.Exit(1)
		}
		totalInstances := response.TotalCount
		totalPages = int(math.Ceil(float64(totalInstances) / float64(pageSize)))
		for _, item := range response.ScalingConfigurations.ScalingConfiguration {

			//Check if imageId is same as the parameter
			//And add to the array list
			if item.ImageId == imageId {
				scalingConfList = append(scalingConfList, item.ScalingConfigurationId)
			}
		}
	}

	return scalingConfList
}
