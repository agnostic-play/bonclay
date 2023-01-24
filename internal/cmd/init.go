package cmd

import (
	"github.com/labstack/gommon/log"
	"gitlab.linkaja.com/be/ditto/internal/config"
	"gitlab.linkaja.com/be/ditto/internal/echo"
	"gitlab.linkaja.com/be/ditto/internal/repository"
	"gitlab.linkaja.com/be/ditto/internal/services"
	"gitlab.linkaja.com/be/ditto/pkg/gorm"
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

	serverExitSignal := echo.RunHttpServer(conf.App.Host, conf.App.Port, conf.App.BaseURL, serviceContainer)

	go func() {
		<-exitSignal // Receive exit signal
		log.Info("disconnecting service dependencies")

		if err := sqlDB.Close(); err != nil {
			log.Error(err)
		}

		//log.Info("finished disconnecting service dependencies")
		exitSignal <- true // Send signal already finish the job
	}()

	return exitSignal, serverExitSignal
}
