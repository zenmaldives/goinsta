package goinsta

import (
	"fmt"
	"sync"

	"github.com/erikdubbelboer/fasthttp"
	"github.com/valyala/bytebufferpool"
)

type reqOptions struct {
	endpoint []byte

	// post data
	postData []byte

	skipStatus bool
	args       *fasthttp.Args
}

func (req *reqOptions) SetEndpointBytes(b []byte) {
	req.endpoint = append(req.endpoint[:0], b...)
}

func (req *reqOptions) SetEndpoint(b string) {
	req.endpoint = append(req.endpoint[:0], b...)
}

func (req *reqOptions) SetDataBytes(b []byte) {
	req.postData = append(req.postData[:0], b...)
}

func (req *reqOptions) SetData(b string) {
	req.postData = append(req.postData[:0], b...)
}

var reqPool = sync.Pool{
	New: func() interface{} {
		return &reqOptions{}
	},
}

func acquireRequest() *reqOptions {
	return reqPool.Get().(*reqOptions)
}

func releaseRequest(r *reqOptions) {
	r.Reset()
	reqPool.Put(r)
}

func (req *reqOptions) Reset() {
	req.endpoint = req.endpoint[:0]
	req.postData = req.postData[:0]
	if req.args != nil {
		fasthttp.ReleaseArgs(req.args)
		req.args = nil
	}
	req.skipStatus = false
}

// TODO: Does the same as sendSimpleRequest
func (insta *Instagram) OptionalRequest(endpoint string, a ...interface{}) (body []byte, err error) {
	req := acquireRequest()
	defer releaseRequest(req)
	req.SetEndpoint(fmt.Sprintf(endpoint, a...))
	return insta.sendRequest(req)
}

func (insta *Instagram) sendSimpleRequest(endpoint string, a ...interface{}) (body []byte, err error) {
	req := acquireRequest()
	defer releaseRequest(req)
	req.SetEndpoint(fmt.Sprintf(endpoint, a...))
	return insta.sendRequest(req)
}

func (insta *Instagram) sendRequest(o *reqOptions) (body []byte, err error) {
	url := bytebufferpool.Get()
	url.WriteString(goInstaAPIUrl)
	defer bytebufferpool.Put(url)

	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	// Setting cookies
	if insta.cookies != nil {
		for _, c := range insta.cookies.Cookies() {
			req.Header.SetCookieBytesKV(c.Key(), c.Value())
		}
	}
	url.Write(o.endpoint)

	if len(o.postData) != 0 {
		req.Header.SetMethod("POST")
		req.SetBody(o.postData)
	} else {
		url.WriteByte('?')
		if o.args != nil {
			url.Write(o.args.QueryString())
		}
	}
	req.SetRequestURIBytes(url.B)

	req.Header.Set("Connection", "close")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie2", "$Version=1")
	req.Header.Set("Accept-Language", "en-US")

	err = insta.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	// if cookie changed or should be setted
	cookie := res.Header.Peek("Set-Cookie")
	if cookie != nil {
		if insta.cookies == nil {
			insta.cookies = &cookies{}
		}
		res.Header.VisitAllCookie(func(key, value []byte) {
			insta.cookies.Set(key, value)
		})
		v := insta.cookies.Peek("csrftoken")
		if len(v) != 0 {
			insta.Info.Token = v
		}
	}

	body = res.Body()
	if res.StatusCode() != 200 && !o.skipStatus {
		switch res.StatusCode() {
		case 400:
			err = ErrLoggedOut
		case 404:
			err = ErrNotFound
		}
	}

	return body, err
}
