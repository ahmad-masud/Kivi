package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	pb "github.com/ahmad-masud/Kivi/proto"
	"google.golang.org/grpc"
)

type kvStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func newKVStore() *kvStore {
	return &kvStore{
		data: make(map[string]string),
	}
}

type server struct {
	pb.UnimplementedKVServer
	store *kvStore
}

func (s *server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	if req.GetKey() == "" {
		return &pb.PutResponse{Ok: false, Error: "empty key"}, nil
	}
	s.store.mu.Lock()
	s.store.data[req.GetKey()] = req.GetValue()
	s.store.mu.Unlock()
	return &pb.PutResponse{Ok: true}, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	if req.GetKey() == "" {
		return &pb.GetResponse{Found: false, Error: "empty key"}, nil
	}
	s.store.mu.RLock()
	val, ok := s.store.data[req.GetKey()]
	s.store.mu.RUnlock()
	if !ok {
		return &pb.GetResponse{Found: false}, nil
	}
	return &pb.GetResponse{Found: true, Value: val}, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if req.GetKey() == "" {
		return &pb.DeleteResponse{Ok: false, Error: "empty key"}, nil
	}
	s.store.mu.Lock()
	_, ok := s.store.data[req.GetKey()]
	if ok {
		delete(s.store.data, req.GetKey())
	}
	s.store.mu.Unlock()
	if !ok {
		return &pb.DeleteResponse{Ok: false, Error: "key not found"}, nil
	}
	return &pb.DeleteResponse{Ok: true}, nil
}

func (s *server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	resp := &pb.ListResponse{}
	for k, v := range s.store.data {
		resp.Pairs = append(resp.Pairs, &pb.KVPair{Key: k, Value: v})
	}
	return resp, nil
}

func main() {
	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}
	s := grpc.NewServer()
	kv := newKVStore()
	pb.RegisterKVServer(s, &server{store: kv})

	go func() {
		log.Printf("gRPC KV server listening on %s", addr)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT or SIGTERM
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down gRPC server")
	s.GracefulStop()
}
