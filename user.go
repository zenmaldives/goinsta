package goinsta

// User is instagram user informations
//
// This datatype is used in requests
type User struct {
	insta *Instagram

	// User objects
	Feed      *UserFeed `json:"-"`
	Followers *Users
	Following *Users

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

// Following sets Followers field in user structure
//
// User.ID is the id of the target user.
// User.Following.NextMaxID is the id of the pagination. If is the first request set to 0.
//
// This function does not get all following. To get all following use #User.AllFollowing.
func (user *User) Following() error {
	userID := getID()
	if userID == "" {
		return ErrNoID
	}
	if user.Following == nil {
		user.Following = NewUsers(insta)
	}

	maxID := user.Following.MaxID
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

	return json.Unmarshal(body, user.Following)
}

// Followers fills User.Followers field.
//
// User.ID is the id of the target user.
// User.Following.NextMaxID is the id of the pagination. If is the first request set to 0.
//
// This function does not get all followers. To get all followers use #User.AllFollowers
func (user *User) Followers() error {
	userID := getID()
	if userID == "" {
		return ErrNoID
	}
	if user.Followers == nil {
		user.Followers = NewUsers(insta)
	}

	maxID := user.Followers.MaxID
	insta := user.insta
	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("friendships/%d/followers/", userID))
	req.args.Set("max_id", maxID)
	req.args.Set("ig_sig_key_version", goInstaSigKeyVersion)
	req.args.Set("rank_token", insta.Info.RankToken)

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, user.Followers)
}
