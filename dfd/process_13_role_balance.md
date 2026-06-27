# Process 13.0: Role & Values Balance

```mermaid
graph TB
    User((User))
    D15[(D15: Roles & Values)]
    
    User -->|13.1 Define org role| P13[13.0 Role Balance]
    User -->|13.2 Define personal value| P13
    
    P13 -->|Store role and value| D15
    
    P13 -->|13.3 Show overlapping circles| User
```

## Data Store: D15 Roles & Values

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users |
| entry_type | VARCHAR(20) | role or value |
| entry_text | TEXT | Entry content |
| created_date | TIMESTAMP | Creation timestamp |
| day_number | INTEGER | Program day (1-56) |
| created_at | TIMESTAMP | Creation timestamp |
