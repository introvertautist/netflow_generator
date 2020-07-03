package main

import (
    "fmt"
    "os"
    "reflect"

    cli "github.com/jawher/mow.cli"
)

// Default arguments
const (
    DefaultCollectorHost      = "0.0.0.0"
    DefaultCollectorPort      = 2056
    DefaultWorkersCount       = 1
    DefaultFlowsCount         = 1
    DefaultSrcNetwork         = "10.99.0.0/30"
    DefaultDstNetwork         = "10.99.0.16/30"
    DefaultSrcPotr            = 2103
    DefaultDstPotr            = 80
    DefaultTrafficVariability = false
    DefaultDelay              = 0
    NetflowV5Type             = 5
    DefaultMaxPackets         = 5
    DefaultMaxPacketZize      = 2048
)

// Print config data to console
func printConfig(config *AppConfig) {
    fmt.Println("Run generator with params:")

    fields := reflect.TypeOf(*config)
    values := reflect.ValueOf(*config)

    for i := 0; i < fields.NumField(); i++ {
        field := fields.Field(i)
        value := values.Field(i)
        fmt.Printf("%20s: %v\n", field.Name, value)
    }
}

/*
Parse common arguments and build config:
Args:
    --collector-host: Netflow collector host
    --collector-port: Netflow collector port
    --workers: Number of workers to generate data
    --flows: Number of flows exported in this packet (1-30)
    --src-network: Network for generating source IP address
    --dst-network: Network for generating destination IP address
    --src-port: TCP/UDP source port
    --dst-port: TCP/UDP destination port
    --traffic-variability: Use different traffic type (SSH, SNMP, HTTP, etc.)
    --delay: Delay between requests in milliseconds. If 0 - without delay
    --max-packets: Max count of generating pakets
    --max-packets-size: Max packet size in bytes
*/
func commonArgs(cmd *cli.Cmd) *AppConfig {
    cmd.Spec = "[--collector-host=<collector-host>] " +
        "[--collector-port=<collector-port>] " +
        "[--workers=<workers>] " +
        "[--flows=<flows>] " +
        "[--src-network=<src-network>] " +
        "[--dst-network=<dst-network>] " +
        "[--src-port=<src-port>] " +
        "[--dst-port=<dst-port>] " +
        "[--traffic-variability=<traffic-variability>] " +
        "[--delay=<delay>] " +
        "[--max-packets=<max-packets>] " +
        "[--max-packets-size=<max-packets-size>] "

    var config AppConfig

    cmd.StringOptPtr(&config.CollectorHost, "collector-host", DefaultCollectorHost, "Collector host")
    cmd.IntOptPtr(&config.CollectorPort, "collector-port", DefaultCollectorPort, "Collector port")

    cmd.IntOptPtr(&config.Workers, "workers", DefaultWorkersCount, "Number of workers to generate data")
    cmd.IntOptPtr(&config.Flows, "flows", DefaultFlowsCount, "Number of flows exported in this packet (1-30)")

    cmd.StringOptPtr(&config.Source, "src-network", DefaultSrcNetwork, "Network for generating source IP address")
    cmd.StringOptPtr(&config.Destination, "dst-network", DefaultDstNetwork, "Network for generating destination IP address")

    cmd.IntOptPtr(&config.SourcePort, "src-port", DefaultSrcPotr, "TCP/UDP source port")
    cmd.IntOptPtr(&config.DestinationPort, "dst-port", DefaultDstPotr, "TCP/UDP destination port")

    cmd.BoolOptPtr(&config.TrafficVariability, "traffic-variability", DefaultTrafficVariability, "Use different traffic type (SSH, SNMP, HTTP, etc.)")
    cmd.IntOptPtr(&config.Delay, "delay", DefaultDelay, "Delay between requests in milliseconds. If 0 - without delay")

    cmd.IntOptPtr(&config.MaxPackets, "max-packets", DefaultMaxPackets, "Max count of generating pakets")
    cmd.IntOptPtr(&config.MaxPacketSize, "max-packets-size", DefaultMaxPacketZize, "Max packet size in bytes")

    return &config
}

/*
Prepearing to netflow generation
Parse common arguments and generate data
*/
func prepareNetflowV5(cmd *cli.Cmd) {
    config := commonArgs(cmd)

    cmd.Action = func() {
        fmt.Println(Info("Netflow v5 generator start"))

        config.Version = NetflowV5Type
        printConfig(config)

        NetflowV5Generator(config)
    }
}

/*
Generate netflow traffic

Available arguments:
    netflow:
        Usage: genaratectl netflow VERSION [OPTIONS]
        Available VERSION:
            v5 - Netflow version 5
*/
func main() {
    app := cli.App("statctl", "Genarate statistic")

    app.Command("netflow", "Genarate Netflow statistic", func(config *cli.Cmd) {
        config.Command("v5", "Netflow v5", prepareNetflowV5)
    })

    err := app.Run(os.Args)

    if err != nil {
        fmt.Println(Error("Got fatal error: ", err))
        os.Exit(1)
    }
}
