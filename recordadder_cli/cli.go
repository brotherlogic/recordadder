package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/recordadder/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(&utils.DiscoveryClientResolverBuilder{})
}

func main() {
	conn, err := grpc.Dial("discovery:///recordadder", grpc.WithInsecure(), grpc.WithBalancerName("my_pick_first"))
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
		var cost = addFlags.Int("cost", 0, "Cost of the record")
		var folder = addFlags.Int("folder", 0, "Goal folder for the record")
		var resetFoldr = addFlags.Int("resetfolder", 0, "Reset folder for the record")
		var year = addFlags.Int("year", time.Now().Year(), "Year for accounting purposes")

		if err := addFlags.Parse(os.Args[2:]); err == nil {
			res, err := client.AddRecord(ctx, &pb.AddRecordRequest{Cost: int32(*cost), Id: int32(*id), Folder: int32(*folder), ResetFolder: int32(*resetFoldr), AccountingYear: int32(*year)})
			if err != nil {
				log.Fatalf("Error on Add Record: %v", err)
			}
			fmt.Printf("Expected to be added on %v\n", time.Unix(res.ExpectedAdditionDate, 0))
		}
	}
}
