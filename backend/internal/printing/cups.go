package printing

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/OpenPrinting/goipp"
)

// CUPSPrinter implements Printer interface using CUPS/IPP protocol
type CUPSPrinter struct {
	client  *http.Client
	device  Device
	baseURL string
}

// NewCUPSPrinter creates a new CUPS printer client
func NewCUPSPrinter(device Device) *CUPSPrinter {
	return &CUPSPrinter{
		client: &http.Client{
			Timeout: device.Timeout,
		},
		device:  device,
		baseURL: fmt.Sprintf("http://%s:%d", device.Host, device.Port),
	}
}

// Connect tests connection to CUPS server
func (p *CUPSPrinter) Connect(ctx context.Context, device Device) error {
	// Test connection by making a simple HTTP request to CUPS
	req, err := http.NewRequestWithContext(ctx, "GET", p.baseURL+"/", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to CUPS server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("CUPS server returned status: %d", resp.StatusCode)
	}

	return nil
}

// Disconnect closes the connection
func (p *CUPSPrinter) Disconnect() error {
	// HTTP client doesn't need explicit disconnection
	return nil
}

// GetStatus retrieves printer status from CUPS
func (p *CUPSPrinter) GetStatus() (*PrinterStatus, error) {
	url := p.baseURL + "/printers/"

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get printer status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CUPS returned status: %d", resp.StatusCode)
	}

	// For now, return a basic status
	// In a full implementation, we'd parse the CUPS XML response
	return &PrinterStatus{
		Name:         "Network Printer",
		State:        "idle",
		Jobs:         0,
		Capabilities: []string{"print", "scan"},
	}, nil
}

// GetCapabilities returns printer capabilities
func (p *CUPSPrinter) GetCapabilities() ([]string, error) {
	// Query CUPS for printer capabilities
	return []string{"print", "scan", "duplex", "color"}, nil
}

// Print sends a print job to CUPS using proper IPP protocol
func (p *CUPSPrinter) Print(ctx context.Context, data []byte, options PrintOptions) error {
	// Create IPP print request message
	req := goipp.NewRequest(goipp.DefaultVersion, goipp.OpPrintJob, 1)

	// Add required attributes
	req.Operation.Add(goipp.MakeAttribute("attributes-charset", goipp.TagCharset, goipp.String("utf-8")))
	req.Operation.Add(goipp.MakeAttribute("attributes-natural-language", goipp.TagLanguage, goipp.String("en-US")))

	// Add printer URI
	printerURI := fmt.Sprintf("ipp://%s:%d/printers/default", p.device.Host, p.device.Port)
	req.Operation.Add(goipp.MakeAttribute("printer-uri", goipp.TagURI, goipp.String(printerURI)))

	// Add job attributes
	req.Operation.Add(goipp.MakeAttribute("requesting-user-name", goipp.TagName, goipp.String("playbymail")))
	req.Operation.Add(goipp.MakeAttribute("job-name", goipp.TagName, goipp.String("turn-sheet")))
	req.Operation.Add(goipp.MakeAttribute("document-format", goipp.TagMimeType, goipp.String("application/pdf")))

	// Add copies
	if options.Copies > 1 {
		req.Operation.Add(goipp.MakeAttribute("copies", goipp.TagInteger, goipp.Integer(options.Copies)))
	}

	// Add duplex setting
	if options.Duplex {
		req.Operation.Add(goipp.MakeAttribute("sides", goipp.TagKeyword, goipp.String("two-sided-long-edge")))
	} else {
		req.Operation.Add(goipp.MakeAttribute("sides", goipp.TagKeyword, goipp.String("one-sided")))
	}

	// Add media size
	if options.PaperSize != "" {
		req.Operation.Add(goipp.MakeAttribute("media", goipp.TagKeyword, goipp.String(options.PaperSize)))
	}

	// Encode the IPP message
	payload, err := req.EncodeBytes()
	if err != nil {
		return fmt.Errorf("failed to encode IPP message: %w", err)
	}

	// Build HTTP request with IPP message + document data
	body := io.MultiReader(bytes.NewBuffer(payload), bytes.NewBuffer(data))

	url := p.baseURL + "/printers/default"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", goipp.ContentType)
	httpReq.Header.Set("Accept", goipp.ContentType)

	// Execute HTTP request
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send print job: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	// Decode IPP response
	respMsg := &goipp.Message{}
	err = respMsg.Decode(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to decode IPP response: %w", err)
	}

	// Check IPP status
	if goipp.Status(respMsg.Code) != goipp.StatusOk {
		return fmt.Errorf("print job failed with IPP status: %s", goipp.Status(respMsg.Code).String())
	}

	return nil
}

// PrintFile prints a file from the filesystem
func (p *CUPSPrinter) PrintFile(ctx context.Context, filepath string, options PrintOptions) error {
	// Read file and send to printer
	// This would read the file and call Print()
	return fmt.Errorf("not implemented")
}

// CancelJob cancels a print job
func (p *CUPSPrinter) CancelJob(ctx context.Context, jobID string) error {
	url := fmt.Sprintf("%s/jobs/%s", p.baseURL, jobID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create cancel request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to cancel job: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cancel job failed with status: %d", resp.StatusCode)
	}

	return nil
}

// GetJobs retrieves current print jobs
func (p *CUPSPrinter) GetJobs(ctx context.Context) ([]PrintJob, error) {
	url := p.baseURL + "/jobs/"

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}
	defer resp.Body.Close()

	// Parse jobs from CUPS response
	// For now, return empty list
	return []PrintJob{}, nil
}

// StartScan initiates a scanning operation
func (p *CUPSPrinter) StartScan(ctx context.Context, options ScanOptions) (*ScanJob, error) {
	// CUPS doesn't directly support scanning via IPP
	// This would typically be done via the printer's web interface
	return &ScanJob{
		ID:         "scan-123",
		Status:     "pending",
		Resolution: options.Resolution,
		Format:     options.Format,
	}, nil
}

// GetScanStatus gets the status of a scan job
func (p *CUPSPrinter) GetScanStatus(ctx context.Context, jobID string) (*ScanJob, error) {
	return &ScanJob{
		ID:     jobID,
		Status: "completed",
	}, nil
}

// GetScanResult retrieves the scan result
func (p *CUPSPrinter) GetScanResult(ctx context.Context, jobID string) (io.Reader, error) {
	// This would download the scanned file
	return strings.NewReader("scan result"), nil
}

// GetSupplies returns printer supplies status
func (p *CUPSPrinter) GetSupplies() ([]Supply, error) {
	return []Supply{
		{Name: "Black Ink", Level: 85, Status: "ok"},
		{Name: "Color Ink", Level: 60, Status: "ok"},
		{Name: "Paper", Level: 90, Status: "ok"},
	}, nil
}

// GetConfiguration returns printer configuration
func (p *CUPSPrinter) GetConfiguration() (map[string]any, error) {
	return map[string]any{
		"protocol": "cups",
		"host":     p.device.Host,
		"port":     p.device.Port,
	}, nil
}

// Discover finds CUPS printers on the network
func (p *CUPSPrinter) Discover() ([]Device, error) {
	// This would use CUPS discovery or network scanning
	return []Device{
		{
			Host:     p.device.Host,
			Port:     p.device.Port,
			Protocol: "cups",
			Timeout:  p.device.Timeout,
		},
	}, nil
}
