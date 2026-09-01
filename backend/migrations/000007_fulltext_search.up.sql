BEGIN;

-- Add a generated tsvector column for fulltext search on title + description.
-- Using 'simple' config to handle Cyrillic and mixed-language content without
-- needing pg_catalog language packs. coalesce() handles NULL descriptions.
ALTER TABLE todoapp.tasks
ADD COLUMN search_vector tsvector
GENERATED ALWAYS AS (
    to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(description, ''))
) STORED;

CREATE INDEX idx_tasks_search_vector ON todoapp.tasks USING GIN (search_vector);

COMMIT;
