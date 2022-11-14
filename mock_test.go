package httpmocking

import (
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"
)

func TestFixedResponseRoundTripper(t *testing.T) {
	const fixedResponse = "fixed response"

	frt := &FixedResponseRoundTripper{RespBytes: []byte(fixedResponse)}
	c := &http.Client{
		Transport: frt,
	}
	uRL, err := url.Parse("http://test.com")

	if err != nil {
		t.Fatal("Should parse URL", err)
	}

	request := &http.Request{
		Method: http.MethodGet,
		URL:    uRL,
	}

	const header = "Header"

	request.Header = map[string][]string{header: {"headervalue"}}
	response, err := c.Do(request)

	if err != nil {
		t.Fatal("Should not err out", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Should return code 200 - actual: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("Should not err out", err)
	}

	if string(body) != fixedResponse {
		t.Fatal("Should give a fixed response")
	}

	if frt.Request.Header[header] == nil {
		t.Fatal("Should correctly record request header")
	}

	if len(frt.Request.Header[header]) != 1 {
		t.Fatal("Should correctly record request header")
	}

	if frt.Request.Header[header][0] != "headervalue" {
		t.Fatal("Should correctly record request header")
	}
}

func TestBadRequestRoundTripper(t *testing.T) {
	const fixedResponse = "fixed response"

	brt := &BadRequestRoundTripper{RespBytes: []byte(fixedResponse)}
	c := &http.Client{
		Transport: brt,
	}
	uRL, err := url.Parse("http://test.com")

	if err != nil {
		t.Fatal("Should parse URL", err)
	}

	request := &http.Request{
		Method: http.MethodGet,
		URL:    uRL,
	}

	const header = "Header"

	const headerValue = "headervalue"

	request.Header = map[string][]string{header: {headerValue}}
	response, err := c.Do(request)

	if err != nil {
		t.Fatal("Should not err out", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("Should return code 400 - actual: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("Should not err out", err)
	}

	if string(body) != fixedResponse {
		t.Fatal("Should give a fixed response")
	}

	if brt.Request.Header[header] == nil {
		t.Fatal("Should correctly record request header")
	}

	if len(brt.Request.Header[header]) != 1 {
		t.Fatal("Should correctly record request header")
	}

	if brt.Request.Header[header][0] != headerValue {
		t.Fatal("Should correctly record request header")
	}
}

func TestWaitGroupRoundTripper(t *testing.T) {
	const fixedResponse = "fixed response"

	wrt := &WaitGroupRoundTripper{RespBytes: []byte(fixedResponse), WG: &sync.WaitGroup{}}
	wrt.WG.Add(2)

	c := &http.Client{
		Transport: wrt,
	}
	uRL, err := url.Parse("http://test.com")

	if err != nil {
		t.Fatal("Should parse URL", err)
	}

	request := &http.Request{
		Method: http.MethodGet,
		URL:    uRL,
	}

	const header = "Header"

	const headerValue = "headervalue"

	request.Header = map[string][]string{header: {headerValue}}

	var response, response2 *http.Response

	var err2 error

	go func() {
		response, err = c.Do(request)
		defer response.Body.Close()
	}()
	go func() {
		response2, err2 = c.Do(request)
		defer response2.Body.Close()
	}()

	wrt.WG.Wait()

	if err != nil {
		t.Fatal("Should not err out", err)
	}

	if err2 != nil {
		t.Fatal("Should not err out", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Should return code 200 - actual: %d", response.StatusCode)
	}

	if response2.StatusCode != http.StatusOK {
		t.Fatalf("Should return code 200 - actual: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("Should not err out", err)
	}

	if string(body) != fixedResponse {
		t.Fatal("Should give a fixed response")
	}

	if wrt.Request.Header[header] == nil {
		t.Fatal("Should correctly record request header")
	}

	if len(wrt.Request.Header[header]) != 1 {
		t.Fatal("Should correctly record request header")
	}

	if wrt.Request.Header[header][0] != headerValue {
		t.Fatal("Should correctly record request header")
	}
}
