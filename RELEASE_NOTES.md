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

**Theme**: AI Integration, Game Subscriptions & Turn Sheet Art

### 🎯 Major Features

**Turn Sheet Background Images**
- Game-level background image upload for join game turn sheets
- Location-level background image upload for location choice turn sheets
- Location images fall back to game-level image if not set
- Turn sheet preview with PDF display in modal
- Image validation (WebP/PNG/JPEG, 1MB max, dimension constraints)

**Game Subscription & Join Workflow**
- Game subscription processing with approval workflow
- Join game turn sheet generation and download
- Subscription-based access control (Designer/Manager roles)

**Turn Sheet Scanning**
- Replaced OCR-based scanning with OpenAI vision integration

**Starting Locations**
- Starting location concept for adventure games
- Starting location selection in frontend studio

### 🎨 User Interface

**Game Designer Studio**
- Turn sheet image upload in game and location edit modals
- Turn sheet preview button in resource tables
- Shared sidebar layout with game context menu

**Game Management**
- Turn sheet download and upload interface
- Resource table layout for game instances

**Layout & Navigation**
- New SidebarLayout component for consistent layouts
- Inline action buttons replacing dropdown menus
- Modal z-index fixes for mobile devices

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

