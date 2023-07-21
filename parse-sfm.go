package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strconv"
)

type SfmPose struct {
	PoseId string `json:"poseId"`
	Pose   struct {
		Transform struct {
			Rotation []string `json:"rotation"`
			Center   []string `json:"center"`
		} `json:"transform"`
		Locked string `json:"locked"`
	} `json:"pose"`
}

type SfmIntrinsic struct {
	IntrinsicId                  string   `json:"intrinsicId"`
	Width                        string   `json:"width"`
	Height                       string   `json:"height"`
	SensorWidth                  string   `json:"sensorWidth"`
	SensorHeight                 string   `json:"sensorHeight"`
	SerialNumber                 string   `json:"serialNumber"`
	Type                         string   `json:"type"`
	InitializationMode           string   `json:"initializationMode"`
	InitialFocalLength           string   `json:"initialFocalLength"`
	FocalLength                  string   `json:"focalLength"`
	PixelRatio                   string   `json:"pixelRatio"`
	PixelRatioLocked             string   `json:"pixelRatioLocked"`
	PrincipalPoint               []string `json:"principalPoint"`
	DistortionInitializationMode string   `json:"distortionInitializationMode"`
	DistortionParams             []string `json:"distortionParams"`
	Locked                       string   `json:"locked"`
}

type SfmViewMetadata struct {
	ImageDescription string `json:"ImageDescription"`
	JpegSubsampling  string `json:"jpeg:subsampling"`
	OiioColorSpace   string `json:"oiio:ColorSpace"`
}

type SfmView struct {
	ViewId      string          `json:"viewId"`
	PoseId      string          `json:"poseId"`
	FrameId     string          `json:"frameId"`
	IntrinsicId string          `json:"intrinsicId"`
	ResectionId string          `json:"resectionId"`
	Path        string          `json:"path"`
	Width       string          `json:"width"`
	Height      string          `json:"height"`
	Metadata    SfmViewMetadata `json:"metadata"`
}

type SfmFile struct {
	Version         []string       `json:"version"`
	FeaturesFolders []string       `json:"featuresFolders"`
	MatchesFolders  []string       `json:"matchesFolders"`
	Views           []SfmView      `json:"views"`
	Intrinsics      []SfmIntrinsic `json:"intrinsics"`
	Poses           []SfmPose      `json:"poses"`
}

type Pose struct {
	Frame     int `json:"frame"`
	Transform struct {
		Rotation [9]float64 `json:"rotation"`
		Center   [3]float64 `json:"center"`
	} `json:"transform"`
}

type Output struct {
	Data []Pose `json:"data"`
}

func main() {

	args := os.Args
	if len(args) != 3 {
		fmt.Println("Invalid input")
		fmt.Println("Usage: parse-sfm <filename.sfm> <output.json>")
		os.Exit(1)
	}

	sfmBytes, err := os.ReadFile(args[1])
	if err != nil {
		fmt.Println(err)
	}

	var sfmData SfmFile

	err = json.Unmarshal(sfmBytes, &sfmData)
	if err != nil {
		fmt.Println(err)
	}

	var output Output
	for _, view := range sfmData.Views {
		frame, err := strconv.Atoi(view.FrameId)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var pose Pose
		pose.Frame = frame

		var poseIndex = -1
		for i := range sfmData.Poses {
			if sfmData.Poses[i].PoseId == view.PoseId {
				poseIndex = i
				break
			}
		}
		if poseIndex == -1 {
			fmt.Println("Cannot find pose")
			continue
		}

		for i := 0; i < 9; i++ {
			val, err := strconv.ParseFloat(sfmData.Poses[poseIndex].Pose.Transform.Rotation[i], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			pose.Transform.Rotation[i] = val
		}

		for i := 0; i < 3; i++ {
			val, err := strconv.ParseFloat(sfmData.Poses[poseIndex].Pose.Transform.Center[i], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			pose.Transform.Center[i] = val
		}

		output.Data = append(output.Data, pose)
	}

	outputBytes, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile(args[2], outputBytes, fs.ModeAppend)
}
