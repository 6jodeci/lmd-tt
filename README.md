## Тестовое задание Lamoda

API предоставляет возможность управления запасами продуктов на складе

### Технологии:
- Go
- PostgreSQL
- Docker
- Dbdocs

### Библиотеки:
- gin-gonic/gin 1.8.2
- swaggo/swag v1.3.3
- sirupsen/logrus 1.9.0
- ilyakaznacheev/cleanEnv 1.4.2

### Диаграмма Базы Данных:
![Untitled](https://user-images.githubusercontent.com/65400970/219599773-fb08868d-00cd-4e3c-baab-d231532da420.png)

### Документация Базы Данных:
- https://dbdocs.io/6jodeci/lamoda_test

### SwaggerUI:
![image](https://user-images.githubusercontent.com/65400970/219967009-707d8dd7-9335-40f5-b83d-63f3668de439.png)

### Запуск приложения локально (выполнять по порядку!):
- git clone https://github.com/6jodeci/lmd-tt 
- перейти в папку configs и переименовать файл example-env.txt в app.env
- make postgres
- make createdb
- make migratecreate
- make migrateup
- make migrateup
- cd app/cmd 
- go run main.go