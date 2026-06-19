package reviewer

import (
	pb "Arboris/generated/go"
	"Arboris/go_server/config"
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToPython(envVar *config.Config) (pb.ReviewServiceClient, *grpc.ClientConn, error) {
	host := envVar.PyServer.Host
	port := envVar.PyServer.Port

	target := host + ":" + port

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		slog.Error("Couldn't establish connection with python server", "ERROR", err)
		return nil, nil, err
	}

	client := pb.NewReviewServiceClient(conn)
	return client, conn, nil
}

func GenerateEmbeddings(client pb.ReviewServiceClient, text string, model string) (*pb.GenerateEmbeddingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	req := &pb.GenerateEmbeddingRequest{
		Text:      text,
		ModelName: model,
	}

	res, err := client.GenerateEmbedding(ctx, req)
	if err != nil {
		slog.Error("Unable to reach python client", "Service", "GenerateEmbeddings", "ERROR", err)
		return nil, err
	}

	return res, nil
}
