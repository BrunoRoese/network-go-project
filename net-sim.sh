#!/bin/bash

set -e

# Get default network interface
IFACE=$(route get default | grep interface | awk '{print $2}')

start_netem() {
    echo "[+] Setting up UDP network simulation on interface: $IFACE"

    # Clear any previous rules
    sudo dnctl -q flush
    sudo pfctl -F all -f /etc/pf.conf

    # Configure the pipe: 200ms latency, 5% packet loss, 1Mbit/s bandwidth
    sudo dnctl pipe 1 config delay 200ms plr 0.05 bw 1Mbit/s

    # Create a temporary pf rule that applies only to UDP traffic
    echo "dummynet out proto udp from any to any pipe 1" | sudo pfctl -ef -

    echo "[✓] UDP network simulation active: 200ms delay, 5% loss, 1Mbit/s limit"
}

stop_netem() {
    echo "[+] Stopping network simulation..."

    sudo pfctl -F all -f /etc/pf.conf
    sudo dnctl -q flush

    echo "[✓] Network simulation removed"
}

case "$1" in
    start)
        start_netem
        ;;
    stop)
        stop_netem
        ;;
    *)
        echo "Usage: $0 {start|stop}"
        exit 1
        ;;
esac

