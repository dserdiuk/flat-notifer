package main

import (
	"github.com/dserdiuk/flat-notifier/internal/notifier"
	"github.com/dserdiuk/flat-notifier/internal/service"
	"github.com/dserdiuk/flat-notifier/internal/source"
	"os"
)

func main() {
	var sources []source.Source

	n := notifier.NewTelegramNotifier(os.Getenv("TG_TOKEN"), os.Getenv("TG_CHAT_ID"))
	myHomeSource := source.NewMyHomeSource(os.Getenv("MYHOME_QUERY_STR"))

	sources = append(sources, myHomeSource)
	s := service.NewCheckService(sources, n)
	go s.Start()

	forever := make(chan bool)
	<-forever
}
