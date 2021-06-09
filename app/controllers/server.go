package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"todo_app/app/models"
	"todo_app/config"
)

// 画面を表示するための処理
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

// クッキーの取得
func session(w http.ResponseWriter, r *http.Request) (sess models.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = models.Session{UUID: cookie.Value}
		if ok, _ := sess.CheckSession(); !ok {
			err = fmt.Errorf("Invalid session")

		}
	}
	return sess, err
}

// URL解析
var validPath = regexp.MustCompile("^/todos/(edit|update|delete)/([0-9]+)$")

// 激むず
// 「http.ResponseWriter, *http.Request, int」を引数とする関数を引数として"parseURL"メソッドに渡す
// HandlerFunc型の値が戻り値となる
func parseURL(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ex) pathが "/todos/edit/1" の場合
		// リクエストのパスと"validPath"のパスでマッチした部分をスライスで取得
		q := validPath.FindStringSubmatch(r.URL.Path)
		// q には{"/todos/edit/1" , "edit", "1" } が入ってる
		// fmt.Println(q)
		if q == nil {
			// マッチしないのでNotFound
			http.NotFound(w, r)
			return
		}
		// q[2]にはidが入っているので id を数値型に変換する
		qi, err := strconv.Atoi(q[2])
		if err != nil {
			// 数値に変換できないのでエラーとなりNotFound
			http.NotFound(w, r)
			return
		}

		// 引数で渡した関数を実行する
		fn(w, r, qi)
	}
}

// URLの設定
func StartMainServer() error {
	files := http.FileServer(http.Dir(config.Config.Static))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	http.HandleFunc("/", top)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/todos", index)
	http.HandleFunc("/todos/new", todoNew)
	http.HandleFunc("/todos/save", todoSave)
	http.HandleFunc("/todos/edit/", parseURL(todoEdit))
	http.HandleFunc("/todos/update/", parseURL(todoUpdate))
	http.HandleFunc("/todos/delete/", parseURL(todoDelete))
	return http.ListenAndServe(":"+config.Config.Port, nil)

	// sqlite3用
	// return http.ListenAndServe(":"+config.Config.Port, nil)

	// heroku Up用
	port := os.Getenv("PORT")
	return http.ListenAndServe(":"+port, nil)

}
