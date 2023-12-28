package main

import (
	"encoding/json"
	"os"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mridulganga/dlt-manager/pkg/db"
	"github.com/mridulganga/dlt-manager/pkg/mqttlib"
	"github.com/mridulganga/dlt-manager/pkg/view"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	MQTT_HOST = "MQTT_HOST"
	MQTT_PORT = "MQTT_PORT"
	MONGO     = "MONGO"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	mqttHost := os.Getenv(MQTT_HOST)
	mqttPort, _ := strconv.Atoi(os.Getenv(MQTT_PORT))
	mongo := os.Getenv(MONGO)
	dbName := os.Getenv("DB_NAME")
	logrus.Infof("mqtt host %s port %v", mqttHost, mqttPort)

	lastActiveLoadTestId := ""

	m := mqttlib.NewMqtt(mqttHost, mqttPort)
	go m.Connect()
	m.WaitUntilConnected()

	d, err := db.NewDatabase(mongo, dbName)
	if err != nil {
		panic(err)
	}

	m.Sub("manager", func(client mqtt.Client, message mqtt.Message) {
		data := db.NGHeartbeat{}
		json.Unmarshal(message.Payload(), &data)

		switch data.Action {
		case "ng_update":
			logrus.Info("processing ng_update")
			isNGHealthy := data.NodeGroupStatus == "healthy"

			// update ng health db collection
			err := d.UpdateNodeGroupHealth(data.NodeGroupID, isNGHealthy)
			if err != nil {
				logrus.Errorf("error while UpdateNodeGroupHealth %v", err.Error())
			}

			// update node list if ng healthy
			if isNGHealthy {
				d.UpdateNodeGroup(data.NodeGroupID, bson.M{"nodes": data.Nodes})
			}

			// check if lt active
			if data.IsLoadTestActive {
				// get lt results and put in db
				nodeUpdates := db.NodeUpdates{}
				json.Unmarshal([]byte(data.NodeUpdates), &nodeUpdates)
				d.PushLoadTestResult(data.LoadTestId, nodeUpdates)
				lastActiveLoadTestId = data.LoadTestId
			}

			// check if lt inactive
			if !data.IsLoadTestActive && lastActiveLoadTestId != "" {
				lt, _ := d.GetLoadTestByID(lastActiveLoadTestId)
				if lt.Status == "running" {
					d.UpdateLoadTest(lastActiveLoadTestId, bson.M{"status": "complete"})
				}
				// 	consolidate results and push to result collection
				ltSummary, _ := d.FetchLoadTestResults(lastActiveLoadTestId)
				d.CreateLoadTestSummary(ltSummary)
			}

		}

	})

	vi := view.NewView(d, m)

	r := gin.New()
	r.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/"),
		gin.Recovery(),
	)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "running",
		})
	})

	g := r.Group("/api")

	g.GET("/ngs", vi.ListNodeGroups)
	g.GET("/ngs/:id", vi.GetNodeGroup)
	g.PATCH("/ngs/:id", vi.UpdateLoadTest)
	g.PUT("/ngs", vi.CreateNodeGroup)
	g.DELETE("/ngs/:id", vi.DeleteNodeGroup)

	g.GET("/loadtests", vi.ListLoadTests)
	g.GET("/loadtests/:id", vi.GetLoadTest)
	g.PATCH("/loadtests/:id", vi.UpdateLoadTest)
	g.PUT("/loadtests", vi.CreateLoadTest)
	g.DELETE("/loadtests/:id", vi.DeleteLoadTest)

	g.PUT("/loadtests/stop", vi.StopLoadTest)
	g.GET("/loadtests/:id/results", vi.GetLoadTestResults)

	r.Run() // listen and serve on 0.0.0.0:8080
}
