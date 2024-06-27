package Repository

import (
	Nbj "github.com/nbj/go-collections"
	"gorm.io/gorm"
	"log"
)

var defaultConfiguration *Config

type Repository[T any] struct {
	connection *gorm.DB
	model      *T
}

type Config struct {
	DatabaseConnection *gorm.DB
}

func Of[T any](config ...Config) *Repository[T] {
	// Check if it is impossible to instantiate a repository instance and bail
	if len(config) == 0 && defaultConfiguration == nil {
		log.Fatal("No configuration passed to Repository and no default configuration available.")
	}

	// Create the repository instance
	var repository Repository[T]

	// Assign the appropriate configuration
	repository.applyConfiguration(defaultConfiguration)

	if len(config) > 0 {
		repository.applyConfiguration(&config[0])
	}

	// Return the newly created repository
	return &repository
}

// SetDefaultConfig Sets the default configuration repositories will be instantiated with
func SetDefaultConfig(config *Config) {
	defaultConfiguration = config
}

// applyConfiguration Assigns a configuration to the repository instance
func (repository *Repository[T]) applyConfiguration(config *Config) {
	repository.connection = config.DatabaseConnection
}

// All Gets a collection of all entries in the repository.
// Returns nil if query fails
func (repository *Repository[T]) All() *Nbj.Collection[T] {
	var entries []T

	if result := repository.connection.Find(&entries); result.Error != nil {
		return nil
	}

	return Nbj.Collect(entries)
}

// Query Takes a closure containing a gorm query, executes it and
// returns the result as a collection of entries. Returns nil if
// query fails
func (repository *Repository[T]) Query(closure func(query *gorm.DB) *gorm.DB) *Nbj.Collection[T] {
	var entries []T

	query := closure(repository.connection)

	if result := query.Find(&entries); result.Error != nil {
		return nil
	}

	return Nbj.Collect(entries)
}

func (repository *Repository[T]) First(closures ...func(query *gorm.DB) *gorm.DB) *T {
	var entry T

	// If no closures are passed to the method
	if len(closures) == 0 {
		if result := repository.connection.First(&entry); result.Error != nil {
			return nil
		}

		return &entry
	}

	// Apply all closures to the query
	query := repository.connection

	for _, closure := range closures {
		query = closure(repository.connection)
	}

	if result := query.First(&entry); result.Error != nil {
		return nil
	}

	return &entry
}
