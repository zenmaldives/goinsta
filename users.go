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
