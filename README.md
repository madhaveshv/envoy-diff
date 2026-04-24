# envoy-diff

> A CLI tool to diff and audit environment variable changes across deployment stages.

---

## Installation

```bash
go install github.com/yourusername/envoy-diff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envoy-diff.git
cd envoy-diff
go build -o envoy-diff .
```

---

## Usage

Compare environment variables between two deployment stages:

```bash
envoy-diff --from staging.env --to production.env
```

**Example output:**

```
~ DATABASE_URL   changed
+ NEW_FEATURE_FLAG  added
- DEPRECATED_KEY    removed
```

### Flags

| Flag       | Description                          | Default  |
|------------|--------------------------------------|----------|
| `--from`   | Source environment file or stage     | required |
| `--to`     | Target environment file or stage     | required |
| `--format` | Output format (`text`, `json`, `md`) | `text`   |
| `--strict` | Exit with non-zero code on any diff  | `false`  |

### CI Integration

```bash
envoy-diff --from .env.staging --to .env.production --strict
```

Use `--strict` in CI pipelines to fail builds when unexpected environment changes are detected.

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)