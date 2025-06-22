# Task Scheduler

- [English](#English)
    1. [Introduction](#introduction)
    2. [Requirement](#requirement)
    3. [How to Run](#how-to-run)
        - [Docker](#docker)
        - [Bare-metal](#bare-metal)
    4. [API Endpoints](#api-endpoints)
    5. [Key Design Decisions](#key-design-decisions)
- [Русский](#Русский)
    1. [Введение](#введение)
    2. [Требования](#требования)
    3. [Как запустить](#как-запустить)
    4. [API Эндпоинты](#api-эндпоинты)
    5. [Ключевые проектные решения](#ключевые-проектные-решения)

---

## English

### Introduction

This project presents a robust and efficient in-memory task management API built with Go (Golang). Designed to handle long-running, asynchronous operations, this service provides a simple yet powerful interface for submitting tasks, tracking their lifecycle, and retrieving their results. Adhering to the principle of "no third-party services," all task data is meticulously managed and stored directly within the service's memory, showcasing a strong command of concurrent programming and resource management in Go.

The core functionality revolves around simulating I/O-bound tasks that typically complete within 3 to 5 minutes. The API allows for real-time status updates, calculation of processing durations, and the graceful handling of both successful completions and simulated failures. Leveraging Go's powerful concurrency primitives and the Gin web framework, the service ensures responsiveness and scalability while maintaining a minimal operational footprint, including its deployment within a highly optimized scratch Docker image.

This solution serves as a foundational component for systems requiring reliable background task orchestration without external dependencies, demonstrating a clean architecture suitable for future expansion and integration.

### Requirement

* [Go](https://go.dev/dl/) 1.24+ (for building and development or running on bare-metal)
* [Docker](https://docs.docker.com/get-docker/) (for running the containerized application)
* Make (to start the container easier)

### How To Run

##### 1. Clone the repository:
```bash
    git clone https://github.com/Eassen496/Task-Scheduler
    cd [workmate-go]
```

#### Docker

##### 2. Build and Run Containers using `Makefile`:
Use the `make` command to orchestrate Docker Compose:
```bash
    make up
```
This command will build the Docker image (if not already built) and start the container in detached mode.

##### 3. View Logs (optional):
To view the service logs, use:
```bash
    make logs
```
##### 4. Stop and Clean Up (optional):
To stop and remove containers, networks, images, and volumes associated with the project:
```bash
    make clean
```
The server will be available at `http://localhost:8080`.

#### Bare-metal

##### 2. Build the application:

To build the application:
```bash
    cd ./app/src
    go mod download
    go build . -o taskScheduler
```
This commands will build the application.



##### 3. Start the app:
Launch the application:
```bash
    ./taskScheduler
```
### API Endpoints

#### 1. Create a New Task

* **Method:** `POST`
* **Path:** `/`
* **Description:** Initiates a new simulated long-running task. The task will immediately begin processing in the background.
* **Example Request (using `curl`):**
    ```bash
    curl -X POST http://localhost:8080/
    ```
* **Example Response (HTTP 201 Created):**
    ```json
  {
      "id": "d4ed1190-f5dc-4449-ae6b-599e18952ddd",
      "status": "pending",
      "creationTime": "2025-06-21T16:07:53.323050002Z"
    }
    ```

#### 2. Get Task Status and Details

* **Method:** `GET`
* **Path:** `/:id`
* **Description:** Retrieves the current status, creation time, processing duration, and results (if completed) or error (if failed/deleted) of a specific task.
* **Example Request (using `curl`):**
    ```bash
    curl http://localhost:8080/d4ed1190-f5dc-4449-ae6b-599e18952ddd
    ```
* **Example Response (HTTP 200 OK - In Progress):**
    ```json
    {
        "id": "d4ed1190-f5dc-4449-ae6b-599e18952ddd",
        "status": "in_progress",
        "creationTime": "2025-06-21T04:30:00Z",
        "startTime": "2025-06-21T04:30:00Z",
        "processingDuration": 110
    }
    ```
* **Example Response (HTTP 200 OK - Completed):**
    ```json
    {
      "id": "d4ed1190-f5dc-4449-ae6b-599e18952ddd",
      "status": "completed",
      "creationTime": "2025-06-21T16:07:53.323050002Z",
      "startTime": "2025-06-21T16:07:53.323930335Z",
      "completionTime": "2025-06-21T16:10:53.33098996Z",
      "processingDuration": 180,
      "result": "Task completed successfully!"
    }                         
    ```
* **Example Response (HTTP 404 Not Found):**
    ```json
    {
        "message": "Task not found"
    }
    ```

#### 3. Delete a Task

* **Method:** `DELETE`
* **Path:** `/:id`
* **Description:** Marks an existing task for deletion. The background processing for the task will be signaled to terminate gracefully.
* **Example Request (using `curl`):**
    ```bash
    curl -X DELETE http://localhost:8080/d4ed1190-f5dc-4449-ae6b-599e18952ddd
    ```
* **Example Response (HTTP 200 OK):**
    ```json
    {
        "message": "Task 'd4ed1190-f5dc-4449-ae6b-599e18952ddd' marked for deletion. It will cease processing shortly."
    }
    ```

### Key Design Decisions

* **Concurrency with Goroutines and Mutexes:** Long-running tasks are executed in separate goroutines to prevent blocking the main HTTP server. A `sync.RWMutex` is used to ensure safe concurrent access to the in-memory task map.
* **In-Memory Storage:** All task data is stored in a `map[uuid.UUID]*Task` in the application's memory, adhering to the requirement of no external databases or queues.
* **Simulated Task Logic:** Tasks simulate their duration using `time.Sleep` and randomly simulate success or failure, updating their status accordingly.
* **Graceful Deletion:** Tasks can be marked for deletion via the API, and their respective goroutines are designed to detect this signal and terminate gracefully.
* **Structured Logging (`log/slog`):** Utilizes Go's modern structured logging for clear, machine-readable output, aiding in monitoring and debugging.
* **Gin Framework:** Used for simplified HTTP routing and request/response handling

---

## Русский

### Введение

Данный проект представляет собой надежный и эффективный API для управления задачами, работающий в оперативной памяти**, построенный на Go (Golang). Разработанный для обработки долгосрочных асинхронных операций, этот сервис предоставляет простой, но мощный интерфейс для отправки задач, отслеживания их жизненного цикла и получения результатов. Придерживаясь принципа "без сторонних сервисов", все данные задач тщательно управляются и хранятся непосредственно в памяти сервиса, демонстрируя уверенное владение параллельным программированием и управлением ресурсами в Go.

Основная функциональность проекта сосредоточена на симуляции операций с интенсивным вводом-выводом, которые обычно выполняются от 3 до 5 минут. API позволяет получать обновления статуса в реальном времени, рассчитывать продолжительность обработки и корректно обрабатывать как успешные завершения, так и симулированные сбои. Используя мощные примитивы конкурентности Go и веб-фреймворк Gin, сервис обеспечивает отзывчивость и масштабируемость при минимальном операционном следе, включая его развертывание в высокооптимизированном Docker-образе на основе scratch.

Это решение служит основополагающим компонентом для систем, требующих надежной оркестровки фоновых задач без внешних зависимостей, демонстрируя чистую архитектуру, подходящую для будущего расширения и интеграции.

### Требования

* [Go](https://go.dev/dl/) 1.24+ (для сборки и разработки или запуска на "голом железе")
* [Docker](https://docs.docker.com/get-docker/) (для запуска контейнерного приложения)
* Make (для более простого запуска контейнера)

### Как запустить

#### 1. Клонировать репозиторий:
```bash
    git clone https://github.com/Eassen496/workmate-go
    cd [workmate-go]
```

#### Docker

##### 2. Собрать и запустить контейнеры с использованием `Makefile`:
Используйте команду `make` для управления Docker Compose:
```bash
    make up
```
Эта команда соберет образ Docker (если он еще не собран) и запустит контейнер в фоновом режиме.

##### 3. Просмотр логов (необязательно):
Для просмотра логов сервиса используйте:
```bash
    make logs
```

##### 4. Остановка и очистка (необязательно):
Для остановки и удаления контейнеров, сетей, образов и томов, связанных с проектом:
```bash
    make clean
```
Сервер будет доступен по адресу `http://localhost:8080`.

#### На "голом железе"

##### 2. Собрать приложение:
Для сборки приложения:
```bash
    cd ./app/src
    go mod download
    go build . -o taskScheduler
```
Эти команды соберут приложение.

##### 3. Запустить приложение:
Запустите приложение:
```bash
    ./taskScheduler
```

### API Эндпоинты

#### 1. Создать новую задачу

* **Метод:** `POST`
* **Путь:** `/`
* **Описание:** Инициирует новую симулированную долгосрочную задачу. Задача немедленно начнет обработку в фоновом режиме.
* **Пример запроса (с использованием `curl`):**
    ```bash
    curl -X POST http://localhost:8080/
    ```
* **Пример ответа (HTTP 201 Created):**
    ```json  
    {
      "id": "d4ed1190-f5dc-4449-ae6b-599e18952ddd",
      "status": "pending",
      "creationTime": "2025-06-21T16:07:53.323050002Z"
    }
    ```

#### 2. Получить статус и детали задачи

* **Метод:** `GET`
* **Путь:** `/:id`
* **Описание:** Получает текущий статус, время создания, продолжительность обработки и результаты (если завершена) или ошибку (если не удалось/удалена) конкретной задачи.
* **Пример запроса (с использованием `curl`):**
    ```bash
    curl http://localhost:8080/d4ed1190-f5dc-4449-ae6b-599e18952ddd
    ```
* **Пример ответа (HTTP 200 OK - В процессе):**
    ```json
    {
        "id": "d4ed1190-f5dc-4449-ae6b-599e18952ddd",
        "status": "in_progress",
        "creationTime": "2025-06-21T04:30:00Z",
        "startTime": "2025-06-21T04:30:00Z",
        "processingDuration": 110
    }
    ```
* **Пример ответа (HTTP 200 OK - Завершена):**
    ```json
    {
      "id": "d4ed1190-f5dc-4449-ae6b-599e18952ddd",
      "status": "completed",
      "creationTime": "2025-06-21T16:07:53.323050002Z",
      "startTime": "2025-06-21T16:07:53.323930335Z",
      "completionTime": "2025-06-21T16:10:53.33098996Z",
      "processingDuration": 180,
      "result": "Task completed successfully!"
    }    
    ```
* **Пример ответа (HTTP 404 Not Found):**
    ```json
    {
        "message": "Task not found"
    }
    ```

#### 3. Удалить задачу

* **Метод:** `DELETE`
* **Путь:** `/:id`
* **Описание:** Помечает существующую задачу для удаления. Фоновая обработка задачи будет сигнализирована о необходимости корректного завершения.
* **Пример запроса (с использованием `curl`):**
    ```bash
    curl -X DELETE http://localhost:8080/d4ed1190-f5dc-4449-ae6b-599e18952ddd
    ```
* **Пример ответа (HTTP 200 OK):**
    ```json
    {
        "message": "Task 'd4ed1190-f5dc-4449-ae6b-599e18952ddd' marked for deletion. It will cease processing shortly."
    }
    ```

### Ключевые проектные решения

* **Конкурентность с использованием goroutine и mutex:** Долгосрочные задачи выполняются в отдельных goroutine, чтобы предотвратить блокировку основного HTTP-сервера. Для обеспечения безопасного одновременного доступа к карте задач, находящейся в памяти, используется `sync.RWMutex`.
* **Хранение данных в памяти:** Все данные задач хранятся в `map[uuid.UUID]*Task` в памяти приложения, что соответствует требованию об отсутствии внешних баз данных или очередей.
* **Логика симуляции задач:** Задачи симулируют свою продолжительность с помощью `time.Sleep` и случайным образом симулируют успех или неудачу, соответствующим образом обновляя свой статус.
* **Корректное удаление:** Задачи могут быть помечены для удаления через API, и соответствующие goroutine предназначены для обнаружения этого сигнала и корректного завершения работы.
* **Структурированное логирование (`log/slog`):** Используется современное структурированное логирование Go для четкого, машиночитаемого вывода, что облегчает мониторинг и отладку.
* **Фреймворк Gin:** Используется для упрощения маршрутизации HTTP и обработки запросов/ответов.
