package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	// Define command line flags
	ip := flag.String("ip", "localhost", "IP address of the target Prometheus Node Exporter")
	port := flag.Int("port", 9100, "Port on which the Node Exporter is listening")
	flag.Parse()

	// Construct the URL from the command line arguments
	url := fmt.Sprintf("http://%s:%d/metrics", *ip, *port)

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error creating request: ", err)
	}

	// Send the request and get a response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request to server: ", err)
	}
	defer resp.Body.Close()

	// Use a bufio.Scanner to read the response body line by line
	scanner := bufio.NewScanner(resp.Body)
	// Regex to match Linux device metrics, adjust according to specific needs
	deviceRegex := regexp.MustCompile(`node_(disk|filesystem|network|thermal_zone|uname)_.*`)
	// Map to store unique device information categorized by type
	deviceInfo := make(map[string]map[string]bool)

	fmt.Println("Grouped Linux Device Metrics:")
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line contains device related metrics
		if deviceRegex.MatchString(line) {
			// Determine the category of the metric
			category := ""
			if strings.Contains(line, "disk") {
				category = "Disk"
			} else if strings.Contains(line, "filesystem") {
				category = "Filesystem"
			} else if strings.Contains(line, "network") {
				category = "Network"
			} else if strings.Contains(line, "thermal_zone") {
				category = "Thermal"
			} else if strings.Contains(line, "uname") {
				category = "System Info"
			}

			// Extract the information within curly braces
			start := strings.Index(line, "{")
			end := strings.Index(line, "}")
			if start != -1 && end != -1 && end > start {
				info := line[start : end+1]
				if deviceInfo[category] == nil {
					deviceInfo[category] = make(map[string]bool)
				}
				deviceInfo[category][info] = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading response body: ", err)
	}

	// Print out the grouped device information
	for category, infos := range deviceInfo {
		fmt.Printf("\n%s:\n", category)
		for info := range infos {
			fmt.Println(info)
		}
	}
}
