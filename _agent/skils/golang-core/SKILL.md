# Skill: Golang Core

This skill defines the foundational standards for developing in Golang within the agency. It focuses on readability, documentation, and error handling.

## Standards

### 1. Code Documentation
- **Godoc**: Every exported function, type, and constant MUST have a Godoc comment. 
- **Format**: Follow the official Godoc format (e.g., `// FunctionName does X...`).
- **Explanation**: Comments should explain **WHY** the code exists and any tricky logic, not just **WHAT** it does (the code should be self-documenting for "what").

### 2. Error Handling
- **Explicit Checks**: Check every return error. Never ignore errors unless explicitly justified.
- **Wrapping**: Use `fmt.Errorf("context: %w", err)` to provide context to errors while preserving the original error for type checking.
- **Sentinel Errors**: Define custom error types or sentinel errors in the `domain` layer if they represent business rules.

### 3. Naming Conventions
- **Interfaces**: Usually end in `-er` (e.g., `Reader`, `Writer`, `Processor`).
- **Packages**: Short, concise, and lowercase. Avoid `util` or `common`.
- **Variables**: Use short names for small scopes (e.g., `i`, `r`) and descriptive names for larger scopes.

## Examples

### Good Documentation
```go
// Fprint formats using the default formats for its operands and writes to w.
// It exists to allow flexible output redirection for formatted strings.
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
    // implementation
}
```

### Good Error Handling
```go
data, err := os.ReadFile(path)
if err != nil {
    // We wrap the error to provide context about which file failed
    return nil, fmt.Errorf("failed to read configuration file %s: %w", path, err)
}
```
