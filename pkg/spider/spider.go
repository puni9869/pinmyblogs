package spider

import (
	"codeberg.org/readeck/go-readability/v2"
	"context"
	"encoding/json"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"net/http"
	u "net/url"
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

	baseURL, _ := u.Parse(url.WebLink)
	article, err := readability.FromReader(resp.Body, baseURL)
	if err != nil {
		log.WithFields(map[string]any{
			"parseWeblink": "failed",
			"webLink":      url.WebLink,
		}).Info("Failed to fetch the url")
	}
	url.Title = article.Title()
	url.Summary = article.Excerpt()

	db := database.Db()
	db.Model(&models.Url{}).Where("id = ?", url.ID).Updates(map[string]interface{}{
		"title":   url.Title,
		"comment": url.Comment,
		"summary": url.Summary,
	})

	log.Info("Updating weblink metadata after fetching.", url.WebLink)
}

func FetchAndUpdateURL(url *models.Url) {
	log := logger.NewLogger()

	// Context with timeout (stronger than client timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 9 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.WebLink, nil)
	if err != nil {
		log.WithFields(map[string]any{
			"webLink": url.WebLink,
			"error":   err.Error(),
		}).Error("Failed to create HTTP request")
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(map[string]any{
			"webLink": url.WebLink,
			"error":   err.Error(),
		}).Error("Failed to fetch the url")
		return
	}
	defer func() { _ = resp.Body.Close() }()

	// Store status code
	comment := map[string]string{
		"statusCode": strconv.Itoa(resp.StatusCode),
	}

	jsonStr, err := json.Marshal(comment)
	if err != nil {
		log.WithError(err).Warn("Failed to marshal comment JSON")
	}
	url.Comment = string(jsonStr)

	baseURL, err := u.Parse(url.WebLink)
	if err != nil {
		log.WithFields(map[string]any{
			"webLink": url.WebLink,
			"error":   err.Error(),
		}).Warn("Failed to parse base URL")
		return
	}

	article, err := readability.FromReader(resp.Body, baseURL)
	if err != nil {
		log.WithFields(map[string]any{
			"webLink": url.WebLink,
			"error":   err.Error(),
		}).Warn("Failed to parse article content")
		return
	}

	url.Title = article.Title()
	url.Summary = article.Excerpt()

	db := database.Db()
	err = db.Model(&models.Url{}).
		Where("id = ?", url.ID).
		Updates(map[string]interface{}{
			"title":   url.Title,
			"comment": url.Comment,
			"summary": url.Summary,
		}).Error

	if err != nil {
		log.WithFields(map[string]any{
			"urlID": url.ID,
			"error": err.Error(),
		}).Error("Failed to update URL metadata in DB")
		return
	}

	log.WithField("webLink", url.WebLink).
		Info("Successfully updated weblink metadata")
}
