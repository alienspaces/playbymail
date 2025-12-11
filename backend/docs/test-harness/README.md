# Test Harness Documentation

## Overview

The test harness is a comprehensive testing framework for the play-by-mail game system that automates test data lifecycle management. It provides declarative configuration for creating complex game scenarios with accounts, games, game instances, and all related adventure game entities.

## Key Features

- **Declarative Configuration**: Define test data using `DataConfig` structures
- **Automatic Cleanup**: Handles setup and teardown of all test data
- **Entity Relationships**: Manages complex relationships between accounts, games, locations, items, creatures, and characters
- **Workflow Simulation**: Supports turn processing, join game workflows, and scan data simulation
- **Image Handling**: Manages game images and media assets for testing
- **Transaction Management**: Provides flexible transaction handling for different test scenarios

## Quick Start

### Essential Setup Pattern

```go
func TestExample(t *testing.T) {
    // Create harness with default configuration
    th := testutil.NewTestHarness(t)
    require.NotNil(t, th, "NewTestHarness returns without error")
    
    // Setup test data
    _, err := th.Setup()
    require.NoError(t, err, "Test data setup returns without error")
    
    // Ensure cleanup happens
    defer func() {
        err = th.Teardown()
        require.NoError(t, err, "Test data teardown returns without error")
    }()
    
    // Your test logic here
    // Access created data via th.Data references
}
```

### Minimal Working Example

```go
package example_test

import (
    "testing"
    
    "github.com/stretchr/testify/require"
    
    "gitlab.com/alienspaces/playbymail/internal/harness"
    "gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

func TestMinimalExample(t *testing.T) {
    // Initialize harness
    th := testutil.NewTestHarness(t)
    require.NotNil(t, th, "NewTestHarness returns without error")
    
    // Setup test data with automatic cleanup
    _, err := th.Setup()
    require.NoError(t, err, "Setup returns without error")
    defer func() {
        err = th.Teardown()
        require.NoError(t, err, "Teardown returns without error")
    }()
    
    // Test your functionality
    accountRec, err := th.Data.GetAccountRecByRef(harness.AccountOneRef)
    require.NoError(t, err, "GetAccountRecByRef returns without error")
    require.NotEmpty(t, accountRec.ID, "Account has valid ID")
    
    gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
    require.NoError(t, err, "GetGameRecByRef returns without error")
    require.Equal(t, accountRec.ID, gameRec.AccountID, "Game belongs to account")
}
```

## Core Components

- **Testing**: Main harness struct with setup/teardown methods
- **DataConfig**: Configuration structure defining what test data to create
- **Data**: Storage for created records with reference-based lookup
- **Reference System**: String-based keys for easy entity retrieval

### Basic Test Structure

#### 1. Harness Initialization

```go
// Default configuration (creates accounts, games, locations, items)
th := testutil.NewTestHarness(t)

// Custom configuration
config := harness.DefaultDataConfig()
// Modify config as needed...
th := testutil.NewTestHarnessWithConfig(t, config)
```

#### Configuration Customization

```go
// Start with default config and customize
config := harness.DefaultDataConfig()

// Add additional accounts
config.AccountConfigs = append(config.AccountConfigs, harness.AccountConfig{
    Ref:  "custom-account",
    Name: "Custom Test Account",
    Type: "player",
})

// Customize game configuration
config.GameConfigs[0].Name = "Custom Adventure Game"
config.GameConfigs[0].Description = "A customized test game"

// Add custom locations
config.AdventureGameLocationConfigs = append(config.AdventureGameLocationConfigs, 
    harness.AdventureGameLocationConfig{
        Ref:         "custom-location",
        GameRef:     harness.GameOneRef,
        Name:        "Custom Location",
        Description: "A test location",
    })

// Create harness with custom config
th := testutil.NewTestHarnessWithConfig(t, config)
```

#### 2. Setup and Teardown

```go
// Setup creates all configured test data
_, err := th.Setup()
require.NoError(t, err, "Setup returns without error")

// Teardown cleans up all created data
defer func() {
    err = th.Teardown()
    require.NoError(t, err, "Teardown returns without error")
}()
```

#### 3. Data Access

```go
// Access created records using references
accountRec, err := th.Data.GetAccountRecByRef(harness.AccountOneRef)
require.NoError(t, err, "GetAccountRecByRef returns without error")

gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
require.NoError(t, err, "GetGameRecByRef returns without error")

locationRec, err := th.Data.GetAdventureGameLocationRecByRef(harness.GameLocationOneRef)
require.NoError(t, err, "GetAdventureGameLocationRecByRef returns without error")
```

