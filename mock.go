package httpmocking

import (
	"bytes"
	"io"
	"net/http"
	"sync"
)

type FixedResponseRoundTripper struct {
	Request   *http.Request
	response  *http.Response
	RespBytes []byte
}

func (f *FixedResponseRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	f.Request = r
	byteReader := io.NopCloser(bytes.NewReader(f.RespBytes))
	f.response = &http.Response{}
	f.response.Body = byteReader
	f.response.StatusCode = http.StatusOK

	return f.response, nil
}

type BadRequestRoundTripper struct {
	Request   *http.Request
	response  *http.Response
	RespBytes []byte
}

func (b *BadRequestRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	b.Request = r
	byteReader := io.NopCloser(bytes.NewReader(b.RespBytes))
	b.response = &http.Response{}
	b.response.Body = byteReader
	b.response.StatusCode = http.StatusBadRequest

	return b.response, nil
}

type WaitGroupRoundTripper struct {
	Request   *http.Request
	response  *http.Response
	RespBytes []byte
	WG        *sync.WaitGroup
}

func (w *WaitGroupRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	w.Request = r
	byteReader := io.NopCloser(bytes.NewReader(w.RespBytes))
	w.response = &http.Response{}
	w.response.Body = byteReader
	w.response.StatusCode = http.StatusOK
	w.WG.Done()

	return w.response, nil
}
