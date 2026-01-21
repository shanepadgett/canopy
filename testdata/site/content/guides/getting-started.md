---
{
  "title": "Getting Started",
  "description": "Learn how to set up your first Canopy site",
  "weight": 1
}
---

This guide walks you through setting up your first Canopy site.

## Prerequisites

- Go 1.21 or later
- A text editor

## Installation

Build Canopy from source:

```bash
go install github.com/shanepadgett/canopy/cmd/canopy@latest
```

## Create a new site

```bash
mkdir my-site && cd my-site
canopy init
```

## Build and serve

```bash
canopy serve
```

Your site is now running at <http://localhost:8080>.
