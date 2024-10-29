# Image Previewer

![GitHub Actions](https://img.shields.io/github/actions/workflow/status/devv4n/image-previewer/ci.yml?branch=main)
![GitHub Release](https://img.shields.io/github/release/devv4n/image-previewer.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/devv4n/image-previewer)
![Go Version](https://img.shields.io/badge/Go-1.23%2B-blue)

## Описание

Image Previewer - это проект, который предоставляет возможность предварительного просмотра изображений, обрабатывая запросы на изменение размера изображений и их кэширование.

## Установка

1. Убедитесь, что у вас установлен [Docker](https://www.docker.com/get-started).
2. Склонируйте репозиторий:

   ```bash
   git clone https://github.com/devv4n/image-previewer.git
   cd image-previewer
   ```

3. Запустите приложение:

   ```bash
   make run
   ```
   
## Конфигурация

Путь до конфигурационного файла указывается через флаг `-c` Default: `.config.json`

   ```bash
    go run ./cmd/image-previewer/main.go -c=.example-cfg.json
   ```

## Тестирование

Чтобы запустить интеграционные тесты, используйте:

```bash
make integration-test
```

## Использование

После запуска приложения, вы можете отправлять HTTP-запросы для изменения размера изображений. Пример запроса:

```
GET http://localhost:8080/fill/<width>/<height>/<image-url>
```
