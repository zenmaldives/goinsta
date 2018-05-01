package goinsta

import (
	"encoding/json"
	"fmt"
)

// FriendShip represents friendship request
type FriendShip struct {
	Following       bool   `json:"following"`
	FollowedBy      bool   `json:"followed_by"`
	Status          string `json:"status"`
	IsPrivate       bool   `json:"is_private"`
	IsMutingReel    bool   `json:"is_muting_reel"`
	OutgoingRequest bool   `json:"outgoing_request"`
	IsBlockingReel  bool   `json:"is_blocking_reel"`
	Blocking        bool   `json:"blocking"`
	IncomingRequest bool   `json:"incoming_request"`
}

// User is instagram user informations
type User struct {
	insta *Instagram

	StatusResponse

	// User objects
	Feed       *UserFeed   `json:"-"`
	Followers  *Users      `json:"-"`
	Following  *Users      `json:"-"`
	FriendShip *FriendShip `json:"friendship_status"`
	Story      *Story      `json:"-"`
	Threads    *Threads    `json:"-"`

	// Json objects and user data
	Username                   string `json:"username"`
	Biography                  string `json:"biography"`
	HasAnonymousProfilePicture bool   `json:"has_anonymouse_profile_picture"`
	ProfilePictureID           string `json:"profile_pic_id"`
	ProfilePictureURL          string `json:"profile_pic_url"`
	FullName                   string `json:"full_name"`
	ID                         int64  `json:"pk"`
	IDStr                      string `json:"-"`
	IsVerified                 bool   `json:"is_verified"`
	IsPrivate                  bool   `json:"is_private"`
	IsFavorite                 bool   `json:"is_favorite"`
	IsUnpublished              bool   `json:"is_unpublished"`
	IsBusiness                 bool   `json:"is_business"`
	ExternalLynxURL            string `json:"external_lynx_url"`
	MediaCount                 int    `json:"media_count"`
	AutoExpandChaining         bool   `json:"auto_expand_chaining"`
	FollowingCount             int    `json:"following_count"`
	FollowerCount              int    `json:"follower_count"`
	ExternalURL                string `json:"external_url"`
	HdProfilePicVersions       []struct {
		Height int    `json:"height"`
		Width  int    `json:"width"`
		URL    string `json:"url"`
	} `json:"hd_profile_pic_versions"`
	UserTagsCount       int `json:"usertags_count"`
	HdProfilePicURLInfo struct {
		Height int    `json:"height"`
		Width  int    `json:"width"`
		URL    string `json:"url"`
	} `json:"hd_profile_pic_url_info"`
	GeoMediaCount int  `json:"geo_media_count"`
	HasChaining   bool `json:"has_chaining"`
}

// NewUser returns new username structure
func NewUser(insta *Instagram) *User {
	user := &User{
		insta: insta,
	}
	user.Feed = NewFeed(user)
	user.Following = NewUsers(user, false)
	user.Followers = NewUsers(user, true)
	return user
}

// Get gets user information using User.ID.
func (user *User) Get() error {
	data, err := user.insta.prepareData(make(map[string]interface{}))
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest(req)
	req.SetEndpoint(fmt.Sprintf("users/%s/info/", user.ID))
	req.SetData(generateSignature(data))

	body, err := user.insta.sendRequest(req)
	if err == nil {
		err = json.Unmarshal(body, user)
	}

	return err
}

// GetByName gets user data using User.Username
func (user *User) GetByName() error {
	body, err := user.insta.sendSimpleRequest(
		fmt.Sprintf("users/%s/usernameinfo/", user.Username),
	)
	if err == nil {
		err = json.Unmarshal(body, user)
	}

	return err
}

func (user *User) getID() string {
	var userID string

	if user.ID != 0 {
		userID = fmt.Sprintf("%d", user.ID)
	} else {
		if user.IDStr == "" {
			return ""
		}
		userID = user.IDStr
	}
	return userID
}

