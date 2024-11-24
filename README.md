# music-storage

music storage written in golang


# Инструкция по запуску

1. Клонируем репозиторий:

`
git clone git@github.com:cutlery47/music-storage.git
`

2. После успешного копирования репозитория, в корневой папке будет лежать файл .env с переменными окружения по умоланию. 
Пользователь может изменять этот файл, однако этого НЕ рекоммендуется делать :) 


3. В корневой директории проекта прописываем:

`
make up_build
`

В результате выполнения данной команды у нас сбилдятся образы и запустятся контейнеры.


4) Запускаем приложение:

При первом запуске пользователю настоятельно рекоммендуется ознакомиться с API приложения. Для этого, достаточно обратиться к http://localhost:8080/swagger/ при помощи браузера (если использовать настройки по умолчанию).

Тут же можно пощупать и остальные эндпойнты, благо делать это супер приятно благодаря интуитивному UI Open API.
