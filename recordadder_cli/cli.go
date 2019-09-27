package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/recordadder/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	host, port, err := utils.Resolve("recordadder", "recordadder-cli")
	if err != nil {
		log.Fatalf("Unable to reach recordadder: %v", err)
	}
	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewAddRecordServiceClient(conn)
	ctx, cancel := utils.BuildContext("recordader-cli", "recordadder")
	defer cancel()

	switch os.Args[1] {
	case "add":
		addFlags := flag.NewFlagSet("AddRecords", flag.ExitOnError)
		var id = addFlags.Int("id", -1, "Id of the record to add")
		var cost = addFlags.Int("cost", -1, "Cost of the record")
		var folder = addFlags.Int("folder", -1, "Goal folder for the record")

		if err := addFlags.Parse(os.Args[2:]); err == nil {
			res, err := client.AddRecord(ctx, &pb.AddRecordRequest{Cost: int32(*cost), Id: int32(*id), Folder: int32(*folder)})
			if err != nil {
				log.Fatalf("Error on Add Record: %v", err)
			}
			fmt.Printf("Expected to be added on %v\n", time.Unix(res.ExpectedAdditionDate, 0))
		}
	}
}
