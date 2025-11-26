# Go Style Guide for pr-faq-validator

This document defines the Go coding standards for this repository.
Standards are derived from Effective Go, the Go standard library, Google Go Style Guide,
and Uber Go Style Guide.

## References

- [Effective Go](https://go.dev/doc/effective_go) - Primary official idioms
- [Go Standard Library](https://pkg.go.dev/std) - Real-world gold standard
- [Google Go Style Guide](https://google.github.io/styleguide/go/)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

## Package Structure

```
cmd/                    # Main applications (one per binary)
internal/              # Private packages (not importable externally)
  parser/              # Document parsing and scoring
  llm/                 # LLM integration
  ui/                  # Terminal UI components
  prompts/             # Prompt management
pkg/                   # Public packages (importable)
```

## Naming Conventions

### Packages
- **lowercase, single word** preferred: `parser`, `llm`, `ui`
- Avoid `utils`, `common`, `helpers` - be specific

### Variables and Functions
- **camelCase** for unexported: `targetPath`, `repoList`
- **PascalCase** for exported: `ParsePRFAQ`, `NewLogger`
- **Avoid stuttering**: `parser.Parse` not `parser.ParseDocument`
- **Acronyms**: `URL`, `HTTP`, `ID` (all caps), or `url`, `http`, `id` (all lower)

### Interfaces
- **Single-method interfaces**: Name after the method with -er suffix
  - `Reader`, `Writer`, `Scanner`
- **Multi-method interfaces**: Descriptive noun

### Receivers
- **Short, 1-2 letters**, consistent across type's methods
  - `func (l *Logger) Info(...)` not `func (logger *Logger) Info(...)`
- **Pointer receivers** for mutable state; value receivers for immutable

## Error Handling

### Error Wrapping
```go
// Good: Add context with fmt.Errorf and %w
if err != nil {
    return fmt.Errorf("failed to read file %s: %w", path, err)
}

// Bad: Losing context
if err != nil {
    return err
}
```

### Error Variables
```go
// Define package-level sentinel errors
var (
    ErrNotFound     = errors.New("not found")
    ErrInvalidInput = errors.New("invalid input")
)
```

### Never Panic in Library Code
- Use `panic` only in `main()` or `init()` for unrecoverable setup errors
- Always return errors from functions

## Functions and Methods

### Function Length
- Target: **≤50 lines** per function
- Maximum: **100 lines** (refactor if approaching)
- Single responsibility principle

### Parameters
- **≤5 parameters** - use options struct if more needed
- Context first: `func Foo(ctx context.Context, ...)`
- Error always last return: `func Foo() (Result, error)`

### Named Return Values
- Use for documentation purposes in short functions
- Avoid in functions >10 lines (naked returns obscure flow)

## Concurrency

### Goroutine Safety
```go
// Good: Clear ownership, explicit synchronization
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()
data = append(data, item)

// Bad: Race condition
go func() {
    data = append(data, item) // Unsynchronized
}()
```

### Worker Pools
- Use `sync.WaitGroup` for coordinating goroutines
- Use buffered channels for work queues
- Always handle context cancellation

## Testing

### Table-Driven Tests
```go
func TestParsePath(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"empty", "", "", true},
        {"valid", "/tmp/foo", "/tmp/foo", false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParsePath(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParsePath() error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("ParsePath() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Coverage
- Target: **≥80%** for core logic
- Focus on behavior, not implementation details
- Use `testdata/` directory for test fixtures

## Documentation

### Package Comments
```go
// Package parser provides document parsing for PR-FAQ documents
// and quality analysis. It supports various header formats and styles.
package parser
```

### Function Comments
```go
// ParsePRFAQ extracts structured sections from a PR-FAQ document.
// It handles various header formats (H1/H2) and naming conventions.
// Returns an error if the document is empty or malformed.
func ParsePRFAQ(content string) (*Document, error)
```

## File Organization

### File Naming
- `snake_case.go` for multi-word files: `quality_analysis.go`
- `*_test.go` for tests
- `doc.go` for package documentation (if needed)

### File Length
- Target: **≤400 lines** per file
- Split large files by responsibility

## Linting

All code must pass:
```bash
golangci-lint run
```

Key linters enabled:
- `errcheck` - Check error returns
- `govet` - Static analysis
- `staticcheck` - Advanced checks
- `gosec` - Security issues
- `gofmt` - Formatting
- `goimports` - Import organization

