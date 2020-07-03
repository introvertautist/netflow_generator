package main

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "math/rand"
    "net"
    "os"
    "sync"
    "time"
)

func createNFlowHeader(flowCount int, flowSequence uint32, startTime int64) NetflowHeaderV5 {
    timeNow := time.Now().UnixNano()
    timeNowSec := timeNow / int64(time.Second)
    timeNowMSec := timeNow - timeNowSec*int64(time.Second)
    sysUptime := uint32((timeNow-startTime)/int64(time.Millisecond)) + 1000

    h := new(NetflowHeaderV5)
    h.Version = 5
    h.FlowCount = uint16(flowCount)
    h.SysUptime = sysUptime
    h.UnixSec = uint32(timeNowSec)
    h.UnixMsec = uint32(timeNowMSec)
    h.FlowSequence = flowSequence
    h.EngineType = 1
    h.EngineID = 0
    h.SampleInterval = 0

    return *h
}

func getNetflowPayloadV5(config *AppConfig, srcIPList *[]string, dstIPList *[]string, flowSequence *uint32) []NetflowPayloadV5 {
    payload := make([]NetflowPayloadV5, config.Flows)

    for i := 0; i < config.Flows; i++ {
        srcAddr := GetRandonIP(*srcIPList)
        dstAddr := GetRandonIP(*dstIPList)
        numPackets := RandomUint32(config.MaxPackets)
        numOctets := RandomUint32(config.MaxPacketSize)
        srcPort := uint16(config.SourcePort)
        dstPort := uint16(config.DestinationPort)

        payload[i].SrcIP = IPtoUint32(srcAddr)
        payload[i].DstIP = IPtoUint32(dstAddr)
        payload[i].NextHopIP = IPtoUint32("172.199.15.1") //FIXME: To options
        payload[i].SrcPort = srcPort
        payload[i].DstPort = dstPort
        payload[i].SnmpInIndex = RandomUint16(65535)  //FIXME: To options
        payload[i].SnmpOutIndex = RandomUint16(65535) //FIXME: To options
        payload[i].NumPackets = numPackets
        payload[i].NumOctets = numOctets
        payload[i].SysUptimeStart = rand.Uint32()
        payload[i].SysUptimeEnd = rand.Uint32()
        payload[i].Padding1 = 0
        payload[i].IPProtocol = 6
        payload[i].IPTos = 0
        payload[i].SrcPrefixMask = uint8(rand.Intn(32))
        payload[i].DstPrefixMask = uint8(rand.Intn(32))
        payload[i].Padding2 = 0

        totalBytes := numPackets * numOctets
        srcAddrFmt := fmt.Sprintf("%s:%d", srcAddr, srcPort)
        dstAddrFmt := fmt.Sprintf("%s:%d", dstAddr, dstPort)

        fmt.Printf(Debug(
            fmt.Sprintf("\r%10d %20s %20s %5d %8d", *flowSequence, srcAddrFmt, dstAddrFmt, config.Flows, totalBytes)))

    }

    return payload
}

func buildNFlowPayload(data NetflowV5) bytes.Buffer {
    buffer := new(bytes.Buffer)

    err := binary.Write(buffer, binary.BigEndian, &data.Header)
    if err != nil {
        fmt.Printf(Error("\nWriting netflow header failed: %s\n", err))
    }

    for _, record := range data.Payload {
        err := binary.Write(buffer, binary.BigEndian, &record)
        if err != nil {
            fmt.Printf(Error("\nWriting netflow record failed: %s\n", err))
        }
    }

    return *buffer
}

func worker(connection *net.UDPConn, config *AppConfig, flowSequence *uint32, startTime int64, srcIPRange *[]string, dstIPRange *[]string, wg *sync.WaitGroup, m *sync.Mutex, results chan<- int) {
    for {
        m.Lock()
        *flowSequence = *flowSequence + 1
        m.Unlock()

        data := new(NetflowV5)
        header := createNFlowHeader(config.Flows, *flowSequence, startTime)
        payload := []NetflowPayloadV5{}

        payload = getNetflowPayloadV5(config, srcIPRange, dstIPRange, flowSequence)

        data.Header = header
        data.Payload = payload

        // buffer := buildNFlowPayload(*data)

        // _, err := connection.Write(buffer.Bytes())
        // if err != nil {
        //     fmt.Printf(Error("\nError connecting to the collector: %s\n", err))
        //     os.Exit(1)
        // }

        if config.Delay != 0 {
            time.Sleep(time.Duration(config.Delay) * time.Millisecond)
        }
    }

    wg.Done()
}

/*
NetflowV5Generator - Generate netflow v5 trafic
Connect to collector and send data
Run n goroutins, where n is config.Workers count
*/
func NetflowV5Generator(config *AppConfig) {
    collectorAddr := fmt.Sprintf("%s:%d", config.CollectorHost, config.CollectorPort)
    udpAddr, err := net.ResolveUDPAddr("udp", collectorAddr)

    if err != nil {
        fmt.Printf(Error("Resolve error: %s\n", err))
        os.Exit(1)
    }

    connection, err := net.DialUDP("udp", nil, udpAddr)
    if err != nil {
        fmt.Printf(Error("Error connecting to the collector: %s\n", err))
        os.Exit(1)
    }

    // Use for detect device uptime
    startTime := time.Now().UnixNano()

    // Use for count packets
    var flowSequence uint32 = 1

    fmt.Printf(Debug2(
        fmt.Sprintf("%10s %20s %20s %5s %8s\n", "TotalRows", "Src", "Dst", "Flows", "Pkg.Size")))

    rand.Seed(time.Now().Unix())

    // Get IP ranges
    srcIPRange := GetIPRange(config.Source)
    dstIPRange := GetIPRange(config.Destination)

    var wg sync.WaitGroup
    var m sync.Mutex

    results := make(chan int, config.Workers)
    for wk := 1; wk <= config.Workers; wk++ {
        wg.Add(1)
        go worker(connection, config, &flowSequence, startTime, &srcIPRange, &dstIPRange, &wg, &m, results)
    }

    for {
        <-results
    }

}
