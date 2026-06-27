-- ============================================================
-- Sepur Ravan Psychology App - PostgreSQL Database Schema
-- Version: 1.0
-- Database: PostgreSQL 14+
-- ============================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- SCHEMA: public - Main Application Tables
-- ============================================================

-- ------------------------------------------------------------
-- Table: users - User accounts and profiles
-- ------------------------------------------------------------
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    registration_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login TIMESTAMP WITH TIME ZONE,
    login_count INTEGER DEFAULT 1,
    agreement_accepted BOOLEAN DEFAULT FALSE,
    agreement_accepted_date TIMESTAMP WITH TIME ZONE,
    cloud_sync_enabled BOOLEAN DEFAULT FALSE,
    do_not_disturb_enabled BOOLEAN DEFAULT FALSE,
    dnd_start_time TIME,
    dnd_end_time TIME,
    android_version VARCHAR(10),
    app_version VARCHAR(10),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT valid_phone CHECK (phone_number ~ '^\+?[0-9]{10,15}$')
);

COMMENT ON TABLE users IS 'کاربران اپلیکیشن - User accounts';
COMMENT ON COLUMN users.agreement_accepted IS 'پذیرش میثاق‌نامه دیجیتال - Digital agreement acceptance';
COMMENT ON COLUMN users.cloud_sync_enabled IS 'فعال‌سازی همگام‌سازی ابری - Cloud sync preference';

-- ------------------------------------------------------------
-- Table: user_settings - User preferences and settings
-- ------------------------------------------------------------
CREATE TABLE user_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    notification_enabled BOOLEAN DEFAULT TRUE,
    vibration_enabled BOOLEAN DEFAULT TRUE,
    language VARCHAR(5) DEFAULT 'fa',
    font_size VARCHAR(10) DEFAULT 'medium',
    theme VARCHAR(20) DEFAULT 'calm',
    crisis_alert_threshold INTEGER DEFAULT 7,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id)
);

COMMENT ON TABLE user_settings IS 'تنظیمات کاربر - User preferences';

-- ------------------------------------------------------------
-- Table: daily_calendar - 56-day program calendar tracking
-- ------------------------------------------------------------
CREATE TABLE daily_calendar (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    day_number INTEGER NOT NULL CHECK (day_number BETWEEN 1 AND 56),
    calendar_date DATE NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    activities_completed JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id, day_number)
);

COMMENT ON TABLE daily_calendar IS 'تقویم 56 روزه - 56-day program calendar';
COMMENT ON COLUMN daily_calendar.activities_completed IS 'فعالیت‌های تکمیل شده روز - Completed daily activities';

CREATE INDEX idx_daily_calendar_user ON daily_calendar(user_id);
CREATE INDEX idx_daily_calendar_date ON daily_calendar(calendar_date);

-- ------------------------------------------------------------
-- Table: emotion_triangle_interactions - Emotion triangle usage
-- ------------------------------------------------------------
CREATE TABLE emotion_triangle_interactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    interaction_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    side_clicked VARCHAR(20) NOT NULL CHECK (side_clicked IN ('thought', 'body', 'behavior')),
    thought_accounts_viewed JSONB,
    vibration_triggered BOOLEAN DEFAULT FALSE,
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE emotion_triangle_interactions IS 'تعاملات مثلث هیجان - Emotion triangle interactions';
COMMENT ON COLUMN emotion_triangle_interactions.side_clicked IS 'ضلع انتخاب شده: thought/f فکر, body/بدن, behavior/رفتار';

CREATE INDEX idx_emotion_triangle_user ON emotion_triangle_interactions(user_id);

-- ------------------------------------------------------------
-- Table: stress_events - Work-related stress event registrations
-- ------------------------------------------------------------
CREATE TABLE stress_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    event_timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    situation_type VARCHAR(50) NOT NULL,
    situation_description TEXT,
    intensity_level INTEGER NOT NULL CHECK (intensity_level BETWEEN 1 AND 10),
    location VARCHAR(100),
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE stress_events IS 'موقعیت‌های استرس‌زا شغلی - Work stress events';
COMMENT ON COLUMN stress_events.intensity_level IS 'شدت استرس 1-10 - Stress intensity 1-10';

CREATE INDEX idx_stress_events_user ON stress_events(user_id);
CREATE INDEX idx_stress_events_date ON stress_events(event_timestamp);

