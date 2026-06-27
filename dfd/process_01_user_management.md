# Process 1.0: User Management & Authentication

```mermaid
graph TB
    User((User))
    D1[(D1: Users)]
    
    User -->|1.1 Initial registration| P1[1.0 User Management]
    User -->|1.2 Login| P1
    User -->|1.3 Update profile| P1
    
    P1 -->|Store user data| D1
    D1 -->|Load data| P1
    
    P1 -->|1.4 Show login history| User
    P1 -->|1.5 Permission grants| User
```

## Data Store: D1 Users

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| phone_number | VARCHAR(20) | Unique phone number |
| registration_date | TIMESTAMP | Registration timestamp |
| last_login | TIMESTAMP | Last login timestamp |
| login_count | INTEGER | Total login count |
| agreement_accepted | BOOLEAN | Digital agreement status |
| cloud_sync_enabled | BOOLEAN | Cloud sync preference |
| do_not_disturb_enabled | BOOLEAN | DND mode status |
| created_at | TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | Last update timestamp |