// Follow follows specified user using his/her id.
//
// If this function does not return any error User.FriendShip will be replaced
// with the new relationship with this user.
func (user *User) Follow() error {
	userID := user.getID()
	if userID == "" {
		return ErrNoID
	}

	resp := followResponse{}
	data, err := insta.prepareData(
		map[string]interface{}{
			"user_id": userID,
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest(req)
	req.SetEndpoint(fmt.Sprintf("friendships/create/%s/", userID))
	req.SetData(generateSignature(data))

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}
	if user.FriendShip == nil {
		user.FriendShip = &user.FriendShip{}
	}

	return json.Unmarshal(body, user)
}

// Unfollow unfollows specified user
//
// If this function does not return any error User.FriendShip will be replaced
// with the new relationship with this user.
func (user *User) Unfollow() error {
	userID := user.getID()
	if userID == "" {
		return ErrNoID
	}

	data, err := insta.prepareData(
		map[string]interface{}{
			"user_id": userID,
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest(req)
	req.SetEndpoint(fmt.Sprintf("friendships/destroy/%s/", userID))
	req.SetData(generateSignature(data))

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}
	if user.FriendShip == nil {
		user.FriendShip = &user.FriendShip{}
	}

	return json.Unmarshal(body, user)
}

// Block blocks user using user ID
func (insta *Instagram) Block() error {
	userID := user.getID()
	if userID == "" {
		return ErrNoID
	}

	data, err := insta.prepareData(
		map[string]interface{}{
			"user_id": userID,
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest(req)
	req.SetEndpoint(fmt.Sprintf("friendships/block/%s/", userID))
	req.SetData(generateSignature(data))

	_, err = insta.sendRequest(req)
	// TODO: What this function returns?
	return err
}

// Unblock unblocks instagram user
func (insta *Instagram) Unblock() error {
	userID := user.getID()
	if userID == "" {
		return ErrNoID
	}

	data, err := insta.prepareData(
		map[string]interface{}{
			"user_id": userID,
		},
	)
	if err != nil {
		return []byte{}, err
	}

	req := acquireRequest(req)
	defer releaseRequest(req)
	req.SetEndpoint(fmt.Sprintf("friendships/unblock/%s/", userID))
	req.SetData(generateSignature(data))

	_, err = insta.sendRequest(req)
	// TODO: What this functions returns?
	return err
}

// Stories gets all available Instagram stories for the given user id
//
// If the is no error in return statement the stories will be saved in User.Story
func (user *User) Stories() error {
	userID := user.getID()
	if userID == "" {
		return ErrNoID
	}

	// TODO: Probably we have to adapt this request data to user object...
	body, err := insta.sendSimpleRequest("feed/user/%s/reel_media/", userID)
	if err == nil {
		if user.Story == nil {
			user.Story = &Story{}
		}
		err = json.Unmarshal(body, user.Story)
	}

	return err
}

// FriendShip
func (user *User) FriendShip() error {
	userID := user.getID()
	if userID == "" {
		return ErrNoID
	}

	data, err := insta.prepareData(
		map[string]interface{}{
			"user_id": userID,
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest(req)
	req.SetEndpoint(fmt.Sprintf("friendships/show/%s/", userID))
	req.SetData(generateSignature(data))

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}
	if user.FriendShip == nil {
		user.FriendShip = &user.FriendShip{}
	}
	return json.Unmarshal(bytes, user)
}

// Send sends direct message to user.
//
// Send uses user's ID. Result is stored in User.Threads
func (user *User) Send(message string) error {
	recipients, err := json.Marshal([][]string{{recipient}})
	if err != nil {
		return err
	}

	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	w := multipart.NewWriter(b)
	w.SetBoundary(insta.Informations.UUID)
	w.WriteField("recipient_users", string(recipients))
	w.WriteField("client_context", insta.Informations.UUID)
	w.WriteField("thread_ids", `["0"]`)
	w.WriteField("text", message)
	w.Close()

	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.Header.SetMultipartFormBoundaryBytes(b.Bytes())
	req.Header.SetMethod("POST")
	req.SetRequestURI(goInstaAPIUrl + "direct_v2/threads/broadcast/text/")

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en_US")
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Connection", "keep-alive")

	if insta.cookies != nil {
		for _, c := range insta.cookies.Cookies() {
			req.Header.SetCookieBytesKV(c.Key(), c.Value())
		}
	}

	err = client.Do(req, res)
	if err != nil {
		return err
	}

	if res.StatusCode() != 200 {
		return fmt.Errorf("HTTP returned status: %d", res.StatusCode())
	}

	r := threadsResponse{}
	err = json.Unmarshal(body, &r)
	if err == nil {
		user.Threads = &r
	}
	return err
}
