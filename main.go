package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Open the input file
	inputFile, err := os.Open("Whatsapp Chat with eProd Solutions.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	// Create the output file
	outputFile, err := os.Create("austinndauwa_messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Create a scanner to read the input file line by line
	scanner := bufio.NewScanner(inputFile)

	var currentMessage strings.Builder
	var currentDateTime string
	isAustinMessage := false
	messageCount := 0

	// Function to write the current message to the output file
	writeMessage := func() {
		if isAustinMessage && currentMessage.Len() > 0 {
			messageCount++
			_, err := fmt.Fprintf(outputFile, "Message %d:\n%s\n%s\n\n", messageCount, currentDateTime, strings.TrimSpace(currentMessage.String()))
			if err != nil {
				log.Fatal(err)
			}
		}
		currentMessage.Reset()
		isAustinMessage = false
		currentDateTime = ""
	}

	// Iterate through each line in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Check if this line looks like the start of a new message
		if strings.Contains(line, " - ") {
			// Write the previous message if it exists
			writeMessage()

			// Split the line into timestamp and content
			parts := strings.SplitN(line, " - ", 2)
			if len(parts) == 2 {
				currentDateTime = parts[0] // Capture the timestamp

				// Check if this is a message from austinndauwa
				senderAndMessage := strings.SplitN(parts[1], ": ", 2)
				if len(senderAndMessage) == 2 && strings.EqualFold(strings.TrimSpace(senderAndMessage[0]), "austinndauwa") {
					isAustinMessage = true
					currentMessage.WriteString(senderAndMessage[1])
					currentMessage.WriteString("\n")
				}
			}
		} else if isAustinMessage {
			// If we're in an Austin message, append this line
			currentMessage.WriteString(line)
			currentMessage.WriteString("\n")
		}
	}

	// Write the last message if it exists
	writeMessage()

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Extraction complete. %d messages from austinndauwa saved to austinndauwa_messages.txt\n", messageCount)
}
