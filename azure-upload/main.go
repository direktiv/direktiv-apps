package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

type AzureBlobUpload struct {
	Container  string `json:"container"`
	Account    string `json:"storage-account"`
	AccountKey string `json:"storage-account-key"`
	Data       string `json:"data"`
	Blobname   string `json:"upload-name"`
}

type AzureBlobOutput struct {
	BlobURL string `json:"blob-url"`
}

const code = "com.azure-upload.error"

func main() {
	direktivapps.StartServer(AzureUpload)
}

func AzureUpload(w http.ResponseWriter, r *http.Request) {
	obj := new(AzureBlobUpload)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	ctx := context.Background()

	os.Setenv("AZURE_STORAGE_ACCOUNT", obj.Account)
	os.Setenv("AZURE_STORAGE_ACCESS_KEY", obj.AccountKey)

	credential, err := azblob.NewSharedKeyCredential(obj.Account, obj.AccountKey)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	decoded, err := base64.StdEncoding.DecodeString(obj.Data)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", obj.Account, obj.Container))
	containerURL := azblob.NewContainerURL(*u, p)
	blobURL := containerURL.NewBlockBlobURL(obj.Blobname)

	err = ioutil.WriteFile(obj.Blobname, decoded, 0700)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	file, err := os.Open(obj.Blobname)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	defer file.Close()

	contentType, err := GetFileContentType(file)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	file.Seek(0, 0)

	// Upload file
	resp, err := azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: contentType,
		},
	})
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	httpR := resp.Response()
	defer httpR.Body.Close()

	data, err := ioutil.ReadAll(httpR.Body)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, data)
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
