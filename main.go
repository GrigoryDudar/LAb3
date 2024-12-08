package main

import (
	"fmt"
	"sync"
	
)

// UserCounter зберігає кількість відвідувань для кожного користувача
type UserCounter struct {
	mu       sync.Mutex
	counters map[string]int
}

// NewUserCounter створює новий UserCounter
func NewUserCounter() *UserCounter {
	return &UserCounter{
		counters: make(map[string]int),
	}
}

// Increment збільшує кількість відвідувань користувача
func (uc *UserCounter) Increment(userID string) int {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.counters[userID]++
	return uc.counters[userID]
}

func main() {
	// Створюємо об'єкт для відстеження кількості відвідувань
	counter := NewUserCounter()

	// Функція для генерації привітання
	greet := func(userID string) {
		// Збільшуємо кількість відвідувань
		visits := counter.Increment(userID)
		// Формуємо і друкуємо привітання
		fmt.Printf("Привіт, %s! Це ваш %d-й візит!\n", userID, visits)
	}

	// Моделюємо одночасну обробку запитів
	users := []string{"Alice", "Bob", "Alice", "Charlie", "Bob", "Alice"}
	var wg sync.WaitGroup

	for _, user := range users {
		wg.Add(1)
		go func(user string) {
			defer wg.Done()
			greet(user)
		}(user)
	}

	// Чекаємо завершення всіх горутин
	wg.Wait()
}
