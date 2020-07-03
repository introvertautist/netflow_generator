# Statistic generator
A utility that allows you to generate statistics/traffic data for testing. Currently supports only Netflow. In next releases will be added:
- Netflow v9
- Netflow v10
- sFlow
- Graphite

## Building
Statistic generator is a Golang project. In order to build a project from source you need Golang. Follow [this](https://golang.org/doc/install) instructions to install language. Then execute:
```bash
go build -o ./bin/statctl .
```

## Common usage
Run `statctl --help` to show all available commands
```bash
Usage: statctl COMMAND [arg...]

Genarate statistic
               
Commands:      
  netflow      Genarate Netflow statistic
```

## Netflow generator
Currently supports only the Netflow v5 version. In future releases will be added versions 9 and 10.  

### Netflow v5
To generate Netflow v5 data, you should run the `statctl` executable file with the argument `netflow v5`. 
```bash
statctl netflow v5
```
**Usage:**  
```
statctl netflow v5 [--collector-host=<collector-host>]  
                   [--collector-port=<collector-port>] 
                   [--workers=<workers>] 
                   [--flows=<flows>] 
                   [--src-network=<src-network>] 
                   [--dst-network=<dst-network>] 
                   [--src-port=<src-port>] 
                   [--dst-port=<dst-port>] 
                   [--traffic-variability=<traffic-variability>] 
                   [--delay=<delay>] [--max-packets=<max-packets>] 
                   [--max-packets-size=<max-packets-size>]
```

**Options:**  
```
--collector-host        Collector host (default "0.0.0.0")
--collector-port        Collector port (default 2056)
--workers               Number of workers to generate data (default 1)
--flows                 Number of flows exported in this packet (1-30)(default  1)

--src-network           Network for generating source IP address (default "10.99.0.0/30")

--dst-network           Network for generating destination IP address (default "10.99.0.16/30")

--src-port              TCP/UDP source port (default 2103)
--dst-port              TCP/UDP destination port (default 80)
--traffic-variability   Use different traffic type (SSH, SNMP, HTTP, etc.)
--delay                 Delay between requests in milliseconds. If 0 - without delay (default 0)

--max-packets           Max count of generating pakets (default 5)
--max-packets-size      Max packet size in bytes (default 2048)
```

