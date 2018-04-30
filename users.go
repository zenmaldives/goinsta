package goinsta

import (
	"encoding/json"
	"fmt"

	"github.com/erikdubbelboer/fasthttp"
)

// Users represents instagram list of users
type Users struct {
	user  *User
	insta *Instagram

	// followers
	f bool

	StatusResponse
	BigList   bool   `json:"big_list"`
	Users     []User `json:"users"`
	PageSize  int    `json:"page_size"`
	NextMaxID string `json:"next_max_id"`
}

// NewUsers returns a list of users.
//
// If followers it's true when Get is called it will be load followers.
// In other case it will be load following.
func NewUsers(user *User, followers bool) *Users {
	users := &Users{
		user:  user,
		insta: user.insta,
		f:     followers,
	}
	return users
}

// Next does not keep last data.
// To keep last data use AddNext.
func (users *Users) Next() bool {
	// TODO
	return false
}

// Get gets following or follower values into Users structure
//
// User.ID is the id of the target user.
// User.Following/Followers.NextMaxID is the id of the pagination. If is the first request set to 0.
//
// This function does not get all following/followers. To get all following use #Users.All.
func (users *Users) Get() (err error) {
	switch {
	case users.f:
		err = users.followers()
	default:
		err = users.following()
	}
	return err
}

func (users *Users) following() error {
	user := users.user
	userID := user.getID()
	if userID == "" {
		return ErrNoID
	}

	maxID := users.NextMaxID
	insta := user.insta
	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("friendships/%s/following/", userID))
	req.args.Set("max_id", maxID)
	req.args.Set("ig_sig_key_version", goInstaSigKeyVersion)
	req.args.Set("rank_token", insta.Info.RankToken)

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, users)
}

func (users *Users) followers() error {
	user := users.user
	userID := user.getID()
	if userID == "" {
		return ErrNoID
	}

	maxID := users.NextMaxID
	insta := user.insta
	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("friendships/%s/followers/", userID))
	req.args.Set("max_id", maxID)
	req.args.Set("ig_sig_key_version", goInstaSigKeyVersion)
	req.args.Set("rank_token", insta.Info.RankToken)

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, users)
}

// AllFollowing ...
func (users *Users) All() (err error) {
	switch {
	case users.f:
		err = users.followers()
	default:
		err = users.following()
	}
	return err
}

func (users *Users) allFollowing() (err error) {
	userList := users.Users
	for {
		if err = users.following(); err != nil {
			break
		}

		userList = append(userList, users.Users...)

		if users.BigList {
			break
		}
	}
	return
}

func (users *Users) allFollowers() (err error) {
	userList := users.Users
	for {
		if err = users.followers(); err != nil {
			break
		}

		userList = append(userList, users.Users...)

		if users.BigList {
			break
		}
	}
	return
}
