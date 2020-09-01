package apply_test

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/raba-jp/primus/exec"
	"github.com/raba-jp/primus/executor"
	"github.com/raba-jp/primus/executor/apply"
	"github.com/spf13/afero"
)

type MockRoundTripper struct {
	http.RoundTripper
	Fn func(req *http.Request) *http.Response
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Fn(req), nil
}

func MockHttpClient(fn func(req *http.Request) *http.Response) *http.Client {
	return &http.Client{
		Transport: &MockRoundTripper{Fn: fn},
	}
}

var in = new(bytes.Buffer)
var out = new(bytes.Buffer)
var errout = new(bytes.Buffer)
var fs = afero.NewMemMapFs()
var execIF exec.Interface

func NewTestExecutor() executor.Executor {
	client := MockHttpClient(func(req *http.Request) *http.Response {
		buf := bytes.NewBufferString("test file")
		body := ioutil.NopCloser(buf)
		return &http.Response{
			Body: body,
		}
	})

	return apply.NewApplyExecutorWithArgs(in, out, errout, execIF, fs, client)
}

func Reset() {
	in.Reset()
	out.Reset()
	errout.Reset()
	fs = afero.NewMemMapFs()
	execIF = nil
}
