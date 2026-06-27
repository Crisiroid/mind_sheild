# Process 10.0: Mind Court

```mermaid
graph TB
    User((User))
    D10[(D10: Negative Thoughts)]
    D11[(D11: Court Evidence)]
    
    User -->|10.1 Log supporting evidence| P10[10.0 Mind Court]
    User -->|10.2 Log contradicting evidence| P10
    User -->|10.3 Request guide help| P10
    
    P10 -->|Store evidence| D11
    D10 -->|Related negative thought| P10
    
    P10 -->|10.4 Show scale model| User
    P10 -->|10.5 Generate alternative thought| User
```

## Data Store: D11 Court Evidence

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users |
| negative_thought_id | UUID | FK to negative_thoughts |
| supporting_evidence | TEXT | Supporting evidence |
| contradicting_evidence | TEXT | Contradicting evidence |
| guide_helper_used | BOOLEAN | Guide helper activated |
| alternative_thought | TEXT | Generated alternative |
| created_date | TIMESTAMP | Creation timestamp |
| day_number | INTEGER | Program day (1-56) |
| created_at | TIMESTAMP | Creation timestamp |
