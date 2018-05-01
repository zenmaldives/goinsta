package goinsta

import (
	"github.com/erikdubbelboer/fasthttp"
)

// ClientInfo ...
type ClientInfo struct {
	Username string
	// TODO: Is safe to store password in memory?
	Password string
	// TODO: Allow user change this fields?
	DeviceID  string
	UUID      string
	RankToken string
	Token     string
	PhoneID   string
}

type cookies map[string]*fasthttp.Cookie

func (ck *cookies) Set(key, value []byte) {
	ks := b2s(key)
	c, ok := (*ck)[ks]
	if !ok {
		c = fasthttp.AcquireCookie()
	}
	c.SetKeyBytes(key)
	c.SetValueBytes(value)
	(*ck)[ks] = c
}

func (ck *cookies) SetCookies(cks []*fasthttp.Cookie) {
	for _, c := range cks {
		(*ck)[b2s(c.Key())] = c
	}
}

func (ck *cookies) Release() {
	for k, c := range *ck {
		fasthttp.ReleaseCookie(c)
		delete(*ck, k)
	}
}

func (ck *cookies) Peek(v string) string {
	c := (*ck)[v]
	if c != nil {
		return b2s(c.Value())
	}
	return ""
}

func (ck *cookies) Cookies() map[string]*fasthttp.Cookie {
	return *ck
}

// Instagram ....
type Instagram struct {
	Logged bool
	Info   ClientInfo

	// Current is current user (logged user)
	Current ProfileData `json:"user,logged_in_user"`

	// DialFunc allows user to use proxy function.
	// See also: https://godoc.org/github.com/erikdubbelboer/fasthttp#Client.Dial
	DialFunc fasthttp.DialFunc

	client  *fasthttp.Client
	cookies *cookies

	// Instagram objects
	User    *User
	Media   *Media
	Search  *Search
	Explore *Explore
	Inbox   *Inbox
	Tag     *Tag

	StatusResponse
}
