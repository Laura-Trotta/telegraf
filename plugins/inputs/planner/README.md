# Planner Plugin

This input plugin reads one file once per day and processes its metrics, changing the year, day and month to the current ones. When it has finished reading all the metrics in every file, it starts over. 

The files containing the metrics must be located in the folder specified in the configuration, and they will be read in alphabetical order. 
Another directory must be configured where the plugin will save its "plan", that is a .json file to keep track of the plugin's work. 

It can be used to simulate having new metrics every day, recycling the same files. 

#Configuration

[[inputs.planner]]

## Directory containing the files to be read
directory = ""

## Directory where the plan will be saved
plandirectory = ""

## The dataformat to be read from files
## Each data format has its own unique set of configuration options, read
## more about them here:
## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
data_format = "influx"
