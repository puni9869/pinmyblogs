package spider

import (
	"encoding/json"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"golang.org/x/net/html"
	"net/http"
	"strconv"
	"time"
)

func ScrapeUrl(url *models.Url) {
	log := logger.NewLogger()
	var err error
	c := http.Client{
		Timeout: 5 * time.Second, // 5 seconds
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 9 {
				return http.ErrUseLastResponse
			}
			return nil // Follow the redirect
		},
	}

	resp, err := c.Get(url.WebLink)
	if err != nil {
		log.WithFields(map[string]any{
			"fetchUrl": "failed",
			"webLink":  url.WebLink,
		}).Info("Failed to fetch the url")
		return
	}
	defer resp.Body.Close()
	_, err = html.Parse(resp.Body)
	if err != nil {
		log.WithFields(map[string]any{
			"parseWeblink": "failed",
			"webLink":      url.WebLink,
		}).Info("Failed to fetch the url")
		return
	}
	comment := make(map[string]string)
	comment["statusCode"] = strconv.Itoa(resp.StatusCode)
	jsonStr, err := json.Marshal(comment)

	url.Comment = string(jsonStr)
	db := database.Db()
	db.Updates(url)
	log.Info("Updating weblink metadata after fetching.", url.WebLink)
}