-- ------------------------------------------------------------
-- Table: body_tension_maps - Body tension mapping
-- ------------------------------------------------------------
CREATE TABLE body_tension_maps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    mapping_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    body_regions JSONB NOT NULL,
    overall_intensity INTEGER CHECK (overall_intensity BETWEEN 1 AND 10),
    severity_color VARCHAR(20),
    notes TEXT,
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE body_tension_maps IS 'نقشه تنش بدنی - Body tension maps';
COMMENT ON COLUMN body_tension_maps.body_regions IS 'نواحی بدن با شدت - Body regions with intensity (JSON)';
COMMENT ON COLUMN body_tension_maps.severity_color IS 'zرد تا قرمز تیره - Yellow to dark red';

CREATE INDEX idx_body_tension_user ON body_tension_maps(user_id);

-- ------------------------------------------------------------
-- Table: breathing_sessions - Mindful breathing exercise sessions
-- ------------------------------------------------------------
CREATE TABLE breathing_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_start TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    session_end TIMESTAMP WITH TIME ZONE,
    duration_seconds INTEGER,
    breathing_pattern VARCHAR(50),
    is_completed BOOLEAN DEFAULT FALSE,
    calendar_ticked BOOLEAN DEFAULT FALSE,
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE breathing_sessions IS 'جلسات تنفس آگاهانه - Mindful breathing sessions';

CREATE INDEX idx_breathing_user ON breathing_sessions(user_id);

-- ------------------------------------------------------------
-- Table: cognitive_error_games - Cognitive error detection games
-- ------------------------------------------------------------
CREATE TABLE cognitive_error_games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    game_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    scenario_id INTEGER NOT NULL,
    scenario_type VARCHAR(50),
    score INTEGER,
    is_correct BOOLEAN,
    time_taken_seconds INTEGER,
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE cognitive_error_games IS 'بازی تشخیص خطای شناختی - Cognitive error detection games';

CREATE INDEX idx_cognitive_games_user ON cognitive_error_games(user_id);

-- ------------------------------------------------------------
-- Table: mental_musts - Mental "musts" backpack entries
-- ------------------------------------------------------------
CREATE TABLE mental_musts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    must_text TEXT NOT NULL,
    created_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_released BOOLEAN DEFAULT FALSE,
    released_date TIMESTAMP WITH TIME ZONE,
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE mental_musts IS 'بایدهای ذهنی کوله‌پشتی - Mental musts backpack';

CREATE INDEX idx_mental_musts_user ON mental_musts(user_id);

-- ------------------------------------------------------------
-- Table: negative_thoughts - Negative thought radar entries
-- ------------------------------------------------------------
CREATE TABLE negative_thoughts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    thought_text TEXT NOT NULL,
    situation TEXT,
    cognitive_error_type VARCHAR(50),
    impact_level INTEGER CHECK (impact_level BETWEEN 1 AND 10),
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE negative_thoughts IS 'رادار افکار منفی - Negative thoughts radar';
COMMENT ON COLUMN negative_thoughts.cognitive_error_type IS 'نوع خطای شناختی - Cognitive error type';
COMMENT ON COLUMN negative_thoughts.impact_level IS 'میزان اثر بر عملکرد 1-10 - Impact on performance 1-10';

CREATE INDEX idx_negative_thoughts_user ON negative_thoughts(user_id);
CREATE INDEX idx_negative_thoughts_date ON negative_thoughts(recorded_at);

-- ------------------------------------------------------------
-- Table: mind_court_evidence - Mind court evidence records
-- ------------------------------------------------------------
CREATE TABLE mind_court_evidence (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    negative_thought_id UUID REFERENCES negative_thoughts(id) ON DELETE CASCADE,
    supporting_evidence TEXT,
    contradicting_evidence TEXT,
    guide_helper_used BOOLEAN DEFAULT FALSE,
    alternative_thought TEXT,
    created_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE mind_court_evidence IS 'شواهد دادگاه ذهن - Mind court evidence';
COMMENT ON COLUMN mind_court_evidence.alternative_thought IS 'فکر جایگزین منطقی - Logical alternative thought';

CREATE INDEX idx_mind_court_user ON mind_court_evidence(user_id);

-- ------------------------------------------------------------
-- Table: conflict_exercises - Work conflict scenario exercises
-- ------------------------------------------------------------
CREATE TABLE conflict_exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    scenario_id INTEGER NOT NULL,
    practice_count INTEGER DEFAULT 1,
    last_practice_date TIMESTAMP WITH TIME ZONE,
    performance_score INTEGER,
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id, scenario_id)
);

COMMENT ON TABLE conflict_exercises IS 'سناریوهای تعارض کاری - Work conflict exercises';

