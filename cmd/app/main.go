package main

import (
	"context"
	"log"
	"os"
	"time"

	"datatom/internal"
	"datatom/internal/adapter/pg"
	"datatom/internal/adapter/rmq"
	"datatom/internal/api"
	"datatom/internal/grpc"
	"datatom/internal/handlers"
	"datatom/internal/migrations"
	"datatom/internal/rest"
	"datatom/internal/routines"
	pkgpg "datatom/pkg/db/pg"
	pkglog "datatom/pkg/log"
	pkgrmq "datatom/pkg/message_broker/rmq"
	"github.com/go-co-op/gocron"
	"golang.org/x/sync/errgroup"
)

var Version string

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
	info.SetVersion(Version)

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
	poolCC, err := pkgpg.NewPoolConfig(pkgpg.Config{
		Address:      c.PostgresAddress,
		Port:         c.PostgresPort,
		User:         c.PostgresUser,
		Password:     c.PostgresPassword,
		DatabaseName: c.PostgresDBName,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	dbCtx, dbCancel := context.WithTimeout(context.Background(), time.Second*10)
	repo, err := pg.NewRepository(dbCtx, pg.Config{
		ConnectConfig: poolCC,
		Logger:        l,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	dbCancel()
	if err != nil {
		l.Fatal(err.Error())
	}
	defer repo.Close()
	l.Info("repository configured")

	amqConn, err := pkgrmq.NewConnection(pkgrmq.ConnectionConfig{
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
	defer amqConn.Close()
	l.Infoln("RMQ connection established")

	publisher, err := pkgrmq.NewPublisher(pkgrmq.PublisherConfig{
		Logger:   l,
		Conn:     amqConn,
		Exchange: c.DWExchange,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	defer publisher.Close()
	l.Infoln("RMQ publisher configured")

	broker, err := rmq.NewBroker(rmq.Config{
		Publisher: publisher,
		Logger:    l,
	})

	refTypeManager, err := api.NewRefTypeManager(api.RefTypeConfig{
		Repository: repo,
		Broker:     broker,
		Timeout:    time.Second,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	l.Info("reference types manager configured")

	dbManager, err := api.NewDBManager(api.DBConfig{Repository: repo})
	if err != nil {
		l.Fatal(err.Error())
	}
	l.Info("database manager configured")

	recordManager, err := api.NewRecordManager(api.RecordConfig{
		Repository: repo,
		Broker:     broker,
		Timeout:    time.Second,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	l.Info("records manager configured")

	propertyManager, err := api.NewPropertyManager(api.PropertyConfig{
		Repository: repo,
		Broker:     broker,
		Timeout:    time.Second,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	l.Info("properties manager configured")

	valueManager, err := api.NewValueManager(api.ValueConfig{
		Repository: repo,
		Broker:     broker,
		Timeout:    time.Second,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	l.Info("values manager configured")

	changedDataManager, err := api.NewChangedDataManager(api.ChangedDataConfig{
		Repository: repo,
		Timeout:    time.Second * 5,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	l.Info("change data manager configured")

	storedConfigsManager, err := api.NewStoredConfigManager(api.StoredConfigsConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		l.Fatal(err.Error())
	}
	l.Info("stored configs manager configured")

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
		consumer := pkgrmq.NewConsumer(pkgrmq.ConsumerConfig{
			Logger:    l,
			Conn:      amqConn,
			Queue:     c.RMQConsumeQueue,
			QueueArgs: pkgrmq.NewQueueArgs().AsClassic().AddDLEArg(c.RMQDLE),
			Handler: handlers.NewConsumeHandler(handlers.ConsumeHandlerConfig{
				Logger:          l,
				Timeout:         time.Second * 2,
				ValueManager:    valueManager,
				PropertyManager: propertyManager,
				RecordManager:   recordManager,
			}),
		})
		if err := consumer.Consume(); err != nil {
			l.Errorf("messages consuming error: %s", err)
			return err
		}
		return nil
	})

	s := gocron.NewScheduler(time.UTC)

	if amqConn != nil && c.DWExchange != "" {
		if _, err := s.Every(10).Second().SingletonMode().Do(routines.NewSendChangedDataRoutine(routines.SendChangedDataConfig{
			Logger:               l,
			ReferenceTypeManager: refTypeManager,
			RecordManager:        recordManager,
			PropertyManager:      propertyManager,
			ValueManager:         valueManager,
			ChangedDataManager:   changedDataManager,
			StoredConfigsManager: storedConfigsManager,
			DBManager:            dbManager,
			Exchange:             c.DWExchange,
			RoutingKeys:          []string{c.DWRoutingKey},
		})); err != nil {
			l.Fatalf("add routine job error: %s", err)
		}
	}
	s.StartAsync()
	l.Infof("routines are running")

	l.Infof("%s service is up", internal.ServiceName)

	if err := g.Wait(); err != nil {
		l.Fatal(err.Error())
	}
}
