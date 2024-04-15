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
		for _, s := range generateRandomWord() {
			res = append(res, s)
		}
	}
	return res
}
