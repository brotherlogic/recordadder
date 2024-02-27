package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brotherlogic/goserver/utils"

	pb "github.com/brotherlogic/recordadder/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	ctx, cancel := utils.ManualContext("recordader-cli", time.Second*10)
	defer cancel()

	conn, err := utils.LFDialServer(ctx, "recordadder")
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewAddRecordServiceClient(conn)

	switch os.Args[1] {
	case "add":
		addFlags := flag.NewFlagSet("AddRecords", flag.ExitOnError)
		var id = addFlags.Int("id", -1, "Id of the record to add")
		var cost = addFlags.Int("cost", 0, "Cost of the record")
		var folder = addFlags.Int("folder", 0, "Goal folder for the record")
		var resetFoldr = addFlags.Int("resetfolder", 0, "Reset folder for the record")
		var year = addFlags.Int("year", time.Now().Year(), "Year for accounting purposes")
		var location = addFlags.String("location", "", "Purchase Location")
		var parents = addFlags.Bool("parents", false, "")

		if err := addFlags.Parse(os.Args[2:]); err == nil {
			if *id > 0 && *cost > 0 && *folder > 0 {
				res, err := client.AddRecord(ctx, &pb.AddRecordRequest{WasParents: *parents, PurchaseLocation: *location, Cost: int32(*cost), Id: int32(*id), Folder: int32(*folder), ResetFolder: int32(*resetFoldr), AccountingYear: int32(*year)})
				if err != nil {
					log.Fatalf("Error on Add Record: %v", err)
				}
				fmt.Printf("Expected to be added on %v\n", time.Unix(res.ExpectedAdditionDate, 0))
			} else {
			log.Printf("Missing data: %v, %v, %v", *id, *cost, *folder)
			}
		} else {
		log.Printf("Error in add: %v", err)
		}
	case "list":
		res, err := client.ListQueue(ctx, &pb.ListQueueRequest{})
		if err != nil {
			log.Fatalf("Error in listing: %v", err)
		}
		cost := int32(0)
		for _, entry := range res.GetRequests() {
			fmt.Printf("%v - %v: %v\n", entry.GetId(), entry.GetCost(), entry)
			cost += entry.GetCost()
		}
		fmt.Printf("Total Cost = %v\n", cost)
	case "here":
		addFlags := flag.NewFlagSet("AddRecords", flag.ExitOnError)
		var id = addFlags.Int("id", -1, "Id of the record to add")

		if err := addFlags.Parse(os.Args[2:]); err == nil {
			if *id > 0 {
				_, err := client.UpdateRecord(ctx, &pb.UpdateRecordRequest{Id: int32(*id), Available: true})
				fmt.Printf("Updated: %v\n", err)
			}
		}
	case "delete":
		addFlags := flag.NewFlagSet("AddRecords", flag.ExitOnError)
		var id = addFlags.Int("id", -1, "Id of the record to add")

		if err := addFlags.Parse(os.Args[2:]); err == nil {
			if *id > 0 {
				_, err := client.DeleteRecord(ctx, &pb.DeleteRecordRequest{Id: int32(*id)})
				fmt.Printf("Deleted: %v\n", err)
			}
		}
	case "proc":
		addFlags := flag.NewFlagSet("AddRecords", flag.ExitOnError)
		var id = addFlags.String("id", "", "Name of the type")

		if err := addFlags.Parse(os.Args[2:]); err == nil {
			if *id != "" {
				_, err := client.ProcAdded(ctx, &pb.ProcAddedRequest{Type: *id})
				fmt.Printf("Processed: %v\n", err)
			}
		}
	}

}
