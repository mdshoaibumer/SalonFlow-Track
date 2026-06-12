DROP TABLE IF EXISTS license_notifications;

-- SQLite does not support DROP COLUMN directly; these columns are additive and safe to leave.
-- If needed, recreate the table without them.
