package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/williamchanrico/aliawan/config"
	"github.com/williamchanrico/aliawan/ecs"
	"github.com/williamchanrico/aliawan/ess"
	"github.com/williamchanrico/aliawan/slb"
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
	fmt.Println("====== v1.3                              ========")
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

func disableRTCAndSyncTime() error {
	_, err := exec.LookPath("timedatectl")
	if err != nil {
		return err
	}
	err = exec.Command("timedatectl", "set-local-rtc", "0").Run()
	if err != nil {
		return err
	}

	_, err = exec.LookPath("ntpd")
	if err != nil {
		return err
	}
	err = exec.Command("ntpdate", "-u", "time.google.com").Run()
	if err != nil {
		return err
	}

	return nil
}

func slbCommand(cfg *config.Config) {
	fmt.Println("Disabling RTC in local tz and sync system time")
	if err := disableRTCAndSyncTime(); err != nil {
		fmt.Println(err)
	}

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
		fmt.Printf("could not send request AddInstanceToVServerGroup to alibaba: %v\n", err)
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
		fmt.Printf("Error while replacing scaling group config %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All feature using image %s (%s) has been replaced to use image %s (%s)\n", *flagOldName, oldImageID, *flagNewName, newImageID)

	fmt.Println("Deleting junk image...")
	junkImageNameSuffix := "-should-be-deleted"
	junkImageID := ecsClient.GetImageIdByName(*flagOldName + junkImageNameSuffix)
	err = ecsClient.DeleteImageByID(junkImageID)
	if err != nil {
		fmt.Printf("Looks like there's no junk images: %v\n", err)
	}

	if oldImageID != "" {
		fmt.Println("Changing old image name to junk name")
		err = ecsClient.ChangeImageName(oldImageID, *flagOldName+junkImageNameSuffix)
		if err != nil {
			fmt.Printf("Error while change old image name %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Change new image name, to become old image name")
	err = ecsClient.ChangeImageName(newImageID, *flagOldName)
	if err != nil {
		fmt.Printf("Error while change new image name %v\n", err)
		os.Exit(1)
	}

	if *flagDeleteOld && oldImageID != "" {
		fmt.Println("Delete Old Image Defined, will delete old image...")
		err = ecsClient.DeleteImageByID(oldImageID)
		if err != nil {
			fmt.Printf("Error while deleting old image %v\n", err)
			os.Exit(1)
		}
	}
}