CREATE INDEX idx_conflict_ex_user ON conflict_exercises(user_id);

-- ------------------------------------------------------------
-- Table: mood_tracker - Mood and activity tracking
-- ------------------------------------------------------------
CREATE TABLE mood_tracker (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    activity_id UUID,
    activity_name VARCHAR(100),
    mood_before INTEGER NOT NULL CHECK (mood_before BETWEEN 1 AND 10),
    mood_after INTEGER NOT NULL CHECK (mood_after BETWEEN 1 AND 10),
    activity_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    notes TEXT,
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE mood_tracker IS 'ردیاب خلق و فعالیت - Mood and activity tracker';
COMMENT ON COLUMN mood_tracker.activity_name IS 'تماس با دوست، مطالعه، چای، ورزش، موسیقی';

CREATE INDEX idx_mood_tracker_user ON mood_tracker(user_id);

-- ------------------------------------------------------------
-- Table: roles_and_values - Organizational roles and personal values
-- ------------------------------------------------------------
CREATE TABLE roles_and_values (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    entry_type VARCHAR(20) NOT NULL CHECK (entry_type IN ('role', 'value')),
    entry_text TEXT NOT NULL,
    created_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE roles_and_values IS 'نقش سازمانی و ارزش‌های فردی - Organizational roles and personal values';

CREATE INDEX idx_roles_values_user ON roles_and_values(user_id);

-- ------------------------------------------------------------
-- Table: sky_thoughts - Thought sky cloud entries
-- ------------------------------------------------------------
CREATE TABLE sky_thoughts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    thought_text TEXT NOT NULL,
    cloud_swiped BOOLEAN DEFAULT FALSE,
    swiped_at TIMESTAMP WITH TIME ZONE,
    created_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE sky_thoughts IS 'افکار آسمان - Sky thoughts (clouds)';

CREATE INDEX idx_sky_thoughts_user ON sky_thoughts(user_id);

-- ------------------------------------------------------------
-- Table: mindful_timers - Mindful activity timer sessions
-- ------------------------------------------------------------
CREATE TABLE mindful_timers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    timer_start TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    timer_end TIMESTAMP WITH TIME ZONE,
    duration_seconds INTEGER,
    vibration_reminders_count INTEGER DEFAULT 0,
    is_completed BOOLEAN DEFAULT FALSE,
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE mindful_timers IS 'تایمر فعالیت آگاهانه - Mindful activity timers';

CREATE INDEX idx_mindful_timers_user ON mindful_timers(user_id);

-- ------------------------------------------------------------
-- Table: acceptance_exercises - Active acceptance vs surrender exercises
-- ------------------------------------------------------------
CREATE TABLE acceptance_exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    video_watched BOOLEAN DEFAULT FALSE,
    watched_at TIMESTAMP WITH TIME ZONE,
    understanding_level INTEGER CHECK (understanding_level BETWEEN 1 AND 10),
    notes TEXT,
    day_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE acceptance_exercises IS 'تمرینات پذیرش فعال - Active acceptance exercises';

CREATE INDEX idx_acceptance_user ON acceptance_exercises(user_id);

-- ------------------------------------------------------------
-- Table: weekly_reports - Weekly progress reports
-- ------------------------------------------------------------
CREATE TABLE weekly_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    week_number INTEGER NOT NULL CHECK (week_number BETWEEN 1 AND 8),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    stress_events_count INTEGER DEFAULT 0,
    avg_stress_intensity DECIMAL(5,2),
    breathing_sessions_count INTEGER DEFAULT 0,
    negative_thoughts_count INTEGER DEFAULT 0,
    body_tension_maps_count INTEGER DEFAULT 0,
    mood_improvement_score DECIMAL(5,2),
    activities_distribution JSONB,
    progress_percentage DECIMAL(5,2),
    generated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id, week_number)
);

COMMENT ON TABLE weekly_reports IS 'گزارش وضعیت هفتگی - Weekly progress reports';
COMMENT ON COLUMN weekly_reports.activities_distribution IS 'توزیع فعالیت‌ها - Activity distribution (pie chart data)';

CREATE INDEX idx_weekly_reports_user ON weekly_reports(user_id);

-- ============================================================
-- SCHEMA: admin_panel - Admin Panel Tables
-- ============================================================
CREATE SCHEMA IF NOT EXISTS admin_panel;

