package goinsta

// User is instagram user informations
//
// This datatype is used in requests
type User struct {
	insta *Instagram

	// User objects
	Feed *UserFeed `json:"-"`

	// Json objects and user data
	Username                   string `json:"username"`
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
}

func NewUser(insta *Instagram) User {
	user := &User{
		insta: insta,
	}
	user.Feed = NewFeed(user)
	return user
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

// UserTaggedFeed - Returns the feed for medua a given user is tagged in
func (user *User) UserTaggedFeed(userID, maxID int64, minTimestamp string) (UserTaggedFeedResponse, error) {
	resp := UserTaggedFeedResponse{}
	maxid := ""
	if maxID != 0 {
		maxid = string(maxID)
	}

	body, err := insta.sendRequest(&reqOptions{
		Endpoint: fmt.Sprintf("usertags/%d/feed/", userID),
		Query: map[string]string{
			"max_id":         maxid,
			"rank_token":     insta.Informations.RankToken,
			"min_timestamp":  minTimestamp,
			"ranked_content": "true",
		},
	})
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(body, &resp)

	return resp, err
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
