# Installation Session Repositories

## Available Implementations

### MemorySessionRepository (Recommended for MVP)
- **Status**: ✅ Fully functional
- **Use case**: Development, testing, MVP
- **Pros**:
  - Fast
  - Thread-safe
  - No external dependencies
  - Full support for all operations
- **Cons**:
  - Data lost on restart
  - Not suitable for production with multiple instances

### SQLiteSimpleSessionRepository
- **Status**: ✅ Fully functional
- **Use case**: Single-instance production with persistence needs
- **Features**:
  - Complete CRUD operations (Save, FindByID, List)
  - Domain factory method for aggregate reconstruction (`ReconstructInstallationSession`)
  - DTO pattern for proper serialization/deserialization
  - WAL mode for better concurrency
  - Indexed queries for performance
- **Pros**:
  - Persistent storage across restarts
  - Zero external dependencies (embedded database)
  - Full test coverage
  - Proper DDD encapsulation maintained
- **Cons**:
  - Single-instance only (not suitable for horizontal scaling)
  - File-based storage (backup required for data safety)

## Implementation Details

The SQLite repository uses a DTO (Data Transfer Object) pattern to bridge the gap between domain aggregates and persistence:

1. **Domain → Storage**: `toStorageModel()` converts domain objects to DTOs with public fields
2. **Storage → Domain**: `fromStorageModel()` reconstructs domain objects using proper constructors
3. **Factory Method**: `ReconstructInstallationSession()` validates invariants during reconstruction

This approach maintains DDD principles while enabling persistence without compromising encapsulation.

## Recommendation

### For Development/Testing
Use `MemorySessionRepository`:
- Fast and simple
- No setup required
- Perfect for unit tests

### For Production
Use `SQLiteSimpleSessionRepository`:
- Persistent storage across restarts
- Suitable for single-server deployments
- Full feature parity with memory repository
- Consider implementing backups for data safety

### For High-Scale Production
Consider future implementations:
- PostgreSQL repository for multi-instance deployments
- Event sourcing for complete audit trail
- CQRS with read models for query optimization

## Configuration

Current server uses `MemorySessionRepository` by default:

```go
// cmd/server/main.go or internal/cli/cmd/server.go
sessionRepo := repository.NewMemorySessionRepository()
```

To switch to SQLite (once reconstruction is implemented):
```go
sessionRepo, err := repository.NewSQLiteSimpleSessionRepository("./gohan.db")
if err != nil {
    log.Fatal(err)
}
defer sessionRepo.Close()
```
