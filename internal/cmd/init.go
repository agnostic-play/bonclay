package cmd

import (
	"github.com/agnostic-play/ditoo/internal/config"
	"github.com/agnostic-play/ditoo/internal/echo"
	"github.com/agnostic-play/ditoo/internal/repository"
	"github.com/agnostic-play/ditoo/internal/services"
	"github.com/agnostic-play/ditoo/pkg/gorm"
	"github.com/labstack/gommon/log"
)

func Init(conf *config.Config) (chan bool, chan bool) {

	var (
		exitSignal = make(chan bool)
		db, sqlDB  = gorm.Init(gorm.ConfigDB{
			DebugLog:            conf.Database.DebugLog,
			Host:                conf.Database.Host,
			User:                conf.Database.User,
			Port:                conf.Database.Port,
			Pass:                conf.Database.Pass,
			DatabaseName:        conf.Database.DatabaseName,
			PoolMaxIdleConn:     conf.Database.PoolMaxIdleConn,
			PoolMaxOpenConn:     conf.Database.PoolMaxOpenConn,
			PoolMaxConnLifetime: conf.Database.PoolMaxConnLifetime,
		})
	)
	repoContainer := repository.NewRepoContainerGorm(db, sqlDB, conf)
	serviceContainer := services.NewServiceContainer(repoContainer, conf)

	serverExitSignal := echo.RunHttpServer(conf, serviceContainer)

	go func() {
		<-exitSignal // Receive exit signal
		log.Info("disconnecting service dependencies")

		if err := sqlDB.Close(); err != nil {
			log.Error(err)
		}

		log.Info("finished disconnecting service dependencies")
		exitSignal <- true // Send signal already finish the job
	}()

	return exitSignal, serverExitSignal
}
