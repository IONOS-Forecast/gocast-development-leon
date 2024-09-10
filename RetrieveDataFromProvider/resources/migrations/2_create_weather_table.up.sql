CREATE TABLE IF NOT EXISTS weather_records
(
    id                           SERIAL PRIMARY KEY,
    timestamp                    TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    source_id                    INTEGER NOT NULL,
    precipitation                NUMERIC(10,4),
    pressure_msl                 NUMERIC(10,4),
    sunshine                     NUMERIC(10,4),
    temperature                  NUMERIC(10,4),
    wind_direction               INTEGER,
    wind_speed                   NUMERIC(10,4),
    cloud_cover                  NUMERIC(10,4),
    dew_point                    NUMERIC(10,4),
    relative_humidity            INTEGER,
    visibility                   INTEGER,
    wind_gust_direction          INTEGER,
    wind_gust_speed              NUMERIC(10,4),
    condition                    TEXT,
    precipitation_probability    INTEGER,
    precipitation_probability_6h INTEGER,
    solar                        NUMERIC(10,4),
    icon                         TEXT,
    city                         TEXT,

    FOREIGN KEY (city) REFERENCES cities (name)
);