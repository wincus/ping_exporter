package pingexporter

import (
	"github.com/go-ping/ping"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &collector{}

type collector struct {
	PacketsRecv *prometheus.Desc
	PacketsSent *prometheus.Desc

	RTT        *prometheus.Desc
	PacketLoss *prometheus.Desc
	Jitter     *prometheus.Desc

	pingers []*ping.Pinger
}

func NewCollector(pingers []*ping.Pinger) prometheus.Collector {
	return &collector{
		PacketsRecv: prometheus.NewDesc(
			"ping_packets_recv",
			"Number of packets received.",
			[]string{"host"},
			nil,
		),

		PacketsSent: prometheus.NewDesc(
			"ping_packets_sent",
			"Number of packets sent.",
			[]string{"host"},
			nil,
		),

		RTT: prometheus.NewDesc(
			"ping_rtt",
			"Running average of the RTT.",
			[]string{"host"},
			nil,
		),

		PacketLoss: prometheus.NewDesc(
			"ping_packet_loss",
			"Percentage of packet loss.",
			[]string{"host"},
			nil,
		),

		Jitter: prometheus.NewDesc(
			"ping_jitter",
			"RTT jitter.",
			[]string{"host"},
			nil,
		),

		pingers: pingers,
	}
}

// Describe implements prometheus.Collector.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.PacketsRecv,
		c.PacketsSent,
		c.RTT,
		c.PacketLoss,
		c.Jitter,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect implements prometheus.Collector.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	for _, pinger := range c.pingers {
		stats := pinger.Statistics()

		ch <- prometheus.MustNewConstMetric(c.PacketsSent, prometheus.CounterValue, float64(stats.PacketsSent), pinger.Addr())
		ch <- prometheus.MustNewConstMetric(c.PacketsRecv, prometheus.CounterValue, float64(stats.PacketsRecv), pinger.Addr())
		ch <- prometheus.MustNewConstMetric(c.RTT, prometheus.GaugeValue, float64(stats.AvgRtt.Milliseconds()), pinger.Addr())
		ch <- prometheus.MustNewConstMetric(c.PacketLoss, prometheus.GaugeValue, float64(stats.PacketLoss), pinger.Addr())
		ch <- prometheus.MustNewConstMetric(c.Jitter, prometheus.GaugeValue, float64(stats.StdDevRtt.Milliseconds()), pinger.Addr())
	}
}
