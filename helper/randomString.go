package helper

import "math/rand"

var letterChoice string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(length int) string {
	randString := ""
	for i := 0; i < length; i++ {
		randString += string(letterChoice[rand.Intn(len(letterChoice))])
	}
	return randString
}
