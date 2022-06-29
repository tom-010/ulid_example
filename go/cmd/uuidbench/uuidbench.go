package main

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofrs/uuid"
)

type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	m.ID = uuid
	return nil
}

type Employee struct {
	Model
	FirstName        string
	LastName         string
	Age              uint
	EmploymentStatus string
	Street           string
	HouseNumber      string
	City             string
	PostalCode       string
}

type Patient struct {
	Model
	FirstName  string
	LastName   string
	EmployeeID uuid.UUID
	Employee   Employee
}

type Session struct {
	gorm.Model
	Notes     string
	PatientID uuid.UUID
	Patient   Patient
}

func createOneSample(db *gorm.DB) {
	employee := Employee{
		FirstName:        "firstname",
		LastName:         "lastname",
		Age:              12,
		EmploymentStatus: "employment status asdf asdf asdf asdf asdf ewgf dsaf asdfjsadlkfj asdjf akdsjf ajdsf lasjdf ajsdfl jasdlkfj asldjf lasdjf lajdsf lkjdsaf ölkjsdflk jasldfj aldsjflkajfdslasdjf ljadslf jlasdjf ljadsf kja",
		Street:           "street",
		HouseNumber:      "123",
		City:             "city",
		PostalCode:       "1234",
	}
	db.Create(&employee)
	for i := 0; i < 3; i++ {
		patient := Patient{
			FirstName:  "first-name",
			LastName:   "lastname",
			EmployeeID: employee.ID,
		}
		db.Create(&patient)
		for j := 0; j < 4; j++ {
			db.Create(&Session{
				Notes:     "aslkdf lajsfdöljsadflk jaslfdkj alskdfj lakdsjf lkajfds ljsadfl jaldskfj alsdjf alsjfd ksajfdl kjdsafl jasldfj alkdsfj lksajfd lkajfds lkjdsaflk jasdflkj lkdsafj adsjf lasdfasd lfdsaf jaldsjf kjdsafölk jsadölfj asljf ajf jfdsa ljdsafl jsadflj asdjf aslkdjf lsajfd l",
				PatientID: patient.ID,
			})
		}
	}
}

func main() {
	dsn := "host=localhost user=test password=test dbname=test port=5432 sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&Employee{}, &Patient{}, &Session{})
	if err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	// cleanup
	db.Unscoped().Where("1=1").Delete(&Session{})
	db.Unscoped().Where("1=1").Delete(&Patient{})
	db.Unscoped().Where("1=1").Delete(&Employee{})

	// create some records
	start := time.Now()
	n := 10000
	for i := 0; i < n; i++ {
		createOneSample(db)
	}
	elapsed := time.Since(start)
	log.Printf("Running %d times took %s", n, elapsed)
}
