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

Текст: капс, tolower, раскладки, кодировки, шифрования
Фото: шакал, искажения
Счетчики
Генератор случайных чисел

$env:GOOS = 'linux'; $env:GOARCH = 'arm64'; $env:CGO_ENABLED = '0'; go build -o lrn .

*/
