# Process 15.0: Mindful Activity Timer

```mermaid
graph TB
    User((User))
    D17[(D17: Timers)]
    
    User -->|15.1 Start timer| P15[15.0 Mindful Timer]
    User -->|15.2 Receive reminder vibration| P15
    
    P15 -->|Store session| D17
    
    P15 -->|15.3 Periodic vibration| User
```

## Data Store: D17 Timers

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users |
| timer_start | TIMESTAMP | Timer start time |
| timer_end | TIMESTAMP | Timer end time |
| duration_seconds | INTEGER | Timer duration |
| vibration_reminders_count | INTEGER | Vibration count |
| is_completed | BOOLEAN | Completion status |
| day_number | INTEGER | Program day (1-56) |
| created_at | TIMESTAMP | Creation timestamp |
