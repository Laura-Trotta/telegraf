# Planner Plugin

This input plugin reads one file once per day and processes its metrics, changing the year, day and month to the current ones. When it has finished reading all the metrics in every file, it starts over. 

#Configuration

  [[inputs.planner]]

  ## Directory containing the files to be modified
  directory = "/home/ltrotta/Desktop/provaplugin/archiveddata"

  ## Directory containing configuration file
  confdirectory = "/home/ltrotta/Desktop/provaplugin/configplanner"

  ## The dataformat to be read from files
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
  data_format = "influx"