-- ------------------------------------------------------------
-- Table: admin_panel.admin_users - Admin user accounts
-- ------------------------------------------------------------
CREATE TABLE admin_panel.admin_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    role_id UUID,
    is_active BOOLEAN DEFAULT TRUE,
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE admin_panel.admin_users IS 'کاربران پنل ادمین - Admin panel users';

-- ------------------------------------------------------------
-- Table: admin_panel.roles - Admin roles
-- ------------------------------------------------------------
CREATE TABLE admin_panel.roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    permissions JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE admin_panel.roles IS 'نقش‌های ادمین - Admin roles';

-- ------------------------------------------------------------
-- Table: admin_panel.user_reports - Anonymized user reports
-- ------------------------------------------------------------
CREATE TABLE admin_panel.user_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_type VARCHAR(50) NOT NULL,
    report_date DATE NOT NULL,
    total_users INTEGER DEFAULT 0,
    active_users INTEGER DEFAULT 0,
    avg_engagement_score DECIMAL(5,2),
    crisis_alerts_count INTEGER DEFAULT 0,
    anonymized_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE admin_panel.user_reports IS 'گزارش‌های ناشناس کاربران - Anonymized user reports';
COMMENT ON COLUMN admin_panel.user_reports.anonymized_data IS 'داده‌های آماری بدون شناسه - Statistical data without identifiers';

CREATE INDEX idx_user_reports_date ON admin_panel.user_reports(report_date);

-- ------------------------------------------------------------
-- Table: admin_panel.system_logs - System activity logs
-- ------------------------------------------------------------
CREATE TABLE admin_panel.system_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    log_type VARCHAR(50) NOT NULL,
    log_message TEXT,
    user_id UUID,
    severity VARCHAR(20) DEFAULT 'info',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE admin_panel.system_logs IS 'لاگ‌های سیستم - System logs';

CREATE INDEX idx_system_logs_type ON admin_panel.system_logs(log_type);
CREATE INDEX idx_system_logs_date ON admin_panel.system_logs(created_at);

-- ============================================================
-- VIEWS - Useful queries for reporting and analytics
-- ============================================================

-- ------------------------------------------------------------
-- View: user_progress_summary - Complete user progress
-- ------------------------------------------------------------
CREATE VIEW user_progress_summary AS
SELECT 
    u.id AS user_id,
    u.phone_number,
    u.registration_date,
    u.last_login,
    dc.day_number AS current_day,
    dc.is_completed AS today_completed,
    COUNT(DISTINCT dc.id) FILTER (WHERE dc.is_completed = TRUE) AS days_completed,
    COUNT(DISTINCT se.id) AS total_stress_events,
    COUNT(DISTINCT bs.id) AS total_breathing_sessions,
    COUNT(DISTINCT nt.id) AS total_negative_thoughts,
    ROUND(AVG(se.intensity_level), 2) AS avg_stress_intensity,
    ROUND(
        (COUNT(DISTINCT dc.id) FILTER (WHERE dc.is_completed = TRUE)::DECIMAL / 56) * 100, 
        2
    ) AS progress_percentage
FROM users u
LEFT JOIN daily_calendar dc ON u.id = dc.user_id
LEFT JOIN stress_events se ON u.id = se.user_id
LEFT JOIN breathing_sessions bs ON u.id = bs.user_id
LEFT JOIN negative_thoughts nt ON u.id = nt.user_id
GROUP BY u.id, dc.day_number, dc.is_completed;

COMMENT ON VIEW user_progress_summary IS 'خلاصه پیشرفت کاربر - User progress summary';

-- ------------------------------------------------------------
-- View: weekly_activity_stats - Weekly activity statistics
-- ------------------------------------------------------------
CREATE VIEW weekly_activity_stats AS
SELECT 
    user_id,
    DATE_TRUNC('week', created_at) AS week_start,
    COUNT(DISTINCT stress_events.id) AS stress_events_count,
    COUNT(DISTINCT breathing_sessions.id) AS breathing_sessions,
    COUNT(DISTINCT negative_thoughts.id) AS negative_thoughts,
    COUNT(DISTINCT body_tension_maps.id) AS body_maps,
    COUNT(DISTINCT mood_tracker.id) AS mood_entries,
    ROUND(AVG(mood_tracker.mood_after - mood_tracker.mood_before), 2) AS avg_mood_improvement
FROM users
LEFT JOIN stress_events ON users.id = stress_events.user_id
LEFT JOIN breathing_sessions ON users.id = breathing_sessions.user_id
LEFT JOIN negative_thoughts ON users.id = negative_thoughts.user_id
LEFT JOIN body_tension_maps ON users.id = body_tension_maps.user_id
LEFT JOIN mood_tracker ON users.id = mood_tracker.user_id
GROUP BY user_id, DATE_TRUNC('week', created_at);

