# Go Flexible Scheduler

A flexible, production-ready task scheduler built with Go, Fiber, and GORM. Supports both one-time scheduled tasks and recurring tasks with minute-level precision.

## Features

### Core Scheduling
- **One-time Tasks**: Schedule tasks for specific date/time execution
- **Recurring Tasks**: Advanced recurring patterns with minute-level precision
- **HTTP Task Execution**: Execute HTTP requests (GET, POST, PUT, DELETE, etc.)
- **Webhook Support**: Optional webhook notifications after task completion
- **Retry Mechanism**: Automatic retry on task failures
- **Multi-user Support**: JWT-based authentication and user isolation

### Recurring Task Patterns
- **Daily**: Execute every day at specific time (e.g., every day at 07:30)
- **Hourly**: Execute every hour at specific minute (e.g., every hour at minute 30)
- **Minutely**: Execute every X minutes (e.g., every 5 minutes)
- **Weekly**: Execute weekly on specific day and time (e.g., every Monday at 09:00)
- **Monthly**: Execute monthly on specific date and time (e.g., every 15th at 14:30)
- **Custom**: Support for cron-like expressions (extensible)

### Task Management
- **Real-time Tracking**: Tasks move through pending → running → completed/failed states
- **Process ID Tracking**: Unique process IDs for task correlation
- **Automatic Cleanup**: Configurable cleanup of old completed/failed tasks
- **Performance Metrics**: Request duration and webhook duration tracking

## Quick Start

### Prerequisites
- Go 1.24+
- MySQL database
- Environment configuration

### Installation

1. Clone the repository:
```bash
git clone https://github.com/abilfida/go-flexible-scheduler.git
cd go-flexible-scheduler
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Configure your environment variables in `.env`:
```env
DB_DSN=user:password@tcp(localhost:3306)/scheduler_db?charset=utf8mb4&parseTime=True&loc=Local
APP_PORT=3000
APP_TIMEZONE=Asia/Jakarta
JWT_SECRET=your-secret-key
CLEANUP_INTERVAL_HOURS=24
```

4. Run the application:
```bash
go run main.go
```

### Using Docker

```bash
docker-compose up -d
```

## API Documentation

### Authentication

#### Sign Up
```http
POST /api/v1/signup
Content-Type: application/json

{
  "username": "user123",
  "email": "user@example.com",
  "password": "password123"
}
```

#### Sign In
```http
POST /api/v1/signin
Content-Type: application/json

{
  "username": "user123",
  "password": "password123"
}
```

### One-time Tasks

#### Create Task
```http
POST /api/v1/tasks
Authorization: Bearer <token>
Content-Type: application/json

