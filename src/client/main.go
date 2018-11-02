package main
import (
    "context"
    "log"
    "fmt"
    pb "../../proto"
    "google.golang.org/grpc"
)

func main() {
    connect,err := grpc.Dial("127.0.0.1:19003",grpc.WithInsecure())
    if err != nil {
        log.Fatal("client connection error:", err)
    }
    defer connect.Close()
    client := pb.NewPersonClient(connect)
    message := &pb.GetMessage{TargetType:1}
    res, err2 := client.GetPerson(context.TODO(),message)
    if err2 != nil {
        panic(err2)
    }
    fmt.Println(res)
}