## Documentation Structure

- [Configuration](configuration.md) - DataConfig examples and patterns
- [Workflows](workflows.md) - Game simulation and turn processing
- [Examples](examples/) - Runnable code examples

## Common Patterns

### Account Types

- **Designer**: Creates and manages games
- **Manager**: Manages game instances and processing
- **Player**: Participates in games with characters

### Game Configuration

- **Adventure Games**: Locations, items, creatures, characters with relationships
- **Game Instances**: Running games with specific players and state
- **Turn Processing**: Automated workflow simulation with scan data

### Reference System

Use string references to link entities and retrieve created records:

```go
// Define references in configuration
AccountOneRef = "account-one"
GameOneRef = "game-one"

// Retrieve created records
accountRec, err := th.Data.GetAccountRecByRef(AccountOneRef)
gameRec, err := th.Data.GetGameRecByRef(GameOneRef)
```

#### Common Reference Constants

```go
// Accounts
harness.AccountOneRef   // Designer account
harness.AccountTwoRef   // Manager account  
harness.AccountThreeRef // Player account

// Games
harness.GameOneRef // Default adventure game

// Locations
harness.GameLocationOneRef // Starting location
harness.GameLocationTwoRef // Secondary location

// Items and Creatures
harness.GameItemOneRef     // Default item
harness.GameCreatureOneRef // Default creature
harness.GameCharacterOneRef // Default character

// Game Instances
harness.GameInstanceOneRef   // Configured instance
harness.GameInstanceCleanRef // Clean instance (no data)
```

#### Available Data Getter Methods

```go
// Account records
accountRec, err := th.Data.GetAccountRecByRef(ref)

// Game records
gameRec, err := th.Data.GetGameRecByRef(ref)

// Adventure game entities
locationRec, err := th.Data.GetAdventureGameLocationRecByRef(ref)
itemRec, err := th.Data.GetAdventureGameItemRecByRef(ref)
creatureRec, err := th.Data.GetAdventureGameCreatureRecByRef(ref)
characterRec, err := th.Data.GetAdventureGameCharacterRecByRef(ref)

// Game instances
instanceRec, err := th.Data.GetGameInstanceRecByRef(ref)

// Access all records of a type
allAccounts := th.Data.AccountRecs
allGames := th.Data.GameRecs
allLocations := th.Data.AdventureGameLocationRecs
// ... etc for other entity types
```

## Advanced Usage

### Error Handling

```go
// Always check errors from harness operations
_, err := th.Setup()
require.NoError(t, err, "Setup should not fail")

// Handle missing references gracefully
accountRec, err := th.Data.GetAccountRecByRef("invalid-ref")
if err != nil {
    t.Logf("Expected error for invalid reference: %v", err)
}
```

### Parallel Test Support

```go
func TestParallelExample(t *testing.T) {
    t.Parallel() // Enable parallel execution
    
    // Each test gets isolated data via unique identifiers
    th := testutil.NewTestHarness(t)
    // ... rest of test setup
}
```

### Transaction Handling

```go
// For tests that need transaction control
th := testutil.NewTestHarness(t)
th.ShouldCommitData = false // Keep data in transaction

tx, err := th.Setup()
require.NoError(t, err, "Setup returns transaction")

// Your test logic...

// Explicitly commit or rollback
if shouldCommit {
    err = tx.Commit(context.TODO())
    require.NoError(t, err, "Commit succeeds")
} else {
    err = tx.Rollback(context.TODO())
    require.NoError(t, err, "Rollback succeeds")
}
```

## Troubleshooting

### Common Setup Issues

- **Missing dependencies**: Ensure all required packages are imported
- **Reference errors**: Use predefined constants from `harness` package
- **Cleanup failures**: Always use `defer` for teardown calls
- **Parallel conflicts**: Each test gets unique data automatically

### Debug Data Creation

```go
// Log created data for debugging
th := testutil.NewTestHarness(t)
_, err := th.Setup()
require.NoError(t, err, "Setup returns without error")

// Check what was created
t.Logf("Created accounts: %d", len(th.Data.AccountRecs))
t.Logf("Created games: %d", len(th.Data.GameRecs))
```

### Performance Tips

- Use `harness.GameInstanceCleanRef` for tests that don't need full game setup
- Minimize configuration complexity for faster test execution
- Consider custom configurations for specific test scenarios

## Best Practices

- Use unique references for all entities
- Configure relationships through reference strings
- Leverage default configurations and extend as needed
- Handle transactions appropriately for your test scenario
- Clean up properly using the teardown mechanism
