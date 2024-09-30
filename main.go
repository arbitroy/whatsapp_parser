package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	// Prompt for the year to parse
	var year string
	fmt.Print("Enter the year to parse (e.g., 2024): ")
	fmt.Scanln(&year)

	// Prompt for user selection
	var userOption string
	fmt.Print("Select user option (1: All users, 2: Austin, 3: Custom username): ")
	fmt.Scanln(&userOption)

	var targetUser string
	switch userOption {
	case "1":
		targetUser = ""  // Empty string means all users
	case "2":
		targetUser = "austinndauwa"
	case "3":
		fmt.Print("Enter the username to parse: ")
		fmt.Scanln(&targetUser)
	default:
		log.Fatal("Invalid option selected")
	}

	// Open the input file
	inputFile, err := os.Open("Whatsapp Chat with eProd Solutions.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	// Create the output file
	outputFileName := fmt.Sprintf("whatsapp_messages_%s_%s.txt", year, strings.ReplaceAll(targetUser, " ", "_"))
	if targetUser == "" {
		outputFileName = fmt.Sprintf("whatsapp_messages_%s_all_users.txt", year)
	}
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Create a scanner to read the input file line by line
	scanner := bufio.NewScanner(inputFile)

	var currentMessage strings.Builder
	var currentDateTime string
	var currentMonth string
	var isRelevantMessage bool
	messageCount := 0
	monthlyMessageCount := 0

	// Function to write the current message to the output file
	writeMessage := func() {
		if isRelevantMessage && currentMessage.Len() > 0 {
			messageCount++
			monthlyMessageCount++
			_, err := fmt.Fprintf(outputFile, "Message %d:\n%s\n%s\n\n", monthlyMessageCount, currentDateTime, strings.TrimSpace(currentMessage.String()))
			if err != nil {
				log.Fatal(err)
			}
		}
		currentMessage.Reset()
		isRelevantMessage = false
		currentDateTime = ""
	}

	// Regular expression to match and clean the date string
	dateRegex := regexp.MustCompile(`(\d{1,2}/\d{1,2}/\d{2},\s+\d{1,2}:\d{2}).*?(am|pm)`)

	// Function to parse the date string
	parseDate := func(dateStr string) (time.Time, error) {
		// Use regex to extract and clean the date string
		matches := dateRegex.FindStringSubmatch(dateStr)
		if len(matches) < 3 {
			return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
		}

		// Reconstruct the date string in the correct format
		cleanDateStr := fmt.Sprintf("%s %s", matches[1], matches[2])

		// Parse the cleaned date string
		return time.Parse("2/1/06, 3:04 pm", cleanDateStr)
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

				// Parse the date to get the month and year
				date, err := parseDate(currentDateTime)
				if err != nil {
					log.Printf("Error parsing date: %v for date string: %s", err, currentDateTime)
					continue
				}

				// Check if the message is from the specified year
				if date.Format("2006") != year {
					continue
				}

				newMonth := date.Format("January 2006")

				// If the month has changed, write a new month header
				if newMonth != currentMonth {
					if currentMonth != "" {
						_, err = fmt.Fprintf(outputFile, "\nTotal messages for %s: %d\n\n", currentMonth, monthlyMessageCount)
						if err != nil {
							log.Fatal(err)
						}
					}
					_, err = fmt.Fprintf(outputFile, "--- %s ---\n\n", newMonth)
					if err != nil {
						log.Fatal(err)
					}
					currentMonth = newMonth
					monthlyMessageCount = 0
				}

				// Check if this is a relevant message based on user selection
				senderAndMessage := strings.SplitN(parts[1], ": ", 2)
				if len(senderAndMessage) == 2 {
					sender := strings.TrimSpace(senderAndMessage[0])
					isRelevantMessage = targetUser == "" || strings.EqualFold(sender, targetUser)
					if isRelevantMessage {
						currentMessage.WriteString(fmt.Sprintf("%s: %s\n", sender, senderAndMessage[1]))
					}
				}
			}
		} else if isRelevantMessage {
			// If we're in a relevant message, append this line
			currentMessage.WriteString(line)
			currentMessage.WriteString("\n")
		}
	}

	// Write the last message if it exists
	writeMessage()

	// Write the final month's total
	if currentMonth != "" {
		_, err = fmt.Fprintf(outputFile, "\nTotal messages for %s: %d\n", currentMonth, monthlyMessageCount)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Extraction complete. %d messages saved to %s\n", messageCount, outputFileName)
}