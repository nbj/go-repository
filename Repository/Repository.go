package Repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/nbj/go-collections/Collection"
	"github.com/nbj/go-support/Support"
	"gorm.io/gorm"
	"log"
)

var defaultConfiguration *Config

type WithRelationships interface {
	With() []string
}

type Repository[T any] struct {
	connection  *gorm.DB
	model       *T
	query       *gorm.DB
	latestError error
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

// GetLatestError
// Returns the latest error recorded
func (repository *Repository[T]) GetLatestError() error {
	return repository.latestError
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

// Query
// Shorthand for starting a new query builder
func (repository *Repository[T]) Query() *QueryBuilder[T] {
	var builder QueryBuilder[T]

	builder.query = repository.connection

	return &builder
}

// All
// Gets a collection of all entries in the repository.
// Returns nil if query fails
func (repository *Repository[T]) All() *Collection.Collection[T] {
	var entries []T

	query := repository.connection
	query = repository.applyRelationships(query)

	if result := query.Find(&entries); result.Error != nil {
		repository.latestError = result.Error
		panic("Repository[All]: " + result.Error.Error())
	}

	return Collection.Collect(entries)
}

// Create
// Creates a new database entry
func (repository *Repository[T]) Create(value T) *T {
	if result := repository.connection.Create(&value); result.Error != nil {
		repository.latestError = result.Error
		panic("Repository[Create]: " + result.Error.Error())
	}

	return &value
}

// Update
// Updates an existing database entry with values from map.
// Returns true if successful, false if not
func (repository *Repository[T]) Update(id uuid.UUID, values any) error {
	query := repository.connection.
		Model(repository.model).
		Where("id = ?", id).
		Updates(values)

	if query.Error != nil {
		repository.latestError = query.Error

		return repository.latestError
	}

	if query.RowsAffected == 0 {
		repository.latestError = errors.New("Repository[Update]: No rows were affected")

		return repository.latestError
	}

	return nil
}

// GormQuery
// Takes a closure containing a gorm query, executes it and
// returns the result as a collection of entries. Returns nil if
// query fails
func (repository *Repository[T]) GormQuery(closure func(query *gorm.DB) *gorm.DB) *Collection.Collection[T] {
	var entries []T

	query := closure(repository.connection)
	query = repository.applyRelationships(query)

	if result := query.Find(&entries); result.Error != nil {
		repository.latestError = result.Error
		panic("Repository[GormQuery]: " + result.Error.Error())
	}

	return Collection.Collect(entries)
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
			repository.latestError = result.Error
			panic("Repository[First]: " + result.Error.Error())
		}

		return &entry
	}

	// Apply all closures to the query
	for _, closure := range closures {
		query = closure(query)
	}

	if result := query.First(&entry); result.Error != nil {
		repository.latestError = result.Error
		panic("Repository[First]: " + result.Error.Error())
	}

	return &entry
}

// Transaction
// Performs a closure as a database transaction
func Transaction(closure func(transactionConfig Config) error) error {
	// We start by creating the transaction
	transaction := defaultConfiguration.DatabaseConnection.Begin()

	// Create a transaction config to use for repositories inside
	// the closure housing the transaction
	transactionConfig := Config{
		DatabaseConnection: transaction,
	}

	// Pass config to closure and execute query. If any errors
	// are produced, we roll back the transaction
	if err := closure(transactionConfig); err != nil {
		transaction.Rollback()

		return err
	}

	// If everything went well we commit the transaction
	transaction.Commit()

	return nil
}
