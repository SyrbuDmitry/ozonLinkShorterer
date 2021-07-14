# ozonLinkShorterer
# Задание на стажировку Ozon
  Cервис сокращения ссылок.
  
## Принцип работы сервиса

Короткие ссылки основываются на id записи в БД, переведённой в систему счисления с алфавитом [A-Za-z0-9]. 
Для получения id оригинальной ссылки в БД короткая ссылка(число в системе счисления [A-Za-z0-9]) переводится в десятичное число.

## API

1. На вход поступает длинная ссылка, возвращается сокращённая ссылка\
Request:\
POST /short {"url": "long-url-here"}\
Response:\
{"url": "short-url-here"}

2. На вход поступает сокращённая ссылка, возвращается полная ссылка\
Request:\
POST /long {"url": "short-url-here"}\
Response:\
{"url": "long-url-here"}

3. В случае сбоя работы возвращается код ошибки.

### Пример запросов

```bash

$ curl --request POST \
    --data '{"url": "long-url-here"}' \
    http://localhost:8000/short

   
  {"url":"Bn"}

$ curl --request POST \
    --data '{"url": "Bn"}' \
    http://localhost:8000/long
    
  {"url":"long-url-here"}
```
## Запуск

С помощью Docker:

```bash
$ git clone https://github.com/SyrbuDmitry/ozonLinkShorterer.git
$ cd ~/ozonLinkShorterer
$ docker-compose build
$ docker-compose up
```

'Вручную':
```bash
$ git clone https://github.com/SyrbuDmitry/ozonLinkShorterer.git
$ cd ~/ozonLinkShorterer/cmd/ozonLinkShorterer
$ go run .
```

