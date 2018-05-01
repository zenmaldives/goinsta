package goinsta

import (
	"encoding/json"
	"fmt"

	"github.com/erikdubbelboer/fasthttp"
)

// MediaInfo contains media information
type MediaInfo struct {
	media *Media

	Status              string `json:"status"`
	NumResults          int    `json:"num_results"`
	AutoLoadMoreEnabled bool   `json:"auto_load_more_enabled"`
	Items               []Item `json:"items"`
	MoreAvailable       bool   `json:"more_available"`
	CommentLikesEnabled bool   `json:"comment_likes_enabled"`
}

// MediaComments struct for get array of comments of a media
type MediaComments struct {
	media *Media

	StatusResponse
	NextMaxID           string            `json:"next_max_id"`
	CommentLikesEnabled bool              `json:"comment_likes_enabled"`
	Comments            []CommentResponse `json:"comments"`
}

// MediaLikers struct for get array of users that like a media
type MediaLikers struct {
	media *Media

	StatusResponse
	UserCount int    `json:"user_count"`
	Users     []User `json:"users"`
}

// Media is media object representation
//
// Do not use concurrently
type Media struct {
	insta *Instagram

	ID string

	Info     *MediaInfo
	Comments *MediaComments
	Likers   *MediaLikers
}

// NewMedia returns new media
func NewMedia(insta *Instagram) *Media {
	media := &Media{insta: insta}
	media.Info = &MediaInfo{media: media}
	media.Comments = &MediaComments{media: media}
	media.Likers = &MediaLikers{media: media}
	return media
}

// Get gets all media data (Likes, Likers, Comments)
//
// id can be "" if Media.ID have been setted before.
func (media *Media) Get(id string) {
	if id != "" {
		media.ID = id
	}
	media.Comments.Get()
	media.Likers.Get()
	media.Info.Get()
}

// GetAsync does the same as Get but using goroutines
func (media *Media) GetAsync(id string) {
	if id != "" {
		media.ID = id
	}
	go media.Comments.Get()
	go media.Likers.Get()
	go media.Info.Get()
}

// Get collect comments.
//
// You can use Get once again to paginate.
func (comments *MediaComments) Get() error {
	media := comments.media
	insta := media.insta
	mediaID := media.ID
	maxID := comments.NextMaxID

	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("media/%s/comments", mediaID))
	req.args.Set("max_id", maxID)

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, comments)
	return err
}

// Get fills likers structure.
//
// This structure cannot be paginated.
func (likers *MediaLikers) Get() error {
	body, err := likers.media.insta.sendSimpleRequest("media/%s/likers/", likers.media.ID)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, likers)
}

// Get fills Info with media info
func (info *MediaInfo) Get() error {
	media := info.media
	insta := media.insta
	mediaID := media.ID

	req := acquireRequest()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("media/%s/info", mediaID))

	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": mediaID,
		},
	)
	if err != nil {
		return err
	}
	req.SetData(generateSignature(data))

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, info)
}

// Like gives like to media
func (media *Media) Like() error {
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": media.ID,
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest(req)
	req.SetEndpoint(fmt.Sprintf("media/%s/like/", media.ID))
	req.SetData(generateSignature(data))

	_, err = insta.sendRequest(req)
	return err
}

// Unlike gives unlike to media
func (media *Media) Unlike() error {
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": mediaID,
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest()
	req.SetEndpoint(fmt.Sprintf("media/%s/unlike/", media.ID))
	req.SetData(generateSignature(data))

	_, err = insta.sendRequest(req)
	return err
}

// DisableComments()
func (media *Media) DisableComments() error {
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": mediaID,
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest()
	req.SetEndpoint(fmt.Sprintf("media/%s/disable_comments/", media.ID))
	req.SetData(generateSignature(data))

	_, err = insta.sendRequest(req)
	return err
}

// EnableComments
func (media *Media) EnableComments() error {
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": mediaID,
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest()
	req.SetEndpoint(fmt.Sprintf("media/%s/enable_comments/", media.ID))
	req.SetData(generateSignature(data))

	_, err = insta.sendRequest(req)
	return err
}

// Edit
func (media *Media) Edit(caption string) error {
	data, err := insta.prepareData(
		map[string]interface{}{
			"caption_text": caption,
		},
	)
	if err != nil {
		return
	}

	req := acquireRequest()
	defer releaseRequest()
	req.SetEndpoint(fmt.Sprintf("media/%s/edit_media/", media.ID))
	req.SetData(generateSignature(data))

	_, err = insta.sendRequest(req)
	return err
}

// Delete
func (media *Media) Delete() error {
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": media.ID,
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest()
	req.SetEndpoint(fmt.Sprintf("media/%s/delete/", media.ID))
	req.SetData(generateSignature(data))

	_, err = insta.sendRequest(req)
	return err
}

// Tag
// TODO:
func (insta *Instagram) RemoveSelfTag(mediaID string) ([]byte, error) {
	// TODO: Probably there are any Tag object in media.
	data, err := insta.prepareData(make(map[string]interface{}))
	if err != nil {
		return
	}

	req := acquireRequest()
	defer releaseRequest()
	req.SetEndpoint(fmt.Sprintf("media/%s/remove/", media.ID))
	req.SetData(generateSignature(data))

	_, err = insta.sendRequest(req)
	return err
}

// Comment
func (media *Media) Comment(text string) error {
	data, err := insta.prepareData(
		map[string]interface{}{
			"comment_text": text,
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest()
	req.SetEndpoint(fmt.Sprintf("media/%s/comment/", media.ID))
	req.SetData(generateSignature(data))

	_, err = insta.sendRequest(req)
	return err
}

// TODO
func (insta *Instagram) DeleteComment(mediaID, commentID string) ([]byte, error) {
	data, err := insta.prepareData()
	if err != nil {
		return []byte{}, err
	}

	return insta.sendRequest(&reqOptions{
		Endpoint: fmt.Sprintf("media/%s/comment/%s/delete/", mediaID, commentID),
		PostData: generateSignature(data),
	})
}
