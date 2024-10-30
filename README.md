# Gocast
Gocast is an application for the weather, written in golang. You can use query to get weather-data.

|arg|environment|default|description|
|-|-|-|-|
|--weather-api-url|WEATHER_API_URL|https://api.brightsky.dev/weather?lat=52&lon=7.6&date=2024-08-21|URL of a weather forecast data provider|
| --weather-geo-url|GEO_API_URL|http://api.openweathermap.org/geo/1.0/direct?q=Berlin&limit=1&appid=|URL of a geocoding data provider|
|--geo-api-key|GEO_API_KEY|*********4d62|API-Key for a geocoding data provider|
|--req-aft-min|REQ_AFT_MIN|30|Weather Data Requests per minute|
|--fdb-user|FDB_USER|forecast|User of the database|
|--fdb-password|FDB_PASSWORD|forecast|Password to the database|
|--fdb-database|FDB_DATABASE|forecast|Name of the database|
|--fdb-address|FDB_ADDRESS|localhost:5544|Adress to the database|
---
### Grafana Dashboards


**Dashboard**
![Dashboard](resources/screenshots/Dashboard.png)
**Average Dashboard**
![Average Dashboard 1](resources/screenshots/Average-Dashboard1.png)
![Average Dashboard 2](resources/screenshots/Average-Dashboard2.png)

---
### How to run
To run the application you need to start the database first, for that you can use `make startdb`, then you need to use `go run main.go` or `make run2` to start the application.
To stop the application you need to close it with `CTRL+C`.
\
Afterwards you can also close the database using `make stopdb`.

You can run the application with `make run` or `docker-compose up`.
If you want to Pause it you can use `CTRL+C`. When you want to completely stop the application you need to close it with
\
`make stop` or `docker-compose down`.

### Example usages
Use `date` and `city` query parameters to get data for a certain date and city.
\
Both parameters are required!
```bash
curl "http://localhost:8080/direct?date=2024-08-25&city=Berlin"

curl "http://localhost:8080/direct?date=2010-12-10&city=Hamburg"

curl "http://localhost:8080/direct?date=2014-05-31&city=Berlin"
```
### Example with example output
You receive data in JSON format.
```bash
curl "http://localhost:8080/direct?date=2024-10-20&city=hamburg"
```
```json
// Output: (Reduced to 1 hour)
{
  "weather": [
    {
      "timestamp": "2024-10-20 00:00:00+00",
      "source_id": 134968,
      "precipitation": 0,
      "pressure_msl": 1019.7,
      "sunshine": 0,
      "temperature": 11.3,
      "wind_direction": 140,
      "wind_speed": 6.1,
      "cloud_cover": 100,
      "dew_point": 11.3,
      "relative_humidity": 100,
      "visibility": 15680,
      "wind_gust_direction": 150,
      "wind_gust_speed": 9.4,
      "condition": "dry",
      "precipitation_probability": 0,
      "precipitation_probability_6h": 0,
      "solar": 0,
      "icon": "cloudy",
      "City": "hamburg"
    }
  ]
}
```