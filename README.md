# LaureneAssistantBot
Помощник с различными возможностями (текст, фото, всякое).
### Docker
```
docker build --platform linux/amd64 -t laurenebot:latest -t laurenebot:1.x.x .
docker run -d -v /path/files/:/app/files/ --net=host --name lrn laurenebot
```

### Environment
обязательные - *
* LOGLVL (panic, fatal, error, warn or warning, info, debug, trace. По дефолту info)
* NAMEDB* (Database name)
* MONGOURL* (Ссылка для подключения mongodb)
* TOKENTG* (telegram bot api token)
* USERLIST*,ADMINLIST*,NOTIFYLIST*,ERRORLIST* (user IDs - "id,id,id")
* LOC (локация для времени, смотреть tzdata)
* PINGPORT (Порт для проверки работоспособности бота, например UptimeRobot. Пример ссылки по которой будет доступ - "http://[ip]:6975/pingLaurene", отсутствие PINGPORT - сервер для пинга не запуститься.)

Пример:  
LOGLVL=INFO  
TOKENTG=19209:AAFSsiJY  
MONGOURL=mongodb://log:pass@127.0.0.1:27017  
NAMEDB=loren  
USERLIST=123456789,352536  
ADMINLIST=123456789,352536  
NOTIFYLIST=123456789,352536  
ERRORLIST=123456789,352536  
LOC=Europe/Moscow  
PINGPORT=6975  

### Папки

```
files/          (папка и рабочее место бота)
    cfg.env     (конфиг)
    logrus.log  (файл логов)
    temp/       (временные файлы)
```