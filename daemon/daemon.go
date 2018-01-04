package daemon

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/itantgo/api/db"
	"github.com/itantgo/api/model"
	"github.com/itantgo/api/ui"
	"net"
)

type Config struct {
	ListenSpec string

	Db db.Config
	UI ui.Config
}

func Run( cfg *Config) error  {

	log.Fatalf("Starting Http on: %s\n", cfg.ListenSpec)

	db, err := db.initDB(cfg.Db);
		if err != nil{
			log.Printf("Error initialize database: %v\n", err)
			return err
		}
	m := model.New(db)

	l, err := net.Listen("tcp", cfg.ListenSpec)
	if err != nil {
		log.Printf("Error creating listener %v\n", err)
		return  err
	}

	ui.Start(cfg.UI, m, l)
	waitForSignal()
	return nil
}

func waitForSignal()  {
	ch := make( chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <- ch
	log.Printf("Got signal: %v, exiting. ", s)
}