package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/laWiki/translation/config"
	"github.com/laWiki/translation/database"
	"github.com/laWiki/translation/model"
	"github.com/laWiki/translation/utils"
)

// Health check

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Translate solicita la traducción de un texto dado un idioma objetivo
func Translate(w http.ResponseWriter, r *http.Request) {

	var TranslationService = utils.NewTranslationService(config.App)

	var request struct {
		SourceID   string `json:"sourceId"`
		SourceLang string `json:"sourceLang"`
		TargetLang string `json:"targetLang"`
		SourceText string `json:"sourceText"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if request.SourceID == "" || request.SourceLang == "" || request.TargetLang == "" || request.SourceText == "" {
		http.Error(w, "Faltan campos requeridos", http.StatusBadRequest)
		return
	}

	// Traducir el texto utilizando el servicio de traducción
	translatedText, err := TranslationService.TranslateText(request.SourceText, request.TargetLang)
	if err != nil {
		http.Error(w, "Error al traducir el texto", http.StatusInternalServerError)
		return
	}

	// Crear la traducción
	newTranslation := model.Translation{
		SourceID:       request.SourceID,
		SourceLang:     request.SourceLang,
		TargetLang:     request.TargetLang,
		SourceText:     request.SourceText,
		TranslatedText: translatedText,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Insertar la traducción en la base de datos
	_, err = database.TranslationCollection.InsertOne(r.Context(), newTranslation)
	if err != nil {
		http.Error(w, "Error al guardar la traducción", http.StatusInternalServerError)
		return
	}

	// Devolver la traducción creada
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTranslation)
}
