package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rifki192/alicloud-image-overwriter/config"
	"github.com/rifki192/alicloud-image-overwriter/ecs"
	"github.com/rifki192/alicloud-image-overwriter/ess"
)

func main() {
	var err error

	flagOldName := flag.String("oldname", "", "Old Image Name")
	flagNewName := flag.String("newname", "", "New Image Name")
	flagDeleteOld := flag.Bool("deleteold", false, "Delete Old Image")
	flag.Parse()

	if *flagNewName == "" {
		fmt.Println("Please provide new image name with --newname")
		os.Exit(1)
	}

	if *flagOldName == "" {
		fmt.Println("Please provide old image name with --oldname")
		os.Exit(1)
	}

	fmt.Println("=============================================")
	fmt.Println("======ALIBABA CLOUD IMAGES OVERWRITER========")
	fmt.Println("======  for replacing image used by  ========")
	fmt.Println("======  any resources on ali cloud   ========")
	fmt.Println("=============================================")
	fmt.Println()

	cfg := config.LoadConfig()

	ecsClient := ecs.New(cfg)
	oldImageID := ecsClient.GetImageIdByName(*flagOldName)
	fmt.Printf("Will replace image %s (%s)\n", *flagOldName, oldImageID)
	newImageID := ecsClient.GetImageIdByName(*flagNewName)
	fmt.Printf("With image %s (%s)\n", *flagNewName, newImageID)

	essClient := ess.New(cfg)
	err = essClient.ReplaceScalingConfigurationsWithImageId(oldImageID, newImageID)
	if err != nil {
		fmt.Printf("Error while replacing scaling group config %v \n", err)
		os.Exit(1)
	}
	fmt.Println("All feature using image %s (%s) has been replaced to use image %s (%s)", *flagOldName, oldImageID, *flagNewName, newImageID)

	if *flagDeleteOld {
		fmt.Println("Delete Old Image Defined, will delete old image...")
		err = ecsClient.DeleteImageByID(oldImageID)
		fmt.Printf("Error while deleting old image %v \n", err)
		os.Exit(1)
	}

	fmt.Println()
	os.Exit(0)
}
