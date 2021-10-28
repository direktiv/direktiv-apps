package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
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

const credFile = "/creds"
const code = "com.imagerecognition.error"

func ImageRecognition(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	obj := new(VisionAPIRecognition)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	err = ioutil.WriteFile(credFile, []byte(obj.ServiceAccountKey), 0777)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	visionClient, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	img := vision.NewImageFromURI(obj.URL)

	resp, err := visionClient.DetectSafeSearch(ctx, img, nil)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

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

	direktivapps.Respond(w, detailData)
}

func main() {
	direktivapps.StartServer(ImageRecognition)
}
