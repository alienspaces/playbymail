# Turn Sheet Testing Guide

## Overview

Turn sheet tests are organized to test both **generation** and **scanning** functionality for each turn sheet type. Tests are implemented in separate files for each turn sheet type (e.g., `adventure_game_location_choice_test.go`).

## Current Structure

```
turn_sheet/
  adventure_game_location_choice.go          # Location choice processor
  adventure_game_location_choice_test.go      # Location choice tests
  adventure_game_join_game.go                 # Join game processor
  turn_sheet.go                              # Base processor
  turn_sheet_processor.go                     # Document processor interface
  types.go                                    # Shared types
  testdata/
    adventure_game_location_choice_turn_sheet_scan.jpg  # Test scan image
    adventure_game_location_choice_turn_sheet.pdf        # Generated PDF
```

## Test Structure Pattern

Each turn sheet processor has two main test functions:

### 1. GenerateTurnSheet Tests

**Test function**: `Test{ProcessorName}_GenerateTurnSheet`

Tests PDF generation from turn sheet data structures:
- Validates data structure (required fields, format)
- Generates PDF from template data
- Uses `testutil.NewDefaultDependencies(t)` for test harness
- Returns PDF data for output

Example:
```go
func TestLocationChoiceProcessor_GenerateTurnSheet(t *testing.T) {
    tests := []struct {
        name        string
        data        any
        expectError bool
        errorMsg    string
    }{
        {
            name: "given valid LocationChoiceData when generating turn sheet then PDF is generated successfully",
            data: &turn_sheet.LocationChoiceData{
                TurnSheetTemplateData: turn_sheet.TurnSheetTemplateData{
                    GameName:      convert.Ptr("Test Adventure"),
                    GameType:      convert.Ptr("adventure"),
                    TurnNumber:    convert.Ptr(1),
                    AccountName:   convert.Ptr("Test Player"),
                    TurnSheetCode: convert.Ptr("TEST123"),
                },
                LocationName: "Starting Location",
                LocationOptions: []turn_sheet.LocationOption{...},
            },
            expectError: false,
        },
    }
    // ... test implementation
}
```

### 2. ScanTurnSheet Tests

**Test function**: `Test{ProcessorName}_ScanTurnSheet`

Tests OCR extraction and parsing from scanned images:
- Extracts text from scanned images using OCR
- Parses player choices from extracted text
- Validates scanned data structure
- Uses test images from `testdata/` directory

Example:
```go
func TestLocationChoiceProcessor_ScanTurnSheet(t *testing.T) {
    tests := []struct {
        name                  string
        imageDataFn           func() ([]byte, error)
        sheetData             any
        expectError           bool
        expectedTurnSheetCode string
        expectedChoices       []string
    }{
        {
            name: "given real scanned turn sheet image when scanning then turn sheet code and location choices are extracted correctly",
            imageDataFn: func() ([]byte, error) {
                return os.ReadFile("testdata/adventure_game_location_choice_turn_sheet_scan.jpg")
            },
            sheetData: map[string]any{
                "locations": []any{...},
            },
            expectError:           false,
            expectedTurnSheetCode: "ABC123XYZ",
            expectedChoices:       []string{"dark_tower"},
        },
    }
    // ... test implementation
}
```

### 3. PDF Generation for Manual Testing

**Test function**: `TestGenerate{Type}PDFForPrinting`

Generates PDF files for physical testing and printing:
- Creates realistic test data
- Generates PDF to testdata directory
- Provides instructions for printing and scanning
- Used for OCR development and validation

## Running Tests

### Run All Backend Tests
```bash
cd playbymail
./tools/test-backend
```

### Run Only Turn Sheet Tests
```bash
cd playbymail/backend
go test ./internal/turn_sheet/... -v
```

### Run Specific Test
```bash
cd playbymail/backend
go test ./internal/turn_sheet -run TestLocationChoiceProcessor_GenerateTurnSheet -v
```

