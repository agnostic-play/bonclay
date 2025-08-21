package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/agnostic-play/ditoo/internal/config"
	"gorm.io/gorm"
)

type RepoContainer interface {
	SquadRepoInterface
	CollectionRepoInterface
	EndpointRepoInterface
	ScenarioRepoInterface
	EndpointScenarioRepoInterface
	CustomVariableRepoInterface
	Begin()
	Rollback() error
	Commit() error
}

func NewRepoContainerGorm(db *gorm.DB, sqlDb *sql.DB, config *config.Config) RepoContainer {
	return &repoContainerGorm{
		db:     db,
		sqlDB:  sqlDb,
		config: config,
	}
}

type repoContainerGorm struct {
	sqlDB   *sql.DB
	db      *gorm.DB
	transDB *gorm.DB
	config  *config.Config
}

func (cont repoContainerGorm) write(ctx context.Context) *gorm.DB {
	if cont.transDB != nil {
		return cont.transDB.WithContext(ctx)
	}
	return cont.db.WithContext(ctx)
}

func (cont repoContainerGorm) Begin() {
	cont.transDB = cont.db.Begin()
}

func (cont repoContainerGorm) Rollback() error {
	if cont.transDB == nil {
		return fmt.Errorf("transaction has not started")
	}
	cont.transDB.Rollback()
	return nil
}

func (cont repoContainerGorm) Commit() error {
	if cont.transDB == nil {
		return fmt.Errorf("transaction has not started")
	}
	cont.transDB.Commit()
	return nil
}
