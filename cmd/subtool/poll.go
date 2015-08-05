package main

import (
	"encoding/gob"
	"fmt"
	"github.com/Symantec/Dominator/proto/sub"
	"net/rpc"
	"os"
	"time"
)

func pollSubcommand(client *rpc.Client, args []string) {
	var err error
	clientName := fmt.Sprintf("%s:%d", *subHostname, *subPortNum)
	sleepDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", *interval))
	for iter := 0; *numPolls < 0 || iter < *numPolls; iter++ {
		if iter > 0 {
			time.Sleep(sleepDuration)
		}
		if client == nil {
			client, err = rpc.DialHTTP("tcp", clientName)
			if err != nil {
				fmt.Printf("Error dialing\t%s\n", err)
				os.Exit(1)
			}
		}
		var request sub.PollRequest
		var reply sub.PollResponse
		err = client.Call("Subd.Poll", request, &reply)
		if err != nil {
			fmt.Printf("Error calling\t%s\n", err)
			os.Exit(1)
		}
		if *newConnection {
			client.Close()
			client = nil
		}
		fs := reply.FileSystem
		if fs == nil {
			fmt.Println("No FileSystem pointer")
		} else {
			fs.RebuildPointers()
			if *debug {
				fs.DebugWrite(os.Stdout, "")
			} else {
				fmt.Print(fs)
			}
			if *file != "" {
				f, err := os.Create(*file)
				if err != nil {
					fmt.Printf("Error creating: %s\t%s\n", *file, err)
					os.Exit(1)
				}
				encoder := gob.NewEncoder(f)
				encoder.Encode(fs)
				f.Close()
			}
		}
	}
}