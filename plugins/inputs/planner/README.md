# Planner Plugin

This input plugin reads files from a given folder and processes the metrics they contain. The files must have numbers as their name, as they will be read n days after the plugin is first started, where n is the file name. When a file is processed, the timestamp of its metrics will be changed to match the year, month and day to the date when they are processed. When it has finished reading all the metrics in every file, it starts over. 

The tags of the metrics can also be modified to custom values using configuration.

The files containing the metrics must be located in the folder specified in the configuration.
Another directory must be configured where the plugin will save its "plan", that is a .json file to keep track of the plugin's work. 

It can be used to simulate having new metrics when you want, recycling the same files. 

## Configuration

```toml
[[inputs.planner]]

## Directory containing the files to be read
directory = ""

## Directory where the plan will be saved
plandirectory = ""

## List of tags to be modified. 
## Example: "tag1=newtag1,tag2=newtag2"
tagslist = ""

## The dataformat to be read from files
## Each data format has its own unique set of configuration options, read
## more about them here:
## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
data_format = "influx"
```