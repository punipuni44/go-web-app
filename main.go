package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var tmpl *template.Template

func main() {
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("template parse error: %v", err)
	}

	http.HandleFunc("/", formHandler)
	http.HandleFunc("/menu", menuHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/game/result", gameResultHandler)
	http.HandleFunc("/settings", settingsHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// フォームを表示するハンドラ
func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "form.html", nil)
}

// メニュー
func menuHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ許可
	if r.Method != http.MethodPost {
		http.Error(w, "不正アクセスです", http.StatusMethodNotAllowed)
		return
	}

	// フォームの値を取得
	username := r.FormValue("username")

	// 入力に応じてレスポンスを変える
	if username == "" {
		http.Error(w, "名前が入力されていません", http.StatusBadRequest)
		return
	}

	data := struct {
		Username string
	}{Username: username}

	tmpl.ExecuteTemplate(w, "menu.html", data)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	rand.Seed(time.Now().UnixNano())
	answer := rand.Intn(3) + 1

	data := struct {
		Username string
		Answer   int
	}{Username: username, Answer: answer}

	tmpl.ExecuteTemplate(w, "game.html", data)
}

func gameResultHandler(w http.ResponseWriter, r *http.Request) {
	// POSTリクエストのみ許可
	if r.Method != http.MethodPost {
		http.Error(w, "不正アクセスです", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	number := r.FormValue("number")
	answer := r.FormValue("answer")
	result := "はずれ((+_+))"
	if number == answer {
		result = "あたり(-.-)"
	}

	data := struct {
		Username string
		Number   string
		Answer   string
		Result   string
	}{Username: username, Number: number, Answer: answer, Result: result}

	tmpl.ExecuteTemplate(w, "result.html", data)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	data := struct{ Username string }{Username: username}
	tmpl.ExecuteTemplate(w, "settings.html", data)
}
