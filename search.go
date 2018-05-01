package goinsta

type Search struct {
	insta *Instagram
}

func NewSearch(insta *Instagram) {
	search := &Search{
		insta: insta,
	}
	return search
}

// SearchLocation return search location by lat & lng & search query in instagram
func (insta *Instagram) SearchLocation(lat, lng, search string) (SearchLocationResponse, error) {
	if lat == "" || lng == "" {
		return SearchLocationResponse{}, fmt.Errorf("lat & lng must not be empty")
	}

	query := map[string]string{
		"rank_token":     insta.Informations.RankToken,
		"latitude":       lat,
		"longitude":      lng,
		"ranked_content": "true",
	}

	if search != "" {
		query["search_query"] = search
	} else {
		query["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	}
	body, err := insta.sendRequest(&reqOptions{
		Endpoint: "location_search/",
		Query:    query,
	})

	if err != nil {
		return SearchLocationResponse{}, err
	}

	resp := SearchLocationResponse{}
	err = json.Unmarshal(body, &resp)
	return resp, err
}

// GetLocationFeed return location feed data by locationID in Instagram
func (insta *Instagram) GetLocationFeed(locationID int64, maxID string) (LocationFeedResponse, error) {
	body, err := insta.sendRequest(&reqOptions{
		Endpoint: fmt.Sprintf("feed/location/%d/", locationID),
		Query: map[string]string{
			"max_id": maxID,
		},
	})
	if err != nil {
		return LocationFeedResponse{}, err
	}

	resp := LocationFeedResponse{}
	err = json.Unmarshal(body, &resp)
	return resp, err
}

func (insta *Instagram) SearchUsername(query string) (SearchUserResponse, error) {
	result := SearchUserResponse{}
	body, err := insta.sendRequest(&reqOptions{
		Endpoint: "users/search/",
		Query: map[string]string{
			"ig_sig_key_version": goInstaSigKeyVersion,
			"is_typeahead":       "true",
			"query":              query,
			"rank_token":         insta.Informations.RankToken,
		},
	})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)

	return result, err
}

func (insta *Instagram) SearchTags(query string) (SearchTagsResponse, error) {
	result := SearchTagsResponse{}
	body, err := insta.sendRequest(&reqOptions{
		Endpoint: "tags/search/",
		Query: map[string]string{
			"is_typeahead": "true",
			"rank_token":   insta.Informations.RankToken,
			"q":            query,
		},
	})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)

	return result, err
}

func (insta *Instagram) SearchFacebookUsers(query string) ([]byte, error) {
	return insta.sendRequest(&reqOptions{
		Endpoint: "fbsearch/topsearch/",
		Query: map[string]string{
			"query":      query,
			"rank_token": insta.Informations.RankToken,
		},
	})
}
