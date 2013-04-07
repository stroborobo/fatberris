package fatberris

import (
	"github.com/Knorkebrot/m3u"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"fmt"
	"net/url"
	"strings"
)

type Song struct {
	Artist string
	Title string
	File string
}

type Moods struct {
	Sundaychillsession []Song
	Downtempo []Song
	Uptempo []Song
}

const (
	URL_PREFIX string = "http://fatberris.com/music/"	// no ssl, sorry
	FEEDS string = "list.json"
)

func GetMoods(moodArgs []string) (string, error) {
	if len(moodArgs) == 0 {
		return "", errors.New("No moods requested")
	}

	var mix bool
	for _, m := range moodArgs {
		if m != "chill" && m != "up" &&
		   m != "down" && m != "mix" {
			return "", fmt.Errorf("Unknown mood: %s\n", m)
		}
		if m == "mix" {
			mix = true
		}
	}
	if mix {
		moodArgs = []string{"chill", "up", "down"}
	}

	resp, err := http.Get(URL_PREFIX + FEEDS)
	if err != nil {
		return "", err
	}
	jsonData, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	var moods Moods
	err = json.Unmarshal(jsonData, &moods)
	if err != nil {
		// In Go version <1.1 null values in JSON
		// were only allowed for types that are able
		// to be represended as null in Go.
		// This is fixed in 1.1:
		// https://codereview.appspot.com/6759043/
		return "", err
	}

	length := 0
	for _, m := range moodArgs {
		if m == "chill" {
			length += len(moods.Sundaychillsession)
		}
		if m == "up" {
			length += len(moods.Downtempo)
		}
		if m == "down" {
			length += len(moods.Uptempo)
		}
	}

	pls := make(m3u.Playlist, 0, length)
	for _, m := range moodArgs {
		var songs []Song
		if m == "chill" {
			songs = moods.Sundaychillsession
		}
		if m == "up" {
			songs = moods.Uptempo
		}
		if m == "down" {
			songs = moods.Downtempo
		}
		for i, _ := range songs {
			var song m3u.Song
			title := make([]string, 0, 3)
			if songs[i].Artist != "" {
				title = append(title, songs[i].Artist)
			}
			if songs[i].Title != "" {
				if len(title) > 0 {
					title = append(title, "-")
				}
				title = append(title, songs[i].Title)
			}
			song.Title = strings.Join(title, " ")
			url, err := url.Parse(URL_PREFIX + songs[i].File)
			if err != nil {
				continue
			}
			song.Path = url.String()
			pls = append(pls, song)
		}
	}

	out, err := pls.String()
	if err != nil {
		return "", err
	}

	return out, nil
}
