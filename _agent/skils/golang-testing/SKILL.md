# Skill: Golang Testing

This skill defines the standards for writing tests in Golang, ensuring reliability and high quality of deliverables.

## Testing Standards

**ATTENTION**
We must be use the packages **testfy** for testing.
Examples:
- `"github.com/stretchr/testify/assert"`
- `"github.com/stretchr/testify/require"`
- `"github.com/stretchr/testify/mock"`
- `"github.com/stretchr/testify/suite"`

### 1. Naming Convention
- Test functions must follow the pattern: `Test[FunctionName]_[Scenario]_[ExpectedResult]`.
- Example: `TestCalculateTotal_MultipleItems_Success`.

### 2. Test Types
- **Unit Tests**: Focus on a single function or method. Use mocks for dependencies.
- **Integration Tests**: Focus on the interaction between multiple components (e.g., Service + Database).
- **Functional/E2E Tests**: Test the entire flow from an external perspective (e.g., API requests).

### 3. Table-Driven Tests
- Use table-driven tests for functions with multiple input scenarios to ensure clean and maintainable test code.

### 4. Mocking
- Generate mocks for interfaces in the `domain` layer to test services in isolation.

## Examples

### Table-Driven Test
```go
func TestSum(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 1, 2, 3},
        {"negative numbers", -1, -2, -3},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Sum(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Sum(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```
