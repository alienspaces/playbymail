# Release Notes

## Overview

Here's what's been happening with Play by Mail — from the early foundations in September 2025, through turn sheet scanning and AI integration, to a major accounts and subscriptions update in early 2026.

---

## 🎯 What's Next

- End-game win conditions (defeat a creature, retrieve an item, or reach a location)
- Closed testing tester invite flow
- Additional demo game content

---

## 🗓️ March 2026 (late)

**Theme**: Interactive Location Objects & Game Instance Controls

### 🎯 Major Features

**Interactive Location Objects**
- Locations can now contain interactive objects — chests, levers, altars, wells, gates, and more
- Each object supports a set of actions (inspect, touch, open, break, pull, and many others)
- Interacting with an object can change its state, give you an item, damage or heal your character, reveal hidden objects, or unlock blocked paths
- Objects with multiple states (e.g. locked → unlocked → open) respond differently depending on their current condition
- Available objects and their actions appear on your location choice turn sheet alongside paths — choosing an action means you stay in place and interact rather than move
- Objects can be hidden until triggered by another effect elsewhere in the game

### 🎨 Improvements

**Per-Instance Turn Duration**
- Game managers can now set a turn duration (in hours) when creating a game instance, controlling how long players have between turns
- Turn duration can be edited inline while the instance is still in the 'created' state

**Delete Cancelled Game Instances**
- Game managers can now permanently delete instances that have been cancelled
- A delete button appears in the completed instances table for cancelled instances

**Edit Published Games**
- Game designers can now edit a published game directly from the Studio without needing to unpublish it first

**Location Link Travel Log Preview**
- The Studio location links view now shows a live preview of how the travel log entry will read for a link, including traversal description

**Creature Encounter Turn Sheet**
- The character panel in combat sheets now shows your attack and defence stats, along with the name of your equipped weapon and armour
- Creature health is now correctly capped at its maximum value

### 🐛 Bug Fixes

- Fixed location link return paths having different names to their forward-direction counterparts
- Fixed creature health displaying above the creature's maximum health value
- Fixed action radio button labels appearing below the button instead of beside it in encounter sheets

---

## 🗓️ March 2026

**Theme**: Player Game Experience & Adventure Game Foundations

### 🎯 Major Features

**Public Game Catalog**
- Players can now browse all published games from a dedicated catalog page
- Join a game directly from the catalog with a single click

**Player App & Turn Sheet Viewer**
- A dedicated player experience for viewing and completing turn sheets
- Access your turn sheets securely via a link in your notification email — no password required
- Navigate between sheets, fill them in, and submit — all from the browser

**Join Game Flow**
- Players sign in and complete a join game turn sheet to enter a game
- Choose your preferred delivery method (email or post) when joining

**Unified Link Requirements**
- Game designers can now control which paths players can see and traverse, based on items they carry, items they have equipped, or creatures they have defeated
- Locked paths display an atmospheric description so players know something blocks the way — without revealing exactly what
- Hidden paths are omitted from the turn sheet entirely until conditions are met

**Link Requirements Studio UI**
- New view in the Game Designer Studio for creating and editing link requirements
- Supports item-based and creature-based conditions for both visibility and traversal

### 🎨 Improvements

**Simplified Turn Sheet UI**
- Form data is now cached automatically as you move between sheets
- A single Submit button saves and submits all sheets at once — no manual Save or Mark Ready steps

**Delivery Method Selection**
- Players choose how they receive future turn sheets (email or post) during the join flow
- Games configured with a single delivery method skip the choice automatically

**Turn Sheet Ordering**
- Adventure game sheets now appear in a consistent order: location choice first, then inventory management
- Combat sheets will appear in the correct place once the monster encounter feature launches

**Equipment & Inventory**
- Items at a location are now auto-picked-up before equipping, matching the intent of the turn sheet
- Item actions on the inventory sheet use radio buttons instead of checkboxes, making choices mutually exclusive

**Account Contacts**
- Contact detail fields (address, phone) are now optional unless you choose postal delivery

### 🐛 Bug Fixes

- Fixed equipment slot errors when equipping non-weapon items (e.g. armour)
- Fixed turn sheet background images missing from location choice and inventory sheets
- Fixed duplicate turn processing jobs being queued for the same game instance
- Fixed the player turn sheet view being too narrow on some screens
- Fixed join game submission failing due to a missing account ID