COMMENT ON VIEW weekly_activity_stats IS 'آمار فعالیت هفتگی - Weekly activity statistics';

-- ============================================================
-- FUNCTIONS & TRIGGERS
-- ============================================================

-- ------------------------------------------------------------
-- Function: Update updated_at timestamp
-- ------------------------------------------------------------
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to tables with updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_settings_updated_at BEFORE UPDATE ON user_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_conflict_exercises_updated_at BEFORE UPDATE ON conflict_exercises
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ------------------------------------------------------------
-- Function: Detect crisis thoughts (self-harm indicators)
-- ------------------------------------------------------------
CREATE OR REPLACE FUNCTION detect_crisis_thoughts()
RETURNS TRIGGER AS $$
BEGIN
    -- Check for crisis keywords in negative thoughts
    IF NEW.thought_text ILIKE '%خودکشی%' OR 
       NEW.thought_text ILIKE '%مرگ%' OR
       NEW.thought_text ILIKE '%کشتن خود%' THEN
        -- Insert into crisis alerts table (should be created)
        INSERT INTO admin_panel.system_logs (log_type, log_message, user_id, severity)
        VALUES ('crisis_alert', 'Potential crisis thought detected', NEW.user_id, 'critical');
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER crisis_detection_trigger 
    AFTER INSERT ON negative_thoughts
    FOR EACH ROW EXECUTE FUNCTION detect_crisis_thoughts();

-- ============================================================
-- INDEXES - Performance optimization
-- ============================================================

-- Composite indexes for common queries
CREATE INDEX idx_daily_calendar_user_day ON daily_calendar(user_id, day_number);
CREATE INDEX idx_stress_events_user_date ON stress_events(user_id, event_timestamp);
CREATE INDEX idx_negative_thoughts_user_date ON negative_thoughts(user_id, recorded_at);
CREATE INDEX idx_breathing_user_date ON breathing_sessions(user_id, session_start);
CREATE INDEX idx_mood_tracker_user_date ON mood_tracker(user_id, activity_date);

-- JSONB indexes
CREATE INDEX idx_body_tension_regions ON body_tension_maps USING GIN (body_regions);
CREATE INDEX idx_calendar_activities ON daily_calendar USING GIN (activities_completed);

-- ============================================================
-- INITIAL DATA
-- ============================================================

-- Insert default admin roles
INSERT INTO admin_panel.roles (role_name, description, permissions) VALUES
('super_admin', 'Full system access', '{"all": true}'),
('content_manager', 'Manage content and exercises', '{"content": true, "reports": true}'),
('viewer', 'View reports only', '{"reports": true}');

-- ============================================================
-- SAMPLE QUERIES for Tracking & Analytics
-- ============================================================

/*
-- 1. User engagement tracking (last 7 days)
SELECT 
    DATE(created_at) AS activity_date,
    COUNT(DISTINCT user_id) AS active_users,
    COUNT(*) AS total_interactions
FROM (
    SELECT user_id, created_at FROM stress_events
    UNION ALL
    SELECT user_id, created_at FROM breathing_sessions
    UNION ALL
    SELECT user_id, created_at FROM negative_thoughts
) AS all_activities
WHERE created_at >= NOW() - INTERVAL '7 days'
GROUP BY DATE(created_at)
ORDER BY activity_date;

-- 2. Most common stress situations
SELECT 
    situation_type,
    COUNT(*) AS occurrence_count,
    ROUND(AVG(intensity_level), 2) AS avg_intensity
FROM stress_events
GROUP BY situation_type
ORDER BY occurrence_count DESC;

-- 3. Body tension heat map data
SELECT 
    jsonb_object_keys(body_regions) AS body_region,
    COUNT(*) AS tension_count,
    ROUND(AVG(overall_intensity), 2) AS avg_intensity
FROM body_tension_maps
GROUP BY body_region
ORDER BY tension_count DESC;

-- 4. Weekly progress for specific user
SELECT * FROM weekly_reports 
WHERE user_id = 'USER_UUID_HERE'
ORDER BY week_number;

-- 5. Cognitive error type distribution
SELECT 
    cognitive_error_type,
    COUNT(*) AS frequency,
    ROUND(AVG(impact_level), 2) AS avg_impact
FROM negative_thoughts
WHERE cognitive_error_type IS NOT NULL
GROUP BY cognitive_error_type
ORDER BY frequency DESC;
*/
