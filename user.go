package goinsta

import (
	"encoding/json"
	"fmt"
)

// User is instagram user informations
type User struct {
	insta *Instagram

	// User objects
	Feed      *UserFeed `json:"-"`
	Followers *Users    `json:"-"`
	Following *Users    `json:"-"`

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

func (user *User) GetByUser() error {
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
