package ipmi

import (
	"fmt"
	"os/exec"
	"strings"
)

// ExecCommand will execute commands using the ipmitool.
func (c *Client) ExecCommand(command string) (string, error) {
	args := fmt.Sprintf("ipmitool -I lanplus -H %s -U %s -P %s %s", c.BMCName, c.BMCUsername, c.BMCPassword, command)
	cmd := exec.Command("/bin/sh", "-c", args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %w", output, err)
	}

	str := string(output)

	return strings.TrimSpace(str), nil
}

func parseSelEvents(s string) []SelInfo {
	var res [][]string
	fields := strings.Split(s, "\n")
	for i := 0; i < len(fields); i++ {

		// start of a new id
		if strings.Contains(fields[i], "SEL Record ID") {
			var selEntry []string
			for j := i; j < len(fields); j++ {
				// In case of end of entry:
				// 1. we get to the start of another record.
				// 2. we get to the end of the
				if (strings.Contains(fields[j], "SEL Record ID") && i != j) || j+1 == len(fields) {
					res = append(res, selEntry)
					break
				}

				fields[j] = strings.TrimSpace(fields[j])
				selEntry = append(selEntry, fields[j])
			}
		}
	}

	var selEvents []SelInfo

	for _, sels := range res {
		var selData SelInfo
		for _, sel := range sels {
			if len(strings.SplitN(sel, ":", 2)) != 2 {
				continue
			}

			data := strings.SplitN(sel, ":", 2)
			data[0] = strings.TrimSpace(data[0])
			data[1] = strings.TrimSpace(data[1])

			switch data[0] {
			case "SEL Record ID":
				selData.RecordID = data[1]
			case "Record Type":
				selData.RecordType = data[1]
			case "Timestamp":
				selData.Timestamp = data[1]
			case "Generator ID":
				selData.GeneratorID = data[1]
			case "EvM Revision":
				selData.EvmRevision = data[1]
			case "Sensor Type":
				selData.SensorType = data[1]
			case "Sensor Number":
				selData.SensorNumber = data[1]
			case "Event Type":
				selData.EventType = data[1]
			case "Event Direction":
				selData.EventDirection = data[1]
			case "Event Data":
				selData.EventData = data[1]
			case "Description":
				selData.Description = data[1]
			}
		}
		selEvents = append(selEvents, selData)
	}

	return selEvents
}
