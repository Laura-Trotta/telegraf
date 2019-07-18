# Archiver Plugin

This plugin saves metrics to files dividing them by metric name, host, year and day.
Given a path, a metric will be stored in:
	
given-directory/metric-name/host/year/day-of-year

# COnfiguration

[[outputs.archiver]]

  ## Directory to write to
  directory = ""

 
  ## Data format to output.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_OUTPUT.md
  data_format = "influx"
