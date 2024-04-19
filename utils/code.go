package utils

import (
	"fmt"
	"math/rand"
)

// generate a length of 4 random character and number
// with first and third character is a letter
// and second and fourth character is a number
func GenUniqueCode() string {
	letter := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	number := "0123456789"

	letter1 := letter[rand.Intn(len(letter))]
	number1 := number[rand.Intn(len(number))]
	letter2 := letter[rand.Intn(len(letter))]
	number2 := number[rand.Intn(len(number))]

	return fmt.Sprintf("%c%c%c%c", letter1, number1, letter2, number2)
}
