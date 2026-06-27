# Process 2.0: Agreement & Program Initiation

```mermaid
graph TB
    User((User))
    D1[(D1: Users)]
    D2[(D2: Agreements)]
    
    User -->|2.1 Accept agreement| P2[2.0 Agreement Mgmt]
    
    P2 -->|Store signed agreement| D2
    D1 -->|User data| P2
    
    P2 -->|2.2 Show roadmap| User
    P2 -->|2.3 Activate Day 1| User
```

## Data Store: D2 Agreements

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users |
| agreement_text | TEXT | Agreement content |
| accepted | BOOLEAN | Acceptance status |
| accepted_date | TIMESTAMP | Acceptance timestamp |
| created_at | TIMESTAMP | Creation timestamp |
