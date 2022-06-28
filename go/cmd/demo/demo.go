package main

import (
	"math/rand"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/oklog/ulid/v2"
)

type Model struct {
	ID        ulid.ULID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())

	uuid, err := ulid.New(ms, entropy)
	if err != nil {
		return err
	}
	m.ID = uuid
	return nil
}

type Product struct {
	Model
	Code  string
	Price uint
}

func main() {
	dsn := "host=localhost user=test password=test dbname=test port=5432 sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	db.Create(&Product{
		Code:  "D42",
		Price: 100,
	})
}
