# CRUD_Books_Go

## DB framework
https://github.com/jackc/pgx


## Migrations

### запуск
миграции на создание таблиц
```bash
migrate -path ./migrations -database "$DATABASE_URL" up
```

миграции на удаление таблиц
```bash
migrate -path ./migrations -database "$DATABASE_URL" down
```
```bash
$DATABASE_URL = "pgx://user:password@localhost:5432/database_name"
```

Cпец символы необходимо закодировать. Например:
* Пробел → %20
* @ → %40
* : → %3A
* $ → %24
* & → %26


### Незавершенные миграции

Еcли произошла ошибка при миграции (dirty) можно использовать

# чтобы отметить 1-ю миграцию как «успешно применённую» и сбросить флаг dirty:
```bash
migrate -path ./migrations \
        -database "pgx://user:password@localhost:5432/database_name?sslmode=disable" \
        force 1
```
или

чтобы откатиться «до начала» (до версии 1) и сбросить dirty:
```bash
migrate -path ./migrations \
        -database "pgx://user:password@localhost:5432/database_name?sslmode=disable" \
        force 0
```

сброс схемы миграций
```bash
migrate -path ./migrations \
        -database "pgx://user:password@localhost:5432/database_name?sslmode=disable" \
        drop
```
