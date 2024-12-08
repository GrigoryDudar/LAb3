package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

// Структура для зберігання книг
type BookStore struct {
	mu    sync.Mutex
	Books *[]string
}

func main() {
	// Ініціалізуємо сховище для книг
	books := &BookStore{
		Books: &[]string{},
	}

	// Реєструємо маршрути
	http.HandleFunc("/", books.HandleForm)
	http.HandleFunc("/books", books.HandleBooks)

	// Запускаємо сервер
	log.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Обробник для відображення форми та збереження даних
func (b *BookStore) HandleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Зчитуємо книгу з форми
		book := r.FormValue("book")
		if book != "" {
			// Додаємо книгу до списку за допомогою mutex
			b.mu.Lock()
			*b.Books = append(*b.Books, book)
			b.mu.Unlock()
		}
		http.Redirect(w, r, "/books", http.StatusSeeOther)
		return
	}

	// Відображення HTML-форми
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Улюблені книги</title>
	</head>
	<body>
		<h1>Введіть вашу улюблену книгу</h1>
		<form method="POST" action="/">
			<input type="text" name="book" placeholder="Назва книги" required>
			<button type="submit">Додати</button>
		</form>
		<a href="/books">Переглянути всі книги</a>
	</body>
	</html>`
	w.Write([]byte(tmpl))
}

// Обробник для відображення списку книг
func (b *BookStore) HandleBooks(w http.ResponseWriter, r *http.Request) {
	// Створюємо HTML для відображення списку книг
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Улюблені книги</title>
	</head>
	<body>
		<h1>Список улюблених книг</h1>
		<ul>
			{{range .}}
				<li>{{.}}</li>
			{{else}}
				<p>Список поки що порожній.</p>
			{{end}}
		</ul>
		<a href="/">Повернутися до форми</a>
	</body>
	</html>`

	// Генеруємо HTML з шаблону
	t, err := template.New("books").Parse(tmpl)
	if err != nil {
		http.Error(w, "Помилка шаблону", http.StatusInternalServerError)
		return
	}

	// Читаємо дані зі списку за допомогою mutex
	b.mu.Lock()
	defer b.mu.Unlock()
	t.Execute(w, *b.Books)
}
