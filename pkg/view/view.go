package view

import (
	"github.com/gin-gonic/gin"
	"github.com/mridulganga/dlt-manager/pkg/db"
	"github.com/mridulganga/dlt-manager/pkg/mqttlib"
	"go.mongodb.org/mongo-driver/bson"
)

type View struct {
	d *db.DB
	m mqttlib.MqttClient
}

func NewView(database *db.DB, mqttClient mqttlib.MqttClient) View {
	return View{
		d: database,
		m: mqttClient,
	}
}

func (v View) CreateNodeGroup(c *gin.Context) {
	ng := db.NodeGroup{}
	c.BindJSON(&ng)
	result, err := v.d.CreateNodeGroup(&ng)
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func (v View) GetNodeGroup(c *gin.Context) {
	id := c.Param("id")
	result, err := v.d.GetNodeGroupByID(id)
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func (v View) UpdateNodeGroup(c *gin.Context) {
	id := c.Param("id")
	ng := bson.M{}
	c.BindJSON(&ng)

	result, err := v.d.UpdateNodeGroup(id, ng)
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func (v View) DeleteNodeGroup(c *gin.Context) {
	id := c.Param("id")
	err := v.d.DeleteNodeGroup(id)
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(200, map[string]string{"status": "ok"})
}

func (v View) ListNodeGroups(c *gin.Context) {
	results, err := v.d.ListNodeGroup()
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(200, results)
}

func (v View) CreateLoadTest(c *gin.Context) {
	lt := db.LoadTest{}
	c.BindJSON(&lt)
	result, err := v.d.CreateLoadTest(&lt)
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// get all nodegroups
	nodegroups, err := v.d.ListNodeGroup()
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// trigger load test in all node groups
	for _, ng := range *nodegroups {
		v.m.Publish(ng.Topic, map[string]any{
			"action":       "start_loadtest",
			"load_test_id": result.ID,
			"plugin_data":  lt.Logic,
			"duration":     lt.Duration,
			"tps":          lt.TPS,
		})
	}

	c.JSON(200, result)
}

func (v View) GetLoadTest(c *gin.Context) {
	id := c.Param("id")
	result, err := v.d.GetLoadTestByID(id)
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func (v View) UpdateLoadTest(c *gin.Context) {
	id := c.Param("id")
	lt := bson.M{}
	c.BindJSON(&lt)

	result, err := v.d.UpdateLoadTest(id, lt)
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func (v View) DeleteLoadTest(c *gin.Context) {
	id := c.Param("id")
	err := v.d.DeleteNodeGroup(id)
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(200, map[string]string{"status": "ok"})
}

func (v View) ListLoadTests(c *gin.Context) {
	results, err := v.d.ListLoadTest()
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(200, results)
}

func (v View) StopLoadTest(c *gin.Context) {
	// get all nodegroups
	nodegroups, err := v.d.ListNodeGroup()
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// trigger stpo load test in all node groups
	for _, ng := range *nodegroups {
		v.m.Publish(ng.Topic, map[string]any{
			"action": "stop_loadtest",
		})
	}

	c.JSON(200, map[string]string{"status": "stopping"})
}

func (v View) GetLoadTestResults(c *gin.Context) {
	id := c.Param("id")
	result, err := v.d.FetchLoadTestResults(id)
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}
