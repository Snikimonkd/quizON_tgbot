ALTER TABLE games ADD COLUMN openning_time timestamptz NOT NULL DEFAULT now();
