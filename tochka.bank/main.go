package main

import "fmt"

// task 1. Заменить символ точку на ы
func main() {
	s := "Привет."
	// s == "Приветы"

	// база
	newStr := ""
	for _, r := range s {
		if r == '.' {
			newStr += "ы"
		} else {
			newStr += string(r)
		}
	}

	fmt.Println(newStr)

	// Оптимальное решение
	result := []rune(s)
	fmt.Println(string(result[:len(result)-1]) + "ы") // либо через WriteByte

	// task 2. Написать, что выведется
	// test
	// nil
	task2()
}

type MyStruct struct {
	Value string
}

type MyInterface interface {
	GetValue() string
}

func (m *MyStruct) GetValue() string {
	return m.Value
}

func ToNil(obj MyInterface) {
	if obj == nil {
		fmt.Printf("obj = nil")
		return
	}
	obj = nil
}

func task2() {
	myStruct := &MyStruct{Value: "test"}

	ToNil(myStruct)
	fmt.Printf("myStruct: %+v \n", myStruct)

	myStruct = nil

	ToNil(myStruct)
	fmt.Printf("myStruct: %+v \n", myStruct)
}
