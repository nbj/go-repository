package Feature

import (
	"github.com/nbj/go-repository/Repository"
	"github.com/nbj/go-repository/Tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_query_builder_can_get_a_collection_of_all_entries(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	collection := repository.Builder().Get()

	// Assert
	assert.Equal(t, 5, collection.Count())
}

func Test_query_builder_can_get_a_collection_of_entries_using_where(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	collection := repository.Builder().
		Where("value", "Value [1]").
		OrWhere("value", "Value [3]").
		Get()

	// Assert
	assert.Equal(t, 2, collection.Count())
	assert.Equal(t, "Value [1]", collection.First().Value)
	assert.Equal(t, "Value [3]", collection.Last().Value)
}

func Test_query_builder_can_get_the_first_entry(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	entry := repository.Builder().First()

	// Assert
	assert.Equal(t, "Value [1]", entry.Value)
}

func Test_query_builder_can_get_the_first_entry_using_where(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	entry := repository.Builder().
		Where("value", "Value [3]").
		First()

	// Assert
	assert.Equal(t, "Value [3]", entry.Value)
}

func Test_query_builder_can_get_the_first_entry_or_fail(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	defer func() {
		if recover() != nil {
			// Makes sure t.Fail() i never reached if a panic occurs
		}
	}()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	repository.Builder().
		Where("value", "this-does-not-exist").
		FirstOrFail()

	// Assert
	t.Fail()
}

func Test_query_builder_can_order_entries(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	entries := repository.Builder().
		OrderBy("value", "desc").
		Get()

	// Assert
	valueA := entries.Shift()
	valueB := entries.Shift()
	valueC := entries.Shift()
	valueD := entries.Shift()
	valueE := entries.Shift()

	assert.Equal(t, "Value [5]", valueA.Value)
	assert.Equal(t, "Value [4]", valueB.Value)
	assert.Equal(t, "Value [3]", valueC.Value)
	assert.Equal(t, "Value [2]", valueD.Value)
	assert.Equal(t, "Value [1]", valueE.Value)
}

func Test_query_builder_can_check_if_entry_exists(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	resultA := repository.Builder().
		Where("value", "Value [4]").
		Exists()

	resultB := repository.Builder().
		Where("value", "this-does-not-exist").
		Exists()

	// Assert
	assert.True(t, resultA)
	assert.False(t, resultB)
}
