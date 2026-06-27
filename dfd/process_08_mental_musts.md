# Process 8.0: Mental Musts Backpack

```mermaid
graph TB
    User((User))
    D9[(D9: Mental Musts)]
    
    User -->|8.1 Write mental must| P8[8.0 Backpack]
    User -->|8.2 View stones| P8
    
    P8 -->|Store mental must| D9
    
    P8 -->|8.3 Show backpack image| User
```

## Data Store: D9 Mental Musts

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users |
| must_text | TEXT | Mental must content |
| created_date | TIMESTAMP | Creation timestamp |
| is_released | BOOLEAN | Released status |
| released_date | TIMESTAMP | Release timestamp |
| day_number | INTEGER | Program day (1-56) |
| created_at | TIMESTAMP | Creation timestamp |
