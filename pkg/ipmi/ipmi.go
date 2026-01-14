package ipmi

import (
	"fmt"
	"strings"
)

type ClientOpts struct {
	BMCUsername string
	BMCPassword string
	BMCName     string
}

type Client struct {
	ClientOpts
}

func NewClient(opts ClientOpts) *Client {
	return &Client{
		ClientOpts: opts,
	}
}

// IsUp checks if the bmc is up.
func (c *Client) IsUp() (bool, error) {
	output, err := c.ExecCommand("power status")
	if err != nil {
		return false, err
	}

	if strings.Contains(output, "on") {
		return true, nil
	}
	return false, nil
}

func (c *Client) SelList() ([]SelInfo, error) {

	// Should return such result:
	// SEL Record ID          : 0001
	//  Record Type           : 02
	//  Timestamp             : 04/13/2023 11:50:14
	//  Generator ID          : 0020
	//  EvM Revision          : 04
	//  Sensor Type           : Event Logging Disabled
	//  Sensor Number         : a0
	//  Event Type            : Sensor-specific Discrete
	//  Event Direction       : Assertion Event
	//  Event Data            : 02ffff
	//  Description           : Log area reset/cleared

	output, err := c.ExecCommand("-v sel elist")
	if err != nil {
		return nil, fmt.Errorf("failed getting sel list: %w", err)
	}

	return parseSelEvents(output), nil
}

// ClearSel will clean the sel event log, so we won't get duplicated events
func (c *Client) ClearSel() error {
	output, err := c.ExecCommand("sel clear")
	if err != nil {
		return fmt.Errorf("clear sel record %s: %w", string(output), err)
	}
	return nil
}
