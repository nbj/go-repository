package Repository

import (
	"fmt"
	"github.com/nbj/go-collections/Collection"
	"github.com/nbj/go-support/Support"
	"gorm.io/gorm"
)

type QueryBuilder[T any] struct {
	query  *gorm.DB
	model  *T
	orders []string
}

func (builder *QueryBuilder[T]) With(query string, args ...any) *QueryBuilder[T] {
	builder.query = builder.query.Preload(query, args)

	return builder
}

func (builder *QueryBuilder[T]) Where(query any, args ...any) *QueryBuilder[T] {
	builder.query = builder.query.Where(query, args...)

	return builder
}

func (builder *QueryBuilder[T]) OrWhere(query any, args ...any) *QueryBuilder[T] {
	builder.query = builder.query.Or(query, args...)

	return builder
}

func (builder *QueryBuilder[T]) OrderBy(column string, direction string) *QueryBuilder[T] {
	builder.query = builder.query.Order(fmt.Sprintf("%s %s", column, direction))

	return builder
}

// Exists
// Checks if the query find any results
func (builder *QueryBuilder[T]) Exists() bool {
	var entries []T
	var result *gorm.DB

	if result = builder.query.First(&entries); result.Error != nil {
		return false
	}

	return result.RowsAffected != 0
}

// Get
// Executes the query and get a collection containing all results
func (builder *QueryBuilder[T]) Get() *Collection.Collection[T] {
	var entries []T

	builder.applyRelationships()

	if result := builder.query.Find(&entries); result.Error != nil {
		return nil
	}

	return Collection.Collect(entries)
}

// First
// Executes the query and fetches the first result
func (builder *QueryBuilder[T]) First() *T {
	var entry T

	builder.applyRelationships()

	if result := builder.query.First(&entry); result.Error != nil || result.RowsAffected == 0 {
		return nil
	}

	return &entry
}

// FirstOrFail
// Executes the query and fetches the first result or dies trying
func (builder *QueryBuilder[T]) FirstOrFail() *T {
	entry := builder.First()

	if entry == nil {
		panic("FirstOrFail: entry not found")
	}

	return entry
}

// applyRelationships
// Applies any relationships set with the With() function on models
func (builder *QueryBuilder[T]) applyRelationships() {
	if Support.Implements[WithRelationships](builder.model) {
		relationships := Support.Cast[WithRelationships](builder.model).With()

		for _, relationship := range relationships {
			builder.query = builder.query.Preload(relationship)
		}
	}
}
