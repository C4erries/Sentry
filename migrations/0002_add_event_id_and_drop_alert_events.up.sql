
BEGIN;

-- 1. Добавляем в alerts колонку event_id (одно событие, породившее алерт)
ALTER TABLE IF EXISTS alerts
  ADD COLUMN IF NOT EXISTS event_id UUID;

-- 2. Добавим индекс для нового поля
CREATE INDEX IF NOT EXISTS idx_alerts_event_id ON alerts(event_id);

-- 3. Удаляем связующую таблицу alert_events целиком
DROP TABLE IF EXISTS alert_events;

COMMIT;
