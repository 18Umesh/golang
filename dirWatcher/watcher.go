package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jinzhu/gorm"
	"io"
	"strings"
)

type Watcher struct {
	config     *Config
	db         *gorm.DB
	watcher    *fsnotify.Watcher
	stopSignal chan struct{}
}

func NewWatcher(config *Config, db *gorm.DB) *Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	return &Watcher{
		config:     config,
		db:         db,
		watcher:    watcher,
		stopSignal: make(chan struct{}),
	}
}

func (w *Watcher) Start() {
	go w.processEvents()

	err := filepath.Walk(w.config.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			w.watchFile(path)
			w.processFile(path)
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking through the directory: %v", err)
	}
}

func (w *Watcher) Stop() {
	close(w.stopSignal)
	w.watcher.Close()
}

func (w *Watcher) processEvents() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			w.handleEvent(event)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		case <-w.stopSignal:
			return
		}
	}
}

func (w *Watcher) watchFile(path string) {
    err := w.watcher.Add(path)
    if err != nil {
        log.Printf("Error adding %s to watcher: %v", path, err)
    } else {
        log.Printf("Watching file: %s", path)
    }
}

func (w *Watcher) processFileCreation(filePath string) {
    log.Printf("Processing file creation: %s", filePath)

    // Check if there is an ongoing task run, and if not, create a new one
    var taskRun TaskRuns
    w.db.Order("start_time desc").First(&taskRun, "status = 'in progress'")
    if taskRun.ID == 0 {
        taskRun = TaskRuns{
            StartTime: time.Now(),
            EndTime:   time.Time{}, // You can set it to the zero value for now
            Runtime:   0,            // You can calculate this once the task is completed
            FilesAdded:   []string{},   // Initialize as empty; update when files are added
            FilesDeleted: []string{},   // Initialize as empty; update when files are deleted
            MagicStringHits: 0,
            Status: "in progress",
        }
        w.db.Create(&taskRun)
    }

    // Simulate updating the database (replace this with actual database logic)
    taskRun.FilesAdded = append(taskRun.FilesAdded, filePath)
    w.db.Save(&taskRun)

    // Update task status to "success" if everything is processed successfully
    taskRun.Status = "success"
    w.db.Save(&taskRun)

    log.Printf("File creation processed: %s", filePath)
}


func (w *Watcher) processFile(filePath string) {
    log.Printf("Processing file: %s", filePath)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Error getting file info for %s: %v", filePath, err)
		return
	}

	if fileInfo.IsDir() {
		log.Printf("Skipping directory: %s", filePath)
		return
	}

    // Open the file
    file, err := os.Open(filePath)
    if err != nil {
        log.Printf("Error opening file %s: %v", filePath, err)
        return
    }
    defer file.Close()

    // Read the content of the file
    content, err := io.ReadAll(file)
    if err != nil {
        log.Printf("Error reading file %s: %v", filePath, err)
        return
    }

    // Log content of the file
    log.Printf("File content: %s", content)

    // Count occurrences of the magic string (case-insensitive)
    magicString := strings.ToLower(w.config.MagicString)
    contentStr := strings.ToLower(string(content))
    count := strings.Count(contentStr, magicString)

    // Initialize TaskRuns struct with start time and status
    taskRun := TaskRuns{
		StartTime:       time.Now(),
		EndTime:         time.Time{}, // You can set it to the zero value for now
		Runtime:         0,            // You can calculate this once the task is completed
		FilesAdded:      []string{},   // Initialize as empty; update when files are added
		FilesDeleted:    []string{},   // Initialize as empty; update when files are deleted
		MagicStringHits: count,
		Status:          "in progress",
	}
	
    // Begin a transaction
    tx := w.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Save results to the database in a transaction
    if err := tx.Create(&taskRun).Error; err != nil {
        log.Printf("Error creating TaskRun: %v", err)
        tx.Rollback()
        return
    }

    // Update task status to "success" if everything is processed successfully
    taskRun.Status = "success"
    if err := tx.Save(&taskRun).Error; err != nil {
        log.Printf("Error updating TaskRun: %v", err)
        tx.Rollback()
        return
    }

    // Commit the transaction
    tx.Commit()

    // Now, taskRun should have updated values
    log.Printf("TaskRun: %+v", taskRun)

}

func (w *Watcher) processFileDeletion(filePath string) {
    log.Printf("Processing file deletion: %s", filePath)

    // Check if there is an ongoing task run, and if not, create a new one
    var taskRun TaskRuns
    w.db.Order("start_time desc").First(&taskRun, "status = 'in progress'")
    if taskRun.ID == 0 {
        taskRun = TaskRuns{
            StartTime: time.Now(),
            EndTime:   time.Time{}, // You can set it to the zero value for now
            Runtime:   0,            // You can calculate this once the task is completed
            FilesAdded:   []string{},   // Initialize as empty; update when files are added
            FilesDeleted: []string{},   // Initialize as empty; update when files are deleted
            MagicStringHits: 0,
            Status: "in progress",
        }
        w.db.Create(&taskRun)
    }

    // Simulate updating the database (replace this with actual database logic)
    taskRun.FilesDeleted = append(taskRun.FilesDeleted, filePath)
    w.db.Save(&taskRun)

    // Update task status to "success" if everything is processed successfully
    taskRun.Status = "success"
    w.db.Save(&taskRun)

    log.Printf("File deletion processed: %s", filePath)
}

func (w *Watcher) handleEvent(event fsnotify.Event) {
    switch {
    case event.Op&fsnotify.Write == fsnotify.Write:
        // File modified
        log.Printf("File modified: %s", event.Name)
        w.processFile(event.Name)

    case event.Op&fsnotify.Create == fsnotify.Create:
        // New file added
        log.Printf("New file added: %s", event.Name)
        w.watchFile(event.Name)
        w.processFileCreation(event.Name)

    case event.Op&fsnotify.Remove == fsnotify.Remove:
        // File deleted
        log.Printf("File deleted: %s", event.Name)
        w.processFileDeletion(event.Name)
    }
}
