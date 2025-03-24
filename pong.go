package main
 
import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)
 
func runNmapScan(target string) (map[string][]string, error) {
	fmt.Printf("\nRunning Nmap Scan for %s\n", target)
 
	cmd := exec.Command("nmap", "-A", "-T5", "--min-rate", "1000", "--max-retries", "1", target, "-oG", "-")
 
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
 
		// Regex to match host information (e.g., "Host: 192.168.1.1 (hostname)")
		hostRegex := regexp.MustCompile(`Host: ([\d\.]+) \((.*?)\)`)
		hostMatches := hostRegex.FindStringSubmatch(line)
 
		// Regex to match open ports (e.g., "80/open/tcp//http")
		portsRegex := regexp.MustCompile(`(\d+)/(\w+)/tcp//([\w-]+)`)
		portMatches := portsRegex.FindAllStringSubmatch(line, -1)
 
		// Regex to match OS details (e.g., "OS details: Linux 3.x, Windows 10, etc.")
		osRegex := regexp.MustCompile(`OS: (.+)`) // Captures OS and all details after it
		osMatches := osRegex.FindStringSubmatch(line)
 
		// If a new host is found, initialize it
		if len(hostMatches) > 0 {
			ip := hostMatches[1]
			hostname := hostMatches[2]
			fullID := fmt.Sprintf("IP: %s (%s)", ip, hostname)
			lastHost = fullID
 
			if _, exists := hostMap[fullID]; !exists {
				hostMap[fullID] = []string{"Ports:"}
			}
		}
 
		// Add port information if applicable
		if lastHost != "" && len(portMatches) > 0 {
			for _, match := range portMatches {
				port := match[1]
				state := match[2]
				service := match[3]
				hostMap[lastHost] = append(hostMap[lastHost], fmt.Sprintf("    %s/%s/%s", port, state, service))
			}
		}
 
		// Attach OS info to the last seen host
		if lastHost != "" && len(osMatches) > 1 {
			hostMap[lastHost] = append(hostMap[lastHost], "OS:", "    "+osMatches[1])
		}
	}
 
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Nmap output: %v", err)
	}
 
	return hostMap, nil
}
 
// -------------------------------------------------------------------------------------------------------
// getActiveHosts performs a Ping Scan (nmap -sP) to find active hosts.
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
				activeHosts = append(activeHosts, parts[1]) // Store just the IP address
			}
		}
	}
 
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("error waiting for Ping Scan: %v", err)
	}
 
	return activeHosts, nil
}
 
func scanActiveHost(host string, wg *sync.WaitGroup, ch chan map[string][]string) {
	defer wg.Done()
 
	// Run Nmap scan for the active host
	hostData, err := runNmapScan(host)
	if err != nil {
		fmt.Printf("Error running Nmap for host %s: %v\n", host, err)
		return
	}
 
	// Send results to channel
	ch <- hostData
}
 
func main() {
	target := "10.0.2.0/24"
 
	// Get active hosts
	activeHosts, err := getActiveHosts(target)
	if err != nil {
		fmt.Println("\nError getting active hosts:", err)
		return
	}
 
	fmt.Println("\nActive Hosts:")
	for _, host := range activeHosts {
		fmt.Println(host)
	}
 
	// WaitGroup to manage goroutines
	var wg sync.WaitGroup
	// Channel to collect Nmap scan results
	ch := make(chan map[string][]string, len(activeHosts))
 
	// Run full Nmap scan on discovered hosts concurrently
	for _, host := range activeHosts {
		wg.Add(1)
		go scanActiveHost(host, &wg, ch)
	}
 
	// Wait for all goroutines to finish
	wg.Wait()
	close(ch)
 
	// Collect and print the results
	fmt.Println("\nHOST INFO ========================================================")
	for hostData := range ch {
		for host, ports := range hostData {
			fmt.Println(host)
			for _, portInfo := range ports {
				fmt.Println(" ", portInfo)
			}
			fmt.Println()
		}
	}
}
