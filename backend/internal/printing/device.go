package printing

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Device represents a network printer/scanner
type Device struct {
	Host     string
	Port     int
	Protocol string // "cups", "raw", "http"
	Timeout  time.Duration
}

// PrinterStatus represents the current status of a printer
type PrinterStatus struct {
	Name        string
	State       string // "idle", "printing", "error", "offline"
	Jobs        int
	Supplies    []Supply
	Capabilities []string
}

// Supply represents printer supplies (ink, paper, etc.)
type Supply struct {
	Name   string
	Level  int // percentage
	Status string
}

// ScanJob represents a scanning operation
type ScanJob struct {
	ID          string
	Status      string // "pending", "scanning", "completed", "error"
	Resolution  int    // DPI
	Format      string // "pdf", "jpeg", "png"
	Destination string // where to save the scan
}

// PrintJob represents a print job
type PrintJob struct {
	ID       string
	Status   string
	Pages    int
	Created  time.Time
	Document string
}

// PrintOptions defines printing parameters
type PrintOptions struct {
	Copies     int
	Duplex     bool
	Color      bool
	Resolution int
	PaperSize  string
	MediaType  string
}

// ScanOptions defines scanning parameters
type ScanOptions struct {
	Resolution  int    // DPI
	Format      string // "pdf", "jpeg", "png"
	ColorMode   string // "color", "grayscale", "monochrome"
	PaperSize   string
	Destination string // where to save
}

// Printer interface defines methods for printer operations
type Printer interface {
	// Device discovery
	Discover() ([]Device, error)
	
	// Connection management
	Connect(ctx context.Context, device Device) error
	Disconnect() error
	
	// Status and capabilities
	GetStatus() (*PrinterStatus, error)
	GetCapabilities() ([]string, error)
	
	// Printing operations
	Print(ctx context.Context, data []byte, options PrintOptions) error
	PrintFile(ctx context.Context, filepath string, options PrintOptions) error
	CancelJob(ctx context.Context, jobID string) error
	GetJobs(ctx context.Context) ([]PrintJob, error)
	
	// Scanning operations
	StartScan(ctx context.Context, options ScanOptions) (*ScanJob, error)
	GetScanStatus(ctx context.Context, jobID string) (*ScanJob, error)
	GetScanResult(ctx context.Context, jobID string) (io.Reader, error)
	
	// Device management
	GetSupplies() ([]Supply, error)
	GetConfiguration() (map[string]interface{}, error)
}

// Scanner interface for scanning operations
type Scanner interface {
	StartScan(ctx context.Context, options ScanOptions) (*ScanJob, error)
	GetScanStatus(ctx context.Context, jobID string) (*ScanJob, error)
	GetScanResult(ctx context.Context, jobID string) (io.Reader, error)
	GetScanCapabilities() ([]string, error)
}

// Factory function to create appropriate printer/scanner
func NewDevice(device Device) (Printer, error) {
	switch device.Protocol {
	case "cups":
		return NewCUPSPrinter(device), nil
	case "raw":
		return NewRawPrinter(device), nil
	case "http":
		return NewHTTPPrinter(device), nil
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", device.Protocol)
	}
}

// Default device configuration
func DefaultDevice() Device {
	return Device{
		Host:     "192.168.1.11",
		Port:     631,
		Protocol: "cups",
		Timeout:  30 * time.Second,
	}
} 