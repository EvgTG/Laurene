package main

import (
	"context"
)

func main() {
	app := New()
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	defer app.Stop(context.Background())
}

/*
TODO:
// /info от юзеров анон чата, выдавать дату рождения
tm := time.Now().Unix() - (694 * 24 * 60 * 60) - (21 * 60 * 60) - (51 * 60)
fmt.Println(time.Unix(tm, 0).String())

Счетчики
Генератор случайных чисел

$env:GOOS = 'linux'; $env:GOARCH = 'arm64'; $env:CGO_ENABLED = '0'; go build -o lrn .

*/
