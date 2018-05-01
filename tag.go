package goinsta

type Tag struct {
	insta *Instagram
	ID    string

	Related *RelatedTag
	Feed    *FeedTag
}

// NewTag
func NewTag(insta *Instagram) *Tag {
	tag := Tag{insta: insta}
	tag.Releated = &RelatedTag{}
	tag.Feed = &FeedTag{}
	return tag
}

// GetRelated can get related tags by tags in instagram
//
// GetRelated uses Tag.ID
func (tag *Tag) GetRelated() error {
	insta := tag.insta

	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)
	req.SetEndpoint(fmt.Sprintf("tags/%s/related", tag.ID))
	req.args.Set("visited", fmt.Sprintf(`[{"id":"%s","type":"hashtag"}]`, tag.ID))
	req.args.Set("related_types", `["hashtag"]`)

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}
	if tag.Related == nil {
		tag.Related = &RelatedTag{}
	}

	return json.Unmarshal(body, tag.Releated)
}

// GetFeed search by tags in instagram
func (tag *Tag) GetFeed() error {
	insta := tag.insta

	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)
	req.SetEndpoint(fmt.Sprintf("feed/tag/%s/", tag))
	req.args.Set("rank_token", insta.Info.RankToken)
	req.args.Set("ranked_content", "true")

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}
	if tag.Feed == nil {
		tag.Feed = &FeedTag{}
	}

	return json.Unmarshal(body, tag.Feed)
}
