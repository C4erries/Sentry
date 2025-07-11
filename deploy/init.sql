-- events table
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    user_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    ip INET,
    geo_country TEXT,
    data JSONB
);

CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id);
CREATE INDEX IF NOT EXISTS idx_events_event_type ON events(event_type);
CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at);

CREATE TABLE IF NOT EXISTS rules (
    id TEXT PRIMARY KEY,         -- example: 'login_storm', 'geo_switching'
    description TEXT NOT NULL    -- человекочитаемое описание, example: 'Login storm detection'
);

INSERT INTO rules (id, description) VALUES
  ('login_storm', 'Login storm detection'),
  ('geo_switching', 'Geographical switch anomaly');

CREATE TABLE IF NOT EXISTS levels (
    id TEXT PRIMARY KEY,         -- example: 'warning', 'critical', 'info'
    priority INTEGER NOT NULL    -- для сортировки / фильтрации
);

INSERT INTO levels (id, priority) VALUES
  ('info', 1),
  ('warning', 2),
  ('critical', 3);

-- alerts table
CREATE TABLE IF NOT EXISTS alerts (
    id UUID PRIMARY KEY,
    rule TEXT REFERENCES rules(id) ON DELETE RESTRICT,
    level TEXT REFERENCES levels(id) ON DELETE RESTRICT,
    detected_at TIMESTAMPTZ NOT NULL,
    data JSONB,
);

-- relation table
CREATE TABLE IF NOT EXISTS alert_events (
    alert_id UUID REFERENCES alerts(id) ON DELETE CASCADE,
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    PRIMARY KEY (alert_id, event_id)
);
