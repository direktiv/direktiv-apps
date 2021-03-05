package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"google.golang.org/api/option"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

type VisionAPIRecognition struct {
	ServiceAccountKey string `json:"serviceAccountKey"`
	URL               string `json:"url"`
}

type Details struct {
	SafeForWork bool                `json:"safeForWork"`
	Racy        visionpb.Likelihood `json:"racyLikelihood"`
	Adult       visionpb.Likelihood `json:"adultLikelihood"`
	Violence    visionpb.Likelihood `json:"violenceLikelihood"`
}

const credFile = "/tmp/creds"

func main() {

	g := direktivapps.ActionError{
		ErrorCode:    "com.imagerecognition.error",
		ErrorMessage: "",
	}
	var err error
	ctx := context.Background()
	obj := new(VisionAPIRecognition)

	direktivapps.ReadIn(obj, g)
	fmt.Println("a")

	err = ioutil.WriteFile(credFile, []byte(obj.ServiceAccountKey), 0777)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	fmt.Println("b")

	visionClient, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	img := vision.NewImageFromURI(obj.URL)
	fmt.Println("c")

	resp, err := visionClient.DetectSafeSearch(ctx, img, nil)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}
	fmt.Println("d")

	var sfw bool
	if resp.GetRacy() == visionpb.Likelihood_VERY_LIKELY || resp.GetRacy() == visionpb.Likelihood_LIKELY || resp.GetAdult() == visionpb.Likelihood_VERY_LIKELY ||
		resp.GetViolence() == visionpb.Likelihood_VERY_LIKELY || resp.GetAdult() == visionpb.Likelihood_LIKELY || resp.GetViolence() == visionpb.Likelihood_LIKELY {
		sfw = false
	} else {
		sfw = true
	}

	detailData, err := json.Marshal(&Details{
		SafeForWork: sfw,
		Adult:       resp.GetAdult(),
		Violence:    resp.GetViolence(),
		Racy:        resp.GetRacy(),
	})

	direktivapps.WriteOut(detailData, g)
}
