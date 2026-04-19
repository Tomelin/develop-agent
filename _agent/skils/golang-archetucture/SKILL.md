# Skill: Golang Architecture

This skill formalizes the architectural layers for backend development as defined in the agency's Playbook. It ensures structured code and clear separation of concerns.

## Project Structure

### 1. Layers
- `cmd/`: Application entry points (e.g., `main.go`). Keep logic here minimal; only bootstrap the app.
- `internal/core/`: The core of the application. Contains business entities, domain services, and interfaces (ports). **No external dependencies allowed here.**
- `internal/core/module_name/entity/`: The entity layer. Contains business entities. **No external dependencies allowed here.**
- `internal/core/module_name/service/`: The service layer. Contains business logic.
- `internal/core/module_name/repository/`: The repository layer. Contains storage logic.  this layer connect with database and other storages if necessary.
- `internal/core/module_name/handlers/`: The handlers layer. Contains handler logic.  The layer contains handler of grpc, http with gin.
- `internal/core/module_name/events/`: The events layer. Contains event handler logic.  The layer contains handler of kafka, rabbitmq.
- `internal/core/module_name/tools/`: The tools layer. Contains tools for the module.
- `internal/core/module_name/module.go`: The module file initialize the application and export the service for other modules.
- `internal/infra/`: Implementation of external dependencies (adapters). E.g., Database repositories, API clients, Email senders.
- `pkg/`: Shared code that can be imported by other projects. Use sparingly.
- `config/`: Configuration mapping and environment setup.
- `api/`: API definitions for grpc and http with gin.
- `build/`: Build definitions for docker, etc.
- `deploy/`: Deployment definitions for kubernetes, etc.
- `scripts/`: Scripts for the application.
- `docs/`: Documentation for the application.
- `test/`: Test files for the application.

**ATTENTION**
When necessary connect on thirty party, example: MongoDB, Redis, Kafka, RabbitMQ, etc.
1. We must be create the category on `pkg/database/mongodb`, `pkg/database/redis`, `pkg/stream/kafka`, `pkg/stream/rabbitmq`, etc.
2. We must be create the interface between provider and business on `internal/infra/database`.
3. We must be create de repository on `internal/core/module_name/repository`.  The repository must be implement the interface on `internal/infra/database`. The repository should not have any dependencies on external providers.


### 2. Dependency Rule
- Dependencies must point **inwards**: 
    - `infra` depends on `domain`.
    - `handler` depends on `domain`.
    - `domain` depends on **nothing** (except standard library).

### 3. Business Logic
- All business rules must reside in `internal/core/module_name/entity`. 
- Use the **Repository Pattern** to abstract data access.

## Implementation Guidelines

### Domain Entity
```go
package entity

// Order represents a customer purchase.
// It is the central entity for the order management domain.
type Order struct {
    ID     string
    Status string
}

// OrderRepository defines the interface for persisting orders.
// Adapters in the infra layer will implement this.
type OrderRepository interface {
    Save(order *Order) error
    GetByID(id string) (*Order, error)
}
```
