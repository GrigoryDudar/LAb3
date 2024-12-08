package main

import (
	"sync"
	"testing"
)

func TestIncrement(t *testing.T) {
	counter := NewUserCounter()

	// Перевірка збільшення лічильника для одного користувача
	if visits := counter.Increment("Alice"); visits != 1 {
		t.Errorf("Очікуваний візит для 'Alice': 1, отримано: %d", visits)
	}

	if visits := counter.Increment("Alice"); visits != 2 {
		t.Errorf("Очікуваний візит для 'Alice': 2, отримано: %d", visits)
	}

	// Перевірка для іншого користувача
	if visits := counter.Increment("Bob"); visits != 1 {
		t.Errorf("Очікуваний візит для 'Bob': 1, отримано: %d", visits)
	}
}

func TestConcurrentIncrement(t *testing.T) {
	counter := NewUserCounter()
	var wg sync.WaitGroup
	userID := "Alice"
	expectedVisits := 1000

	// Запускаємо 1000 горутин для одного користувача
	for i := 0; i < expectedVisits; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment(userID)
		}()
	}

	wg.Wait()

	// Перевірка кількості візитів після завершення горутин
	if visits := counter.Increment(userID); visits != expectedVisits+1 {
		t.Errorf("Очікуваний візит для '%s': %d, отримано: %d", userID, expectedVisits+1, visits)
	}
}
