package lib

import (
	"math/rand"
	"reflect"
	"time"
)

// RandomStringRunes generate a random string
func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func PickRandomValue(values []interface{}) interface{} {
	min := 0
	max := reflect.ValueOf(values).Len()
	randomNumber := rand.Intn(max-min) + min

	return values[randomNumber]
}

func GetRandomValueInRange(min, max int) int {
	randomNumber := rand.Intn(max-min) + min

	return randomNumber
}

func GetRandomValueInRange2(min, max int) int {
	rand.Seed(time.Now().UnixNano() - 1012)
	randomNumber := float64(min) + rand.Float64()*(float64(max)-float64(min))

	return int(randomNumber)
}
