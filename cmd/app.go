package main

import (
	"fmt"
	"github.com/dserdiuk/flat-notifier/internal/notifier"
	"github.com/dserdiuk/flat-notifier/internal/service"
	"github.com/dserdiuk/flat-notifier/internal/source"
	"log"
	"net/http"
	"os"
)

func main() {
	var sources []source.Source

	n := notifier.NewTelegramNotifier(os.Getenv("TG_TOKEN"), os.Getenv("TG_CHAT_ID"))
	myHomeSource := source.NewMyHomeSource(os.Getenv("MYHOME_QUERY_STR"))
	ssSource := source.NewSsSource(os.Getenv("SS_QUERY_STR"))

	sources = append(sources, myHomeSource, ssSource)
	s := service.NewCheckService(sources, n)
	log.Println("Start checking service")
	go s.Start()
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "ok")
	})
	http.ListenAndServe(":8080", nil)
}
