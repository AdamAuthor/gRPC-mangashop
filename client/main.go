package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"strconv"
	"time"
	"worked-gRPC-mangashop/api"
)

const address = "localhost:8080"

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Didn`t connect with %s: %v", address, err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Didn`t close %s: %v", address, err)
		}
	}(conn)

	client := api.NewServiceCRUDClient(conn)

	newManga := api.Manga{
		Id:    121212,
		Name:  "One Piece",
		Genre: []string{"Senen", "Action", "Adventures"},
		Cover: "Hardcover",
		Price: 30,
	}

	updatedManga := api.Manga{
		Id:    123,
		Name:  "Naruto",
		Genre: []string{"Senen", "War", "Fantasy"},
		Cover: "Hardcover",
		Price: 20,
	}

	CreateManga(client, &newManga)
	ReadManga(client, 555)
	ReadAllManga(client)
	UpdateManga(client, &updatedManga)
	ReadAllManga(client)
	DeleteManga(client, 555)
	ReadAllManga(client)
}

func CreateManga(client api.ServiceCRUDClient, manga *api.Manga) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &api.Manga{Id: manga.Id, Name: manga.Name, Genre: manga.Genre, Cover: manga.Cover, Price: manga.Price}
	res, err := client.CreateManga(ctx, req)
	if err != nil {
		log.Fatalf("%v.CreateManga(_) = _, %v", client, err)
	}
	if strconv.FormatInt(res.GetId(), 10) != "" {
		log.Printf("Successful created manga with id: %v", res)
	} else {
		log.Printf("CreateManga Failed")
	}
}

func ReadManga(client api.ServiceCRUDClient, id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &api.Id{Id: id}
	res, err := client.ReadManga(ctx, req)
	if err != nil {
		log.Fatalf("%v.ReadManga(_) = _, %v", client, err)
	}

	log.Printf("Information about manga with id %d: %v", id, res)
}

func ReadAllManga(client api.ServiceCRUDClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &api.Empty{}
	stream, err := client.ReadAllManga(ctx, req)
	if err != nil {
		log.Fatalf("%v.ReadManga(_) = _, %v", client, err)
	}

	counter := 1
	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("%v.ReadAllMovies(_) = _, %v", client, err)
		}
		log.Printf("Manga #%d: %v", counter, row)
		counter++
	}
}

func UpdateManga(client api.ServiceCRUDClient, manga *api.Manga) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &api.Manga{Id: manga.Id, Name: manga.Name, Genre: manga.Genre, Cover: manga.Cover, Price: manga.Price}
	res, err := client.UpdateManga(ctx, req)
	if err != nil {
		log.Fatalf("%v.UpdateManga(_) = _, %v", client, err)
	}
	if res.GetValue() == 1 {
		log.Printf("UpdateManga Success")
	} else {
		log.Printf("UpdateManga Failed")
	}
}

func DeleteManga(client api.ServiceCRUDClient, id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &api.Id{Id: id}
	res, err := client.DeleteManga(ctx, req)
	if err != nil {
		log.Fatalf("%v.DeleteManga(_) = _, %v", client, err)
	}
	if res.GetValue() == 1 {
		log.Printf("DeleteManga Success")
	} else {
		log.Printf("DeleteManga Failed")
	}
}
