# Process 14.0: Thought Sky

```mermaid
graph TB
    User((User))
    D16[(D16: Sky Thoughts)]
    
    User -->|14.1 Type negative thought| P14[14.0 Thought Sky]
    User -->|14.2 Swipe clouds| P14
    
    P14 -->|Store thought| D16
    
    P14 -->|14.3 Cloud animation| User
    P14 -->|14.4 Clouds passing message| User
```

## Data Store: D16 Sky Thoughts

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users |
| thought_text | TEXT | Thought content |
| cloud_swiped | BOOLEAN | Cloud swiped status |
| swiped_at | TIMESTAMP | Swipe timestamp |
| created_date | TIMESTAMP | Creation timestamp |
| day_number | INTEGER | Program day (1-56) |
| created_at | TIMESTAMP | Creation timestamp |
