package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	obspb "github.com/lassenordahl/disaggui/obs/proto"
	"google.golang.org/grpc"
)

func handleStatements(client obspb.CRDBServiceClient, reader *bufio.Reader) {
	for {
		fmt.Print("Enter text: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		input = strings.TrimSpace(input)
		words := strings.Fields(input)
		var capitalizedWords []string
		for _, word := range words {
			if isCapitalized(word) {
				capitalizedWords = append(capitalizedWords, word)
			}
		}

		if len(capitalizedWords) == 0 {
			continue
		}

		// Generate a timestamp for the fingerprint.
		timestamp := time.Now().Format(time.UTC.String())
		fingerprint := &obspb.Fingerprint{
			Input:     strings.Join(capitalizedWords, " "),
			Timestamp: timestamp,
		}

		resp, err := client.ProcessFingerprint(context.Background(), fingerprint)
		if err != nil {
			log.Fatalf("Error calling ProcessFingerprint: %v", err)
		}

		log.Printf("Response from server: %s", resp.GetMessage())
	}
}

func isCapitalized(s string) bool {
	return strings.ToUpper(s) == s
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := obspb.NewCRDBServiceClient(conn)
	reader := bufio.NewReader(os.Stdin)

	handleStatements(client, reader)
}
