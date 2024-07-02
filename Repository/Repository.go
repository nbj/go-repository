package Repository

import (
	Nbj "github.com/nbj/go-collections"
	"github.com/nbj/go-support/Support"
	"gorm.io/gorm"
	"log"
)

var defaultConfiguration *Config

type WithRelationships interface {
	With() []string
}

type Repository[T any] struct {
	connection *gorm.DB
	model      *T
	query      *gorm.DB
}

type Config struct {
	DatabaseConnection *gorm.DB
}

// Of
// Named constructor for creating instances of a repository
func Of[T any](config ...Config) *Repository[T] {
	// Check if it is impossible to instantiate a repository instance and bail
	if len(config) == 0 && defaultConfiguration == nil {
		log.Fatal("No configuration passed to Repository and no default configuration available.")
	}

	// Create the repository instance
	var repository Repository[T]
	repository.model = new(T)

	// Assign the appropriate configuration
	repository.applyConfiguration(defaultConfiguration)

	if len(config) > 0 {
		repository.applyConfiguration(&config[0])
	}

	// Return the newly created repository
	return &repository
}

// SetDefaultConfig
// Sets the default configuration repositories will be instantiated with
func SetDefaultConfig(config *Config) {
	defaultConfiguration = config
}

// applyConfiguration
// Assigns a configuration to the repository instance
func (repository *Repository[T]) applyConfiguration(config *Config) {
	repository.connection = config.DatabaseConnection
}

// applyRelationships
// Applies any relationships set with the With() function on models
func (repository *Repository[T]) applyRelationships(query *gorm.DB) *gorm.DB {
	if Support.Implements[WithRelationships](repository.model) {
		relations := Support.Cast[WithRelationships](repository.model).With()

		for _, relation := range relations {
			query = query.Preload(relation)
		}
	}

	return query
}

// Builder
// Shorthand for starting a new query builder
func (repository *Repository[T]) Builder() *QueryBuilder[T] {
	var builder QueryBuilder[T]

	builder.query = repository.connection

	return &builder
}

// All
// Gets a collection of all entries in the repository.
// Returns nil if query fails
func (repository *Repository[T]) All() *Nbj.Collection[T] {
	var entries []T

	query := repository.connection
	query = repository.applyRelationships(query)

	if result := query.Find(&entries); result.Error != nil {
		return nil
	}

	return Nbj.Collect(entries)
}

// Query
// Takes a closure containing a gorm query, executes it and
// returns the result as a collection of entries. Returns nil if
// query fails
func (repository *Repository[T]) Query(closure func(query *gorm.DB) *gorm.DB) *Nbj.Collection[T] {
	var entries []T

	query := closure(repository.connection)
	query = repository.applyRelationships(query)

	if result := query.Find(&entries); result.Error != nil {
		return nil
	}

	return Nbj.Collect(entries)
}

// First
// Gets the first database entry that matches the queries passed
func (repository *Repository[T]) First(closures ...func(query *gorm.DB) *gorm.DB) *T {
	var entry T

	query := repository.connection
	query = repository.applyRelationships(query)

	// If no closures are passed to the method
	if len(closures) == 0 {
		if result := query.First(&entry); result.Error != nil {
			return nil
		}

		return &entry
	}

	// Apply all closures to the query
	for _, closure := range closures {
		query = closure(query)
	}

	if result := query.First(&entry); result.Error != nil {
		return nil
	}

	return &entry
}
