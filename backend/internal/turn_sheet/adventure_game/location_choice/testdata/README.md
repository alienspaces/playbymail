# Location Choice Test Data

This directory contains test fixtures for location choice turn sheet testing.

## Test Images

When OCR scanning is implemented, add test images here:

- `valid_location_choice.png` - Clean scan with all fields filled correctly
- `poor_quality_scan.png` - Low quality scan to test error handling
- `missing_code_scan.png` - Scan with missing turn sheet code
- `multiple_choices.png` - Scan with multiple location selections

## Expected Results

Add JSON files with expected scan results:

- `valid_location_choice.json` - Expected output for valid scan
- `multiple_choices.json` - Expected output for multiple choice scan

## Usage

Test images should be:
- PNG or JPEG format
- Actual scanned turn sheets or realistic simulations
- Named descriptively to indicate what scenario they test

