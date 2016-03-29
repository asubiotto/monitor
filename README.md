# Monitor for HTTP access logs
## About
Monitor is a program that can be used to monitor an HTTP access log as it is being written to. Every 10 seconds, a ranking of a subset of sections with the highest number of hits ranked by these is displayed in the console. The interval and size of this ranking can be modified by changing values in `monitor/const.go`.
An optional argument, `threshold` is used to alert of traffic spikes. Whenever traffic in a specified time window (default is two minutes, value can be changed in `const.go`) divided by the number of sections hit so far goes above this threshold, an alert occurs. An alert recovery message is printed once this value once again falls below the threshold.
Monitor also provides other useful metrics like requests per second, total sections, and total traffic.

## Usage
```
go get github.com/asubiotto/monitor
go run main.go -path=<path to http access log file> [-threshold=<traffic spike threshold>]
```

## Improvements
### Global tracker
`tracker.go` defines a global object that is accessed both by the reader and the reporter to set and access information. Tracker is initialized only once and the `tracker` global variable is set to point to this Tracker. I would try to find a design pattern so that I would have the minimum amount of information exposed even within the package.

### HTTP Access Log parsing
I came across a file of HTTP access logs (`monitor/epa-http.txt`) that I use for testing the parsing of HTTP access logs. The thing is that even though a general guideline for writing these logs exists, not many logs seem to adhere to it (for example omitting response codes or '-', symbolizing missing information). A more robust access log parsing mechanism is something I would work on.
