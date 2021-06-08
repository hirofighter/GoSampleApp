package controllers

import (
	"log"
	"net/http"
	"todo_app/app/models"
)

func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, nil, "layout", "public_navbar", "signup")
		} else {
			http.Redirect(w, r, "/todos", 302)
		}

	} else if r.Method == "POST" {
		// 入力値を解析
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		// UserのStructに格納
		user := models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			PassWord: r.PostFormValue("password"),
		}
		// DB登録する
		if err := user.CreateUser(); err != nil {
			log.Fatalln(err)
		}
		// Topページにリダイレクト
		http.Redirect(w, r, "/", 302)

	}
}

func login(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, nil, "layout", "public_navbar", "login")

	} else {
		http.Redirect(w, r, "/todos", 302)
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	user, err := models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Println(err)
		// エラー時はログイン画面にリダイレクト
		http.Redirect(w, r, "/login", 302)
	}
	// 入力されたパスワードのチェック
	// 入力値を暗号化してDBの値と比較する
	if user.PassWord == models.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			log.Println(err)
		}

		// クッキーの作成
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID,
			HttpOnly: true,
		}
		// クッキーをセットする
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", 302)
	} else {
		// ログイン失敗時はログイン画面にリダイレクト
		http.Redirect(w, r, "/login", 302)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		log.Println(err)
	}

	if err != http.ErrNoCookie {
		session := models.Session{UUID: cookie.Value}
		session.DeleteSessionByUUID()
		// fmt.Println("delete")
	}
	http.Redirect(w, r, "/login", 302)
	// fmt.Println("redirect")
}
