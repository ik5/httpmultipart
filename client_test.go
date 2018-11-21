package httpmultipart

import (
	"testing"
)

func TestPostRequest(t *testing.T) {
	if testing.Short() {
		return
	}
	address := "http://localhost:8081"
	params := InitParams()
	params.AddString("version", "2")
	params.AddString("language", "he")
	params.AddReadFile("file", "client_test.go")
	resp, err := PostRequest(address, params)
	if err != nil {
		t.Error(err)
	}
	if resp.Body != nil {
		resp.Body.Close()
	}
}
