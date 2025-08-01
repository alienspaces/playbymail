package printing

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Example demonstrates how to use the printing package
func Example() {
	// Create a device configuration for your printer
	device := Device{
		Host:     "192.168.1.11", // Your printer's IP address
		Port:     631,            // CUPS port
		Protocol: "cups",         // Protocol: "cups", "raw", or "http"
		Timeout:  30 * time.Second,
	}

	// Create a printer instance
	printer, err := NewDevice(device)
	if err != nil {
		log.Fatalf("Failed to create printer: %v", err)
	}

	ctx := context.Background()

	// Connect to the printer
	err = printer.Connect(ctx, device)
	if err != nil {
		log.Printf("Failed to connect to printer: %v", err)
		log.Println("This is expected if no printer is available")
		return
	}
	defer printer.Disconnect()

	// Get printer status
	status, err := printer.GetStatus()
	if err != nil {
		log.Printf("Failed to get printer status: %v", err)
	} else {
		fmt.Printf("Printer status: %s\n", status.State)
		fmt.Printf("Printer capabilities: %v\n", status.Capabilities)
	}

	// Example: Print some data
	printData := []byte("Hello, this is a test print job!")
	printOptions := PrintOptions{
		Copies:     1,
		Duplex:     false,
		Color:      false,
		Resolution: 300,
		PaperSize:  "A4",
		MediaType:  "plain",
	}

	err = printer.Print(ctx, printData, printOptions)
	if err != nil {
		log.Printf("Failed to print: %v", err)
	} else {
		fmt.Println("Print job sent successfully")
	}

	// Example: Start a scan (if supported)
	scanOptions := ScanOptions{
		Resolution:  300,
		Format:      "pdf",
		ColorMode:   "color",
		PaperSize:   "A4",
		Destination: "/tmp/scan.pdf",
	}

	scanJob, err := printer.StartScan(ctx, scanOptions)
	if err != nil {
		log.Printf("Failed to start scan: %v", err)
	} else {
		fmt.Printf("Scan job started: %s\n", scanJob.ID)

		// Check scan status
		status, err := printer.GetScanStatus(ctx, scanJob.ID)
		if err != nil {
			log.Printf("Failed to get scan status: %v", err)
		} else {
			fmt.Printf("Scan status: %s\n", status.Status)
		}
	}

	// Get printer supplies
	supplies, err := printer.GetSupplies()
	if err != nil {
		log.Printf("Failed to get supplies: %v", err)
	} else {
		fmt.Println("Printer supplies:")
		for _, supply := range supplies {
			fmt.Printf("  %s: %d%% (%s)\n", supply.Name, supply.Level, supply.Status)
		}
	}
}
