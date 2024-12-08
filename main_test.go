package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestHandleForm(t *testing.T) {
	// Ініціалізуємо BookStore
	store := &BookStore{
		Books: &[]string{},
		mu:    sync.Mutex{},
	}

	// Створюємо POST-запит із даними форми
	formData := "book=TestBook"
	req, err := http.NewRequest("POST", "/", strings.NewReader(formData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Створюємо ResponseRecorder для фіксації відповіді
	rr := httptest.NewRecorder()

	// Викликаємо обробник
	handler := http.HandlerFunc(store.HandleForm)
	handler.ServeHTTP(rr, req)

	// Перевіряємо, чи був редірект на /books
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("Очікував статус %d, але отримав %d", http.StatusSeeOther, status)
	}

	// Перевіряємо, чи книга була додана до сховища
	store.mu.Lock()
	defer store.mu.Unlock()
	if len(*store.Books) != 1 || (*store.Books)[0] != "TestBook" {
		t.Errorf("Очікував книгу 'TestBook', але отримав %+v", *store.Books)
	}
}

func TestHandleBooks(t *testing.T) {
	// Ініціалізуємо BookStore з попередньо доданими книгами
	store := &BookStore{
		Books: &[]string{"Book1", "Book2"},
		mu:    sync.Mutex{},
	}

	// Створюємо GET-запит
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Створюємо ResponseRecorder
	rr := httptest.NewRecorder()

	// Викликаємо обробник
	handler := http.HandlerFunc(store.HandleBooks)
	handler.ServeHTTP(rr, req)

	// Перевіряємо статус відповіді
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Очікував статус %d, але отримав %d", http.StatusOK, status)
	}

	// Перевіряємо, чи відповіді містять назви книг
	responseBody := rr.Body.String()
	if !strings.Contains(responseBody, "Book1") || !strings.Contains(responseBody, "Book2") {
		t.Errorf("Очікував, що відповідь містить 'Book1' і 'Book2', але отримав %s", responseBody)
	}
}
