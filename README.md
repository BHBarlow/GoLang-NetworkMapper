# Nmap Network Scanner in Go

## Description
This Go program automates network scanning using Nmap. It performs the following tasks:

1. **Ping Scan (`nmap -sP`)**: Identifies active hosts in the specified subnet.
2. **Full Nmap Scan (`nmap -A`)**: Runs a comprehensive scan on each active host, gathering information such as open ports, running services, and OS details.
3. **Concurrent Scanning**: Uses Goroutines and WaitGroups to scan multiple hosts in parallel for improved efficiency.
4. **Configuration File Support**: Uses a `config.json` file to store the target subnet and Nmap flags, reducing the need to modify the source code.

## Quick Start (Pre-Configured Binary)
If you have the pre-built binary available, simply run:
```sh
./pong  # Linux/macOS
```
This will automatically use the settings from `config.json` for scanning.

## Installation Requirements for modifications 

### Prerequisites
- Ensure you have **Go** installed on your system. If not, install it:
  ```sh
  sudo apt update && sudo apt install golang -y  # Debian/Ubuntu
  brew install go  # macOS
  ```
- Install **Nmap** (Required for scanning):
  ```sh
  sudo apt install nmap -y  # Debian/Ubuntu
  brew install nmap  # macOS
  ```

## Configuration
The program uses a `config.json` file to store scan settings. Create or modify `config.json` in the same directory as the executable:
```json
{
  "target": "192.168.1.0/24",
  "nmap_flags": ["-A", "-T5", "--min-rate", "1000", "--max-retries", "1"]
}
```
### Configuration Options
- `target`: The IP address or subnet to scan.
- `nmap_flags`: An array of Nmap command-line options to customize scanning.

## Building and Running

1. **Clone the repository (if applicable):**
   ```sh
   git clone <repository-url>
   cd <repository-folder>
   ```
2. **Build the Go program:**
   ```sh
   go build -o pong
   ```
3. **Run the compiled executable:**
   ```sh
   ./pong  # Linux/macOS
   ```

### Overriding Configuration via CLI
You can override the `target` from `config.json` using the command-line flag:
```sh
go run pong.go -target 192.168.1.10
```
If no CLI argument is provided, the program defaults to the `config.json` settings.

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
- Running this script may require **root privileges** on some systems.
- Use responsibly and ensure you have permission before scanning any network.

## License
This project is licensed under the MIT License.

