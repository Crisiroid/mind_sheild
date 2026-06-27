# Process 11.0: Work Conflict Scenarios

```mermaid
graph TB
    User((User))
    D12[(D12: Conflict Exercises)]
    
    User -->|11.1 Select scenario| P11[11.0 Conflict Practice]
    User -->|11.2 Repeated practice| P11
    
    P11 -->|Store exercise progress| D12
    
    P11 -->|11.3 Show feedback| User
```

## Data Store: D12 Conflict Exercises

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users |
| scenario_id | INTEGER | Scenario identifier |
| practice_count | INTEGER | Number of practices |
| last_practice_date | TIMESTAMP | Last practice timestamp |
| performance_score | INTEGER | Performance score |
| day_number | INTEGER | Program day (1-56) |
| created_at | TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | Last update timestamp |
