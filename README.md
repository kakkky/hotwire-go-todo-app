# hotwire-go-todo-app

[hotwire-go](https://github.com/kakkky/hotwire-go) の example TODO アプリ。

同じ機能を 2 通りのテンプレートエンジンで実装している。

- `html-template/` — 標準の `html/template` を利用
- `ah-templ/` — [a-h/templ](https://github.com/a-h/templ) を利用

## 起動

```sh
# html/template 版
go run ./html-template

# a-h/templ 版 (先に templ generate が必要)
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate -path ./ah-templ
go run ./ah-templ
```

`ah-templ/` の `*_templ.go` は生成物なので gitignore しています。

いずれも `http://localhost:8080` で listen する。

## 機能

- 一覧 (`GET /todos`)
- 新規作成 (`GET /todos/new`, `POST /todos`)
- 編集 (`GET /todos/{id}/edit`, `POST /todos/{id}`)
- 削除 (`DELETE /todos/{id}`)

Turbo Frame / Turbo Stream をフォームや行単位の更新に使い、
ページ全体をリロードせずに UI を更新する挙動を確認できる。
