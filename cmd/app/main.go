package main

import (
	"context"
	"datatom/grpc"
	"datatom/internal"
	"datatom/internal/adapter/pg"
	"datatom/internal/api"
	"datatom/internal/migrations"
	pkgpg "datatom/pkg/db/pg"
	"datatom/pkg/log"
	"datatom/rest"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	c := newConfig()
	if err := parse(os.Args[1:], c); err != nil {
		panic(err.Error())
	}
	l, err := log.NewLogger()
	if err != nil {
		panic(err.Error())
	}

	info := internal.NewInfo(c.Title, c.Description)
	info.SetVersion(0, 1, 0)

	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	db, err := pkgpg.NewDB(dbCtx, pkgpg.Config{
		Address:      c.PostgresAddress,
		Port:         c.PostgresPort,
		User:         c.PostgresUser,
		Password:     c.PostgresPassword,
		DatabaseName: c.PostgresDBName,
	})
	if err != nil {
		cancel()
		panic(err.Error())
	}
	cancel()
	defer db.Close()
	l.Info("database connected")

	if err := migrations.UpMigrations(db); err != nil {
		panic(err.Error())
	}

	repo, err := pg.NewRepository(db, l)
	if err != nil {
		panic(err.Error())
	}
	defer repo.CloseConn(db)
	l.Info("repository configured")

	refTypeManager, err := api.NewRefTypeManager(api.RefTypeConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	l.Info("reference types manager configured")

	recordManager, err := api.NewRecordManager(api.RecordConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	l.Info("records manager configured")

	propertyManager, err := api.NewPropertyManager(api.PropertyConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	l.Info("properties manager configured")

	valueManager, err := api.NewValueManager(api.ValueConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	l.Info("values manager configured")

	storedConfigsManager, err := api.NewStoredConfigManager(api.StoredConfigsConfig{
		Repository: repo,
		Timeout:    time.Second,
	})

	dwGRPCConn := grpc.NewConnection(grpc.Config{
		Logger:  l,
		Address: c.DatawayGRPCAddress,
		Port:    c.DatawayGRPCPort,
	})

	restServer, err := rest.NewServer(rest.Config{
		Logger:  l,
		Port:    c.RESTPort,
		Timeout: time.Second * time.Duration(c.RESTTimeoutSec),

		AppInfo: *info,

		RefTypeManager:       refTypeManager,
		RecordManager:        recordManager,
		PropertyManager:      propertyManager,
		ValueManager:         valueManager,
		StoredConfigsManager: storedConfigsManager,

		DatawayGRPCConnection: dwGRPCConn,
	})
	if err != nil {
		panic(err.Error())
	}

	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		err := restServer.Serve()
		l.Errorln("REST server error:", err.Error())
		return err
	})
	l.Infof("REST server listens at port: %d", c.RESTPort)

	l.Infof("%s service is up", internal.ServiceName)

	if err := g.Wait(); err != nil {
		panic(err.Error())
	}
}
