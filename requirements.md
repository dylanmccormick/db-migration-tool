Project Requirements
Core Functionality

Read migration files from a migrations/ directory
Track applied migrations in a schema_migrations table
Migrate up: Apply pending migrations in sequential order
Migrate down: Rollback the last N migrations
Status command: Show which migrations are applied/pending
Support multiple databases: SQLite, PostgreSQL, MySQL

Migration File Format

Naming: {version}_{description}.{up|down}.sql
Example: 001_create_users.up.sql, 001_create_users.down.sql
Version numbers must be sequential integers

migrate up              # Apply all pending migrations
migrate down            # Rollback last migration
migrate down 3          # Rollback last 3 migrations
migrate status          # Show migration status
migrate create <name>   # Create new migration file pair
```

### Error Handling
- Fail fast if migrations are out of order
- Use database transactions (all-or-nothing)
- Clear error messages for common issues

---
