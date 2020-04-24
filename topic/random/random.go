package random

import "math/rand"

var Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomStringCustom(length int, seed []rune) string {
	arr := make([]rune, length)
	for i := 0; i < length; i++ {
		arr[i] = seed[rand.Intn(len(seed))]
	}
	return string(arr)
}

func RandomString(length int) string {
	arr := make([]rune, length)
	for i := 0; i < length; i++ {
		arr[i] = Letters[rand.Intn(len(Letters))]
	}
	return string(arr)
}
