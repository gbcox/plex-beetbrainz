package emby

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"beetbrainz/common"
	"beetbrainz/types"
)

type EmbyItem struct {
	Type        string `json:"Type"`
	Title       string `json:"Name"`
	Parent      string `json:"Album"`
	Grandparent string `json:"Artists"`
}

func (item *EmbyItem) String() string {
	return fmt.Sprintf("%s - %s (%s)", item.Grandparent, item.Title, item.Parent)
}

func (item *EmbyItem) AsMediaItem() *types.MediaItem {
	return &types.MediaItem{
		Artist: item.Grandparent,
		Album:  item.Parent,
		Track:  item.Title,
	}
}

type EmbyAccount struct {
	Title string `json:"Name"`
}

type EmbyRequest struct {
	Event   string      `json:"Event"`
	Account EmbyAccount `json:"User"`
	Item    EmbyItem    `json:"Item"`
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("Request method '%s' is not allowed.", r.Method)
		return
	}

	rq, err := parseRequest(r)
	if err != nil {
		log.Printf("Failed to decode the Emby request: %v", err)
		return
	}

	common.HandleRequest(rq)
}

func parseRequest(r *http.Request) (*common.Request, error) {
	err := r.ParseMultipartForm(16)
	if err != nil {
		return nil, err
	}

	data := []byte(r.FormValue("payload"))
	var embyRequest EmbyRequest
	err = json.Unmarshal(data, &embyRequest)
	if err != nil {
		return nil, err
	}

	rq := &common.Request{
		Event:     toRequestEvent(embyRequest.Event),
		User:      embyRequest.Account.Title,
		MediaType: embyRequest.Item.Type,
		Item:      embyRequest.Item.AsMediaItem(),
	}

	return rq, nil
}

func toRequestEvent(e string) string {
	switch e {
	case "playback.start":
		return "play"
	case "playback.resume":
		return "resume"
	default:
		return "scrobble"
	}
}
