package main

import (
	"fmt"
	"os"
)

func calculate(num1 float64, operator string, num2 float64) float64 {
    switch operator {
    case "+":
        return num1 + num2
    case "-":
        return num1 - num2
    case "*":
        return num1 * num2
    case "/":
        return num1 / num2
    default:
        fmt.Println("Unknown operator:", operator)
        os.Exit(1)
        return 0
    }
}

func main() {
    var num1 float64
    var num2 float64
    var operator string

    fmt.Print("Enter first number: ")
    fmt.Scan(&num1)

    fmt.Print("Enter operator (+, -, *, /): ")
    fmt.Scan(&operator)

    fmt.Print("Enter second number: ")
    fmt.Scan(&num2)

    result := calculate(num1, operator, num2)
    fmt.Printf("%v %v %v = %v\n", num1, operator, num2, result)
}