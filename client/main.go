package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"

	pb "github.com/aagoldingay/commendeer-go/pb"
)

const (
	addr = "localhost:8080"
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	c := pb.NewCommendeerClient(conn)

	loginresp, err := c.LoginUser(context.Background(), &pb.LoginRequest{Username: "admin1", Password: "4dm1n123"})
	fmt.Println(loginresp)

	logoutresp, err := c.LogoutUser(context.Background(), &pb.LogoutRequest{Authcode: loginresp.Authcode})
	fmt.Println(logoutresp)
}
