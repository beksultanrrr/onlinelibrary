ALTER TABLE books ADD CONSTRAINT movies_readtime_check CHECK (readtime >= 0);
ALTER TABLE books ADD CONSTRAINT movies_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
ALTER TABLE books ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);