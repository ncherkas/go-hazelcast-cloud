package common

import (
	"fmt"

	hazelcast "github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/config"
	"github.com/hazelcast/hazelcast-go-client/config/property"
	"github.com/hazelcast/hazelcast-go-client/core"
)

// NewHazelcastClient builds a new Hazelcast client instance
func NewHazelcastClient(clusterName string, password string, token string) (hazelcast.Client, error) {
	cfg := hazelcast.NewConfig()
	cfg.GroupConfig().SetName(clusterName)
	cfg.GroupConfig().SetPassword(password)

	discoveryCfg := config.NewCloudConfig()
	discoveryCfg.SetEnabled(true)
	discoveryCfg.SetDiscoveryToken(token)

	cfg.NetworkConfig().SetCloudConfig(discoveryCfg)

	cfg.SetProperty("hazelcast.client.cloud.url", "https://coordinator.hazelcast.cloud")
	cfg.SetProperty(property.StatisticsEnabled.Name(), "true")
	cfg.SetProperty(property.StatisticsPeriodSeconds.Name(), "1")
	cfg.SetProperty(property.HeartbeatTimeout.Name(), "300000")

	client, err := hazelcast.NewClientWithConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to build a Hazelcast Client instance: %w", err)
	}

	return client, nil
}

func GetMap(hzClient hazelcast.Client, name string) (core.Map, error) {
	m, err := hzClient.GetMap(name)
	if err != nil {
		return nil, fmt.Errorf("Failed to get a Hazelcast map '%v': %w", name, err)
	}
	return m, nil
}
