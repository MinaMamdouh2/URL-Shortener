// Package userdb contains user related CRUD functionality.
// Every core package regardless of it's name has that Core type, that core type represents the API.
// Every single core package is going to have a Store type that is going
// to represent the API related implementing the interface.
// What we are trying to do here is talk to Postgres.
package userdb

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log *zap.SugaredLogger
	db  *gorm.DB
}

// NewStore constructs the api for data access.
func NewStore(log *zap.SugaredLogger, db *gorm.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}
