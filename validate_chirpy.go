package main

import (
	"errors"
	"net/http"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	params := parameters{}

	if err := decodeJson(w, r, &params); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return

	}
	cleaned_body, err := validateProfane(params.Body, []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	type validRes struct {
		Cleaned_body string `json:"cleaned_body"`
	}
	respondWithJSON(w, http.StatusOK, validRes{Cleaned_body: cleaned_body})

}

func validateProfane(content string, banWords []string) (string, error) {
	if len(content) > 140 {
		return "", errors.New("Chirp is too long")
	}
	mWord := make(map[string]bool)
	for _, v := range banWords {
		mWord[strings.ToLower(v)] = true
	}

	wContent := strings.Split(content, " ")
	for i, w := range wContent {
		_, ok := mWord[strings.ToLower(w)]
		if ok {
			wContent[i] = strings.Repeat("*", 4)
		}
	}
	return strings.Join(wContent, " "), nil
}
