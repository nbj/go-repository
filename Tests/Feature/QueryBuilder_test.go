package Feature

import (
	"github.com/nbj/go-repository/Repository"
	"github.com/nbj/go-repository/Tests"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_query_builder_can_get_a_collection_of_all_entries(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	collection := repository.Query().Get()

	// Assert
	assert.Equal(t, 5, collection.Count())
}

func Test_query_builder_can_get_a_collection_of_entries_using_where(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	collection := repository.Query().
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
	entry := repository.Query().First()

	// Assert
	assert.Equal(t, "Value [1]", entry.Value)
}

func Test_query_builder_can_get_the_first_entry_using_where(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	entry := repository.Query().
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
	repository.Query().
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
	entries := repository.Query().
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
	resultA := repository.Query().
		Where("value", "Value [4]").
		Exists()

	resultB := repository.Query().
		Where("value", "this-does-not-exist").
		Exists()

	// Assert
	assert.True(t, resultA)
	assert.False(t, resultB)
}

func Test_query_builder_set_relationships_to_query(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseRelationModel]()

	// Act
	result := repository.Query().
		With("TestCaseModel").
		Get()

	// Assert
	assert.Equal(t, "*Tests.TestCaseModel", reflect.TypeOf(result.First().TestCaseModel).String())
}

func Test_query_builder_can_paginate_all_results(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()
	repository := Repository.Of[Tests.TestCaseRelationModel]()

	// Act
	paginator := repository.Query().Paginate(1, 1, "tests")

	// Assert
	assert.Equal(t, 1, paginator.Page)
	assert.Equal(t, 1, paginator.PerPage)
	assert.Equal(t, 1, paginator.To)
	assert.Equal(t, 1, paginator.From)
	assert.Equal(t, 5, paginator.Total)
	assert.Equal(t, 1, paginator.Items.Count())
}

func Test_query_builder_can_paginate_queried_results(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()
	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	paginator := repository.Query().
		Where("value", "Value [1]").
		Paginate(1, 1, "tests")

	// Assert
	assert.Equal(t, 1, paginator.Page)
	assert.Equal(t, 1, paginator.PerPage)
	assert.Equal(t, 1, paginator.To)
	assert.Equal(t, 1, paginator.From)
	assert.Equal(t, 1, paginator.Total)
	assert.Equal(t, 1, paginator.Items.Count())
}

func Test_query_builder_can_skip_n_entries(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	collection := repository.Query().
		Skip(1).
		Get()

	// Assert
	assert.Equal(t, 4, collection.Count())
	assert.Equal(t, "Value [2]", collection.First().Value)
}

func Test_query_builder_can_take_n_entries(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	collection := repository.Query().
		Take(1).
		Get()

	// Assert
	assert.Equal(t, 1, collection.Count())
	assert.Equal(t, "Value [1]", collection.First().Value)
}

func Test_query_builder_can_skip_n_and_take_n_entries(t *testing.T) {
	// Arrange
	Tests.SetupEnvironment()

	repository := Repository.Of[Tests.TestCaseModel]()

	// Act
	collection := repository.Query().
		Skip(2).
		Take(2).
		Get()

	// Assert
	assert.Equal(t, 2, collection.Count())
	assert.Equal(t, "Value [3]", collection.First().Value)
	assert.Equal(t, "Value [4]", collection.Last().Value)
}
