package main

import (
	"context"
	"errors"
	"flag"
	log "github.com/sirupsen/logrus"
	pb "github.com/smirnoffV/area-svc/pb"
	"google.golang.org/grpc"
	"io"
	"math"
	"net"
	"sync"
)

var (
	serverAddress = flag.String("addr", ":10000", "")
)

func main() {
	log.Infof("Start listening grpc server on %s port", *serverAddress)
	ln, err := net.Listen("tcp", *serverAddress)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterAreaServer(s, &AreaService{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := s.Serve(ln); err != nil {
			log.WithError(err).Error("stop listening server")
			return
		}
	}()

	wg.Wait()
}

type AreaService struct {
	max int64
}

func (s *AreaService) Square(ctx context.Context, req *pb.SquareRequest) (*pb.Result, error) {
	return &pb.Result{
		Area: req.Side * req.Side,
	}, nil
}

func (s *AreaService) Rectangle(ctx context.Context, req *pb.RectangleRequest) (*pb.Result, error) {
	return &pb.Result{
		Area: req.Height * req.Width,
	}, nil
}

func (s *AreaService) Circle(ctx context.Context, req *pb.CircleRequest) (*pb.Result, error) {
	return &pb.Result{
		Area: math.Pi * math.Pow(req.Radius, 2),
	}, nil
}

func (s *AreaService) Max(srv pb.Area_MaxServer) error {
	for {
		select {
		case <-srv.Context().Done():
			return nil
		default:
			req, err := srv.Recv()

			if err == io.EOF {
				return errors.New("error reading from closed connection")
			}

			if err != nil {
				log.WithError(err).Info("receiving message from stream problem")
				continue
			}

			if s.max >= req.Number {
				continue
			}

			s.max = req.Number

			if err := srv.Send(&pb.NumberResponse{}); err != nil {
				log.Error("error sending response")
				continue
			}

			log.Infof("changed max to %v", req.Number)
		}
	}
}
