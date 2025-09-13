# medium-rsjs

[![🐹 Build](https://github.com/orhantugrul/medium-rsjs/actions/workflows/build.yml/badge.svg)](https://github.com/orhantugrul/medium-rsjs/actions/workflows/build.yml)

A lightweight Go API that fetches and parses Medium RSS feeds, providing clean JSON responses for easy consumption.

## What it does

- Fetches Medium RSS feeds by username
- Parses and cleans XML content (removes CDATA, fixes encoding issues)
- Returns structured JSON with normalized data
- Handles multiple date formats and normalizes to RFC3339

## Quick Start

### Local Development

1. **Clone and navigate to the project:**

   ```bash
   git clone https://github.com/orhantugrul/medium-rsjs.git
   cd medium-rsjs
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Run the server:**

   ```bash
   go run src/main.go
   ```

4. **Test the API:**
   ```bash
   curl http://localhost:8080/api/health
   curl http://localhost:8080/api/feed/@username
   ```

### Docker

1. **Build the image:**

   ```bash
   docker build -t medium-rsjs .
   ```

2. **Run the container:**

   ```bash
   docker run -p 8080:8080 medium-rsjs
   ```

3. **Test the API:**
   ```bash
   curl http://localhost:8080/api/health
   ```

## API Endpoints

### Get Medium Feed

```
GET /api/feed/{username}
```

**Example:**

```bash
curl http://localhost:8080/api/feed/@username
```

**Response:**

```json
{
  "title": "User's Medium Feed",
  "link": "https://medium.com/@username",
  "posts": [
    {
      "title": "Article Title",
      "link": "https://medium.com/@username/article-slug",
      "author": "Author Name",
      "published": "2024-01-01T12:00:00Z",
      "content": "Article content...",
      "categories": ["tech", "programming"]
    }
  ]
}
```

## Environment Variables

- `PORT` - Server port (default: 8080)
- `GIN_MODE` - Set to "release" for production
- `TRUSTED_PROXIES` - Comma-separated list of trusted proxy IPs
