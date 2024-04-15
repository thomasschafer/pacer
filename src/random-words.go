package main

import (
	"math/rand"
	"strings"
)

var charSet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@£$%^&*()`-='\\[~_+\"|{},")

const (
	maxLength     = 10
	numberOfWords = 10
)

func generateRandomWord() string {
	length := rand.Intn(maxLength) + 1 // Ensure at least 1 character
	runes := make([]rune, length)
	for i := range runes {
		runes[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(runes)
}

func generateRandomSentence() string {
	res := []string{}
	for i := 0; i < numberOfWords; i++ {
		res = append(res, generateRandomWord())
	}
	return strings.Join(res, " ")
}
