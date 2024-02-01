package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

func initDatabase() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=ssactuarial dbname=newtest password=ssactuarial sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database")
	}

	if err := db.DB().Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&TaskRuns{})

	log.Println("Connected to the database successfully!")

	return db
}

// TaskRun model
type TaskRuns struct {
	gorm.Model
	StartTime       time.Time     `gorm:"not null" json:"start_time"`
	EndTime         time.Time     `gorm:"not null" json:"end_time"`
	Runtime         time.Duration `gorm:"not null" json:"runtime"`
	FilesAdded      pq.StringArray `gorm:"type:text[]" sql:"not null;column:files_added" json:"files_added"`
	FilesDeleted    pq.StringArray `gorm:"type:text[]" sql:"not null;column:files_deleted" json:"files_deleted"`
	MagicStringHits int           `gorm:"not null" json:"magic_string_hits"`
	Status          string        `gorm:"not null;size:20" json:"status"`
}
