package main

import (
	"Arboris/go_server/client/reviewer"
	"Arboris/go_server/config"
	"fmt"
	"log"
	"log/slog"
)

func main() {

	envVar, loadErr := config.LoadEnv()

	if loadErr != nil {
		slog.Error("Unable to load the env variables", "ERROR", loadErr)
	} else {
		slog.Info("Loaded env....")
	}

	pythonClient, conn, reviewServerErr := reviewer.ConnectToPython(envVar)

	if reviewServerErr != nil {
		slog.Error("Unable to connect to python client", "ERROR", reviewServerErr)
	} else {
		slog.Info("Successfully connected to python client")
	}

	defer conn.Close()

	port := envVar.GoServer.Port
	host := envVar.GoServer.Host
	address := host + ":" + port
	log.Println("Host:Port ", address)

	res, err := reviewer.GenerateEmbeddings(pythonClient, "", "")

	if err != nil {
		slog.Error("Unable to reach Generate Embeddings")
	}
	fmt.Printf("Response : %v", res)
}
