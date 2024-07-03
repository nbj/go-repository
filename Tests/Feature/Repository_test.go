package Feature

import (
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
	entries := repository.Query(func(query *gorm.DB) *gorm.DB {
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
	builder := repository.Builder()

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
