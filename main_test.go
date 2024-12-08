package main

import (
	"testing"
)

func TestAddProduct(t *testing.T) {
	cart := NewCart()
	product := &Product{ID: 1, Name: "Молоко", Price: 25.5, Count: 2}

	cart.AddProduct(product)

	if len(cart.Products) != 1 {
		t.Errorf("Очікувана кількість продуктів: 1, отримано: %d", len(cart.Products))
	}

	if cart.Products[1].Count != 2 {
		t.Errorf("Очікувана кількість товару 'Молоко': 2, отримано: %d", cart.Products[1].Count)
	}
}

func TestUpdateProduct(t *testing.T) {
	cart := NewCart()
	product := &Product{ID: 1, Name: "Молоко", Price: 25.5, Count: 2}
	cart.AddProduct(product)

	cart.UpdateProduct(1, 5)

	if cart.Products[1].Count != 5 {
		t.Errorf("Очікувана кількість товару 'Молоко' після оновлення: 5, отримано: %d", cart.Products[1].Count)
	}

	cart.UpdateProduct(2, 3)
	if len(cart.Products) != 1 {
		t.Errorf("Не повинно бути змін для товару, якого немає в кошику")
	}
}

func TestRemoveProduct(t *testing.T) {
	cart := NewCart()
	product := &Product{ID: 1, Name: "Молоко", Price: 25.5, Count: 2}
	cart.AddProduct(product)

	cart.RemoveProduct(1)

	if len(cart.Products) != 0 {
		t.Errorf("Очікувана кількість продуктів після видалення: 0, отримано: %d", len(cart.Products))
	}

	cart.RemoveProduct(2) // Спроба видалити неіснуючий продукт
	if len(cart.Products) != 0 {
		t.Errorf("Кількість продуктів повинна залишитися 0 після видалення неіснуючого продукту")
	}
}

func TestShowCart(t *testing.T) {
	cart := NewCart()
	product1 := &Product{ID: 1, Name: "Молоко", Price: 25.5, Count: 2}
	product2 := &Product{ID: 2, Name: "Хліб", Price: 15.0, Count: 1}
	cart.AddProduct(product1)
	cart.AddProduct(product2)

	cart.ShowCart()
	// Для більш детальної перевірки можна використовувати mocks або захоплення stdout.
}
