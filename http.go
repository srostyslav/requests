package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Request struct {
	Url             string
	Headers, Params map[string]string
	Body            interface{}
	Json            bool

	ResponseBody []byte
	ResponseCode int

	resp   *http.Response
	req    *http.Request
	client *http.Client
}

func (r *Request) init() {
	if r.client == nil {
		r.client = &http.Client{}
	}
}

func (r *Request) SetClient(c *http.Client) {
	r.client = c
}

func (r *Request) send() (err error) {
	if r.Json {
		r.req.Header.Add("content-type", "application/json")
	}

	for k, v := range r.Headers {
		r.req.Header.Add(k, v)
	}

	if r.resp, err = r.client.Do(r.req); err != nil {
		return err
	}

	if r.ResponseBody, err = ioutil.ReadAll(r.resp.Body); err != nil {
		return err
	} else {
		r.resp.Body.Close()
		r.ResponseCode = r.resp.StatusCode
		return nil
	}
}

func (r *Request) Put() (err error) {
	r.init()

	if r.Body != nil {
		if output, err := json.Marshal(&r.Body); err != nil {
			return err
		} else if r.req, err = http.NewRequest(http.MethodPut, r.Url, bytes.NewReader(output)); err != nil {
			return err
		}
	} else if r.req, err = http.NewRequest(http.MethodPut, r.Url, nil); err != nil {
		return err
	}

	r.setQueryParams()

	return r.send()
}

func (r *Request) Patch() (err error) {
	r.init()

	if r.Body != nil {
		if output, err := json.Marshal(&r.Body); err != nil {
			return err
		} else if r.req, err = http.NewRequest(http.MethodPatch, r.Url, bytes.NewReader(output)); err != nil {
			return err
		}
	} else if r.req, err = http.NewRequest(http.MethodPatch, r.Url, nil); err != nil {
		return err
	}

	r.setQueryParams()

	return r.send()
}

func (r *Request) Post() (err error) {
	r.init()

	if r.Body != nil {
		if output, err := json.Marshal(&r.Body); err != nil {
			return err
		} else if r.req, err = http.NewRequest(http.MethodPost, r.Url, bytes.NewReader(output)); err != nil {
			return err
		}
	} else if r.req, err = http.NewRequest(http.MethodPost, r.Url, nil); err != nil {
		return err
	}

	r.setQueryParams()

	return r.send()
}

func (r *Request) Decode(obj interface{}) error {
	return json.Unmarshal(r.ResponseBody, obj)
}

func (r *Request) Try(cnt int, method string) (err error) {
	var f func() error

	switch method {
	case http.MethodPost:
		f = r.Post
	case http.MethodGet:
		f = r.Get
	case http.MethodDelete:
		f = r.Delete
	case http.MethodPatch:
		f = r.Patch
	case http.MethodPut:
		f = r.Put
	default:
		return fmt.Errorf("method %s is not found", method)
	}

	for i := 0; i < cnt; i++ {
		if err = f(); err == nil {
			break
		} else if i+1 < cnt {
			time.Sleep(1 * time.Second)
		}
	}
	return err
}

func (r *Request) setQueryParams() {
	q := r.req.URL.Query()
	for k, v := range r.Params {
		q.Add(k, v)
	}
	r.req.URL.RawQuery = q.Encode()
}

func (r *Request) Get() (err error) {
	r.init()

	if r.req, err = http.NewRequest(http.MethodGet, r.Url, nil); err != nil {
		return err
	}
	r.setQueryParams()

	return r.send()
}

func (r *Request) Delete() (err error) {
	r.init()

	if r.req, err = http.NewRequest(http.MethodDelete, r.Url, nil); err != nil {
		return err
	}
	r.setQueryParams()

	return r.send()
}
