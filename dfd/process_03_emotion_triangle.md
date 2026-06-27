# Process 3.0: Emotion Triangle Interaction

```mermaid
graph TB
    User((User))
    D3[(D3: Triangle Interactions)]
    
    User -->|3.1 Click thought side| P3[3.0 Emotion Triangle]
    User -->|3.2 Click body side| P3
    User -->|3.3 Click behavior side| P3
    
    P3 -->|Log interaction| D3
    
    P3 -->|3.4 Show thought accounts| User
    P3 -->|3.5 Phone vibration| User
```

## Data Store: D3 Triangle Interactions

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users |
| interaction_date | TIMESTAMP | Interaction timestamp |
| side_clicked | VARCHAR(20) | thought/body/behavior |
| thought_accounts_viewed | JSONB | Viewed thought accounts |
| vibration_triggered | BOOLEAN | Vibration activated |
| day_number | INTEGER | Program day (1-56) |
| created_at | TIMESTAMP | Creation timestamp |
