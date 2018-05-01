package goinsta

type Inbox struct {
	insta *Instagram

	DirectList *DirectList
	Pending    *PendingInbox

	Thread *Thread
}

func NewInbox(insta *Instagram) *Inbox {
	inbox := Inbox{insta: insta}
	inbox.DirectList = &DirectList{}
	inbox.Pending = &PendingInbox{}
	inbox.Thread = NewThread(insta)
	return inbox
}

// TODO
func (inbox *Inbox) GetRecent() ([]byte, error) {
	return insta.sendSimpleRequest("direct_share/recent_recipients/")
}

// GetV2 stores inbox messages in Inbox.DirectList
func (inbox *Inbox) GetV2() error {
	body, err := insta.sendSimpleRequest("direct_v2/inbox/?")
	if err != nil {
		return err
	}
	if inbox.DirectList == nil {
		inbox.DirectList = &DirectList{}
	}

	return json.Unmarshal(body, inbox)
}

// PendingRequests stores in Inbox.Pending
func (inbox *Inbox) PendingRequests() error {
	body, err := insta.sendSimpleRequest("direct_v2/pending_inbox/?")
	if err != nil {
		return err
	}
	if inbox.Pending == nil {
		inbox.Pending = &PendingInbox{}
	}

	return json.Unmarshal(body, inbox)
}

// RankedRecipients stores ranked recipients in Inbox.Ranked
func (inbox *Inbox) RankedRecipients() error {
	body, err := insta.sendSimpleRequest("direct_v2/ranked_recipients/?")
	if err != nil {
		return result, err
	}
	if inbox.Ranked == nil {
		inbox.Ranked = &RankedInbox{}
	}

	return json.Unmarshal(body, inbox)
}

type threadsResponse struct {
	StatusResponse
	Threads Threads `json:"threads"`
}

type Threads []Thread

type threadResponse struct {
	StatusResponse
	Thread Thread `json:"thread"`
}

type Thread struct {
	insta *Instagram
	ID    string `json:"-"`

	Named bool `json:"named"`
	Users []struct {
		Username                   string `json:"username"`
		HasAnonymousProfilePicture bool   `json:"has_anonymous_profile_picture"`
		FriendshipStatus           struct {
			Following       bool `json:"following"`
			IncomingRequest bool `json:"incoming_request"`
			OutgoingRequest bool `json:"outgoing_request"`
			Blocking        bool `json:"blocking"`
			IsPrivate       bool `json:"is_private"`
		} `json:"friendship_status"`
		ProfilePicURL string `json:"profile_pic_url"`
		ProfilePicID  string `json:"profile_pic_id"`
		FullName      string `json:"full_name"`
		Pk            int64  `json:"pk"`
		IsVerified    bool   `json:"is_verified"`
		IsPrivate     bool   `json:"is_private"`
	} `json:"users"`
	ViewerID         int64            `json:"viewer_id"`
	MoreAvailableMin bool             `json:"more_available_min"`
	ThreadID         string           `json:"thread_id"`
	ImageVersions2   ImageVersions    `json:"image_versions2"`
	LastActivityAt   int64            `json:"last_activity_at"`
	NextMaxID        string           `json:"next_max_id"`
	Canonical        bool             `json:"canonical"`
	LeftUsers        []interface{}    `json:"left_users"`
	NextMinID        string           `json:"next_min_id"`
	Muted            bool             `json:"muted"`
	Items            []ItemMediaShare `json:"items"`
	ThreadType       string           `json:"thread_type"`
	MoreAvailableMax bool             `json:"more_available_max"`
	ThreadTitle      string           `json:"thread_title"`
	LastSeenAt       struct {
		Num1572292791 struct {
			ItemID    string `json:"item_id"`
			Timestamp string `json:"timestamp"`
		} `json:"1572292791"`
		Num4043092277 struct {
			ItemID    string `json:"item_id"`
			Timestamp string `json:"timestamp"`
		} `json:"4043092277"`
	} `json:"last_seen_at"`
	Inviter struct {
		Username                   string `json:"username"`
		HasAnonymousProfilePicture bool   `json:"has_anonymous_profile_picture"`
		ProfilePicURL              string `json:"profile_pic_url"`
		ProfilePicID               string `json:"profile_pic_id"`
		FullName                   string `json:"full_name"`
		Pk                         int64  `json:"pk"`
		IsVerified                 bool   `json:"is_verified"`
		IsPrivate                  bool   `json:"is_private"`
	} `json:"inviter"`
	Pending bool `json:"pending"`
}

// Get stores requested thread specified by Thread.ID in Thread
func (thread *Thread) Get() error {
	insta := thread.insta
	body, err := insta.sendSimpleRequest("direct_v2/threads/%s/", thread.ID)
	if err != nil {
		return result, err
	}

	r := threadResponse{}
	err = json.Unmarshal(body, &r)
	if err == nil {
		thread = &r.Thread
	}
	return err
}
