package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type request struct {
	Image string `json:"image"` // url of the image
}

var code = "com.blurimage.error"

func getFileContentType(data []byte) (string, error) {

	if len(data) < 512 {
		return "", errors.New("the length of the image is less than 512 bytes unable to check image type")
	}
	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(data)

	return contentType, nil
}

func BlurHandler(w http.ResponseWriter, r *http.Request) {
	var o request
	aid, err := direktivapps.Unmarshal(&o, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Log(aid, "reading image from url")

	resp, err := http.Get(o.Image)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Log(aid, fmt.Sprintf("status %s", resp.Status))

	direktivapps.Log(aid, "detecting content type")
	ct, err := getFileContentType(data[:512])
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Log(aid, fmt.Sprintf("content-type: '%s' detected", ct))

	newImage, err := os.Create(aid)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	defer newImage.Close()

	switch ct {
	case "image/png":
		direktivapps.LogDouble(aid, "decoding png...")
		srcImage, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
		direktivapps.LogDouble(aid, "decoded png")
		dstImage := image.NewNRGBA(srcImage.Bounds())

		graphics.Blur(dstImage, srcImage, &graphics.BlurOptions{StdDev: 10})
		direktivapps.LogDouble(aid, fmt.Sprintf("dst img: %+v", dstImage))

		err = png.Encode(newImage, dstImage)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	case "image/jpeg":
		direktivapps.Log(aid, "decoding jpeg...")

		srcImage, err := jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}

		direktivapps.LogDouble(aid, "decoded jpeg")
		dstImage := image.NewNRGBA(srcImage.Bounds())

		graphics.Blur(dstImage, srcImage, &graphics.BlurOptions{StdDev: 10})

		err = jpeg.Encode(newImage, dstImage, &jpeg.Options{jpeg.DefaultQuality})
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	default:
		direktivapps.RespondWithError(w, code, fmt.Sprintf("'jpeg' and 'png' only supported for blurring."))
		return
	}

	data, err = ioutil.ReadFile(aid)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	direktivapps.Respond(w, data)
}

func main() {
	direktivapps.StartServer(BlurHandler)
}
