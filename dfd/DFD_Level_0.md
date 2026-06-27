# DFD Level 0: Context Diagram

```mermaid
graph TB
    User((User))
    Admin((Admin))
    Cloud[(Cloud Storage)]
    CrisisLine((Crisis Hotline))
    
    System[Sepur Ravan App System]
    
    User -->|Daily interactions| System
    User -->|Stress & tension data| System
    User -->|Negative thoughts| System
    User -->|Activities & mood| System
    User -->|Settings & preferences| System
    
    System -->|Infographics & education| User
    System -->|Notifications & reminders| User
    System -->|Weekly progress reports| User
    System -->|Interactive exercises| User
    System -->|Crisis alerts| User
    
    Admin -->|User management| System
    Admin -->|View reports| System
    Admin -->|Role configuration| System
    
    System -->|Anonymous statistics| Admin
    System -->|Monitoring dashboard| Admin
    
    System <-->|Optional sync| Cloud
    
    System -->|Emergency contact| CrisisLine
    CrisisLine -->|Crisis information| System
```
