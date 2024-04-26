package main

import (
	"fmt"
	"math"
	"time"
)

// Задача 1
// [1,2,3,4,5] -> [1,3,5]
// [2,2,2] -> []
// [1,2,2] -> [1]
func filter(a []int) []int {
	var n int

	for i := 0; i < len(a); i++ {
		if a[i]%2 != 0 {
			a[n] = a[i]
			n++
		}
	}

	return a[:n]
}

// Задача 2

type Limiter interface {
	Allow(key string) bool
}

type tokenParams struct { // Про многопоточку: хочется добавить мьютекс😁
	max        int32
	speed      int32
	count      int32
	accessTime time.Time
}

type Baskets struct {
	basket map[string]tokenParams
}

func NewBaskets() Baskets {
	baskets := Baskets{
		basket: make(map[string]tokenParams),
	}

	return baskets
}

func (b Baskets) Allow(key string) bool {
	if _, exists := b.basket[key]; !exists {
		fmt.Printf("key %s not found", key)
		return false
	}

	creditedTokens := int32(0)
	now := time.Now()

	if !b.basket[key].accessTime.IsZero() {
		delta := now.Sub(b.basket[key].accessTime).Seconds()
		creditedTokens = int32(math.Floor(delta)) * b.basket[key].speed // за < 1 секунды начисляется 0

		if creditedTokens > b.basket[key].max {
			creditedTokens = b.basket[key].max
		}
	} else {
		creditedTokens = b.basket[key].max // первый раз начисляем максимум
	}

	newCount := b.basket[key].count - 1 + creditedTokens
	if newCount < 0 { // не снижаем токены, если токены потрачены
		return false
	}

	b.basket[key] = tokenParams{
		max:        b.basket[key].max,
		speed:      b.basket[key].speed,
		count:      newCount,
		accessTime: now,
	}

	return true
}

func main() {
	fmt.Println(filter([]int{2, 2, 1, 4, 4, 5, 7, 7, 7}))

	basket := NewBaskets()
	basket.basket["key"] = tokenParams{
		max:        5,
		speed:      1,
		count:      5,
		accessTime: time.Time{},
	}
	fmt.Println(basket.Allow("key"))
	fmt.Println(basket.Allow("key"))
	fmt.Println(basket.Allow("key"))
	fmt.Println(basket.Allow("key"))
	fmt.Println(basket.Allow("key"))
}
