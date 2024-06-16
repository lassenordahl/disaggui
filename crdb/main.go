package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	obspb "github.com/lassenordahl/disaggui/obs/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := obspb.NewCRDBServiceClient(conn)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter text: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		words := strings.Fields(input)
		var capitalizedWords []string
		for _, word := range words {
			if isCapitalized(word) {
				capitalizedWords = append(capitalizedWords, word)
			}
		}

		fingerprint := &obspb.Fingerprint{Input: strings.Join(capitalizedWords, " ")}
		resp, err := client.ProcessFingerprint(context.Background(), fingerprint)
		if err != nil {
			log.Fatalf("Error calling ProcessFingerprint: %v", err)
		}

		log.Printf("Response from server: %s", resp.GetMessage())
	}
}

func isCapitalized(word string) bool {
	return word == strings.Title(word)
}
