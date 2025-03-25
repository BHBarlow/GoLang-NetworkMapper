package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

// Config struct to hold target and Nmap flags
type Config struct {
	Target    string   `json:"target"`
	NmapFlags []string `json:"nmap_flags"`
}

// Read configuration from a JSON file
func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func runNmapScan(target string, flags []string) (map[string][]string, error) {
	fmt.Printf("\nRunning Nmap Scan for %s\n", target)

	args := append(flags, target, "-oG", "-")
	cmd := exec.Command("nmap", args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running Nmap: %v", err)
	}

	hostMap := make(map[string][]string)
	var lastHost string

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()

		hostRegex := regexp.MustCompile(`Host: ([\d\.]+) \((.*?)\)`)
		hostMatches := hostRegex.FindStringSubmatch(line)

		portsRegex := regexp.MustCompile(`(\d+)/(\w+)/tcp//([\w-]+)`)
		portMatches := portsRegex.FindAllStringSubmatch(line, -1)

		osRegex := regexp.MustCompile(`OS: (.+)`)
		osMatches := osRegex.FindStringSubmatch(line)

		if len(hostMatches) > 0 {
			ip := hostMatches[1]
			hostname := hostMatches[2]
			fullID := fmt.Sprintf("IP: %s (%s)", ip, hostname)
			lastHost = fullID

			if _, exists := hostMap[fullID]; !exists {
				hostMap[fullID] = []string{"Ports:"}
			}
		}

		if lastHost != "" && len(portMatches) > 0 {
			for _, match := range portMatches {
				port := match[1]
				state := match[2]
				service := match[3]
				hostMap[lastHost] = append(hostMap[lastHost], fmt.Sprintf("    %s/%s/%s", port, state, service))
			}
		}

		if lastHost != "" && len(osMatches) > 1 {
			hostMap[lastHost] = append(hostMap[lastHost], "OS:", "    "+osMatches[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Nmap output: %v", err)
	}

	return hostMap, nil
}

func getActiveHosts(target string) ([]string, error) {
	cmd := exec.Command("nmap", "-sP", target, "-oG", "-")

	outputPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating output pipe for Ping Scan: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting Nmap Ping Scan: %v", err)
	}

	scanner := bufio.NewScanner(outputPipe)
	var activeHosts []string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Host") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				activeHosts = append(activeHosts, parts[1])
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("error waiting for Ping Scan: %v", err)
	}

	return activeHosts, nil
}

func scanActiveHost(host string, flags []string, wg *sync.WaitGroup, ch chan map[string][]string) {
	defer wg.Done()

	hostData, err := runNmapScan(host, flags)
	if err != nil {
		fmt.Printf("Error running Nmap for host %s: %v\n", host, err)
		return
	}

	ch <- hostData
}

func main() {
	// Load configuration from file
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config file:", err)
		return
	}

	target := flag.String("target", "", "Target subnet or IP (e.g., 192.168.1.0/24)")
	flag.Parse()

	if *target == "" {
		*target = config.Target
	}

	if *target == "" {
		fmt.Println("Usage: go run main.go -target <IP/Subnet> or set target in config.json")
		return
	}

	activeHosts, err := getActiveHosts(*target)
	if err != nil {
		fmt.Println("\nError getting active hosts:", err)
		return
	}

	fmt.Println("\nActive Hosts:")
	for _, host := range activeHosts {
		fmt.Println(host)
	}

	var wg sync.WaitGroup
	ch := make(chan map[string][]string, len(activeHosts))

	for _, host := range activeHosts {
		wg.Add(1)
		go scanActiveHost(host, config.NmapFlags, &wg, ch)
	}

	wg.Wait()
	close(ch)

	hostResults := make(map[string][]string)
	for hostData := range ch {
		for host, ports := range hostData {
			hostResults[host] = ports
		}
	}

	jsonData, err := json.MarshalIndent(hostResults, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
