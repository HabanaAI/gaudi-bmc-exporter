package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfiguredServer(t *testing.T) {
	servers := []serverInfo{
		{Hostname: "bmc-01:443", Username: "user", Password: "password"},
	}

	server, ok := configuredServer(servers, "bmc-01:443")
	require.True(t, ok)
	require.Equal(t, servers[0], server)

	for _, target := range []string{"", "attacker:443", "BMC-01:443", "bmc-01:443 "} {
		_, ok := configuredServer(servers, target)
		require.False(t, ok, target)
	}
}
