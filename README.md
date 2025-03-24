# Nmap Network Scanner in Go

## Description
This Go program automates network scanning using Nmap. It performs the following tasks:

1. **Ping Scan (`nmap -sP`)**: Identifies active hosts in the specified subnet.
2. **Full Nmap Scan (`nmap -A`)**: Runs a comprehensive scan on each active host, gathering information such as open ports, running services, and OS details.
3. **Concurrent Scanning**: Uses Goroutines and WaitGroups to scan multiple hosts in parallel for improved efficiency.

## Installation Requirements

### Prerequisites
- Ensure you have **Go** installed on your system. If not, install it:
  ```sh
  sudo apt update && sudo apt install golang -y  # Debian/Ubuntu
  brew install go  # macOS
  winget install -e --id golang.Go  # Windows (with Winget)
  ```
- Install **Nmap** (Required for scanning):
  ```sh
  sudo apt install nmap -y  # Debian/Ubuntu
  brew install nmap  # macOS
  winget install -e --id Nmap.Nmap  # Windows (with Winget)
  ```

## Configuration
- Before running the program, **modify the `target` variable** in `main()` to match your desired IP range:
  ```go
  target := "10.0.2.0/24"  // Change this to match your network
  ```
  This IP range must be hardcoded before running the program.

## Building and Running

1. **Clone the repository (if applicable):**
   ```sh
   git clone <repository-url>
   cd <repository-folder>
   ```
2. **Build the Go program:**
   ```sh
   go build -o nmap_scanner
   ```
3. **Run the compiled executable:**
   ```sh
   ./nmap_scanner
   ```
   (On Windows, run `nmap_scanner.exe` instead.)

## Expected Output
- The program will first list active hosts in the subnet.
- It will then perform a detailed Nmap scan on each host, displaying:
  - Open ports
  - Running services
  - OS details (if detected)

## Example Output
```sh
Active Hosts:
10.0.2.1
10.0.2.5
10.0.2.8

HOST INFO ========================================================
IP: 10.0.2.1 (router)
 Ports:
    80/open/http
    22/open/ssh
 OS:
    Linux 3.x
```

## Notes
- Running this script requires **root privileges** on some systems.
- Use responsibly and ensure you have permission before scanning any network.

## Future Development
- Would like to expand the filtering to grab more important info or any partial info grabbed by the Nmap scan
- Would like to implment changes to make it run faster

## License
This project is licensed under the MIT License.

