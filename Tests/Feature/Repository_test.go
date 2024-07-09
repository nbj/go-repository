package Feature

import (
	"errors"
	"github.com/google/uuid"
	"github.com/nbj/go-repository/Repository"
	"github.com/nbj/go-repository/Tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func Test_a_repository_can_be_instantiated(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	// Act
	repository := Repository.Of[Tests.TestCaseModel]()

	// Assert
	require.Equal(t, "*Repository.Repository[github.com/nbj/go-repository/Tests.TestCaseModel]", reflect.TypeOf(repository).String())
}

func Test_a_repository_can_get_a_collection_of_all_entries(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	entries := repository.All()

	// Assert
	require.Equal(t, 5, entries.Count())
}

func Test_a_repository_can_get_entries_with_relationships(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	entry := repository.First()

	// Assert
	require.Equal(t, 1, len(entry.TestCaseRelationModels))
	require.Equal(t, "*Tests.TestCaseRelationModel", reflect.TypeOf(entry.TestCaseRelationModels[0]).String())
}

func Test_a_repository_can_query_entries(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	entries := repository.GormQuery(func(query *gorm.DB) *gorm.DB {
		return query.
			Where("value IN ?", []string{"Value [2]", "Value [4]"})
	})

	// Assert
	require.Equal(t, 2, entries.Count())
	require.Equal(t, "Value [2]", entries.First().Value)
	require.Equal(t, "Value [4]", entries.Last().Value)
}

func Test_a_repository_can_query_first_matching_conditions(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	entryA := repository.First()
	entryB := repository.First(func(query *gorm.DB) *gorm.DB {
		return query.Where("value = ?", "Value [5]")
	})

	// Assert
	require.Equal(t, "Value [1]", entryA.Value)
	require.Equal(t, "Value [5]", entryB.Value)
}

func Test_a_repository_can_initiate_a_query_builder(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	builder := repository.Query()

	// Assert
	require.Equal(t, "*Repository.QueryBuilder[github.com/nbj/go-repository/Tests.TestCaseModel]", reflect.TypeOf(builder).String())
}

func Test_a_repository_can_create_new_entries(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment(true)

	repository := Repository.Of[Tests.TestCaseModel]()
	assert.Equal(t, 0, repository.All().Count())

	// Act
	repository.Create(Tests.TestCaseModel{
		Value: "Value [NEW]",
	})

	// Assert
	assert.Equal(t, 1, repository.All().Count())
	assert.Equal(t, "Value [NEW]", repository.All().First().Value)
}

func Test_a_repository_can_update_entries(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment(true)

	repository := Repository.Of[Tests.TestCaseModel]()
	assert.Equal(t, 0, repository.All().Count())

	repository.Create(Tests.TestCaseModel{
		Value: "Value [NEW]",
	})

	assert.Equal(t, 1, repository.All().Count())
	assert.Equal(t, "Value [NEW]", repository.All().First().Value)

	// Act
	model := repository.All().First()

	repository.Update(model.Id, map[string]any{
		"value": "Value [UPDATED]",
	})

	// Assert
	assert.Equal(t, 1, repository.All().Count())
	assert.Equal(t, "Value [UPDATED]", repository.All().First().Value)
}

func Test_queries_can_be_performed_as_a_transaction(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment(true)

	repository := Repository.Of[Tests.TestCaseModel]()
	assert.Equal(t, 0, repository.All().Count())

	// Act
	err := Repository.Transaction(func(config Repository.Config) error {
		transaction := Repository.Of[Tests.TestCaseModel](config)

		firstUuid, _ := uuid.NewV7()
		secondUuid, _ := uuid.NewV7()

		transaction.Create(Tests.TestCaseModel{
			Id:    firstUuid,
			Value: "Value [NEW]",
		})

		transaction.Create(Tests.TestCaseModel{
			Id:    secondUuid,
			Value: "Value [MORE-NEW]",
		})

		return nil
	})

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 2, repository.All().Count())
	assert.Equal(t, "Value [NEW]", repository.All().First().Value)
	assert.Equal(t, "Value [MORE-NEW]", repository.All().Last().Value)
}

func Test_if_a_transaction_does_not_return_nil_it_is_rolled_back(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment(true)

	repository := Repository.Of[Tests.TestCaseModel]()
	assert.Equal(t, 0, repository.All().Count())

	// Act
	err := Repository.Transaction(func(config Repository.Config) error {
		transaction := Repository.Of[Tests.TestCaseModel](config)

		firstUuid, _ := uuid.NewV7()
		secondUuid, _ := uuid.NewV7()

		transaction.Create(Tests.TestCaseModel{
			Id:    firstUuid,
			Value: "Value [NEW]",
		})

		transaction.Create(Tests.TestCaseModel{
			Id:    secondUuid,
			Value: "Value [MORE-NEW]",
		})

		return errors.New("this-transaction-will-not-be-commited")
	})

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, "this-transaction-will-not-be-commited", err.Error())
	assert.Equal(t, 0, repository.All().Count())
}
