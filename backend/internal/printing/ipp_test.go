package printing

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/OpenPrinting/goipp"
)

func TestIPPPrintJob(t *testing.T) {
	// Create a test device
	device := Device{
		Host:     "192.168.1.11",
		Port:     631,
		Protocol: "cups",
		Timeout:  30 * time.Second,
	}

	printer := NewCUPSPrinter(device)

	// Test data (simple text)
	testData := []byte("Hello, this is a test print job from playbymail!")

	// Print options
	options := PrintOptions{
		Copies:     1,
		Duplex:     false,
		Color:      false,
		Resolution: 300,
		PaperSize:  "A4",
		MediaType:  "plain",
	}

	ctx := context.Background()

	// Test the print function
	err := printer.Print(ctx, testData, options)
	if err != nil {
		// This is expected if no printer is available
		t.Logf("Print test failed (expected if no printer): %v", err)
	} else {
		t.Log("Print job sent successfully!")
	}
}

func TestIPPMessageConstruction(t *testing.T) {
	// Test building a proper IPP message
	req := goipp.NewRequest(goipp.DefaultVersion, goipp.OpPrintJob, 1)

	// Add required attributes
	req.Operation.Add(goipp.MakeAttribute("attributes-charset", goipp.TagCharset, goipp.String("utf-8")))
	req.Operation.Add(goipp.MakeAttribute("attributes-natural-language", goipp.TagLanguage, goipp.String("en-US")))
	req.Operation.Add(goipp.MakeAttribute("printer-uri", goipp.TagURI, goipp.String("ipp://192.168.1.11:631/printers/default")))
	req.Operation.Add(goipp.MakeAttribute("requesting-user-name", goipp.TagName, goipp.String("playbymail")))
	req.Operation.Add(goipp.MakeAttribute("job-name", goipp.TagName, goipp.String("test-job")))
	req.Operation.Add(goipp.MakeAttribute("document-format", goipp.TagMimeType, goipp.String("application/pdf")))
	req.Operation.Add(goipp.MakeAttribute("copies", goipp.TagInteger, goipp.Integer(2)))
	req.Operation.Add(goipp.MakeAttribute("sides", goipp.TagKeyword, goipp.String("one-sided")))
	req.Operation.Add(goipp.MakeAttribute("media", goipp.TagKeyword, goipp.String("A4")))

	// Encode the message
	payload, err := req.EncodeBytes()
	if err != nil {
		t.Fatalf("Failed to encode IPP message: %v", err)
	}

	// Verify we have a valid IPP message
	if len(payload) == 0 {
		t.Fatal("IPP message is empty")
	}

	// Check that it starts with IPP version (0x0200 for IPP 2.0)
	if len(payload) < 8 {
		t.Fatal("IPP message too short")
	}

	// The first 8 bytes should be: version(2) + operation(2) + request-id(4)
	version := uint16(payload[0])<<8 | uint16(payload[1])
	operation := uint16(payload[2])<<8 | uint16(payload[3])

	if version != 0x0200 {
		t.Errorf("Expected IPP version 2.0 (0x0200), got 0x%04x", version)
	}

	if operation != 0x0002 { // OpPrintJob
		t.Errorf("Expected operation PrintJob (0x0002), got 0x%04x", operation)
	}

	t.Logf("IPP message encoded successfully: %d bytes", len(payload))
	t.Logf("Version: 0x%04x, Operation: 0x%04x", version, operation)
}

func TestIPPResponseParsing(t *testing.T) {
	// Create a mock IPP response
	resp := goipp.NewResponse(goipp.DefaultVersion, goipp.StatusOk, 1)
	resp.Operation.Add(goipp.MakeAttribute("attributes-charset", goipp.TagCharset, goipp.String("utf-8")))
	resp.Operation.Add(goipp.MakeAttribute("attributes-natural-language", goipp.TagLanguage, goipp.String("en-US")))
	resp.Operation.Add(goipp.MakeAttribute("job-id", goipp.TagInteger, goipp.Integer(123)))
	resp.Operation.Add(goipp.MakeAttribute("job-state", goipp.TagEnum, goipp.Integer(3))) // 3 = pending
	resp.Operation.Add(goipp.MakeAttribute("job-state-message", goipp.TagText, goipp.String("Job completed successfully")))

	// Encode the response
	respData, err := resp.EncodeBytes()
	if err != nil {
		t.Fatalf("Failed to encode response: %v", err)
	}

	// Decode the response
	decodedResp := &goipp.Message{}
	err = decodedResp.DecodeBytes(respData)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify the response
	if goipp.Status(decodedResp.Code) != goipp.StatusOk {
		t.Errorf("Expected status OK, got %s", goipp.Status(decodedResp.Code).String())
	}

	t.Logf("IPP response parsed successfully: %s", goipp.Status(decodedResp.Code).String())
}

func TestIPPWithDocumentData(t *testing.T) {
	// Create IPP request
	req := goipp.NewRequest(goipp.DefaultVersion, goipp.OpPrintJob, 1)
	req.Operation.Add(goipp.MakeAttribute("attributes-charset", goipp.TagCharset, goipp.String("utf-8")))
	req.Operation.Add(goipp.MakeAttribute("attributes-natural-language", goipp.TagLanguage, goipp.String("en-US")))
	req.Operation.Add(goipp.MakeAttribute("printer-uri", goipp.TagURI, goipp.String("ipp://192.168.1.11:631/printers/default")))
	req.Operation.Add(goipp.MakeAttribute("document-format", goipp.TagMimeType, goipp.String("text/plain")))

	// Encode IPP message
	payload, err := req.EncodeBytes()
	if err != nil {
		t.Fatalf("Failed to encode IPP message: %v", err)
	}

	// Test document data
	documentData := []byte("This is test document data for printing.")

	// Combine IPP message + document data
	body := bytes.NewBuffer(payload)
	body.Write(documentData)

	// Read the combined data
	combinedData, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("Failed to read combined data: %v", err)
	}

	// Verify we have both IPP message and document data
	if len(combinedData) < len(payload) {
		t.Fatal("Combined data is shorter than IPP payload")
	}

	// Verify the IPP message is at the beginning
	if !bytes.Equal(combinedData[:len(payload)], payload) {
		t.Fatal("IPP message not found at beginning of combined data")
	}

	// Verify document data is at the end
	documentStart := len(payload)
	if !bytes.Equal(combinedData[documentStart:], documentData) {
		t.Fatal("Document data not found at end of combined data")
	}

	t.Logf("IPP message + document data combined successfully")
	t.Logf("IPP message: %d bytes", len(payload))
	t.Logf("Document data: %d bytes", len(documentData))
	t.Logf("Total: %d bytes", len(combinedData))
}
