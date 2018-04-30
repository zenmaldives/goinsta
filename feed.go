package goinsta

// UserFeed contains user feeds
type UserFeed struct {
	user *User

	Status              string `json:"status"`
	NumResults          int    `json:"num_results"`
	AutoLoadMoreEnabled bool   `json:"auto_load_more_enabled"`
	Items               []Item `json:"items"`
	MoreAvailable       bool   `json:"more_available"`
	TotalCount          int    `json:"total_count"`
	RequiresView        bool   `json:"requires_review"`
	NextMaxID           string `json:"next_max_id"`
	MinTimestamp        string `json:"-"`
	// TODO maybe this is photos waiting for review?
	// NewPhotos           []interface{} `json:"new_photos"`
}

// NewFeed returns feed for the given user
func NewFeed(user *User) *UserFeed {
	return &UserFeed{user: user}
}

// SetUser sets new user
func (uf *UserFeed) SetUser(user *User) {
	if user == nil {
		return
	}
	uf.user = user
}

// Reset sets to defaults values UserFeed
func (uf *UserFeed) Reset() {
	uf.Status = ""
	uf.NumResults = 0
	uf.AutoLoadMoreEnabled = false
	uf.Items = nil
	uf.MoreAvailable = false
	uf.NextMaxID = ""
	uf.MinTimestamp = ""
}

// Get sets the Instagram feed for the given user id.
//
// NextMaxID and MinTimestamp can be used for pagination.
// Pagination occurs automatically call by call.
//
// ID or IDStr can be used to interfact with specified user.
func (feed *UserFeed) Get() (err error) {
	userID := user.getID()
	if userID == "" {
		return ErrNoID
	}

	insta := user.insta
	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("feed/user/%s/", userID))
	req.args.Set("max_id", feed.NextMaxID)
	req.args.Set("rank_token", insta.Info.RankToken)
	req.args.Set("min_timestamp", feed.MinTimestamp)
	req.args.Set("ranked_content", "true")

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, feed)
	if err == nil {
		feed.IDStr = strconv.FormatInt(feed.ID, 10)
	}
	return
}

// Latest gets the latest page of users feed.
func (feed *UserFeed) Latest() error {
	feed.Reset()
	return feed.Get()
}

// Tagged sets tagged media in feed structure
func (feed *UserFeed) Tagged() error {
	userID := user.getID()
	if userID == "" {
		return ErrNoID
	}

	insta := feed.user.insta
	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("feed/user/%s/", userID))
	req.args.Set("max_id", feed.NextMaxID)
	req.args.Set("rank_token", insta.Info.RankToken)
	req.args.Set("min_timestamp", feed.MinTimestamp)
	req.args.Set("ranked_content", "true")

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, feed)
	return err
}
