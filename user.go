package goinsta

// User is instagram user informations
//
// This datatype is used in requests
type User struct {
	insta *Instagram

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

// UserFeed - Returns the Instagram feed for the given user id.
// You can use maxID and minTimestamp for pagination, otherwise leave them empty to get the latest page only.
func (user *User) Feed(maxID, minTimestamp string) (resp UserFeedResponse, err error) {
	insta := users.insta
	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("feed/user/%d/", userID))
	req.args.Set("max_id", maxID)
	req.args.Set("rank_token", insta.Info.RankToken)
	req.args.Set("min_timestamp", minTimestamp)
	req.args.Set("ranked_content", "true")

	body, err := insta.sendRequest(req)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(body, &resp)
	if err == nil {
		resp.IDStr = strconv.FormatInt(resp.ID, 10)
	}
	return
}

// LatestFeed gets the latest page of your own Instagram feed.
func (user *User) LatestFeed() (UserFeedResponse, error) {
	return insta.Current.UserFeed("", "")
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
