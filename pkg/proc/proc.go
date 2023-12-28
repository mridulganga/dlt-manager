package proc

import (
	"fmt"
)

// process messages and update db with nodegroup, node and load test data

type Message struct {
	Action          string   `json:"action"`
	LoadTestId      string   `json:"load_test_id"`
	NodeGroupId     string   `json:"ng_id"`
	NodeGroupStatus string   `json:"ng_status"`
	Timestamp       string   `json:"timestamp"`
	Nodes           []string `json:"nodes"`
	NodeUpdates     []string `json:"node_updates"`
}

func Proc(msg Message) error {
	if msg.Action != "ng_update" {
		return fmt.Errorf("invalid action %s", msg.Action)
	}

	// update node group health status
	// update node group node list and their health
	// update isLoadTest active
	// update load test results

	return nil
}
