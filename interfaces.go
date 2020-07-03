package main

// AppConfig - Structure of application config data
type AppConfig struct {
    CollectorHost       string
    CollectorPort       int
    Workers             int
    Flows               int
    Source              string
    Destination         string
    SourcePort          int
    DestinationPort     int
    TrafficVariability  bool
    Version             int
	Delay               int
	MaxPackets			int
	MaxPacketSize		int
}

// NetflowHeaderV5 - Structure of headers Netflow v5
type NetflowHeaderV5 struct {
    Version        uint16
    FlowCount      uint16
    SysUptime      uint32
    UnixSec        uint32
    UnixMsec       uint32
    FlowSequence   uint32
    EngineType     uint8
    EngineID       uint8
    SampleInterval uint16
}

// NetflowPayloadV5 - Structure of Netflow v5 body
type NetflowPayloadV5 struct {
    SrcIP          uint32
    DstIP          uint32
    NextHopIP      uint32
    SnmpInIndex    uint16
    SnmpOutIndex   uint16
    NumPackets     uint32
    NumOctets      uint32
    SysUptimeStart uint32
    SysUptimeEnd   uint32
    SrcPort        uint16
    DstPort        uint16
    Padding1       uint8
    TCPFlags       uint8
    IPProtocol     uint8
    IPTos          uint8
    SrcAsNumber    uint16
    DstAsNumber    uint16
    SrcPrefixMask  uint8
    DstPrefixMask  uint8
    Padding2       uint16
}

// NetflowV5 - Netflow v5 structure
type NetflowV5 struct {
    Header  NetflowHeaderV5
    Payload []NetflowPayloadV5
}
