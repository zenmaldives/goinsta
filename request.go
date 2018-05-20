package goinsta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/erikdubbelboer/fasthttp"
	"github.com/spf13/cast"
)

type reqOptions struct {
	// Endpoint is the request path of instagram api
	Endpoint string

	// IsPost setted to true will send request with POST method.
	//
	// By default this option is false.
	IsPost bool

	// Query is the parameters of the request
	//
	// This parameters are independents of the request method (POST|GET)
	Query map[string]string
}

func (insta *Instagram) sendSimpleRequest(uri string, a ...interface{}) (body []byte, err error) {
	return insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(uri, a...),
		},
	)
}

func (inst *Instagram) sendRequest(o *reqOptions) (body []byte, err error) {
	var args *fasthttp.Args
	method := "GET"
	if o.IsPost {
		method = "POST"
	}

	url := bytebufferpool.Get()
	url.WriteString(goInstaAPIUrl)
	url.WriteString(o.Endpoint)

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	req.Header.SetMethod(method)

	for k, v := range o.Query {
		args.Add(k, v)
	}

	if o.IsPost {
		args.WriteTo(req.BodyWriter())
	} else {
		for k, v := range u.Query() {
			args.Add(k, strings.Join(v, " "))
		}

		url.WriteString("?")
		url.Write(args.QueryString())
	}
	req.SetRequestURIBytes(url.B)

	req.Header.Set("Connection", "close")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie2", "$Version=1")
	req.Header.Set("Accept-Language", "en-US")
	// setting cookie
	inst.j.AddToRequest(req)

	err = inst.c.Do(req, res)
	if err != nil {
		return nil, err
	}
	// getting cookies
	inst.j.ResponseCookies(res)

	body := res.Body()

	switch res.StatusCode() {
	case 200:
	default:
		ierr := instaError{}
		err = json.Unmarshal(body, &ierr)
		if err != nil {
			return nil, fmt.Errorf("Invalid status code: %d", resp.StatusCode)
		}
		return nil, instaToErr(ierr)
	}

	return body, err
}

func (insta *Instagram) prepareData(other ...map[string]interface{}) (string, error) {
	data := map[string]interface{}{
		"_uuid":      insta.uuid,
		"_uid":       insta.Account.ID,
		"_csrftoken": insta.token,
	}
	for i := range other {
		for key, value := range other[i] {
			data[key] = value
		}
	}
	b, err := json.Marshal(data)
	if err == nil {
		return b2s(b), err
	}
	return "", err
}

func (insta *Instagram) prepareDataQuery(other ...map[string]interface{}) map[string]string {
	data := map[string]string{
		"_uuid":      insta.uuid,
		"_csrftoken": insta.token,
	}
	for i := range other {
		for key, value := range other[i] {
			data[key] = cast.ToString(value)
		}
	}
	return data
}
