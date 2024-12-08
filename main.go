package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Product структура для зберігання інформації про товар
type Product struct {
	ID    int
	Name  string
	Price float64
	Count int
}

// Cart структура для зберігання товарів у кошику
type Cart struct {
	Products map[int]*Product
	mu       sync.Mutex
}

// NewCart створює новий кошик
func NewCart() *Cart {
	return &Cart{Products: make(map[int]*Product)}
}

// AddProduct додає товар до кошика
func (c *Cart) AddProduct(product *Product) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if p, exists := c.Products[product.ID]; exists {
		p.Count += product.Count
	} else {
		c.Products[product.ID] = product
	}
	fmt.Printf("Додано: %s (%d шт.)\n", product.Name, product.Count)
}

// UpdateProduct оновлює кількість товару в кошику
func (c *Cart) UpdateProduct(productID, count int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if p, exists := c.Products[productID]; exists {
		p.Count = count
		fmt.Printf("Оновлено: %s, нова кількість: %d\n", p.Name, p.Count)
	} else {
		fmt.Println("Товар не знайдено в кошику!")
	}
}

// RemoveProduct видаляє товар з кошика
func (c *Cart) RemoveProduct(productID int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.Products[productID]; exists {
		delete(c.Products, productID)
		fmt.Println("Товар видалено з кошика.")
	} else {
		fmt.Println("Товар не знайдено в кошику!")
	}
}

// ShowCart виводить вміст кошика
func (c *Cart) ShowCart() {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println("Ваш кошик:")
	if len(c.Products) == 0 {
		fmt.Println("- Кошик порожній.")
		return
	}
	for _, product := range c.Products {
		fmt.Printf("- %s: %d шт. (%.2f грн за шт.)\n", product.Name, product.Count, product.Price)
	}
}

// Головна функція
func main() {
	cart := NewCart()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Система керування кошиком. Використовуйте команди: add, update, remove, show, exit")
	for {
		fmt.Print("\nВведіть команду: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch strings.ToLower(input) {
		case "add":
			fmt.Print("Введіть ID, назву, ціну та кількість через кому (наприклад: 1,Молоко,25.50,2): ")
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			parts := strings.Split(line, ",")
			if len(parts) != 4 {
				fmt.Println("Неправильний формат. Спробуйте ще раз.")
				continue
			}

			id, _ := strconv.Atoi(parts[0])
			name := parts[1]
			price, _ := strconv.ParseFloat(parts[2], 64)
			count, _ := strconv.Atoi(parts[3])

			cart.AddProduct(&Product{ID: id, Name: name, Price: price, Count: count})

		case "update":
			fmt.Print("Введіть ID товару та нову кількість через кому (наприклад: 1,5): ")
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			parts := strings.Split(line, ",")
			if len(parts) != 2 {
				fmt.Println("Неправильний формат. Спробуйте ще раз.")
				continue
			}

			id, _ := strconv.Atoi(parts[0])
			count, _ := strconv.Atoi(parts[1])

			cart.UpdateProduct(id, count)

		case "remove":
			fmt.Print("Введіть ID товару для видалення: ")
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			id, _ := strconv.Atoi(line)

			cart.RemoveProduct(id)

		case "show":
			cart.ShowCart()

		case "exit":
			fmt.Println("Вихід із програми.")
			return

		default:
			fmt.Println("Невідома команда. Спробуйте add, update, remove, show або exit.")
		}
	}
}
