package Tests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nbj/go-repository/Repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

func SetupEnvironment() {
	connection := getSqliteDatabaseConnection()

	seedTestData(connection)

	Repository.SetDefaultConfig(&Repository.Config{DatabaseConnection: connection})
}

func getSqliteDatabaseConnection() *gorm.DB {
	var connection *gorm.DB
	var err error

	if connection, err = gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}); nil != err {
		panic("failed to connect database: " + err.Error())
	}

	modelsToMigrate := []any{
		TestCaseModel{},
	}

	if err = connection.AutoMigrate(modelsToMigrate...); nil != err {
		panic("failed to auto migrate database: " + err.Error())
	}

	return connection
}

func seedTestData(connection *gorm.DB) {
	numberOfEntries := 5

	for entry := 1; entry <= numberOfEntries; entry++ {
		uniqueIdentifier, _ := uuid.NewV7()

		instance := TestCaseModel{
			Id:    uniqueIdentifier,
			Value: fmt.Sprintf("Value [%d]", entry),
		}

		connection.Create(&instance)
	}
}

type TestCaseModel struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;uniqueIndex"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at" gorm:"index;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}
