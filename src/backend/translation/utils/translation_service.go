package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/laWiki/translation/config"
)

// TranslationService maneja las traducciones utilizando la API de DeepL
type TranslationService struct {
	APIKey string
}

// NewTranslationService crea una nueva instancia de TranslationService
func NewTranslationService(cfg config.GlobalConfig) *TranslationService {
	return &TranslationService{
		APIKey: cfg.DeepLAPIKey,
	}
}

// TranslateText traduce el texto al idioma objetivo utilizando DeepL
func (ts *TranslationService) TranslateText(text, targetLang string) (string, error) {
	apiURL := "https://api.deepl.com/v2/translate"
	data := url.Values{}
	data.Set("auth_key", ts.APIKey)
	data.Set("text", text)
	data.Set("target_lang", strings.ToUpper(targetLang))

	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the response
	type deeplResponse struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}

	var result deeplResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if len(result.Translations) > 0 {
		return result.Translations[0].Text, nil
	}

	return "", fmt.Errorf("translation failed")
}