{
  "url": "https://api.example.com/webhook",
  "method": "POST",
  "headers": "{\"Content-Type\":\"application/json\"}",
  "body": "{\"message\":\"Hello World\"}",
  "scheduled_at": "2025-11-07 09:30:00",
  "webhook_url": "https://myapp.com/task-completed"
}
```

#### Get Tasks
```http
GET /api/v1/tasks?status=pending
Authorization: Bearer <token>
```

Status options: `pending`, `running`, `completed`, `failed`

### Recurring Tasks

#### Create Recurring Task

**Daily at 07:30:**
```http
POST /api/v1/recurring
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Daily Morning Report",
  "description": "Generate daily report every morning",
  "url": "https://api.example.com/generate-report",
  "method": "POST",
  "headers": "{\"Content-Type\":\"application/json\"}",
  "body": "{\"type\":\"daily_report\"}",
  "pattern": "daily",
  "time": "07:30",
  "start_date": "2025-11-07T00:00:00+07:00",
  "is_active": true
}
```

**Every hour at minute 30:**
```http
POST /api/v1/recurring
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Hourly Health Check",
  "url": "https://api.example.com/health",
  "method": "GET",
  "pattern": "hourly",
  "time": "00:30",
  "start_date": "2025-11-07T00:00:00+07:00",
  "is_active": true
}
```

**Every 5 minutes:**
```http
POST /api/v1/recurring
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "System Monitor",
  "url": "https://api.example.com/monitor",
  "method": "GET",
  "pattern": "minutely",
  "interval": 5,
  "start_date": "2025-11-07T00:00:00+07:00",
  "is_active": true
}
```

**Weekly (every Monday at 09:00):**
```http
POST /api/v1/recurring
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Weekly Backup",
  "url": "https://api.example.com/backup",
  "method": "POST",
  "pattern": "weekly",
  "time": "09:00",
  "day_of_week": 1,
  "start_date": "2025-11-07T00:00:00+07:00",
  "is_active": true
}
```

**Monthly (15th of each month at 14:30):**
```http
POST /api/v1/recurring
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Monthly Invoice",
  "url": "https://api.example.com/invoice",
  "method": "POST",
  "pattern": "monthly",
  "time": "14:30",
  "day_of_month": 15,
  "start_date": "2025-11-07T00:00:00+07:00",
  "end_date": "2025-12-31T23:59:59+07:00",
  "max_runs": 12,
  "is_active": true
}
```

#### Manage Recurring Tasks

**List recurring tasks:**
```http
GET /api/v1/recurring
Authorization: Bearer <token>
```

**Get specific recurring task:**
```http
GET /api/v1/recurring/{id}
Authorization: Bearer <token>
```

**Update recurring task:**
```http
PUT /api/v1/recurring/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "is_active": false,
  "time": "08:00"
}
```

**Toggle active status:**
```http
POST /api/v1/recurring/{id}/toggle
Authorization: Bearer <token>
```

**Delete recurring task:**
```http
DELETE /api/v1/recurring/{id}
Authorization: Bearer <token>
```

## Recurring Pattern Reference

### Pattern Types

| Pattern | Description | Required Fields | Example |
|---------|-------------|----------------|--------|
| `daily` | Every day at specific time | `time` (HH:MM) | "07:30" |
| `hourly` | Every hour at specific minute | `time` (00:MM) | "00:30" |
| `minutely` | Every X minutes | `interval` (number) | 5 |
| `weekly` | Weekly on specific day/time | `time`, `day_of_week` | "09:00", 1 (Monday) |
| `monthly` | Monthly on specific date/time | `time`, `day_of_month` | "14:30", 15 |
| `custom` | Custom cron expression | `cron_expr` | "0 30 9 * * MON" |

### Day of Week Values
- 0 = Sunday
- 1 = Monday
- 2 = Tuesday
- 3 = Wednesday
- 4 = Thursday
- 5 = Friday
- 6 = Saturday

### Time Format
- Use 24-hour format: "HH:MM"
- Examples: "07:30", "14:45", "23:59"
- For hourly pattern, use "00:MM" where MM is the target minute

### Control Fields

| Field | Type | Description |
|-------|------|-------------|
| `start_date` | datetime | When to start the recurring pattern |
| `end_date` | datetime | Optional end date |
| `max_runs` | integer | Optional maximum number of executions |
| `is_active` | boolean | Enable/disable the recurring task |

## Architecture

### Components

- **Config**: Environment configuration and timezone management
- **Database**: GORM-based MySQL connection and models
- **Migration**: Automatic database schema management
- **Router**: Fiber-based REST API routes
- **Middleware**: JWT authentication and request validation
- **Task Models**: TaskCore, TaskPending, TaskRunning, TaskCompleted, TaskFailed, RecurringTask
- **Schedulers**: 
  - Main scheduler (processes pending tasks every 5 seconds)
  - Recurring scheduler (checks recurring tasks every 30 seconds)
  - Cleanup scheduler (removes old tasks periodically)
- **Executor**: HTTP request execution with retry logic

### Task Lifecycle

1. **One-time Tasks**: User creates → TaskPending → TaskRunning → TaskCompleted/TaskFailed
2. **Recurring Tasks**: User creates RecurringTask → Scheduler spawns TaskPending at scheduled time → continues lifecycle

### Database Tables

- `users`: User authentication
- `tasks_pending`: Tasks waiting to be executed
- `tasks_running`: Tasks currently being executed
- `tasks_completed`: Successfully completed tasks
- `tasks_failed`: Failed tasks (after retries)
- `recurring_tasks`: Recurring task templates and schedules

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|--------|
| `DB_DSN` | MySQL connection string | Required |
| `APP_PORT` | Application port | 3000 |
| `APP_TIMEZONE` | Application timezone | Asia/Jakarta |
| `JWT_SECRET` | JWT signing secret | Required |
| `CLEANUP_INTERVAL_HOURS` | Cleanup interval in hours | 24 |

### Timezone Support

The application supports timezone-aware scheduling. All time calculations are performed in the configured timezone (`APP_TIMEZONE`).

## Performance and Scalability

- **Efficient Polling**: 5-second intervals for pending tasks, 30-second intervals for recurring tasks
- **Database Indexes**: Optimized queries with proper indexing on scheduling fields
- **Concurrent Execution**: Tasks execute in separate goroutines
- **Resource Management**: Automatic cleanup of old task records
- **Process Isolation**: User-based task isolation with JWT authentication

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