---

## 🗓️ December 2025 – February 2026

**Theme**: Accounts, Subscriptions & Permissions

### 🎯 Major Features

**Multiple Users per Account**
- You can now have multiple users under a single account, each with their own email and postal contact details
- Each user signs in independently and has their own session

**Subscription Tiers**
- Choose between Basic and Professional subscriptions for Game Designer, Manager, and Player roles
- Basic Game Designer lets you create up to 10 games; Professional is unlimited
- Your subscription determines what you can access across the platform

**Subscription Management**
- View your active subscriptions, their tier, and limits from your account page
- See your game subscriptions and cancel Player or Manager subscriptions when needed
- Subscribe to published games directly from the management dashboard

**Permissions & Access**
- What you can see and do is now tied to your active subscriptions
- Game Designers, Game Managers, and Players each get access to their relevant features automatically

**Game Validation**
- Before creating a game instance, you can now validate your game to check it's ready — the system reports any issues with severity levels so you know what to fix

**Email Notifications**
- Turn sheets can now be delivered via email notification
- Secure turn sheet access links with time-limited tokens

**Inventory Management Turn Sheets**
- New turn sheet type for managing character inventory
- Characters now have inventory capacity that carries across turns

### 🎨 Improvements

**Game Management Dashboard**
- Browse published games available for subscription
- Subscribe as a Manager with one click
- Games you're already subscribed to are filtered out

**Game Instance Creation**
- Choose delivery methods (email, postal) when creating a game instance
- Clearer form layout with checkboxes and helpful info notices

**Sign-In & Navigation**
- Smoother redirect after signing in to the studio or admin areas
- Game management now opens directly to the turn sheets page

### 🐛 Bug Fixes

- Fixed an issue where downloading a join game turn sheet could fail
- Fixed turn sheet background images not fully loading before PDF generation
- Fixed game creation validation in some edge cases

---

## 🗓️ November 2025

**Theme**: Turn Sheet Art & Game Subscriptions

### 🎯 Major Features

**Turn Sheet Background Images**
- Upload custom background artwork for your join game turn sheets at the game level
- Upload location-specific backgrounds for location choice turn sheets
- If a location doesn't have its own image, it falls back to the game-level background
- Preview your turn sheets as PDFs before sending them out
- Supports WebP, PNG, and JPEG (up to 1MB)

**Game Subscription & Join Workflow**
- Players can now subscribe to your game and go through an approval process
- Generate and download join game turn sheets for approved subscribers

**AI-Powered Turn Sheet Scanning**
- Scanned turn sheets are now read using AI vision instead of basic OCR, giving much more accurate results

**Starting Locations**
- Define which locations players can start in
- Choose starting locations when designing your adventure game

### 🎨 Improvements

**Game Designer Studio**
- Upload turn sheet images directly from the game and location edit screens
- Preview turn sheets from the resource table
- Consistent sidebar layout with a game context menu

**Game Management**
- Download and upload turn sheets from the management interface
- Cleaner resource table layout for game instances

**Layout & Navigation**
- Consistent sidebar layouts across pages
- Inline action buttons instead of dropdown menus
- Fixed modal display issues on mobile devices

---

## 🗓️ October 2025

**Theme**: Turn Sheet Scanning & Production Launch

### 🎯 Major Features

**Turn Sheet Scanning**
- Scan filled-in turn sheets and have the system read player responses automatically
- Handles handwriting quirks and common scan artifacts
- Supports location choice turn sheets with checkbox detection

**Turn Sheet PDF Generation**
- Generate polished PDF turn sheets for your players
- Consistent layout with turn sheet codes positioned for easy scanning
- Join game turn sheets ready for distribution

### 🎨 Improvements

- Removed the hero decoration from the home page for a cleaner look
- Added support email (support@playbymail.games) to the footer

### 🐛 Bug Fixes

- Fixed various scanning accuracy issues
- Fixed deployment and build issues

---

## 🗓️ September 2025

**Theme**: Getting Started

### 🎯 Major Features

- Turn sheet processing system — the foundation for generating, distributing, and scanning turn sheets
- Character-based gameplay with game instance support

### 🎨 Improvements

- Smoother development setup with consistent tooling
- Project documentation and getting started guide

---
