package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"

	"github.com/fogleman/gg"

	_ "image/jpeg"
	_ "image/png"
)

const (
	code = "com.img.watermark.error"
)

type imgInfo struct {
	txt, color string
	img        []byte
}

func checkReqFields(data []byte) (imgInfo, error) {

	var (
		i imgInfo
		m map[string]string
	)

	err := json.Unmarshal(data, &m)
	if err != nil {
		return i, err
	}

	txt, ok := m["text"]
	if !ok {
		return i, fmt.Errorf("field 'text' is missing in payload")
	}
	i.txt = txt

	img, ok := m["img"]
	if !ok {
		return i, fmt.Errorf("field 'img' is missing in payload")
	}

	dec, err := base64.StdEncoding.DecodeString(img)
	if err != nil {
		return i, err
	}
	i.img = dec

	color, ok := m["color"]
	if !ok {
		color = "#000000FF"
	}
	i.color = color

	return i, nil

}

func request(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error watermarking image: %v\n", err)
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	info, err := checkReqFields(body)
	if err != nil {
		log.Printf("error watermarking image: %v\n", err)
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	b, err := watermarkImage(info.img, info.txt, info.color)
	if err != nil {
		log.Printf("error watermarking image: %v\n", err)
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	w.Write(b)

}

func watermarkImage(img []byte, txt, textColor string) ([]byte, error) {

	buf := new(bytes.Buffer)

	im, _, err := image.DecodeConfig(bytes.NewReader(img))
	if err != nil {
		return buf.Bytes(), err
	}

	dc := gg.NewContext(im.Width, im.Height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	m, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return buf.Bytes(), err
	}
	dc.DrawImage(m, 0, 0)

	maxWidth := float64(im.Width) * 0.60
	fs := 200
	for {
		dc.LoadFontFace("/usr/share/fonts/truetype/liberation/LiberationMono-Regular.ttf", float64(fs))
		w, _ := dc.MeasureString(txt)
		if w < maxWidth || fs == 10 {
			break
		}
		fs = fs - 2
	}

	dc.RotateAbout(gg.Radians(325), float64(im.Width/2), float64(im.Height/2))
	dc.SetHexColor(textColor)
	dc.DrawStringAnchored(txt, float64(im.Width/2), float64(im.Height/2), 0.5, 0.5)

	dc.EncodePNG(buf)

	return buf.Bytes(), nil

}

func main() {
	direktivapps.StartServer(request)
}
