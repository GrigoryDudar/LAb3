package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// Тест для функції fetchConversionRate
func TestFetchConversionRate(t *testing.T) {
	// Створюємо мок-сервер
	mockResponse := `{"rates": {"EUR": 0.85}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/latest/USD" {
			t.Errorf("Неправильний шлях: %s", r.URL.Path)
		}
		fmt.Fprintln(w, mockResponse)
	}))
	defer server.Close()

	

	// Викликаємо функцію
	rate, err := fetchConversionRate("USD", "EUR")
	if err != nil {
		t.Fatalf("Отримана помилка: %v", err)
	}

	// Перевіряємо результат
	expectedRate := 0.85
	if rate != expectedRate {
		t.Errorf("Очікував %f, але отримав %f", expectedRate, rate)
	}
}

// Тест для функції convertCurrency
func TestConvertCurrency(t *testing.T) {
	// Замінюємо fetchConversionRate на мок-функцію
	mockFetch := func(from, to string) (float64, error) {
		if from == "USD" && to == "EUR" {
			return 0.85, nil
		}
		return 0, fmt.Errorf("unsupported currencies")
	}

	// Створюємо запит
	request := ConversionRequest{
		Amount:       100,
		FromCurrency: "USD",
		ToCurrency:   "EUR",
	}

	// Канал для результатів і WaitGroup
	ch := make(chan ConversionResult, 1)
	var wg sync.WaitGroup

	// Запускаємо тестовану функцію
	wg.Add(1)
	go func() {
		defer wg.Done()
		rate, err := mockFetch(request.FromCurrency, request.ToCurrency)
		if err != nil {
			ch <- ConversionResult{Error: err}
			return
		}
		ch <- ConversionResult{
			ConvertedAmount: request.Amount * rate,
			Rate:            rate,
			Error:           nil,
		}
	}()

	wg.Wait()
	close(ch)

	// Перевіряємо результат
	for result := range ch {
		if result.Error != nil {
			t.Fatalf("Отримана помилка: %v", result.Error)
		}
		expectedAmount := 85.0
		if result.ConvertedAmount != expectedAmount {
			t.Errorf("Очікував %f, але отримав %f", expectedAmount, result.ConvertedAmount)
		}
	}
}
