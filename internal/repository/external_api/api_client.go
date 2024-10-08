package external_api

import (
	"encoding/json"
	"fmt"
	"log"
	"music_catalog/config"
	"net/http"
	"net/url"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type APIClient interface {
	FetchSongDetails(group, song string) (*SongDetail, error)
}

type ExternalAPIClient struct {
	BaseURL string
}

func NewExternalAPIClient(cfg *config.Config) *ExternalAPIClient {
	return &ExternalAPIClient{
		BaseURL: cfg.ExternalAPIURL, // Используем базовый URL из конфигурации
	}
}

func (client *ExternalAPIClient) FetchSongDetails(group, song string) (*SongDetail, error) {
	encodedGroup := url.QueryEscape(group)
	encodedSong := url.QueryEscape(song)
	url := fmt.Sprintf("%s/info?group=%s&song=%s", client.BaseURL, encodedGroup, encodedSong)
	log.Printf("[INFO] Выполняем запрос к внешнему API: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[ERROR] Ошибка при выполнении запроса к внешнему API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Внешний API вернул статус: %d", resp.StatusCode)
		return nil, fmt.Errorf("внешний API вернул ошибку: %d", resp.StatusCode)
	}

	var songDetail SongDetail
	err = json.NewDecoder(resp.Body).Decode(&songDetail)
	if err != nil {
		log.Printf("[ERROR] Ошибка при декодировании ответа внешнего API: %v", err)
		return nil, err
	}

	return &songDetail, nil
}
