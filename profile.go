package goinsta

// ProfileUser struct is current logged in user profile data
// It's very similar to User struct but have more features
// Gender -> 1 male , 2 female , 3 unknown
type ProfileData struct {
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

func (insta *Instagram) ChangePassword(newpassword string) ([]byte, error) {
	data, err := insta.prepareData(map[string]interface{}{
		"old_password":  insta.Informations.Password,
		"new_password1": newpassword,
		"new_password2": newpassword,
	})
	if err != nil {
		return []byte{}, err
	}
	bytes, err := insta.sendRequest(&reqOptions{
		Endpoint: "accounts/change_password/",
		PostData: generateSignature(data),
	})
	if err == nil {
		insta.Informations.Password = newpassword
	}
	return bytes, err
}

func (insta *Instagram) GetRecentActivity() (RecentActivityResponse, error) {
	result := RecentActivityResponse{}
	bytes, err := insta.sendSimpleRequest("news/inbox/?")
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (insta *Instagram) GetFollowingRecentActivity() (FollowingRecentActivityResponse, error) {
	result := FollowingRecentActivityResponse{}
	bytes, err := insta.sendSimpleRequest("news/?")
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetTrayFeeds - Get all available Instagram stories of your friends
func (insta *Instagram) GetReelsTrayFeed() (TrayResponse, error) {
	bytes, err := insta.sendSimpleRequest("feed/reels_tray/")
	if err != nil {
		return TrayResponse{}, err
	}

	result := TrayResponse{}
	json.Unmarshal([]byte(bytes), &result)

	return result, nil
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
