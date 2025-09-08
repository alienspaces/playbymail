# PDF Generator System

This package provides a comprehensive PDF generation system for play-by-mail games, organized by game type and template structure.

## Architecture

The generator system follows a modular architecture:

- **Core Generator** (`generator.go`): Contains the main PDF generation logic using [Go's html/template](https://pkg.go.dev/html/template) and wkhtmltopdf
- **Game-Specific Generators**: Each game type has its own generator package (e.g., `adventure_game_generator/`)
- **Templates**: Organized by game type in `[game_type]/templates/` directories

## Usage

### Basic Setup

```go
import (
    "gitlab.com/alienspaces/playbymail/internal/generator"
    "gitlab.com/alienspaces/playbymail/internal/generator/adventuregamegenerator"
)

// Create PDF generator
pdfGen := generator.NewPDFGenerator("templates", "output")

// Create game-specific generator
adventureGen := adventuregamegenerator.NewAdventureGameLocationGenerator(logger, domain, pdfGen)
```

### Generating PDFs

```go
// Prepare template data
data := generator.TemplateData{
    // Core context records - always present for every turn
    AccountRec:      accountRecord,      // *account_record.Account
    GameInstanceRec: gameInstanceRecord, // *game_record.GameInstance
    GameRec:         gameRecord,         // *game_record.Game
    
    // Background images
    BackgroundTop:    "images/forest_top.jpg",
    BackgroundMiddle: "images/forest_middle.jpg",
    BackgroundBottom: "images/forest_bottom.jpg",
    
    // Game-specific data - location choice only
    GameData: map[string]any{
        "location_choice": map[string]any{
            "available_locations": []map[string]any{
                {
                    "id":          "mystic_grove",
                    "name":        "Mystic Grove",
                    "description": "A peaceful clearing with ancient trees",
                },
                {
                    "id":          "crystal_caverns",
                    "name":        "Crystal Caverns",
                    "description": "Glowing crystals light the underground passages",
                },
            },
        },
    },
}

// Generate location choice turn sheet PDF
pdfData, err := adventureGen.GenerateLocationChoiceTurnSheet(ctx, data)
if err != nil {
    log.Fatal(err)
}

// Save to file
err = adventureGen.GenerateLocationChoiceTurnSheetToFile(ctx, data, "location_choice_turn_sheet.pdf")
if err != nil {
    log.Fatal(err)
}
```

## Template System

### Template Structure

Templates are organized by game type:

```
internal/generator/
├── adventure_game_generator/
│   └── templates/
│       ├── base.template          # Base template with background structure
│       ├── base.css               # CSS for styling
│       └── location_choice.template  # Location choice turn sheet
└── [other_game_types]/
    └── templates/
        └── [game_specific_templates]
```

### Template Features

- **Background Images**: Support for top, middle, bottom, and overlay images
- **Data-Driven Content**: Dynamic content populated from game data
- **Form Elements**: Checkboxes, radio buttons, text areas for player input
- **Responsive Design**: Optimized for PDF output
- **Game-Specific Layouts**: Different templates for different game mechanics

### Template Data Structure

```go
type TemplateData struct {
    // Core context records - always present for every turn
    AccountRec      *account_record.Account
    GameInstanceRec *game_record.GameInstance
    GameRec         *game_record.Game

    // Background images
    BackgroundTop    string
    BackgroundMiddle string
    BackgroundBottom string

    // Content sections
    Header  map[string]any
    Content map[string]any
    Footer  map[string]any

    // Game-specific data
    GameData map[string]any
}
```

## Adventure Game Templates

### Location Choice Template

For players to choose their next location in adventure games:

- Available locations with descriptions
- Character actions at the chosen location
- Basic inventory management
- Notes section for player input

This is the core template for adventure game turn sheets, allowing players to navigate through the game world.

## Background Image System

The template system supports layered background layouts:

- **Top Background**: 33.33% height at the top (z-index: 0)
- **Middle Background**: Full height overlay (z-index: 1) - can overlay top and bottom
- **Bottom Background**: 33.33% height at the bottom (z-index: 0)

The middle background naturally overlays the top and bottom images based on its size and positioning.

CSS classes:
- `.background-top`
- `.background-middle`
- `.background-bottom`

## Dependencies

- **wkhtmltopdf**: For HTML to PDF conversion
- **Go html/template**: For template processing
- **Go standard library**: For file operations and context handling

## Testing

Run tests with:

```bash
go test ./internal/generator/... -v
```

Tests verify:
- Template file existence
- Generator function calls
- PDF generation (with expected errors for missing templates in test environment)

## Examples

See `example_usage.go` for comprehensive examples of:
- Setting up generators
- Preparing template data
- Generating different types of PDFs
- Error handling

## Future Enhancements

- Support for additional game types
- More template customization options
- Image optimization
- Template validation
- Performance improvements
