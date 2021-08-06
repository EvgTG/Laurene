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
В статус длину альбом менеджера
В либу логов добавить гетер инфы о новых ошибках, каждые 5 минут проверять его
Сделать склейку вправо, квадрат
log.Error узнать что делает

$env:GOOS = 'linux'; $env:GOARCH = 'arm64'; $env:CGO_ENABLED = '0'; go build -o lrn .

*/
