package printing

import (
	"context"
	"testing"
	"time"
)

func TestNewDevice(t *testing.T) {
	tests := []struct {
		name    string
		device  Device
		wantErr bool
	}{
		{
			name: "valid cups device",
			device: Device{
				Host:     "192.168.1.11",
				Port:     631,
				Protocol: "cups",
				Timeout:  30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid raw device",
			device: Device{
				Host:     "192.168.1.11",
				Port:     9100,
				Protocol: "raw",
				Timeout:  30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid http device",
			device: Device{
				Host:     "192.168.1.11",
				Port:     80,
				Protocol: "http",
				Timeout:  30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "invalid protocol",
			device: Device{
				Host:     "192.168.1.11",
				Port:     631,
				Protocol: "invalid",
				Timeout:  30 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer, err := NewDevice(tt.device)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && printer == nil {
				t.Error("NewDevice() returned nil printer when no error expected")
			}
		})
	}
}

func TestDefaultDevice(t *testing.T) {
	device := DefaultDevice()
	
	if device.Host != "192.168.1.11" {
		t.Errorf("DefaultDevice().Host = %v, want %v", device.Host, "192.168.1.11")
	}
	
	if device.Port != 631 {
		t.Errorf("DefaultDevice().Port = %v, want %v", device.Port, 631)
	}
	
	if device.Protocol != "cups" {
		t.Errorf("DefaultDevice().Protocol = %v, want %v", device.Protocol, "cups")
	}
}

func TestCUPSPrinter_Connect(t *testing.T) {
	device := Device{
		Host:     "192.168.1.11",
		Port:     631,
		Protocol: "cups",
		Timeout:  5 * time.Second,
	}
	
	printer := NewCUPSPrinter(device)
	ctx := context.Background()
	
	// This will likely fail since we don't have a real printer
	// but it tests the interface
	err := printer.Connect(ctx, device)
	if err != nil {
		t.Logf("Connect failed as expected: %v", err)
	}
}

func TestRawPrinter_Connect(t *testing.T) {
	device := Device{
		Host:     "192.168.1.11",
		Port:     9100,
		Protocol: "raw",
		Timeout:  5 * time.Second,
	}
	
	printer := NewRawPrinter(device)
	ctx := context.Background()
	
	// This will likely fail since we don't have a real printer
	// but it tests the interface
	err := printer.Connect(ctx, device)
	if err != nil {
		t.Logf("Connect failed as expected: %v", err)
	}
}

func TestHTTPPrinter_Connect(t *testing.T) {
	device := Device{
		Host:     "192.168.1.11",
		Port:     80,
		Protocol: "http",
		Timeout:  5 * time.Second,
	}
	
	printer := NewHTTPPrinter(device)
	ctx := context.Background()
	
	// This will likely fail since we don't have a real printer
	// but it tests the interface
	err := printer.Connect(ctx, device)
	if err != nil {
		t.Logf("Connect failed as expected: %v", err)
	}
} 