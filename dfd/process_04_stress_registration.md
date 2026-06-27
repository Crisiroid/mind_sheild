# Process 4.0: Stress Event Registration

```mermaid
graph TB
    User((User))
    D4[(D4: Stress Events)]
    
    User -->|4.1 Select work situation| P4[4.0 Stress Registration]
    User -->|4.2 Set intensity 1-10| P4
    
    P4 -->|Store stress event| D4
    
    P4 -->|4.3 Confirmation| User
```

## Data Store: D4 Stress Events

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users |
| event_timestamp | TIMESTAMP | Event timestamp |
| situation_type | VARCHAR(50) | Type of stress situation |
| situation_description | TEXT | Situation description |
| intensity_level | INTEGER | Intensity 1-10 |
| location | VARCHAR(100) | Event location |
| day_number | INTEGER | Program day (1-56) |
| created_at | TIMESTAMP | Creation timestamp |
