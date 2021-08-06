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
Описание альбома передавать
Сделать склейку вправо
В статус длину альбом менеджера
log.Error узнать что делает
В либу логов добавить гетер инфы о новых ошибках, каждые 5 минут проверять его

$env:GOOS = 'linux'; $env:GOARCH = 'arm64'; $env:CGO_ENABLED = '0'; go build -o lrn .

*/