### Generate PDF for Manual Testing
```bash
cd playbymail/backend
go test ./internal/turn_sheet -run TestGenerateLocationChoicePDFForPrinting -v
```

## Test Dependencies

Turn sheet tests require:

1. **Test harness**: Uses `testutil.NewDefaultDependencies(t)` for logger and setup
2. **Config**: Loads `TemplatesPath` from `../../templates`
3. **OCR**: Requires Tesseract OCR to be installed for scanning tests
4. **Test images**: Scanned PDFs saved in `testdata/` directory

## Processor Implementation Pattern

Each turn sheet type implements the processor pattern:

```go
type LocationChoiceProcessor struct {
    *BaseProcessor
}

// NewLocationChoiceProcessor creates a new location choice processor
func NewLocationChoiceProcessor(l logger.Logger, cfg *config.Config) *LocationChoiceProcessor {
    return &LocationChoiceProcessor{
        BaseProcessor: NewBaseProcessor(l, cfg),
    }
}

// GenerateTurnSheet generates a location choice turn sheet PDF
func (p *LocationChoiceProcessor) GenerateTurnSheet(ctx context.Context, l logger.Logger, data any) ([]byte, error) {
    // Validate data
    // Generate PDF using generator
}

// ScanTurnSheet scans a location choice turn sheet and extracts player choices
func (p *LocationChoiceProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, imageData []byte, sheetData any) (any, error) {
    // Extract text using OCR
    // Parse choices from text
}
```

## Base Processor Methods

The `BaseProcessor` provides common functionality:

- `GeneratePDF(ctx, templatePath, data)` - Generate PDF from template
- `ExtractTextFromImage(ctx, imageData)` - Extract text using OCR
- `ParseTurnSheetCodeFromImage(ctx, imageData)` - Parse turn sheet code from scanned image
- `ValidateBaseTemplateData(data)` - Validate common template fields

## Test Coverage

### Current Status

✅ **Location Choice Turn Sheet**:
- GenerateTurnSheet tests implemented
- ScanTurnSheet tests implemented
- OCR integration working
- PDF generation working
- Test images available

⏳ **Join Game Turn Sheet**:
- Processor implemented (`adventure_game_join_game.go`)
- No tests yet
- Needs test implementation

## Adding New Turn Sheet Types

When creating a new turn sheet type:

1. **Create processor file**: `{type}_processor.go`
   - Implement `GenerateTurnSheet` method
   - Implement `ScanTurnSheet` method
   - Add data structures for generation and scanning
   - Extend `BaseProcessor` for common functionality

2. **Create test file**: `{type}_test.go`
   - Test `GenerateTurnSheet` with various data scenarios
   - Test `ScanTurnSheet` with scanned images
   - Add PDF generation test for manual testing
   - Use table-driven tests

3. **Add test data**: Update `testdata/` directory
   - Create scanned images for testing
   - Generate PDFs for printing
   - Document test data format

4. **Update types**: Add to `types.go`
   - Add constant for new turn sheet type
   - Define data structures

## Best Practices

1. **Use table-driven tests** - One test function per method with multiple scenarios
2. **Test harness** - Use `testutil.NewDefaultDependencies(t)` for setup
3. **Clear test names** - Use "given/when/then" format for clarity
4. **Validation** - Test both success and error cases
5. **OCR patterns** - Document common OCR patterns for scanned text
6. **Test images** - Use real scanned images for OCR testing
7. **Manual testing** - Generate PDFs for physical testing workflow

## OCR Pattern Matching

The scanner uses regex patterns to extract choices from OCR text:
- Handles OCR artifacts (Q/, O/, Sf for checkboxes)
- Matches against expected location/item names
- Validates against sheet data
- Filters false matches (submit, deadline, turn, etc.)

## Current Status

- ✅ Test structure established
- ✅ Location choice tests implemented
- ✅ OCR integration working
- ✅ PDF generation working
- ⏳ Join game tests pending
- ⏳ Additional turn sheet types pending
