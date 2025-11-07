package spider

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
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
	defer func() { _ = resp.Body.Close() }()
	comment := make(map[string]string)
	comment["statusCode"] = strconv.Itoa(resp.StatusCode)

	jsonStr, _ := json.Marshal(comment)
	url.Comment = string(jsonStr)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.WithFields(map[string]any{
			"parseWeblink": "failed",
			"webLink":      url.WebLink,
		}).Info("Failed to fetch the url")
		return
	}
	url.Title = doc.Find("title").Text()
	db := database.Db()
	db.Model(&models.Url{}).Where("id = ?", url.ID).Updates(map[string]interface{}{
		"title":   url.Title,
		"comment": url.Comment,
	})

	log.Info("Updating weblink metadata after fetching.", url.WebLink)
}
