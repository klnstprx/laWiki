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
		return
	}

	if len(req.Fields) == 0 {
		http.Error(w, "no fields to translate", http.StatusBadRequest)
		return
	}

	// Collect texts and corresponding field names
	texts := []string{}
	fieldNames := []string{}
	for field, text := range req.Fields {
		texts = append(texts, text)
		fieldNames = append(fieldNames, field)
	}

	// Prepare DeepL API request with source_lang set to 'auto' for detection
	apiURL := "https://api.deepl.com/v2/translate"
	data := url.Values{}
	data.Set("auth_key", config.App.DeepLKey)
	data.Set("target_lang", strings.ToUpper(req.TargetLang))
	data.Set("source_lang", "auto")  // Enable source language detection
	data.Set("tag_handling", "html") // Preserve HTML tags in translation
	for _, text := range texts {
		data.Add("text", text)
	}

	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type deepLResp struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
		DetectedSourceLanguage string `json:"detected_source_language"`
	}
	var result deepLResp
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(result.Translations) != len(texts) {
		http.Error(w, "translation count mismatch", http.StatusInternalServerError)
		return
	}

	// Check if detected source language matches target language
	if strings.EqualFold(result.DetectedSourceLanguage, req.TargetLang) {
		http.Error(w, "source language is the same as target language", http.StatusBadRequest)
		return
	}

	// Map translations back to field names
	translated := make(map[string]string)
	for i, translation := range result.Translations {
		translated[fieldNames[i]] = translation.Text
	}

	// Prepare and send response
	respBody := TranslationResponse{
		TranslatedFields:       translated,
		DetectedSourceLanguage: result.DetectedSourceLanguage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respBody)
}
