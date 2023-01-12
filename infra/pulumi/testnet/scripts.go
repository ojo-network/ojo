package main

import (
	"strings"
)

func genStartupScript() string {
	return strings.Join([]string{
		"#!/bin/bash -xeu",
		genBaseSystemUtilitiesInstallScriptSnippet(),
		genGCPOpsAgentInstallScriptSnippet(),
		genCaddyInstallScriptSnippet(),
		"touch /tmp/STARTUP_FINISHED", // so other scripts can know when this is done
	}, "\n")
}

func genCaddyInstallScriptSnippet() string {
	return `# genCaddyInstallScriptSnippet
if [ ! -f "/etc/CADDY_INSTALLED" ]; then
  curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --yes --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
  curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
  apt update
  apt install -y caddy
  touch /etc/CADDY_INSTALLED
fi
`
}

func genGCPOpsAgentInstallScriptSnippet() string {
	return `# genGCPOpsAgentInstallScriptSnippet
if [ ! -f "/etc/GCP_OPS_AGENT_INSTALLED" ]; then
  cd /tmp
  curl -sSO https://dl.google.com/cloudagents/add-google-cloud-ops-agent-repo.sh
  bash add-google-cloud-ops-agent-repo.sh --also-install
  touch /etc/GCP_OPS_AGENT_INSTALLED
fi
`
}

func genBaseSystemUtilitiesInstallScriptSnippet() string {
	return `# genBaseSystemUtilitiesInstallScriptSnippet
# get latest updates
apt update
apt dist-upgrade -y

# install utilities
if [ ! -f "/etc/UTILITIES_INSTALLED" ]; then
    apt install -y \
    ack apt-transport-https bsdmainutils ca-certificates curl debian-keyring debian-archive-keyring iputils-ping jq less lsof nano ncat net-tools nmap sysstat telnet tmux traceroute vim
    touch /etc/UTILITIES_INSTALLED
fi
`
}
