package httpmultipart

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const maxFileSize = 1024 * 1024 * 1024

// Params hold information about how to add multipart content to HTTP client
type Params struct {
	body *bytes.Buffer
	w    *multipart.Writer
}

// InitParams initialize Params struct
func InitParams() Params {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	return Params{body, w}
}

// AddString add/change a string parameter to the Params struct
// If name already exists it will override existed parameter.
func (p *Params) AddString(name, value string) (bool, error) {
	err := p.w.WriteField(name, value)
	if err != nil {
		return false, err
	}
	return true, nil
}

// AddReadFile read a content of a file and set it content as a parameter.
// If something went wrong, an error is returned.
func (p *Params) AddReadFile(name, filename string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return false, err
	}

	stat, err := f.Stat()
	if err != nil {
		return false, err
	}
	size := stat.Size()
	if size > maxFileSize {
		return false, fmt.Errorf("File size %d is too big, max size is : %d", size, maxFileSize)
	}

	part, err := p.w.CreateFormFile(name, filepath.Base(f.Name()))
	if err != nil {
		return false, err
	}
	io.Copy(part, f)
	return true, nil
}

// PostRequest send a request based on method, address and params
func PostRequest(address string, params Params) (*http.Response, error) {
	params.w.Close() // For ending boundary
	req, err := http.NewRequest("POST", address, params.body)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	req.Header.Set("Content-Type", params.w.FormDataContentType())
	return client.Do(req)
}
