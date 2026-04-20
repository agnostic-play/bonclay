package http

import (
	`berlin.allobank.com/tools/bonclay/internal/adapters/echo`
	repository2 `berlin.allobank.com/tools/bonclay/internal/adapters/repositories`
	"berlin.allobank.com/tools/bonclay/internal/config"
	"berlin.allobank.com/tools/bonclay/internal/services"
	"berlin.allobank.com/tools/bonclay/pkg/gorm"
	"github.com/labstack/gommon/log"
)

func Runner(conf *config.Config) (chan bool, chan bool) {
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

	dbClient := repository2.NewDBClient(db)
	repoContainer := repository2.NewRepoContainer(db, conf)
	serviceContainer := services.NewServiceContainer(conf, repoContainer, dbClient)

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
