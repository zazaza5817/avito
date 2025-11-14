# PR Reviewer Assignment Service (Test Task, Fall 2025)
Сервис для автоматического назначения ревьюеров на Pull Request'ы с возможностью управления командами и участниками.

## Технологический стек

- **Язык**: Go 1.22
- **База данных**: PostgreSQL 15

## Запуск через Docker Compose


```bash
# Клонируйте репозиторий
git clone https://github.com/zazaza5817/avito
cd avito

# Создайте .env файл 
cp .env.example .env

# Запустите все сервисы
docker-compose up --build 
# Или
make docker-up

# Сервис будет доступен на http://localhost:8080
```



## Makefile команды

```bash
make docker-up      # Запустить через Docker Compose
make docker-down    # Остановить Docker Compose
make lint           # Запустить линтер
make deps           # Установить зависимости
make fmt            # Форматировать код
make load-test:     # Запустить нагрузочное тестирование
```

## Структура проекта

```
.
├── cmd/
│   └── api/
│       └── main.go                     # Точка входа приложения
├── internal/
│   ├── auth/
│   │   └── jwt.go                      # JWT аутентификация
│   ├── config/
│   │   └── config.go                   # Конфигурация
│   ├── database/
│   │   └── database.go                 # Подключение к БД
│   ├── handlers/
│   │   ├── helpers.go                  # Вспомогательные функции
│   │   ├── pr_handler.go               # HTTP обработчики PR
│   │   ├── team_handler.go             # HTTP обработчики команд
│   │   └── user_handler.go             # HTTP обработчики пользователей
│   ├── middleware/
│   │   └── auth.go                     # Middleware авторизации
│   ├── models/
│   │   ├── errors.go                   # Модели ошибок
│   │   └── models.go                   # Модели данных
│   ├── repository/
│   │   ├── interfaces.go               # Интерфейсы репозиториев
│   │   ├── pr_repository.go            # Репозиторий PR
│   │   ├── team_repository.go          # Репозиторий команд
│   │   └── user_repository.go          # Репозиторий пользователей
│   ├── response/
│   │   └── response.go                 # Структуры ответов
│   └── service/
│       ├── errors.go                   # Ошибки бизнес-логики
│       ├── pr_service.go               # Бизнес-логика PR
│       ├── team_service.go             # Бизнес-логика команд
│       └── user_service.go             # Бизнес-логика пользователей
├── migrations/
│   ├── 000001_init_schema.up.sql       # Миграция схемы вверх
│   ├── 000001_init_schema.down.sql     # Миграция схемы вниз
│   ├── 000002_seed_data.up.sql         # Тестовые данные
│   └── 000002_seed_data.down.sql       # Откат тестовых данных
├── docker-compose.yml                  # Docker Compose конфигурация
├── Dockerfile                          # Multi-stage Docker build
├── k6-load-test.js                     # Нагрузочное тестирование k6
├── Makefile                            # Команды для разработки
├── openapi.yml                         # OpenAPI спецификация
├── .golangci.yml                       # Конфигурация линтера
├── .env.example                        # Пример переменных окружения
├── .gitignore                          # Игнорируемые файлы Git
├── go.mod                              # Go модули
├── go.sum                              # Контрольные суммы зависимостей
└── README.md                           # Этот файл
```

## Линтер

Проект использует `golangci-lint` с настройками:
- **errcheck** - Проверка обработки ошибок
- **govet** - Статический анализ кода
- **ineffassign** - Обнаружение неиспользуемых присваиваний
- **staticcheck** - Расширенная статическая проверка
- **unused** - Обнаружение неиспользуемого кода
- **misspell** - Проверка орфографии в комментариях
- **revive** - Проверка стиля кода 

Запуск линтера:
```bash
make lint
```

## Нагрузочное тестирование

Нагрузочные тесты на основе **k6**:

### Запуск тестов

```bash
make load-test
# Или
k6 run k6-load-test.js
```

### Результаты тестирования

При нагрузке ~571 RPS (40696 запросов за 71 секунду):
- **Средняя Latency**: 4.56ms
- **P95 Latency**: 18.43ms
- **P99 Latency**: 37.20ms
- **Error Rate**: 0%
- **Availability**: 100%

### Производительность по операциям

| Операция | Средняя | P90 | P95 |
|----------|---------|-----|-----|
| Health Check | 0.64ms | 1.29ms | 1.91ms |
| Создание PR | 4.80ms | 12.00ms | 19.17ms |
| Мерж PR | 3.00ms | 6.99ms | 13.74ms |
| Переназначение | 5.00ms | 11.49ms | 17.98ms |
| Добавление в команду | 4.42ms | 8.22ms | 15.66ms |
| Получение команды | 1.69ms | 2.75ms | 6.24ms |
| PR на ревью | 13.53ms | 24.88ms | 33.15ms |
| Изменение статуса | 3.41ms | 9.14ms | 14.99ms |


## Вопросы и принятые решения

### 1. Аутентификация
**Вопрос**: Как реализовать и выдавать AdminToken и UserToken?

**Решение**: Реализовано через JWT токены без хранения ролей в базе данных. Токены содержат claim `is_admin` (boolean), который определяет права доступа. Middleware проверяет наличие и валидность токена, а также соответствие роли требованиям эндпоинта.


## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| DB_HOST | Хост PostgreSQL | localhost |
| DB_PORT | Порт PostgreSQL | 5432 |
| DB_USER | Пользователь БД | postgres |
| DB_PASSWORD | Пароль БД | postgres |
| DB_NAME | Имя БД | pr_reviewer_db |
| DB_SSLMODE | SSL режим | disable |
| SERVER_PORT | Порт сервера | 8080 |
| SERVER_HOST | Хост сервера | 0.0.0.0 |
| ENV | Окружение | development |
| JWT_SECRET | Очкнь секретный ключ JWT | your-secret-key-change-in-production |


## Контакты
Telegram: @booba47

