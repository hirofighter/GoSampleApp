package controllers

import (
	"log"
	"net/http"
	"todo_app/app/models"
)

// Top画面
func top(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, "Hello", "layout", "public_navbar", "top")
	} else {
		http.Redirect(w, r, "/todos", 302)
	}
}

// Todo一覧
func index(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		// ログイン中の場合、ユーザー情報を取得
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		// ユーザー情報からTodoデータを抽出
		todos, _ := user.GetTodosByUser()
		// Todoをユーザー情報に格納
		user.Todos = todos
		generateHTML(w, user, "layout", "private_navbar", "index")
	}
}

// Todo新規作成
func todoNew(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		// 未ログイン
		http.Redirect(w, r, "/login", 302)
	} else {
		// Todo作成ページへ
		generateHTML(w, nil, "layout", "private_navbar", "todo_new")
	}
}

// Todo登録
func todoSave(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		// 未ログイン
		http.Redirect(w, r, "/login", 302)

	} else {
		// パラメータ取得
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		// セッション取得（ログイン確認）
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		// パラメータの"content"に登録されている内容を取得
		content := r.PostFormValue("content")
		// DB登録
		if err := user.CreateTodo(content); err != nil {
			log.Println(err)
		}
		// Todo一覧画面にリダイレクト
		http.Redirect(w, r, "/todos", 302)
	}
}

// Todo編集
func todoEdit(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		// 未ログイン
		http.Redirect(w, r, "/login", 302)
	} else {
		// セッション確認
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		// 選択されたTodoを取得
		t, err := models.GetTodo(id)
		if err != nil {
			log.Println(err)
		}
		// Edit画面を表示
		generateHTML(w, t, "layout", "private_navbar", "todo_edit")
	}
}

// Todo更新
func todoUpdate(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		// 入力フォームの内容を取得
		content := r.PostFormValue("content")
		// Todoのstructに格納
		t := &models.Todo{ID: id, Content: content, UserID: user.ID}
		// DB更新
		if err := t.UpdateTodo(); err != nil {
			log.Println(err)
		}
		// Todo一覧画面にリダイレクト
		http.Redirect(w, r, "/todos", 302)

	}
}

// Todo削除
func todoDelete(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)

	} else {
		_, err := sess.GetUserBySession()
		if err != nil {
			// 未ログイン
			log.Println(err)
		}
		// 該当のTodoデータ取得
		t, err := models.GetTodo(id)
		if err != nil {
			log.Println(err)
		}
		// DB更新（削除）
		if err := t.DeleteTodo(); err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/todos", 302)

	}

}
