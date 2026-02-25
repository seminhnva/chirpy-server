package main

import (
	"errors"
	"log"
	"strings"
)

func validateChirp(s string) (string, error) {
	cleaned_body, err := validateProfane(s, []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	})
	if err != nil {
		log.Printf("Error validate chirps: %v", err)
		return "", err
	}
	return cleaned_body, nil
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
