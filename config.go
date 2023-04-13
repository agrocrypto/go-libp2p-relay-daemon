package relaydaemon

import (
	"encoding/json"
	"os"
	"time"

	relayv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

// Config stores the full configuration of the relays, ACLs and other settings
// that influence behaviour of a relay daemon.
type Config struct {
	Network NetworkConfig
	ConnMgr ConnMgrConfig
	RelayV2 RelayV2Config
	ACL     ACLConfig
	Daemon  DaemonConfig
}

// DaemonConfig controls settings for the relay-daemon itself.
type DaemonConfig struct {
	PprofPort int
}

// NetworkConfig controls listen and annouce settings for the libp2p host.
type NetworkConfig struct {
	ListenAddrs   []string
	AnnounceAddrs []string
}

// ConnMgrConfig controls the libp2p connection manager settings.
type ConnMgrConfig struct {
	ConnMgrLo    int
	ConnMgrHi    int
	ConnMgrGrace time.Duration
}

// RelayV1Config controls activation of V1 circuits and resouce configuration
// for them.
// type RelayV1Config struct {
// 	Enabled   bool
// 	Resources relayv1.Resources
// }

// RelayV2Config controls activation of V2 circuits and resouce configuration
// for them.
type RelayV2Config struct {
	Enabled   bool
	Resources relayv2.Resources
}

// ACLConfig provides filtering configuration to allow specific peers or
// subnets to be fronted by relays. In V2, this specifies the peers/subnets
// that are able to make reservations on the relay. In V1, this specifies the
// peers/subnets that can be contacted through the relays.
type ACLConfig struct {
	AllowPeers   []string
	AllowSubnets []string
}

// DefaultConfig returns a default relay configuration using default resource
// settings and no ACLs.
func DefaultConfig() Config {
	return Config{
		Network: NetworkConfig{
			ListenAddrs: []string{
				"/ip4/0.0.0.0/udp/9095/quic-v1",
				"/ip4/0.0.0.0/udp/9095/quic-v1/webtransport",
				"/ip4/0.0.0.0/tcp/4001",
				"/ip4/0.0.0.0/tcp/9096/ws",
				"/ip4/127.0.0.1/udp/9095/quic-v1/webtransport",
				"/ip4/127.0.0.1/udp/9095/quic-v1/webtransport/certhash/uEiAaP2zrOyYeIFmagpOQg0K_6R4eD6aPxrZBrXzRnsVNUQ/certhash/uEiA4yfEqqYgLIDMaoZFAOEUDjyFL6YYHj3Wc7tf9ll-atg/p2p/12D3KooWEbomrRWemnfMBfgNMHdQmduTYZGmmLzTUXydfwU1iohy",
			},
			AnnounceAddrs: []string{
				"/ip4/127.0.0.1/udp/9095/quic-v1/webtransport",
				// "/ip4/127.0.0.1/udp/9095/quic-v1/webtransport/certhash/uEiAaP2zrOyYeIFmagpOQg0K_6R4eD6aPxrZBrXzRnsVNUQ/certhash/uEiA4yfEqqYgLIDMaoZFAOEUDjyFL6YYHj3Wc7tf9ll-atg/p2p/12D3KooWEbomrRWemnfMBfgNMHdQmduTYZGmmLzTUXydfwU1iohy",
				// "/ip4/127.0.0.1/tcp/9096/ws",
			},
		},
		ConnMgr: ConnMgrConfig{
			ConnMgrLo:    512,
			ConnMgrHi:    768,
			ConnMgrGrace: 2 * time.Minute,
		},
		// RelayV1: RelayV1Config{
		// 	Enabled:   false,
		// 	Resources: relayv1.DefaultResources(),
		// },
		RelayV2: RelayV2Config{
			Enabled:   true,
			Resources: relayv2.DefaultResources(),
		},
		Daemon: DaemonConfig{
			PprofPort: 6060,
		},
	}
}

// LoadConfig reads a relay daemon JSON configuration from the given path.
// The configuration is first initialized with DefaultConfig, so all unset
// fields will take defaults from there.
func LoadConfig(cfgPath string) (Config, error) {
	cfg := DefaultConfig()

	if cfgPath != "" {
		cfgFile, err := os.Open(cfgPath)
		if err != nil {
			return Config{}, err
		}
		defer cfgFile.Close()

		decoder := json.NewDecoder(cfgFile)
		err = decoder.Decode(&cfg)
		if err != nil {
			return Config{}, err
		}
	}

	return cfg, nil
}
