package goinsta

// User is instagram user informations
type User struct {
	insta *Instagram

	// User objects
	Feed      *UserFeed `json:"-"`
	Followers *Users    `json:"-"`
	Following *Users    `json:"-"`

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

// NewUser returns new username structure
func NewUser(insta *Instagram) *User {
	user := &User{
		insta: insta,
	}
	user.Feed = NewFeed(user)
	user.Following = NewUsers(user, false)
	user.Followers = NewUser(user, true)
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
