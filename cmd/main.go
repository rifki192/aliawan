package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/williamchanrico/alicloud-image-overwriter/config"
	"github.com/williamchanrico/alicloud-image-overwriter/ecs"
	"github.com/williamchanrico/alicloud-image-overwriter/ess"
	"github.com/williamchanrico/alicloud-image-overwriter/slb"
)

func main() {
	if len(os.Args[1:]) < 1 {
		fmt.Println("Please provide at least one argument, to see available argument just type -h argument")
		os.Exit(1)
	}

	fmt.Println("=================================================")
	fmt.Println("======    ALIBABA CLOUD CLI WRAPPER      ========")
	fmt.Println("======  another un-official alicloud-cli ========")
	fmt.Println("======      to simplify your task        ========")
	fmt.Println("=================================================")
	fmt.Println()

	cfg := config.LoadConfig()

	switch args := os.Args; args[1] {
	case "images":
		//imagesCommand
		imagesCommand(cfg)
	case "slb":
		//slbCommand
		slbCommand(cfg)
	default:
		fmt.Printf("%s. not defined please see help", args[1])
	}

	fmt.Println()
	os.Exit(0)
}

func slbCommand(cfg *config.Config) {
	var err error
	// Subcommands
	slbCmd := flag.NewFlagSet("slb", flag.ExitOnError)

	flagVGroups := slbCmd.String("vgroupname", "", "VServer Groups Name")
	flagInstanceID := slbCmd.String("instanceid", "", "Instance ID to be added to Vserver Group SLB")
	slbCmd.Parse(os.Args[2:])

	if *flagVGroups == "" {
		fmt.Println("Please provide VGroup Name with --vgroupname")
		os.Exit(1)
	}

	ecsClient := ecs.New(cfg)

	if *flagInstanceID == "" {
		*flagInstanceID = ecsClient.GetInstanceID()
	}

	slbClient := slb.New(cfg)
	err = slbClient.AddInstanceToVServerGroup(*flagVGroups, *flagInstanceID)
	if err != nil {
		log.Printf("could not send request AddInstanceToVServerGroup to alibaba: %v\n", err)
		os.Exit(1)
	}
}

func imagesCommand(cfg *config.Config) {
	var err error

	imagesCmd := flag.NewFlagSet("images", flag.ExitOnError)

	flagOldName := imagesCmd.String("oldname", "", "Old Image Name")
	flagNewName := imagesCmd.String("newname", "", "New Image Name")
	flagDeleteOld := imagesCmd.Bool("deleteold", false, "Delete Old Image")
	imagesCmd.Parse(os.Args[2:])

	if *flagNewName == "" {
		fmt.Println("Please provide new image name with --newname")
		os.Exit(1)
	}

	if *flagOldName == "" {
		fmt.Println("Please provide old image name with --oldname")
		os.Exit(1)
	}

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
	fmt.Printf("All feature using image %s (%s) has been replaced to use image %s (%s) \n", *flagOldName, oldImageID, *flagNewName, newImageID)

	fmt.Println("Change new image name, to become old image name")

	if oldImageID != "" {
		err = ecsClient.ChangeImageName(oldImageID, *flagOldName+"tmp")
		if err != nil {
			fmt.Printf("Error while change old image name %v \n", err)
			os.Exit(1)
		}
	}

	err = ecsClient.ChangeImageName(newImageID, *flagOldName)
	if err != nil {
		fmt.Printf("Error while change new image name %v \n", err)
		os.Exit(1)
	}

	if *flagDeleteOld && oldImageID != "" {
		fmt.Println("Delete Old Image Defined, will delete old image...")
		err = ecsClient.DeleteImageByID(oldImageID)
		if err != nil {
			fmt.Printf("Error while deleting old image %v \n", err)
			os.Exit(1)
		}
	}
}
