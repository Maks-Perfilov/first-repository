package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func romeToInt(s string) int {

	rtoi := map[string]int{"I": 1, "V": 5, "X": 10}
	result := 0
	repeatCounter := 1
	for i := 0; i < len(s)-1; i++ {
		current, next := rtoi[string(s[i])], rtoi[string(s[i+1])]
		if current == next {
			repeatCounter++
			if repeatCounter > 3 || current == rtoi["V"] {
				panic("Недопустимая комбинация римских символов")
			}
		} else {
			if repeatCounter > 1 && next > current {
				panic("Недопустимая комбинация римских символов")
			} else {
				repeatCounter = 1
			}
		}
		if current < next {
			result -= current
		} else {
			result += current
		}
	}
	return result + rtoi[string(s[len(s)-1])]
}

func intToRome(num int) string {
	if num < 0 {
		panic("Римское число должно быть положительным")
	} else if num == 0 {
		return "nulla"
	}

	romeMap := map[int]string{
		100: "C", 90: "XC", 50: "L", 40: "XL", 10: "X",
		9: "IX", 5: "V", 4: "IV", 1: "I",
	}
	keys := []int{100, 90, 50, 40, 10, 9, 5, 4, 1}
	roman := ""
	for _, key := range keys {
		for num >= key {
			roman += romeMap[key]
			num -= key
		}
	}
	return roman
}

func parseExpression(expression string) (int, int, string, bool) {
	re := regexp.MustCompile(`^\s*(\d+|[ivxIVX]+)\s*([+\-*/])\s*(\d+|[ivxIVX]+)\s*$`)
	matches := re.FindStringSubmatch(expression)
	if matches == nil {
		panic("Выражение не является математическим или не соответствует формату")
	}

	num1Str, num2Str := strings.ToUpper(matches[1]), strings.ToUpper(matches[3])
	num1, num2 := convertToInt(num1Str), convertToInt(num2Str)
	isRome := strings.ContainsAny(num1Str, "IVX") || strings.ContainsAny(num2Str, "IVX")

	if (strings.ContainsAny(num1Str, "IVX") && !strings.ContainsAny(num2Str, "IVX")) ||
		(!strings.ContainsAny(num1Str, "IVX") && strings.ContainsAny(num2Str, "IVX")) {
		panic("Нельзя использовать одновременно арабские и римские цифры")
	}
	if num1 > 10 || num1 == 0 || num2 == 0 || num2 > 10 {
		panic("Число не может быть больше 10 или равно 0")
	}

	return num1, num2, matches[2], isRome
}

func convertToInt(numStr string) int {
	if strings.ContainsAny(numStr, "IVX") {
		return romeToInt(numStr)
	}
	num, err := strconv.Atoi(numStr)
	if err != nil {
		panic("Неверные числа")
	}
	return num
}

func calculate(num1, num2 int, operator string) int {
	switch operator {
	case "+":
		return num1 + num2
	case "-":
		return num1 - num2
	case "*":
		return num1 * num2
	case "/":
		if num2 == 0 {
			panic("Деление на ноль")
		}
		return num1 / num2
	default:
		panic("Недопустимый оператор")
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin) //Не знаю зач, это просто нужно для считывания
	for {
		fmt.Print("Введите выражение: ")

		input := make(chan string)
		go func() {
			text, _ := reader.ReadString('\n')
			input <- text
		}()

		timer := time.NewTimer(22 * time.Second)

		select {
		case userInput := <-input:
			fmt.Println("Вы ввели:", userInput) //оставел на всякий, для проверки когда считывание было через scan
			timer.Stop()
			num1, num2, operator, isRoman := parseExpression(userInput)
			result := calculate(num1, num2, operator)
			if isRoman {
				fmt.Printf("Результат: %s\n", intToRome(result))
			} else {
				fmt.Printf("Результат: %d\n", result)
			}
		case <-timer.C:
			fmt.Println("Время вышло! Программа завершена.")
			os.Exit(0)
		}
	}
}
