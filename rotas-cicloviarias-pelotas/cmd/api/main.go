package main

import (
	"flag"
	"log"
	"net/http"

	"rei-julian/internal/mapdata"
	"rei-julian/internal/routing"
	"rei-julian/internal/tracking"
	"rei-julian/internal/transport/httpapi"
)

func main() {
	addr := flag.String("addr", ":8080", "endereco HTTP")
	mapPath := flag.String("map", "maps/city_01.json", "arquivo JSON do mapa")
	staticDir := flag.String("static", "web", "diretorio estatico")
	flag.Parse()

	repo, err := mapdata.Load(*mapPath)
	if err != nil {
		log.Fatal(err)
	}

	routingService := routing.NewService(repo)
	trackingService := tracking.NewService(repo, routingService)
	server := httpapi.NewServer(repo, routingService, trackingService)

	log.Printf("Rei-Julian API em http://localhost%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, server.Routes(*staticDir)))
}
