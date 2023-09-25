package toolkit

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	s := testTools.RandomString(10)
	if len(s) != 10 {
		t.Error("Wrong length random string returned")
	}
}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{
		name:          "Allowed no rename",
		allowedTypes:  []string{"image/jpeg", "image/png", "image/gif"},
		renameFile:    false,
		errorExpected: false,
	}, {
		name:          "Allowed rename",
		allowedTypes:  []string{"image/jpeg", "image/png", "image/gif"},
		renameFile:    true,
		errorExpected: false,
	}, {
		name:          "Not Allowed",
		allowedTypes:  []string{"image/png", "image/gif"},
		renameFile:    false,
		errorExpected: true,
	},
}

func TestTools_Upload(t *testing.T) {
	for _, e := range uploadTests {
		// set up a pipe to avoid buffering
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			// create the form data field 'file'
			part, err := writer.CreateFormFile("file", "./testdata/img.jpg")
			if err != nil {
				t.Error(err)
			}

			f, err := os.Open("./testdata/img.jpg")
			if err != nil {
				t.Error(err)
			}

			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error("error decoding image", err)
			}

			err = jpeg.Encode(part, img, nil)
			if err != nil {
				t.Error("error encoding image", err)
			}
		}()
		// read from the pipe which receives the form data
		request := httptest.NewRequest(http.MethodPost, "/", pr)
		request.Header.Set("Content-Type", writer.FormDataContentType())

		var testTools Tools
		testTools.AllowedFileTypes = e.allowedTypes

		uploadedFiles, err := testTools.UploadFiles(request, "./testdata/uploads/", e.renameFile)
		if err != nil && !e.errorExpected {
			t.Error(err)
		}
		if !e.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName)); os.IsNotExist(err) {
				t.Errorf("%s: expected file does not exist: %s", e.name, err.Error())
			}

			//clean up
			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName))
		}

		if !e.errorExpected && err != nil {
			t.Errorf("%s: error expected but not recieved", e.name)
		}

		wg.Wait()

	}
}

func TestTools_UploadOneFile(t *testing.T) {
	// set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		// create the form data field 'file'
		part, err := writer.CreateFormFile("file", "./testdata/img.jpg")
		if err != nil {
			t.Error(err)
		}

		f, err := os.Open("./testdata/img.jpg")
		if err != nil {
			t.Error(err)
		}

		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			t.Error("error decoding image", err)
		}

		err = jpeg.Encode(part, img, nil)
		if err != nil {
			t.Error("error encoding image", err)
		}
	}()
	// read from the pipe which receives the form data
	request := httptest.NewRequest(http.MethodPost, "/", pr)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	var testTools Tools

	uploadedFile, err := testTools.UploadOneFile(request, "./testdata/uploads/", true)
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFile.NewFileName)); os.IsNotExist(err) {
		t.Errorf("expected file does not exist: %s", err.Error())
	}

	//clean up
	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFile.NewFileName))
}
