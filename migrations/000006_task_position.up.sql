BEGIN;

ALTER TABLE todoapp.tasks
ADD COLUMN position DOUBLE PRECISION DEFAULT EXTRACT(EPOCH FROM NOW()) NOT NULL;

CREATE INDEX idx_tasks_position ON todoapp.tasks (folder_id, position);

COMMIT;
