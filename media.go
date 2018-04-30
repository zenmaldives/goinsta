package goinsta

type Media struct {
	// TODO: Implement
	next  int64
	insta *Instagram

	Comments MediaComments
	Likes    MediaLikes
}

// Media returns comments of a media, input is media ID.
// You can use maxID for pagination, otherwise leave it empty to get the latest page only.
func (com *Comments) Media(mediaID string, maxID string) (resp MediaComments, err error) {
	insta := users.insta
	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("media/%s/comments", mediaID))
	req.args.Set("max_id", maxID)

	body, err := insta.sendRequest(req)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(body, &resp)
	return resp, err
}

// MediaLikers return likers of a media , input is mediaid of a media
func (insta *Instagram) MediaLikers(mediaID string) (MediaLikers, error) {
	body, err := insta.sendSimpleRequest("media/%s/likers/?", mediaID)
	if err != nil {
		return MediaLikersResponse{}, err
	}
	resp := MediaLikersResponse{}
	err = json.Unmarshal(body, &resp)

	return resp, err
}
