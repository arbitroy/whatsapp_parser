package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main(){
	file, err :=  os.Open("Whatsapp Chat with eProd Solutions.txt")

	if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

	// Define a regular expression pattern
    // pattern := `(\d{2}/\d{2}/\d{4}), (\d{2}:\d{2}) - (.*?): (.*)`
	newPattern := `\d{2}/\d{2}/\d{2},\s\d{2}:\d{2}\s*.*austinndauwa.*`

    // Compile the regular expression
    re := regexp.MustCompile(newPattern)

	scanner := bufio.NewScanner(file);

	// Iterate through each line in the file
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			fmt.Printf("Match found: %s\n", line)
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}