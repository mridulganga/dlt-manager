package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Data map[string]any
type NodeUpdates map[string][]Data

type User struct {
	ID        string    `bson:"_id" json:"_id"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"password"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type LoadTest struct {
	ID          string    `bson:"_id" json:"_id,omitempty"`
	Description string    `bson:"description" json:"description"`
	TPS         float64   `bson:"tps" json:"tps"`
	Duration    int       `bson:"duration" json:"duration"`
	Logic       string    `bson:"logic" json:"logic"`
	CreatedBy   string    `bson:"created_by" json:"created_by"`
	StartTime   time.Time `bson:"start_time" json:"start_time"`
	EndTime     time.Time `bson:"end_time" json:"end_time"`
	Status      string    `bson:"status" json:"status"`
}

type NodeGroup struct {
	ID              string    `bson:"_id" json:"_id"`
	Nodes           []string  `bson:"nodes"`
	Topic           string    `bson:"topic"`
	IsHealthy       bool      `bson:"is_healthy"`
	LastHealthCheck time.Time `bson:"last_health_time"`
}

type LoadTestSummary bson.M

type NGHeartbeat struct {
	Action           string   `json:"action"`
	NodeGroupStatus  string   `json:"ng_status"`
	NodeGroupID      string   `json:"ng_id"`
	Nodes            []string `json:"nodes"`
	IsLoadTestActive bool     `json:"isLoadTestActive"`
	Timestamp        string   `json:"timestamp"`
	NodeUpdates      string   `json:"node_updates,omitempty"`
	LoadTestId       string   `json:"load_test_id,omitempty"`
}

type NodeHeartBeat struct {
	Action          string `json:"action"`
	IsTestActive    string `json:"isTestActive"`
	LoadTestID      string `json:"load_test_id"`
	LoadTestResults string `json:"load_test_results"`
	NodeID          string `json:"node_id"`
	NodeStatus      string `json:"node_status"`
	Timestamp       string `json:"timestamp"`
}

type LoadTestEntry struct {
	IsSuccess  string `bson:"isSuccess"`
	LatencyMs  string `bson:"latencyMs"`
	Response   string `bson:"response"`
	StatusCode string `bson:"statusCode"`
}
