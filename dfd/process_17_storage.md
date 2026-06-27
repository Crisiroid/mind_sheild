# Process 17.0: Local Storage & Sync

```mermaid
graph TB
    User((User))
    D19[(D19: Local Storage)]
    Cloud[(Cloud Storage)]
    
    User -->|17.1 Choose sync| P17[17.0 Storage Mgmt]
    
    P17 -->|17.2 Auto local save| D19
    D19 -->|Data| P17
    
    P17 <-->|17.3 Optional sync| Cloud
    
    P17 -->|17.4 Non-replacement warning| User
```

## Data Store: D19 Local Storage

Local SQLite database containing all application tables with sync status fields:

| Field | Type | Description |
|-------|------|-------------|
| sync_status | VARCHAR(20) | pending/synced/conflict |
| last_synced_at | TIMESTAMP | Last sync timestamp |
| device_id | VARCHAR(100) | Device identifier |
| storage_used_mb | INTEGER | Storage usage |
