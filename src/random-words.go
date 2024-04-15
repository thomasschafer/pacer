package main

import (
	"math/rand"
)

var charSet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@£$%^&*()`-='\\[~_+\"|{},")

const (
	maxLength     = 10
	numberOfWords = 10
)

func generateRandomWord() []string {
	length := rand.Intn(maxLength) + 1
	res := make([]string, length)
	for i := range res {
		res[i] = string(charSet[rand.Intn(len(charSet))])
	}
	return res
}

func generateRandomSentence() []string {
	res := []string{}
	for i := 0; i < numberOfWords; i++ {
		res = append(res, generateRandomWord()...)
		if i != numberOfWords-1 {
			res = append(res, " ")
		}
	}
	return res
}
