# Archiver Plugin

This plugin saves metrics to files dividing them by metric name, a specific tag, year and day.
Given a path, a metric will be stored in:
	
given-directory/metric-name/tag/year/day-of-year

## Configuration

```toml
[[outputs.archiver]]

  ## Directory to write to
  directory = ""

  ## Tag 
  tag = ""

 
  ## Data format to output.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_OUTPUT.md
  data_format = "influx"
  ```
