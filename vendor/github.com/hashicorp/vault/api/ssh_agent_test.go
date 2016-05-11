package api

import (
	"fmt"
	"strings"
	"testing"
)

func TestSSH_CreateTLSClient(t *testing.T) {
	// load the default configuration
	config, err := LoadSSHHelperConfig("./test-fixtures/agent_config.hcl")
	if err != nil {
		panic(fmt.Sprintf("error loading agent's config file: %s", err))
	}

	client, err := config.NewClient()
	if err != nil {
		panic(fmt.Sprintf("error creating the client: %s", err))
	}

	// Provide a certificate and enforce setting of transport
	config.CACert = "./test-fixtures/vault.crt"

	client, err = config.NewClient()
	if err != nil {
		panic(fmt.Sprintf("error creating the client: %s", err))
	}
	if client.config.HttpClient.Transport == nil {
		panic(fmt.Sprintf("error creating client with TLS transport"))
	}
}

func TestParseSSHHelperConfig(t *testing.T) {
	config, err := ParseSSHHelperConfig(`
		vault_addr = "1.2.3.4"
`)
	if err != nil {
		t.Fatal(err)
	}

	if config.SSHMountPoint != SSHHelperDefaultMountPoint {
		t.Errorf("expected %q to be %q", config.SSHMountPoint, SSHHelperDefaultMountPoint)
	}
}

func TestParseSSHHelperConfig_missingVaultAddr(t *testing.T) {
	_, err := ParseSSHHelperConfig("")
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "ssh_helper: missing config 'vault_addr'") {
		t.Errorf("bad error: %s", err)
	}
}

func TestParseSSHHelperConfig_badKeys(t *testing.T) {
	_, err := ParseSSHHelperConfig(`
vault_addr = "1.2.3.4"
nope = "bad"
`)
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "ssh_helper: invalid key 'nope' on line 3") {
		t.Errorf("bad error: %s", err)
	}
}