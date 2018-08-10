package ecs

import (
	"fmt"
	"log"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func (c *Client) GetImageIdByName(imageName string) string {
	var imageId string
	request := ecs.CreateDescribeImagesRequest()
	request.PageSize = requests.NewInteger(5)
	request.ImageName = imageName

	response, err := c.client.DescribeImages(request)
	if err != nil {
		// Handle exceptions
		log.Printf("could not send request to alibaba: %s", err)
		os.Exit(1)
	}
	if len(response.Images.Image) > 0 {
		imageId = response.Images.Image[0].ImageId
	}

	return imageId
}

func (c *Client) DeleteImageByID(imageID string) error {
	var err error
	request := ecs.CreateDeleteImageRequest()
	response := ecs.CreateDeleteImageResponse()

	request.ImageId = imageID
	request.Force = requests.NewBoolean(true)
	response, err = c.client.DeleteImage(request)
	if err != nil {
		// Handle exceptions
		return fmt.Errorf("could not send request to alibaba: %s", err)
	}
	if response.IsSuccess() {
		fmt.Printf("Image with ID %s, has been deleted.\n", imageID)
	}
	return nil
}

func (c *Client) ChangeImageName(imageID string, imageName string) error {
	var err error
	request := ecs.CreateModifyImageAttributeRequest()
	response := ecs.CreateModifyImageAttributeResponse()

	request.ImageId = imageID
	request.ImageName = imageName
	response, err = c.client.ModifyImageAttribute(request)
	if err != nil {
		// Handle exceptions
		return fmt.Errorf("could not send request to alibaba: %s", err)
	}
	if response.IsSuccess() {
		fmt.Printf("Image with ID %s, has been renamed to %s.\n", imageID, imageName)
	}
	return nil
}
