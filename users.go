package goinsta

// Users represents instagram users
type Users struct {
	insta  *Instagram
	cursor int64

	Current *User
}

func (user *Users) Next() bool {
	// TODO
}

// UserFollowing return followings of specific user.
//
// UserID is the id of the target user.
// maxID is the id of the pagination. If is the first request set to 0.
//
func (users *Users) Following(userID int64, maxID string) (r UsersResponse, err error) {
	insta := users.insta
	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("friendships/%d/following/", userID))
	req.args.Set("max_id", maxID)
	req.args.Set("ig_sig_key_version", goInstaSigKeyVersion)
	req.args.Set("rank_token", insta.Info.RankToken)

	body, err := insta.sendRequest(req)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(body, &r)
	return r, err
}

// UserFollowers return followers of specific user
// skip maxid with empty string for get first page
func (users *Users) Followers(userID int64, maxID string) (r UsersResponse, err error) {
	insta := users.insta
	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("friendships/%d/followers/", userID))
	req.args.Set("max_id", maxID)
	req.args.Set("ig_sig_key_version", goInstaSigKeyVersion)
	req.args.Set("rank_token", insta.Info.RankToken)

	body, err := insta.sendRequest(req)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(body, &r)
	return r, err
}
