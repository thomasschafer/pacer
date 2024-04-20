package main

import (
	"embed"
	"fmt"
	"math/rand"
	"strings"
)

var charSet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@£$%^&*()`-='\\[~_+\"|{},")

const (
	maxLength     = 8
	numberOfWords = 20 // TODO: allow users to set this
)

type WordType int

const (
	random WordType = iota
	top1000
)

func generateRandomWord() []string {
	length := rand.Intn(maxLength) + 1
	res := make([]string, length)
	for i := range res {
		res[i] = string(charSet[rand.Intn(len(charSet))])
	}
	return res
}

//go:embed words/top1000.txt
var top1000WordFile embed.FS

func getRandomTop1000WordGen() func() []string {
	data, err := top1000WordFile.ReadFile("words/top1000.txt")
	if err != nil {
		panic(fmt.Sprintf("Error opening file: %s", err))
	}
	words := strings.Split(string(data), "\n")

	return func() []string {
		var res []string
		for _, c := range words[rand.Intn(len(words))] {
			res = append(res, string(c))
		}
		return res
	}
}

func getWordGenFunc(wordType WordType) func() []string {
	switch wordType {
	case random:
		return generateRandomWord
	case top1000:
		return getRandomTop1000WordGen()
	}
	panic(fmt.Sprintf("Missing case for wordType %v", wordType))
}

func generateRandomSentence(wordType WordType) []string {
	wordGenFunc := getWordGenFunc(wordType)
	res := []string{}
	for i := 0; i < numberOfWords; i++ {
		res = append(res, wordGenFunc()...)
		if i != numberOfWords-1 {
			res = append(res, " ")
		}
	}
	return res
}
