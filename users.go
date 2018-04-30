package goinsta

// Users represents instagram list of users
type Users struct {
	insta  *Instagram
	cursor int64

	StatusResponse
	BigList   bool   `json:"big_list"`
	Users     []User `json:"users"`
	PageSize  int    `json:"page_size"`
	NextMaxID string `json:"next_max_id"`
}

func NewUsers(insta *Instagram) *Users {
	users := &Users{
		insta: insta,
	}
	return users
}

// Next does not keep last data.
// To keep last data use AddNext.
func (users *Users) Next() bool {
	// TODO
}
