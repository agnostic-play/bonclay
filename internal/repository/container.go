package repository

import (
	"database/sql"
	"gitlab.linkaja.com/be/ditto/internal/config"
	"gitlab.linkaja.com/be/ditto/internal/errs"
	"gorm.io/gorm"
)

type RepoContainer interface {
	Ping() error
}

func NewRepoContainerGorm(db *gorm.DB, sqlDb *sql.DB, config *config.Config) RepoContainer {
	return &repoContainerGorm{
		db:     db,
		sqlDB:  sqlDb,
		config: config,
	}
}

type repoContainerGorm struct {
	sqlDB  *sql.DB
	db     *gorm.DB
	config *config.Config
}

func (cont repoContainerGorm) Ping() error {
	var dbName string
	exec := cont.db.Raw("SELECT DATABASE()").Find(&dbName)
	if exec.Error != nil {
		return exec.Error
	}

	if dbName == "" {
		return errs.ErrNotFound
	}
	return nil
}
