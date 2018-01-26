package main

import (
	"fmt"
	"log"
	"os/exec"
	"encoding/json"
	"os"
	"time"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
)


type peers struct {
	Peers []struct {
		ID        string   `json:"id"`
		Connected bool     `json:"connected"`
		Netaddr   []string `json:"netaddr"`
		Channels  []struct {
			State                    string `json:"state"`
			Owner                    string `json:"owner"`
			ShortChannelID           string `json:"short_channel_id"`
			FundingTxid              string `json:"funding_txid"`
			MsatoshiToUs             int    `json:"msatoshi_to_us"`
			MsatoshiTotal            int    `json:"msatoshi_total"`
			DustLimitSatoshis        int    `json:"dust_limit_satoshis"`
			MaxHtlcValueInFlightMsat uint64  `json:"max_htlc_value_in_flight_msat"`
			ChannelReserveSatoshis   int    `json:"channel_reserve_satoshis"`
			HtlcMinimumMsat          int    `json:"htlc_minimum_msat"`
			ToSelfDelay              int    `json:"to_self_delay"`
			MaxAcceptedHtlcs         int    `json:"max_accepted_htlcs"`
		} `json:"channels"`
	} `json:"peers"`
}
var (
	lightning_cli_path = os.Args[1]
	interval = 10

	lightning_peers = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "lightning",
		Name:      "peers",
		Help:      "Total number of peers",
	})
	lightning_channels = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "lightning",
		Name:      "channels",
		Help:      "Total number of channels",
	})
)

func main() {
	prometheus.MustRegister(lightning_peers)
	prometheus.MustRegister(lightning_channels)

	go pollerLoop()

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":8010", nil)
}

func pollerLoop() {
	for {
		peer_data := []byte(GetPeers())
		peers := peers{}
		total_channels := 0

		jsonErr := json.Unmarshal(peer_data, &peers)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		fmt.Printf("Peers: %v \n", len(peers.Peers))
		for _, peer := range peers.Peers {
			total_channels += len(peer.Channels)
		}
		fmt.Printf("Total channels: %v \n", total_channels)

		f_peers := float64(len(peers.Peers))
		f_channels := float64(total_channels)

		lightning_peers.Set(f_peers)
		lightning_channels.Set(f_channels)

		<-time.After(time.Duration(10) * time.Second)
	}
}

func GetPeers() []byte{
	out, err := exec.Command(lightning_cli_path, "listpeers").Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}
