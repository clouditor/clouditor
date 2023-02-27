# `metrics.json`: Format of metrics (see metric.proto for more details)
| Field | Description |
| ----------- | ----------------- |
| id   | The ID of metric. Must not contain spaces.  |
| name     | A more readable name    |
| description  | A short description |
| category  | The control of the catalog |
| scale  | 1 for nominal data, 2 for ordinal data and 3 for numbers |
| range | Depending on scale: "allowedValues" for nominal, "order" for ordinal and "minMax" for numbers|