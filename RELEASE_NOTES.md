# Release Notes

## Overview

These release notes cover development focused on building a complete turn sheet processing system. The work progressed from foundational architecture in September, to OCR integration in October, to AI-powered scanning in November.

---

## 🎯 What's Next

- Additional turn sheet types (inventory management)
- Real-time turn processing workflows
- Expanded test coverage
- An additional game type!

---

## 🗓️ November 2025

**Theme**: AI Integration, Game Subscriptions & Management UI

### 🎯 Major Features

**Game Subscription & Join Workflow**
- Game subscription processing job worker
- Subscription approval endpoint with email templates
- Join game processor with complete workflow
- Download endpoint for join game turn sheets
- Renamed Collaborator subscription type to Designer for clarity

**Subscription-Based Access Control**
- Games API supports filtering by subscription type
- Studio shows only games where user is Designer
- Management shows only games where user is Manager
- Removed redundant game_administration table

**Starting Locations**
- Starting location concept for adventure games
- Starting location selection in frontend studio

**Turn Sheet Scanning**
- Replaced OCR-based scanning with OpenAI vision integration

**Game Images & Turn Sheet Backgrounds**
- Game image upload API for turn sheet backgrounds
- Turn sheet PDF generation with custom background images
- Session refresh endpoint for extended sessions

### 🎨 User Interface

**Home Page**
- Game genres section added

**Game Designer Studio**
- Subscription-based access control (Designer role)
- Clickable table rows for quick editing
- Prominent game selection button
- Shared sidebar layout with game context menu
- Turn sheet image upload in game edit modal
- Turn sheet preview modal with PDF display

**Game Management**
- Subscription-based access control (Manager role)
- Turn sheet download and upload interface
- Resource table layout for game instances
- Game context menu in sidebar navigation
- Clean flat layout matching Studio design

**Layout & Navigation**
- New SidebarLayout component for consistent layouts
- Game context appears when a game is selected
- Improved navigation with direct edit links in tables
- Replaced dropdown action menus with inline action buttons
- Consistent button sizing across all views
- Unified spacing and layout between Studio and Management sections
- Added h3 title level support with appropriate spacing

---

## 🗓️ October 2025

**Theme**: OCR Integration & Production Deployment

### 🎯 Major Features

**Turn Sheet Scanning & OCR**
- Complete turn sheet scanning with Tesseract OCR integration
- Location choice processor with intelligent pattern matching
- Base processor architecture for extensible turn sheet types
- Handles OCR artifacts (Q/, O/, Sf patterns, cents symbol, checkboxes)
- Priority-based pattern matching with counting heuristics

**Data Structures & PDF Generation**
- Turn sheet data structures with type-safe structs
- Join game turn sheet structure and template
- Fixed footer layout for consistent turn sheet code positioning
- PDF layout optimizations and template improvements

### 🚀 Infrastructure & Deployment

**Heroku & CI/CD**
- Tesseract OCR dependencies for Heroku buildpacks
- Fixed OCR initialization failures in CI/CD pipelines
- Simplified Aptfile and added dedicated Tesseract buildpack
- OCR dependencies added to GitHub Actions
- Enhanced debugging and CI triggers

### 🧪 Testing Infrastructure

**Test Framework**
- Restructured harness for realistic turn testing
- Comprehensive test framework for turn sheet generation and scanning
- Real test image assets for OCR validation
- Created testing guide (TESTING.md) with best practices

### 🎨 User Interface

**Frontend Updates**
- Removed hero decoration from home page
- Added support email (support@playbymail.games) to footer

### 🔧 Technical Improvements

**Code Quality**
- Refactored scanning and generating architecture
- Improved validation, error handling, and code organization
- Code formatting cleanup across codebase

### 🐛 Bug Fixes

- Fixed OCR text cleaning and regex patterns
- Fixed Heroku build dependencies and Tesseract OCR initialization
- Fixed import cycles and test data handling

---

## 🗓️ September 2025

**Theme**: Foundation & Architecture

### 🏗️ Architecture Refactoring

- Interface-based turn processing architecture
- Character-centric approach with game instance record passing
- Established base architecture for turn sheet processing

### 🛠️ Developer Experience

- Node.js version management with `.nvmrc` file
- Updated documentation and development environment
- Pinned npm packages for consistent builds

---

*Generated from git history and codebase analysis*

