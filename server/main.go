package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"worked-gRPC-mangashop/api"
)

const port = ":8080"

var mangas = make([]*api.Manga, 0, 5)

type ServiceCRUD struct {
	api.UnimplementedServiceCRUDServer
}

func (c *ServiceCRUD) CreateManga(ctx context.Context, manga *api.Manga) (*api.Id, error) {
	log.Printf("Received: %v", manga)

	res := api.Id{}
	res.Id = manga.GetId()
	mangas = append(mangas, manga)
	return &res, nil
}

func (c *ServiceCRUD) ReadManga(ctx context.Context, id *api.Id) (*api.Manga, error) {
	log.Printf("Received: %v", id)

	res := &api.Manga{}

	for _, manga := range mangas {
		if manga.GetId() == id.GetId() {
			res = manga
			break
		}
	}

	return res, nil
}

func (c *ServiceCRUD) ReadAllManga(empty *api.Empty, stream api.ServiceCRUD_ReadAllMangaServer) error {
	log.Printf("Received: %v", empty)

	for _, m := range mangas {
		if err := stream.Send(m); err != nil {
			return err
		}
	}
	return nil
}

func (c *ServiceCRUD) UpdateManga(ctx context.Context, manga *api.Manga) (*api.Status, error) {
	log.Printf("Received: %v", manga)

	res := &api.Status{}

	for i, m := range mangas {
		if manga.GetId() == m.GetId() {
			mangas = append(mangas[:i], mangas[i+1:]...)
			manga.Id = m.Id
			mangas = append(mangas, manga)
			res.Value = 1
			break
		}
	}
	return res, nil
}

func (c *ServiceCRUD) DeleteManga(ctx context.Context, id *api.Id) (*api.Status, error) {
	log.Printf("Received: %v", id)

	res := &api.Status{}
	for index, manga := range mangas {
		if manga.GetId() == id.GetId() {
			mangas = append(mangas[:index], mangas[index+1:]...)
			res.Value = 1
			break
		}
	}

	return res, nil
}

func main() {
	initMangas()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("cannot listen to %s: %v", port, err)
	}

	defer func(listen net.Listener) {
		if err := listen.Close(); err != nil {
			log.Fatalf("cannot close %v with error: %v", listen, err)
		}
	}(listener)

	grpcServer := grpc.NewServer()
	serviceCRUD := new(ServiceCRUD)
	api.RegisterServiceCRUDServer(grpcServer, serviceCRUD)

	log.Printf("Server listening at: %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve on %v: %v", listener.Addr(), err)
	}
}

func initMangas() {
	manga1 := &api.Manga{
		Id:    123,
		Name:  "Naruto",
		Genre: []string{"Senen", "Action", "Fantasy"},
		Cover: "Hardcover",
		Price: 20,
	}
	manga2 := &api.Manga{
		Id:    555,
		Name:  "Death Note",
		Genre: []string{"Detective", "Action", "Psychology"},
		Cover: "Softcover",
		Price: 15,
	}

	mangas = append(mangas, manga1, manga2)
}
