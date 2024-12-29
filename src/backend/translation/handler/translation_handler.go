package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/laWiki/translation/config"
)

// TranslationRequest represents the request payload for translation.
type TranslationRequest struct {
	Fields     map[string]string `json:"fields"`
	TargetLang string            `json:"targetLang"`
}

// TranslationResponse represents the response payload with translated fields and detected source language.
type TranslationResponse struct {
	TranslatedFields       map[string]string `json:"translatedFields"`
	DetectedSourceLanguage string            `json:"detectedSourceLanguage"`
}

// HealthCheck responds with a simple OK status.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Translate translates the specified fields to the target language.
// It detects the source language and prevents translation if source and target languages are the same.
func Translate(w http.ResponseWriter, r *http.Request) {
	var req TranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.TargetLang == "" {
		http.Error(w, "missing targetLang param", http.StatusBadRequest)
		return // Added return
	}

	if len(req.Fields) == 0 {
		http.Error(w, "no fields to translate", http.StatusBadRequest)
		return // Added return
	}

	// Collect texts and corresponding field names
	texts := []string{}
	fieldNames := []string{}
	for field, text := range req.Fields {
		texts = append(texts, text)
		fieldNames = append(fieldNames, field)
	}

	// Prepare DeepL API request with source_lang set to 'auto' for detection
	apiURL := "https://api-free.deepl.com/v2/translate"
	data := url.Values{}
	data.Set("auth_key", config.App.DeepLKey)
	data.Set("target_lang", strings.ToUpper(req.TargetLang))
	data.Set("tag_handling", "html") // Preserve HTML tags in translation

	for _, text := range texts {
		data.Add("text", text)
	}

	// Make the HTTP POST request to DeepL API
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		http.Error(w, "Failed to call DeepL API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check if DeepL API responded with success
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, "DeepL API error: "+string(body), resp.StatusCode)
		return
	}

	// Decode DeepL API response
	var deepLResp struct {
		Translations []struct {
			DetectedSourceLanguage string `json:"detected_source_language"`
			Text                   string `json:"text"`
		} `json:"translations"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&deepLResp); err != nil {
		http.Error(w, "Failed to decode DeepL API response", http.StatusInternalServerError)
		return
	}

	// Prepare the TranslationResponse
	translationResp := TranslationResponse{
		TranslatedFields:       make(map[string]string),
		DetectedSourceLanguage: "",
	}

	for i, translation := range deepLResp.Translations {
		field := fieldNames[i]
		translationResp.TranslatedFields[field] = translation.Text

		// Assuming all fields have the same source language
		if translationResp.DetectedSourceLanguage == "" {
			translationResp.DetectedSourceLanguage = translation.DetectedSourceLanguage
		}
	}

	// Encode and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(translationResp); err != nil {
		http.Error(w, "Failed to encode translation response", http.StatusInternalServerError)
		return
	}
}
