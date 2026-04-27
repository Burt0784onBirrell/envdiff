# envdiff

> Compare `.env` files across environments and surface missing or mismatched variables.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff && go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <base-file> <compare-file> [compare-file...]
```

### Example

```bash
envdiff .env.example .env.production
```

**Sample output:**

```
MISSING in .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISMATCHED values:
  - LOG_LEVEL: "debug" (example) vs "info" (production)

✔ All other variables match.
```

### Flags

| Flag | Description |
|------|-------------|
| `--keys-only` | Compare keys only, ignore values |
| `--quiet` | Exit with non-zero status if differences found (CI-friendly) |
| `--json` | Output results as JSON |

---

## Why envdiff?

Keeping `.env` files in sync across environments is error-prone. `envdiff` makes it easy to catch missing secrets or misconfigured variables before they cause issues in production.

---

## Contributing

Pull requests and issues are welcome. Please open an issue before submitting large changes.

---

## License

[MIT](LICENSE) © 2024 yourusername