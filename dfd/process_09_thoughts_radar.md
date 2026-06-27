# Process 9.0: Negative Thoughts Radar

```mermaid
graph TB
    User((User))
    D10[(D10: Negative Thoughts)]
    
    User -->|9.1 Log negative thought| P9[9.0 Thoughts Radar]
    User -->|9.2 Set impact 1-10| P9
    User -->|9.3 Log situation & error type| P9
    
    P9 -->|Store negative thought| D10
    
    P9 -->|9.4 Termite animation| User
    P9 -->|9.5 Incident report UI| User
```

## Data Store: D10 Negative Thoughts

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users |
| thought_text | TEXT | Negative thought content |
| situation | TEXT | Situation context |
| cognitive_error_type | VARCHAR(50) | Error type classification |
| impact_level | INTEGER | Impact 1-10 |
| recorded_at | TIMESTAMP | Recording timestamp |
| day_number | INTEGER | Program day (1-56) |
| created_at | TIMESTAMP | Creation timestamp |
