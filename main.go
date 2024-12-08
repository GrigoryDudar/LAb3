package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type ConversionRequest struct {
	Amount       float64
	FromCurrency string
	ToCurrency   string
}

type ConversionResult struct {
	ConvertedAmount float64
	Rate            float64
	Error           error
}

func fetchConversionRate(from, to string) (float64, error) {
	apiURL := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", from)

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch conversion rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	rates, ok := data["rates"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid response format")
	}

	rate, ok := rates[to].(float64)
	if !ok {
		return 0, fmt.Errorf("conversion rate for %s not found", to)
	}

	return rate, nil
}

func convertCurrency(req ConversionRequest, wg *sync.WaitGroup, ch chan<- ConversionResult) {
	defer wg.Done()

	rate, err := fetchConversionRate(req.FromCurrency, req.ToCurrency)
	if err != nil {
		ch <- ConversionResult{Error: err}
		return
	}

	converted := req.Amount * rate
	ch <- ConversionResult{
		ConvertedAmount: converted,
		Rate:            rate,
		Error:           nil,
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the Currency Converter!")

	// Запитуємо дані у користувача
	fmt.Print("Enter the amount to convert: ")
	amountInput, _ := reader.ReadString('\n')
	amountInput = strings.TrimSpace(amountInput)
	amount, err := strconv.ParseFloat(amountInput, 64)
	if err != nil || amount <= 0 {
		fmt.Println("Invalid amount. Please enter a positive number.")
		return
	}

	fmt.Print("Enter the currency you are converting from (e.g., USD): ")
	fromCurrency, _ := reader.ReadString('\n')
	fromCurrency = strings.ToUpper(strings.TrimSpace(fromCurrency))

	fmt.Print("Enter the currency you are converting to (e.g., EUR): ")
	toCurrency, _ := reader.ReadString('\n')
	toCurrency = strings.ToUpper(strings.TrimSpace(toCurrency))

	// Створюємо запит на конвертацію
	request := ConversionRequest{
		Amount:       amount,
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	}

	var wg sync.WaitGroup
	ch := make(chan ConversionResult, 1)

	// Запускаємо горутину для обробки запиту
	wg.Add(1)
	go convertCurrency(request, &wg, ch)

	// Чекаємо завершення
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Отримуємо результат
	for result := range ch {
		if result.Error != nil {
			fmt.Printf("Error: %v\n", result.Error)
		} else {
			fmt.Printf("Converted Amount: %.2f %s (Rate: %.4f)\n", result.ConvertedAmount, toCurrency, result.Rate)
		}
	}
}
