package main

import (
	"log"
	"math/rand"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofrs/uuid"
	"github.com/oklog/ulid/v2"
)

type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())

	ulid, err := ulid.New(ms, entropy)
	if err != nil {
		return err
	}
	m.ID.UnmarshalBinary(ulid[:])
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

	// cleanup
	db.Unscoped().Where("1=1").Delete(&Product{})

	// create some records
	start := time.Now()
	n := 100
	for i := 0; i < n; i++ {
		db.Create(&Product{
			Code:  "D42",
			Price: 100,
		})
		time.Sleep(2 * time.Millisecond)
	}
	elapsed := time.Since(start)
	log.Printf("Inserting %d datasets took %s", n, elapsed)

	// sort by date and id
	sortedByDate := make([]Product, n)
	db.Order("created_at").Find(&sortedByDate)

	sortedById := make([]Product, n)
	db.Order("id").Find(&sortedById)

	// ensure the order is the same
	for i := 0; i < n; i++ {
		if sortedByDate[i].ID != sortedById[i].ID {
			log.Fatalf("They are not the same at %d", i)
		}
	}
}
