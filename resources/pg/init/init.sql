CREATE TABLE IF NOT EXISTS cities
(
    name TEXT PRIMARY KEY,
    lat  NUMERIC(10,4) NOT NULL,
    lon  NUMERIC(10,4) NOT NULL
);

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

COPY cities(name, lat, lon) FROM '/usr/pgdata/cities.csv' WITH CSV;
COPY weather_records(timestamp, source_id, precipitation, pressure_msl, sunshine, temperature, wind_direction, wind_speed, cloud_cover, dew_point, relative_humidity, visibility, wind_gust_direction, wind_gust_speed, condition, precipitation_probability, precipitation_probability_6h, solar, icon, city) FROM '/usr/pgdata/berlin.csv' WITH CSV;
COPY weather_records(timestamp, source_id, precipitation, pressure_msl, sunshine, temperature, wind_direction, wind_speed, cloud_cover, dew_point, relative_humidity, visibility, wind_gust_direction, wind_gust_speed, condition, precipitation_probability, precipitation_probability_6h, solar, icon, city) FROM '/usr/pgdata/hamburg.csv' WITH CSV;
COPY weather_records(timestamp, source_id, precipitation, pressure_msl, sunshine, temperature, wind_direction, wind_speed, cloud_cover, dew_point, relative_humidity, visibility, wind_gust_direction, wind_gust_speed, condition, precipitation_probability, precipitation_probability_6h, solar, icon, city) FROM '/usr/pgdata/munich.csv' WITH CSV;
