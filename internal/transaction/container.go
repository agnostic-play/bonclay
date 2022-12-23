package repository

import (
	"database/sql"
	"fmt"
	"gitlab.linkaja.com/be/ditto/internal/config"
	"gorm.io/gorm"
)

type RepoContainer interface {
	Ping() error
	InitiateTrans() TransRepo
}
type TransRepo interface {
	Begin()
	Commit()
	Rollback()
}

func NewRepoContainerGorm(db *gorm.DB, sqlDb *sql.DB, config *config.Config) RepoContainer {
	return &repoContainerGorm{
		db:     db,
		sqlDB:  sqlDb,
		config: config,
	}
}
func (r repoContainerGorm) InitiateTrans() TransRepo {
	return
}

type repoContainerGorm struct {
	sqlDB  *sql.DB
	db     *gorm.DB
	config *config.Config
}

func (cont repoContainerGorm) Ping() error {
	var dbName string
	cont.db.Raw("SELECT DATABASE()").Find(&dbName)
	return fmt.Errorf("gagal")
}
