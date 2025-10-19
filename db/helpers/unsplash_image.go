package unsplash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type photoResp struct {
	Urls struct {
		Raw     string `json:"raw"`
		Regular string `json:"regular"`
		Small   string `json:"small"`
	} `json:"urls"`
	Links struct {
		Html             string `json:"html"`
		DownloadLocation string `json:"download_location"`
	} `json:"links"`
	User struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Links    struct {
			Html string `json:"html"`
		} `json:"links"`
	} `json:"user"`
}

type PhotoMeta struct {
	Small            string
	Regular          string
	Raw              string
	PhotoPageURL     string
	AuthorName       string
	AuthorUsername   string
	AuthorProfileURL string
	DownloadLocation string
}

func FetchPhotoMeta(photoID string) (PhotoMeta, error) {
	key := os.Getenv("UNSPLASH_ACCESS_KEY")
	if key == "" {
		return PhotoMeta{}, fmt.Errorf("UNSPLASH_ACCESS_KEY is empty")
	}

	req, _ := http.NewRequest("GET", "https://api.unsplash.com/photos/"+photoID, nil)
	req.Header.Set("Authorization", "Client-ID "+key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return PhotoMeta{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return PhotoMeta{}, fmt.Errorf("unsplash status %d", res.StatusCode)
	}

	var p photoResp
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return PhotoMeta{}, err
	}

	// (opcjonalnie) tracking pobrania – asynchronicznie
	if p.Links.DownloadLocation != "" {
		go func(url string) {
			dl, _ := http.NewRequest("GET", url, nil)
			dl.Header.Set("Authorization", "Client-ID "+key)
			if r, e := http.DefaultClient.Do(dl); r != nil {
				r.Body.Close()
			} else {
				_ = e
			}
		}(p.Links.DownloadLocation)
	}

	return PhotoMeta{
		Small:            p.Urls.Small,
		Regular:          p.Urls.Regular,
		Raw:              p.Urls.Raw,
		PhotoPageURL:     p.Links.Html,
		AuthorName:       p.User.Name,
		AuthorUsername:   p.User.Username,
		AuthorProfileURL: p.User.Links.Html,
		DownloadLocation: p.Links.DownloadLocation,
	}, nil
}
