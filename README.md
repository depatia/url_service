## Основные компоненты

### 1. Domain Layer (`interfaces/`)
- `IRepository` - интерфейс хранилища
- `IHealthCheckerService` - интерфейс сервиса проверки сайтов  
- `IPdfService` - интерфейс генератора PDF

### 2. Delivery Layer (`delivery/`)
- `HTTP Handlers` - обработчики REST API (с применением gin)
- `Swagger Documentation` - автоматическая документация API

### 3. Service Layer (`services/`)
- `HealthCheckerService` - логика проверки доступности сайтов
- `PdfGeneratorService`- генерация PDF отчетов по результатам проверки

### 4. Repository Layer (`repository/`)
- `In-Memory Storage` - хранение результатов проверок в памяти
- потокобезопасные операции с использованием `sync.RWMutex`

## API Endpoints

POST `/api/check-sites`
Проверяет доступность переданных сайтов.  
POST `/api/generate-report`
Генерирует pdf файл по результатам проверки (по уникальному идентификатору проверки).  
GET `/health`
Проверка работоспособности сервера на текущий момент.

## Swagger

[`localhost:8080/swagger`](http://localhost:8080/swagger/index.html)

## Конфигурация

Конфигурация находится в файле config.yaml

## Архитектура

Проект построен по принципам **Clean Architecture** и **Dependency Injection**:
