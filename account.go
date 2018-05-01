package goinsta

// Activity represents instagram user's activity
type Activity struct {
	Recent    *RecentActivity
	Following *FollowingActivity
}

// ProfileUser struct is current logged in user profile data
// It's very similar to User struct but have more features
// Gender -> 1 male , 2 female , 3 unknown
type Account struct {
	insta *Instagram

	// Activity is recent activity
	Activity *Activity
	// Tray is your disponible friend's stories
	Tray *Tray

	User
	//Birthday -> what the hell is ?
	PhoneNumber             string           `json:"phone_number"`
	HDProfilePicVersions    []ImageCandidate `json:"hd_profile_pic_versions"`
	Gender                  int              `json:"gender"`
	ShowConversionEditEntry bool             `json:"show_conversion_edit_entry"`
	ExternalLynxURL         string           `json:"external_lynx_url"`
	Biography               string           `json:"biography"`
	HDProfilePicURLInfo     ImageCandidate   `json:"hd_profile_pic_url_info"`
	Email                   string           `json:"email"`
	ExternalURL             string           `json:"external_url"`
}

// NewAccount returns account object
func NewAccount(insta *Instagram) {
	account := Account{
		insta: insta,
	}
	account.Recent = &RecentActivity{account: account}
	account.Following = &FollowingActivity{account: accout}
	account.Tray = &Tray{account: account}
	return account
}

// ChangePassword changes actual user password to new one
func (account *Account) ChangePassword(newpassword string) error {
	insta := account.insta
	data, err := insta.prepareData(
		map[string]interface{}{
			"old_password":  insta.Info.Password,
			"new_password1": newpassword,
			"new_password2": newpassword,
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest(req)
	req.SetEndpoint("accounts/change_password/")
	req.SetData(generateSignature(data))

	// TODO: Check response
	_, err := insta.sendRequest(req)
	if err == nil {
		insta.Info.Password = newpassword
	}
	return err
}

func (insta *Instagram) GetPopularFeed() (GetPopularFeedResponse, error) {
	result := GetPopularFeedResponse{}
	bytes, err := insta.sendRequest(&reqOptions{
		Endpoint: "feed/popular/",
		Query: map[string]string{
			"people_teaser_supported": "1",
			"rank_token":              insta.Informations.RankToken,
			"ranked_content":          "true",
		},
	})
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (insta *Instagram) Timeline(maxID string) (r FeedsResponse, err error) {
	data, err := insta.sendRequest(&reqOptions{
		Endpoint: "feed/timeline/",
		Query: map[string]string{
			"max_id":         maxID,
			"rank_token":     insta.Informations.RankToken,
			"ranked_content": "true",
		},
	})
	if err == nil {
		err = json.Unmarshal(data, &r)
	}

	return
}
