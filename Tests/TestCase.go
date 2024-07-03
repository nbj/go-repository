package Tests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nbj/go-repository/Repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupEnvironment(noSeed ...bool) {
	connection := getSqliteDatabaseConnection()

	if len(noSeed) == 0 {
		seedTestData(connection)
	}

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
		TestCaseRelationModel{},
	}

	if err = connection.AutoMigrate(modelsToMigrate...); nil != err {
		panic("failed to auto migrate database: " + err.Error())
	}

	return connection
}

func seedTestData(connection *gorm.DB) {
	numberOfEntries := 5
	var instances []*TestCaseModel

	for entry := 1; entry <= numberOfEntries; entry++ {
		instance := makeTestCaseModel(fmt.Sprintf("Value [%d]", entry))
		connection.Create(&instance)
		instances = append(instances, &instance)
	}

	for index, instance := range instances {
		relationInstance := makeTestCaseRelationModel(instance.Id, fmt.Sprintf("Relation Value [%d]", index))
		connection.Create(&relationInstance)
	}
}

func makeTestCaseModel(value string) TestCaseModel {
	uniqueIdentifier, _ := uuid.NewV7()

	instance := TestCaseModel{
		Id:    uniqueIdentifier,
		Value: value,
	}

	return instance
}

func makeTestCaseRelationModel(testCaseModelId uuid.UUID, value string) TestCaseRelationModel {
	uniqueIdentifier, _ := uuid.NewV7()

	instance := TestCaseRelationModel{
		Id:              uniqueIdentifier,
		TestCaseModelId: testCaseModelId,
		Value:           value,
	}

	return instance
}
