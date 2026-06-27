# Process 7.0: Cognitive Error Game

```mermaid
graph TB
    User((User))
    D8[(D8: Game Results)]
    
    User -->|7.1 Select scenario| P7[7.0 Cognitive Game]
    User -->|7.2 Drag & drop| P7
    
    P7 -->|Store game result| D8
    
    P7 -->|7.3 Show feedback| User
    P7 -->|7.4 Scoring| User
```

## Data Store: D8 Game Results

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users |
| game_date | TIMESTAMP | Game timestamp |
| scenario_id | INTEGER | Scenario identifier |
| scenario_type | VARCHAR(50) | Scenario type |
| score | INTEGER | Game score |
| is_correct | BOOLEAN | Correct answer |
| time_taken_seconds | INTEGER | Time to complete |
| day_number | INTEGER | Program day (1-56) |
| created_at | TIMESTAMP | Creation timestamp |
