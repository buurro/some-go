package bluetooth

import (
	"bufio"
	"os/exec"
	"strings"
)

func FetchPairedDevices() ([]device, error) {
	devices := []device{}

	cmd := exec.Command("bluetoothctl", "devices", "Paired")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return devices, err
	}
	if err := cmd.Start(); err != nil {
		return devices, err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), " ", 3)
		if len(parts) < 3 {
			continue
		}
		devices = append(devices, device{
			address: parts[1],
			name:    parts[2],
		})
	}

	return devices, nil
}
