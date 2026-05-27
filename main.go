package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
)

func isUnaryOperator(op string) bool {
	switch op {
	case "sqrt", "log", "sin", "cos", "tan", "!":
		return true
	default:
		return false
	}
}

func calculateUnary(num float64, operator string) (float64, error) {
	switch operator {
	case "sqrt":
		if num < 0 {
			return 0, errors.New("cannot take square root of a negative number")
		}
		return math.Sqrt(num), nil
	case "log":
		if num <= 0 {
			return 0, errors.New("cannot take logarithm of non-positive number")
		}
		return math.Log(num), nil
	case "sin":
		return math.Sin(num), nil
	case "cos":
		return math.Cos(num), nil
	case "tan":
		return math.Tan(num), nil
	case "!":
		if num != float64(int(num)) {
			return 0, errors.New("cannot take factorial of a non-integer")
		}
		n := int(num)
		if n < 0 {
			return 0, errors.New("cannot take factorial of a negative number")
		}
		result := 1
		for i := 1; i <= n; i++ {
			result *= i
		}
		return float64(result), nil
	default:
		return 0, fmt.Errorf("unknown unary operator: %v", operator)
	}
}

func calculateBinary(num1 float64, operator string, num2 float64) (float64, error) {
	switch operator {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	case "*":
		return num1 * num2, nil
	case "/":
		if num2 == 0 {
			return 0, errors.New("cannot divide by zero")
		}
		return num1 / num2, nil
	case "^":
		return math.Pow(num1, num2), nil
	case "%":
		if num2 == 0 {
			return 0, errors.New("cannot modulo by zero")
		}
		return math.Mod(num1, num2), nil	
	default:
		return 0, fmt.Errorf("unknown operator: %v", operator)
	}
}

func formatResult(f float64) string {
	if f == math.Trunc(f) {
		return strconv.FormatFloat(f, 'f', 0, 64)
	}
	return strconv.FormatFloat(f, 'f', 6, 64)
}

func main() {
	fmt.Println("Go Scientific Calculator")
	fmt.Println("------------------------")
	fmt.Println("Binary:  + - * / ^ %")
	fmt.Println("Unary:   sqrt  log  sin  cos  tan")
	fmt.Println("Unary:   ! (e.g. enter 5 then !)")
	fmt.Println("Type 'q' as operator to quit")
	fmt.Println("------------------------")

	for {
		var num1 float64
		var operator string

		fmt.Print("\nFirst number: ")
		_, err := fmt.Scan(&num1)
		if err != nil {
			fmt.Println("Invalid input — enter a number.")
			fmt.Scan(&operator)
			continue
		}

		fmt.Print("Operator: ")
		fmt.Scan(&operator)

		if operator == "q" {
			fmt.Println("Bye!")
			os.Exit(0)
		}

		if isUnaryOperator(operator) {
			result, err := calculateUnary(num1, operator)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			} else {
				fmt.Printf("Result: %s\n", formatResult(result))
				continue
			}
		}

		var num2 float64
		fmt.Print("Second number: ")
		_, err = fmt.Scan(&num2)
		if err != nil {
			fmt.Println("Invalid input — enter a number.")
			fmt.Scan(&operator)
			continue
		}

		result, err := calculateBinary(num1, operator, num2)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		} else {
			fmt.Printf("Result: %s\n", formatResult(result))
		}


	}
}