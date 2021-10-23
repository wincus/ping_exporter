# ping_exporter

[![Go Reference](https://pkg.go.dev/badge/github.com/blainsmith/ping_exporter.svg)](https://pkg.go.dev/github.com/blainsmith/ping_exporter)

Command `ping_exporter` provides a Prometheus exporter for ping metrics such as RTT, packet loss, and jitter to any number of hosts.

## Usage

Available flags for `ping_exporter` include:

```
$ ./ping_exporter -h
Usage of ./ping_exporter:
    -metrics.addr string
        address for ping exporter (default ":9165")
    -metrics.path string
        URL path for surfacing collected metrics (default "/metrics")
    -ping.host value
        host to ping, can be repeated (-ping.host=1.1.1.1 -ping.host=google.com ...)

$ ./ping_exporter -ping.host=1.1.1.1 -ping.host=apple.com -ping.host=google.com
2021/10/23 15:21:21 starting ping exporter on ":9137"
```

The following metrics will be exported:

```
...
# HELP ping_jitter RTT jitter.
# TYPE ping_jitter gauge
ping_jitter{host="1.1.1.1"} 3
ping_jitter{host="apple.com"} 2
ping_jitter{host="google.com"} 2
# HELP ping_packet_loss Percentage of packet loss.
# TYPE ping_packet_loss gauge
ping_packet_loss{host="1.1.1.1"} 0
ping_packet_loss{host="apple.com"} 0
ping_packet_loss{host="google.com"} 0
# HELP ping_packets_recv Number of packets received.
# TYPE ping_packets_recv counter
ping_packets_recv{host="1.1.1.1"} 65
ping_packets_recv{host="apple.com"} 65
ping_packets_recv{host="google.com"} 65
# HELP ping_packets_sent Number of packets sent.
# TYPE ping_packets_sent counter
ping_packets_sent{host="1.1.1.1"} 65
ping_packets_sent{host="apple.com"} 65
ping_packets_sent{host="google.com"} 65
# HELP ping_rtt Running average of the RTT.
# TYPE ping_rtt gauge
ping_rtt{host="1.1.1.1"} 27
ping_rtt{host="apple.com"} 25
ping_rtt{host="google.com"} 28
...
```