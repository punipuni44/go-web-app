package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
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
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Input Form</title>
	</head>
	<body>
		<h1>名前を入力してください</h1>
		<form action="/menu" method="POST">
			<input type="text" name="username">
			<input type="submit" value="送信">
		</form>
	</body>
	</html>
	`

	fmt.Fprint(w, html)
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
		fmt.Fprintln(w, "名前が入力されていません")
		return
	}

	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<body>
		<h2>%sさん、何をしますか？</h2>

		<form action="/game" method="POST">
			<input type="hidden" name="username" value="%s">
			<input type="submit" value="🎮ゲーム">
		</form>

		<form action="/settings" method="GET">
			<input type="submit" value="設定">
		</form>

	</body>
	</html>
	`, username, username)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	rand.Seed(time.Now().UnixNano())
	answer := rand.Intn(3) + 1

	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<body>
		<h2>数あてゲーム</h2>
		<p>1~3 の数字をあててください</p>

		<form action="/game/result" method="POST">
			<input type="hidden" name="username" value="%s">
			<input type="hidden" name="answer" value="%d">
			<input type="number" name="number" min="1" max="3">
			<input type="submit" value="勝負！">
		</form>

		<br>
		<form action="/menu" method="POST">
		    <input type="hidden" name="username" value="%s">
			<input type="submit" value="メニューに戻る">
		</form>
	</body>
	</html>
	`, username, answer, username)
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

	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<body>
		<h2>%sさんの結果</h2>
		<p>あなたの選択：%s</p>
		<p>正解：%s</p>
		<h3>%s</h3>
		<br>
		<form action="/menu" method="POST">
		    <input type="hidden" name="username" value="%s">
			<input type="submit" value="メニューに戻る">
		</form>
	</body>
	</html>
	`, username, number, answer, result)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<body>
		<h2>設定</h2>
		<p>この機能は現在開発中です。</p>
		<br>
		<form action="/menu" method="POST">
		    <input type="hidden" name="username" value="%s">
			<input type="submit" value="メニューに戻る">
		</form>
	</body>
	</html>
	`)
}
