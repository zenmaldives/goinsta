package goinsta

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Comments allows user to interact with media (item) comments.
// You can Add or Delete by index or by user name.
type Comments struct {
	item     *Item
	endpoint string
	err      error

	Items                          []Comment `json:"comments"`
	CommentCount                   int       `json:"comment_count"`
	Caption                        Caption   `json:"caption"`
	CaptionIsEdited                bool      `json:"caption_is_edited"`
	HasMoreComments                bool      `json:"has_more_comments"`
	HasMoreHeadloadComments        bool      `json:"has_more_headload_comments"`
	MediaHeaderDisplay             string    `json:"media_header_display"`
	DisplayRealtimeTypingIndicator bool      `json:"display_realtime_typing_indicator"`
	NextID                         string    `json:"next_max_id"`
	LastID                         string    `json:"next_min_id"`
	Status                         string    `json:"status"`

	//PreviewComments                []Comment   `json:"preview_comments"`
}

func (c *Comments) setValues() {
	for i := range c.Items {
		c.Items[i].inst = c.item.media.instagram()
	}
}

func newComments(item *Item) *Comments {
	c := &Comments{
		item: item,
	}
	return c
}

func (comments Comments) Error() error {
	return comments.err
}

// Disable disables comments in FeedMedia.
//
// See example: examples/media/commentDisable.go
func (comments *Comments) Disable() error {
	switch comments.item.media.(type) {
	case *StoryMedia:
		return fmt.Errorf("Incompatible type. Cannot use Disable() with StoryMedia type")
	default:
	}

	insta := comments.item.media.instagram()
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": comments.item.ID,
		},
	)
	if err != nil {
		return err
	}

	_, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCommentDisable, comments.item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// Enable enables comments in FeedMedia
//
// See example: examples/media/commentEnable.go
func (comments *Comments) Enable() error {
	switch comments.item.media.(type) {
	case *StoryMedia:
		return fmt.Errorf("Incompatible type. Cannot use Enable() with StoryMedia type")
	default:
	}

	insta := comments.item.media.instagram()
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": comments.item.ID,
		},
	)
	if err != nil {
		return err
	}

	_, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCommentEnable, comments.item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// Next allows comment pagination.
//
// This function support concurrency methods to get comments using Last and Next ID
//
// New comments are stored in Comments.Items
func (comments *Comments) Next() bool {
	if comments.err != nil {
		return false
	}

	item := comments.item
	insta := comments.item.media.instagram()
	data, err := insta.prepareData(
		map[string]interface{}{
			"can_support_threading": true,
			"max_id":                comments.NextID,
			"min_id":                comments.LastID,
		},
	)
	if err != nil {
		comments.err = err
		return false
	}

	endpoint := comments.endpoint

	body, err := insta.sendRequest(
		&reqOptions{
			Endpoint: endpoint,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err == nil {
		c := Comments{}
		err = json.Unmarshal(body, &c)
		if err == nil {
			*comments = c
			comments.endpoint = endpoint
			comments.item = item
			if !comments.HasMoreComments || comments.NextID == "" {
				comments.err = ErrNoMore
			}
			comments.setValues()
			return true
		}
	}
	comments.err = err
	return false
}

// Sync prepare Comments to receive comments.
// Use Next to receive comments.
//
// See example: examples/media/commentsSync.go
func (comments *Comments) Sync() {
	endpoint := fmt.Sprintf(urlCommentSync, comments.item.ID)
	comments.endpoint = endpoint
	return
}

// Add push a comment in media.
//
// If parent media is a Story this function will send a private message
// replying the Instagram story.
//
// See example: examples/media/commentsAdd.go
func (comments *Comments) Add(text string) (err error) {
	var url, data string
	item := comments.item
	insta := item.media.instagram()

	var query map[string]string

	switch item.media.(type) {
	case *StoryMedia: // TODO: story causes error
		url = urlReplyStory
		query = insta.prepareDataQuery(
			map[string]interface{}{
				"recipient_users": fmt.Sprintf("[[%d]]", item.User.ID),
				"action":          "send_item",
				"media_id":        item.ID,
				"client_context":  generateUUID(),
				"text":            text,
				"entry":           "reel",
				"reel_id":         item.User.ID,
			},
		)
	case *FeedMedia: // normal media
		url = fmt.Sprintf(urlCommentAdd, item.Pk)
		data, err = insta.prepareData(
			map[string]interface{}{
				"comment_text": text,
			},
		)
		query = generateSignature(data)
	}
	if err != nil {
		return err
	}

	// ignoring response
	_, err = insta.sendRequest(
		&reqOptions{
			Endpoint: url,
			Query:    query,
			IsPost:   true,
		},
	)
	return err
}

// Del deletes comment.
func (comments *Comments) Del(comment *Comment) error {
	insta := comments.item.media.instagram()

	data, err := insta.prepareData()
	if err != nil {
		return err
	}
	id := comment.getid()

	_, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCommentDelete, comments.item.ID, id),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// DelByID removes comment using id.
//
// See example: examples/media/commentsDelByID.go
func (comments *Comments) DelByID(id string) error {
	return comments.Del(&Comment{idstr: id})
}

// DelMine removes all of your comments limited by parsed parameter.
//
// If limit is <= 0 DelMine will delete all your comments.
//
// See example: examples/media/commentsDelMine.go
func (comments *Comments) DelMine(limit int) error {
	i := 0
	if limit <= 0 {
		i = limit - 1
	}
	comments.Sync()

	insta := comments.item.media.instagram()
floop:
	for comments.Next() {
		for _, c := range comments.Items {
			if c.UserID == insta.Account.ID || c.User.ID == insta.Account.ID {
				if i >= limit {
					break floop
				}
				comments.Del(&c)
				i++
			}
		}
	}
	if err := comments.Error(); err != nil && err != ErrNoMore {
		return err
	}
	return nil
}

type Comment struct {
	inst *Instagram

	ID                             int64     `json:"pk"`
	idstr                          string    `json:"-"`
	Text                           string    `json:"text"`
	Type                           int       `json:"type"`
	User                           User      `json:"user"`
	UserID                         int64     `json:"user_id"`
	BitFlags                       int       `json:"bit_flags"`
	ChildCommentCount              int       `json:"child_comment_count"`
	CommentIndex                   int       `json:"comment_index"`
	CommentLikeCount               int       `json:"comment_like_count"`
	ContentType                    string    `json:"content_type"`
	CreatedAt                      int       `json:"created_at"`
	CreatedAtUtc                   int       `json:"created_at_utc"`
	DidReportAsSpam                bool      `json:"did_report_as_spam"`
	HasLikedComment                bool      `json:"has_liked_comment"`
	InlineComposerDisplayCondition string    `json:"inline_composer_display_condition"`
	OtherPreviewUsers              []User    `json:"other_preview_users"`
	PreviewChildComments           []Comment `json:"preview_child_comments"`
	Status                         string    `json:"status"`
}

func (c Comment) getid() string {
	switch {
	case c.ID == 0:
		return c.idstr
	case c.idstr == "":
		return strconv.FormatInt(c.ID, 10)
	}
	return ""
}

// Like likes comment.
func (c *Comment) Like() error {
	data, err := c.inst.prepareData()
	if err != nil {
		return err
	}

	_, err = c.inst.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCommentLike, c.getid()),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// Unlike unlikes comment.
func (c *Comment) Unlike() error {
	data, err := c.inst.prepareData()
	if err != nil {
		return err
	}

	_, err = c.inst.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlCommentUnlike, c.getid()),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}
