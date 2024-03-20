// Package db defines interfaces and implementations for interacting with
// entities in the database. It includes the AuditRepository interface, which outlines
// operations for querying and mutating Audit records, and the SQLAuditRepository
// struct, which provides a concrete implementation of the AuditRepository using GORM.
//
// The SQLAuditRepository implementation leverages GORM to perform CRUD operations
// and complex queries on the database, abstracting the data layer away from the
// service layer to facilitate easier testing and maintenance.
package db

import (
	"fmt"
	"github.com/heather92115/translator/internal/mdl"
	"gorm.io/gorm"
	"log"
)

// AuditRepository defines the operations available for a Audit entity.
type AuditRepository interface {
	FindAuditByID(id int) (*mdl.Audit, error)
	FindAudits(tableName string, duration mdl.Duration, limit int) (*[]mdl.Audit, error)
	CreateAudit(Audit *mdl.Audit) error
}

// SQLAuditRepository provides a GORM-based implementation of the AuditRepository interface.
type SQLAuditRepository struct {
	db *gorm.DB
}

// NewSqlAuditRepository initializes a new SQLAuditRepository with a database connection.
func NewSqlAuditRepository() (repo *SQLAuditRepository, err error) {
	db, err := GetConnection()
	if err != nil {
		return
	}

	repo = &SQLAuditRepository{db: db}

	return
}

// FindAuditByID retrieves a single Audit record from the database using its primary ID.
//
// The function attempts to establish a database connection and then queries the Audit table
// for a record matching the specified ID. It is designed to fetch exactly one record or return
// an error if the record does not exist or in case of a connection or query execution error.
//
// Parameters:
// - id: An integer representing the primary ID of the Audit record to retrieve.
//
// Returns:
//   - *mdl.Audit: A pointer to a Audit struct representing the found record. If no record is found
//     or in case of an error, nil is returned.
//   - error: An error object detailing any issues encountered during the database connection
//     attempt or query execution. Errors could include connection failures, issues executing
//     the query, or the situation where no record is found matching the provided ID.
//     In cases where the operation succeeds and a record is found, nil is returned for the error.
//
// Usage example:
// Audit, err := FindAuditByID(123)
//
//	if err != nil {
//	    log.Printf("An error occurred: %v", err)
//	} else {
//		log.Printf("Retrieved Audit: %+v\n", Audit)
//	}
func (repo *SQLAuditRepository) FindAuditByID(id int) (Audit *mdl.Audit, err error) {

	db, err := GetConnection()
	if err != nil {
		return
	}

	result := db.First(&Audit, id) // `First` method adds `WHERE id = ?` to the query
	if result.Error != nil {
		err = fmt.Errorf("error finding Audit with id %d: %v", id, result.Error)
	}

	return
}

func (repo *SQLAuditRepository) FindAudits(tableName string, duration *mdl.Duration, limit int) (Audits *[]mdl.Audit, err error) {
	db, err := GetConnection()
	if err != nil {
		return
	}

	Audits = &[]mdl.Audit{}

	query := db.Limit(limit)

	// Filter by LearningLangCode
	if len(tableName) > 0 {
		query = query.Where("table_name = ?", tableName)
	}

	if duration != nil {
		query = query.Where("created >= ? and created <= ?", duration.Start, duration.End)
	}

	// Execute the query
	err = query.Find(Audits).Error
	if err != nil {
		log.Printf("Error finding %d Audit records with tableName '%s': %v", limit, tableName, err)
	}

	return
}

// CreateAudit inserts a new Audit record into the database.
// It establishes a database connection, then attempts to insert the provided Audit instance.
// Returns an error if the database connection fails or if the insert operation encounters an error.
func (repo *SQLAuditRepository) CreateAudit(Audit *mdl.Audit) error {
	db, err := GetConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to the db, error: %v", err)
	}

	result := db.Create(Audit)
	if result.Error != nil {
		return result.Error
	}

	return nil
}