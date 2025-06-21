package main

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	tasks      map[uuid.UUID]*Task
	tasksMutex sync.RWMutex
)

func init() {
	tasks = make(map[uuid.UUID]*Task)
}

type Task struct {
	ID                 uuid.UUID     `json:"id"`
	Status             string        `json:"status"`
	CreationTime       time.Time     `json:"creationTime"`
	StartTime          *time.Time    `json:"startTime,omitempty"`
	CompletionTime     *time.Time    `json:"completionTime,omitempty"`
	ProcessingDuration time.Duration `json:"processingDuration,omitempty"`
	Result             string        `json:"result,omitempty"`
	Error              string        `json:"error,omitempty"`
}

func processTask(task *Task) {
	slog.Info("Task routine started", "task_id", task.ID.String())

	startTime := time.Now().UTC()
	taskDuration := time.Duration(rand.Int64N(3)+3) * time.Minute

	shouldFail := rand.IntN(100) < 20

	tasksMutex.Lock()
	task.StartTime = &startTime
	task.Status = "in_progress"
	tasksMutex.Unlock()

	timer := time.NewTimer(taskDuration)
	defer timer.Stop()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timer.C:
			slog.Info("Task timer expired", "task_id", task.ID.String())

			tasksMutex.Lock()
			if task.Status != "deleted" {
				now := time.Now().UTC()
				task.CompletionTime = &now
				task.ProcessingDuration = task.CompletionTime.Sub(*task.StartTime) / 1000000000

				if shouldFail {
					task.Status = "failed"
					task.Error = "Simulated internal processing error during execution."
					slog.Error("Task failed during processing", "task_id", task.ID.String(), "error", task.Error)
				} else {
					task.Status = "completed"
					task.Result = "Task completed successfully!"
					slog.Info("Task completed successfully", "task_id", task.ID.String())
				}
			} else {
				slog.Info("Task completed its duration, but was marked for deletion earlier", "task_id", task.ID.String())
			}
			tasksMutex.Unlock()
			return

		case <-ticker.C:
			tasksMutex.RLock()
			if task.Status == "deleted" {
				tasksMutex.RUnlock()
				slog.Info("Task routine stopping due to deletion request", "task_id", task.ID.String())
				tasksMutex.Lock()
				task.Status = "deleted_by_user"
				task.Error = "Task explicitly deleted by user request."
				tasksMutex.Unlock()
				return
			}
			tasksMutex.RUnlock()

			tasksMutex.Lock()
			task.ProcessingDuration = time.Since(*task.StartTime) / 1000000000
			tasksMutex.Unlock()
		}
	}
}

func postHandler(c *gin.Context) {
	newTask := &Task{
		ID:           uuid.New(),
		CreationTime: time.Now().UTC(),
		Status:       "pending",
	}

	tasksMutex.Lock()
	tasks[newTask.ID] = newTask
	tasksMutex.Unlock()

	go processTask(newTask)

	c.IndentedJSON(http.StatusCreated, newTask)
	slog.Info("POST /tasks: New task created", "task_id", newTask.ID.String())
}

func getHandler(c *gin.Context) {
	requestedIDStr := c.Param("id")
	requestedID, err := uuid.Parse(requestedIDStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		slog.Warn("GET /tasks/:id: Invalid task ID format requested", "id", requestedIDStr, "error", err)
		return
	}

	tasksMutex.RLock()
	task, exists := tasks[requestedID]
	tasksMutex.RUnlock()

	if !exists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		slog.Info("GET /tasks/:id: Task not found", "task_id", requestedIDStr)
		return
	}

	c.IndentedJSON(http.StatusOK, task)
	slog.Info("GET /tasks/:id: Task retrieved", "task_id", requestedIDStr, "status", task.Status)
}

func deleteHandler(c *gin.Context) {
	requestedIDStr := c.Param("id")
	requestedID, err := uuid.Parse(requestedIDStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		slog.Warn("DELETE /tasks/:id: Invalid task ID format requested", "id", requestedIDStr, "error", err)
		return
	}

	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	task, exists := tasks[requestedID]
	if !exists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		slog.Info("DELETE /tasks/:id: Task not found for deletion", "task_id", requestedIDStr)
		return
	}

	task.Status = "deleted"
	task.Error = "Task explicitly deleted by user request."

	c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Task '%s' marked for deletion. It will cease processing shortly.", requestedIDStr)})
	slog.Info("DELETE /tasks/:id: Task marked for deletion", "task_id", requestedIDStr, "current_status", task.Status)
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})))

	router := gin.Default()

	router.POST("/", postHandler)
	router.GET("/:id", getHandler)
	router.DELETE("/:id", deleteHandler)

	slog.Info("Starting HTTP server", "port", 8080)
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
