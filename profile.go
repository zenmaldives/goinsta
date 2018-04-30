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
