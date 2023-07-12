package main

import (
	"context"
	"datatom/grpc"
	"datatom/internal"
	"datatom/internal/adapter/pg"
	"datatom/internal/api"
	"datatom/internal/handlers"
	"datatom/internal/migrations"
	"datatom/pkg/amq"
	pkgpg "datatom/pkg/db/pg"
	pkglog "datatom/pkg/log"
	"datatom/rest"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	versionMajor = 0
	versionMinor = 1
	versionPatch = 0
)

func main() {
	c := newConfig()
	if err := parse(os.Args[1:], c); err != nil {
		log.Fatal(err.Error())
	}
	l, err := pkglog.NewLogger()
	if err != nil {
		log.Fatal(err.Error())
	}

	info := internal.NewInfo(c.Title, c.Description)
	info.SetVersion(versionMajor, versionMinor, versionPatch)

	dbCC, err := pkgpg.NewConnConfig(pkgpg.Config{
		Address:      c.PostgresAddress,
		Port:         c.PostgresPort,
		User:         c.PostgresUser,
		Password:     c.PostgresPassword,
		DatabaseName: c.PostgresDBName,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	if err := migrations.UpMigrations(dbCC); err != nil {
		l.Fatal(err.Error())
	}
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	repo, err := pg.NewRepository(dbCtx, pg.Config{
		ConnectConfig: dbCC,
		Logger:        l,
	})
	cancel()
	if err != nil {
		l.Fatal(err.Error())
	}
	defer repo.Close(context.Background())
	l.Info("repository configured")

	refTypeManager, err := api.NewRefTypeManager(api.RefTypeConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	l.Info("reference types manager configured")

	recordManager, err := api.NewRecordManager(api.RecordConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	l.Info("records manager configured")

	propertyManager, err := api.NewPropertyManager(api.PropertyConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	l.Info("properties manager configured")

	valueManager, err := api.NewValueManager(api.ValueConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		l.Fatal(err.Error())
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
		l.Fatal(err.Error())
	}

	amqConn, err := amq.NewConnection(amq.ConnectionConfig{
		Logger:   l,
		Address:  c.RMQAddress,
		Port:     c.RMQPort,
		User:     c.RMQUser,
		Password: c.RMQPassword,
		VHost:    c.RMQVHost,
	})
	if err != nil {
		l.Fatalf("rmq dial error: %s", err)
	}
	if amqConn != nil {
		defer amqConn.Close()
		l.Infoln("RMQ connection established")
	}

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		err := restServer.Serve()
		l.Errorln("REST server error:", err.Error())
		return err
	})
	l.Infof("REST server listens at port: %d", c.RESTPort)

	g.Go(func() error {
		if amqConn == nil {
			return nil
		}
		if err := amq.RunNewConsumer(amq.ConsumerConfig{
			Logger: l,
			Conn:   amqConn,
			Queue:  c.RMQConsumeQueue,
			Handler: handlers.NewConsumeHandler(handlers.ConsumeHandlerConfig{
				Logger:          l,
				Timeout:         time.Second * 2,
				ValueManager:    valueManager,
				PropertyManager: propertyManager,
			}),
			QueueArgs: amq.NewQueueArgs().AddTypeArg(amqp.QueueTypeClassic).AddDLEArg(c.RMQDLE),
		}); err != nil {
			l.Errorf("rmq run consuming error: %s", err)
			return err
		}
		return nil
	})

	l.Infof("%s service is up", internal.ServiceName)

	if err := g.Wait(); err != nil {
		l.Fatal(err.Error())
	}
}
