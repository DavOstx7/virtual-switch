# Virtual Switch

A Go-based implementation of a virtual network switch. This project aims to provide basic functionality for network packet switching within a virtualized environment.

## Features

- **Virtual Switch:** Simulate network traffic switching between different interfaces.
- **Packet Forwarding:** Forward packets to specified interfaces.
- **Customizable:** Easily configure network behavior as per project needs.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/DavOstx7/virtual-switch.git
   cd virtual-switch
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build the application:
   ```bash
   go build
   ```

## Usage

To run the virtual switch:

```bash
go run main.go
```

## Configuration

You can configure the network interfaces and packet forwarding behavior directly in the source code.

## Future Enhancements

In the future, we plan to add the following features to enhance the functionality of the virtual switch:

- **Improved Frame Source Integration:** Expand support for custom frame sources, including deeper integration with `gopacket` for more versatile packet handling.
- **Traffic Shaping:** Implement traffic shaping features to control the flow of packets, simulate bandwidth limits, and prioritize traffic.
- **Packet Loss Simulation:** Add functionality to simulate packet loss, useful for testing network resilience and recovery mechanisms.
- **GUI for Configuration:** A graphical interface to easily configure network interfaces and forwarding rules.
- **Enhanced Packet Analysis:** Advanced tools for packet inspection and statistics tracking.

Stay tuned for updates!
