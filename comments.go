package goinsta

type Comments struct {
	next  int64
	insta *Instagram
}

// Media returns comments of a media, input is media ID.
// You can use maxID for pagination, otherwise leave it empty to get the latest page only.
func (com *Comments) Media(mediaID string, maxID string) (resp MediaCommentsResponse, err error) {
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
