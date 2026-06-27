# Process 5.0: Body Map & Physical Tension

```mermaid
graph TB
    User((User))
    D5[(D5: Body Maps)]
    
    User -->|5.1 Touch body area| P5[5.0 Body Map]
    User -->|5.2 Set tension intensity| P5
    User -->|5.3 Select color gradient| P5
    
    P5 -->|Store body map| D5
    
    P5 -->|5.4 Show anatomical figure| User
```

## Data Store: D5 Body Maps

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users |
| mapping_date | TIMESTAMP | Mapping timestamp |
| body_regions | JSONB | Body regions with intensity |
| overall_intensity | INTEGER | Overall intensity 1-10 |
| severity_color | VARCHAR(20) | Yellow to dark red |
| notes | TEXT | Additional notes |
| day_number | INTEGER | Program day (1-56) |
| created_at | TIMESTAMP | Creation timestamp |
