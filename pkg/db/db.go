package db

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mridulganga/dlt-manager/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userColl            = "users"
	ngColl              = "nodegroups"
	loadtestColl        = "loadtests"
	loadTestUpdatesColl = "loadtestupdates"
	ltsummaryColl       = "ltsummary"
)

type DBInterface interface{}

// Database - struct
type DB struct {
	client   *mongo.Client
	database string
}

// NewDatabase - new db obj using the connection string
func NewDatabase(connString string, dbName string) (*DB, error) {
	ctx := context.Background()
	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		return nil, fmt.Errorf("error while connecting to db %s", err.Error())
	}

	// create db struct
	db := DB{
		client:   dbClient,
		database: dbName,
	}

	return &db, nil
}

func (d DB) CreateUser(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()
	user.ID = uuid.New().String()

	collection := d.client.Database(d.database).Collection(userColl)
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d DB) GetUserByID(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	collection := d.client.Database(d.database).Collection(userColl)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d DB) GetUserByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	collection := d.client.Database(d.database).Collection(userColl)
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d DB) UpdateUser(id string, update bson.M) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	collection := d.client.Database(d.database).Collection(userColl)
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": update}, opts).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d DB) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := d.client.Database(d.database).Collection(userColl)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (d DB) ListUser() (*[]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := d.client.Database(d.database).Collection(userColl)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := []User{}

	for cursor.Next(ctx) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

func (d DB) CreateLoadTest(loadtest *LoadTest) (*LoadTest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	loadtest.StartTime = time.Now()
	loadtest.ID = uuid.New().String()

	collection := d.client.Database(d.database).Collection(loadtestColl)
	_, err := collection.InsertOne(ctx, loadtest)
	if err != nil {
		return nil, err
	}

	return loadtest, nil
}

func (d DB) GetLoadTestByID(id string) (*LoadTest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var loadtest LoadTest
	collection := d.client.Database(d.database).Collection(loadtestColl)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&loadtest)
	if err != nil {
		return nil, err
	}

	return &loadtest, nil
}

func (d DB) UpdateLoadTest(id string, update bson.M) (*LoadTest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var loadtest LoadTest
	collection := d.client.Database(d.database).Collection(loadtestColl)
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": update}, opts).Decode(&loadtest)
	if err != nil {
		return nil, err
	}

	return &loadtest, nil
}

func (d DB) DeleteLoadTest(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := d.client.Database(d.database).Collection(loadtestColl)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (d DB) ListLoadTest() (*[]LoadTest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := d.client.Database(d.database).Collection(loadtestColl)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	loadtests := []LoadTest{}

	for cursor.Next(ctx) {
		var loadtest LoadTest
		err := cursor.Decode(&loadtest)
		if err != nil {
			return nil, err
		}
		loadtests = append(loadtests, loadtest)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &loadtests, nil
}

func (d DB) CreateNodeGroup(nodegroup *NodeGroup) (*NodeGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nodegroup.ID = uuid.New().String()
	nodegroup.LastHealthCheck = time.Now()

	collection := d.client.Database(d.database).Collection(ngColl)
	_, err := collection.InsertOne(ctx, nodegroup)
	if err != nil {
		return nil, err
	}

	return nodegroup, nil
}

func (d DB) GetNodeGroupByID(id string) (*NodeGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var nodegroup NodeGroup
	collection := d.client.Database(d.database).Collection(ngColl)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&nodegroup)
	if err != nil {
		return nil, err
	}

	return &nodegroup, nil
}

func (d DB) UpdateNodeGroup(id string, update bson.M) (*NodeGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var nodegroup NodeGroup
	collection := d.client.Database(d.database).Collection(ngColl)
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": update}, opts).Decode(&nodegroup)
	if err != nil {
		return nil, err
	}

	return &nodegroup, nil
}

func (d DB) DeleteNodeGroup(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := d.client.Database(d.database).Collection(ngColl)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (d DB) ListNodeGroup() (*[]NodeGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := d.client.Database(d.database).Collection(ngColl)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	nodegroups := []NodeGroup{}

	for cursor.Next(ctx) {
		var nodegroup NodeGroup
		err := cursor.Decode(&nodegroup)
		if err != nil {
			return nil, err
		}
		nodegroups = append(nodegroups, nodegroup)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &nodegroups, nil
}

func (d DB) UpdateNodeGroupHealth(nodeGroupId string, isHealthy bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := d.client.Database(d.database).Collection(ngColl)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": nodeGroupId}, bson.M{"$set": bson.M{
		"is_healthy":       isHealthy,
		"last_health_time": time.Now(),
	}})
	if err != nil {
		return err
	}
	return nil
}

func (d DB) PushLoadTestResult(loadTestId string, result NodeUpdates) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := d.client.Database(d.database).Collection(loadTestUpdatesColl)

	for _, v := range result {
		for _, u := range v {
			nodeUpdate := NodeHeartBeat{}
			utils.DeepCopy(u, &nodeUpdate)

			loadTestResultBase64 := nodeUpdate.LoadTestResults
			loadTestResultsString, _ := base64.StdEncoding.DecodeString(loadTestResultBase64)
			loadTestResults := []string{}
			json.Unmarshal([]byte(loadTestResultsString), &loadTestResults)
			for _, res := range loadTestResults {

				singleResult := bson.M{}
				json.Unmarshal([]byte(res), &singleResult)
				singleResult["load_test_id"] = loadTestId
				singleResult["_id"] = uuid.New().String()
				// add single result to db
				_, err := collection.InsertOne(ctx, singleResult)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d DB) FetchLoadTestResults(loadTestId string) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	loadTest, _ := d.GetLoadTestByID(loadTestId)

	collection := d.client.Database(d.database).Collection(loadTestUpdatesColl)
	cursor, err := collection.Find(ctx, bson.M{"load_test_id": loadTestId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	totalRequestCount := 0
	successCount := 0
	failureCount := 0
	latencySum := 0
	failures := map[string]string{}

	for cursor.Next(ctx) {
		entry := LoadTestEntry{}
		err := cursor.Decode(&entry)
		if err != nil {
			return nil, err
		}
		totalRequestCount = totalRequestCount + 1
		latencyMs, _ := strconv.Atoi(entry.LatencyMs)
		latencySum = latencySum + latencyMs
		if entry.IsSuccess == "true" {
			successCount = successCount + 1
		} else {
			failureCount = failureCount + 1
			failures[entry.StatusCode] = entry.Response
		}
	}

	return map[string]any{
		"load_test_id":   loadTestId,
		"startTime":      loadTest.StartTime,
		"endTime":        loadTest.EndTime,
		"duration":       loadTest.Duration,
		"tps":            loadTest.TPS,
		"totalRequests":  totalRequestCount,
		"successCount":   successCount,
		"failureCount":   failureCount,
		"successPercent": (float64)(successCount / totalRequestCount),
		"avgLatencyMs":   (float64)(latencySum / totalRequestCount),
		"topFailures":    failures,
	}, nil
}

func (d DB) CreateLoadTestSummary(ltsummary LoadTestSummary) (LoadTestSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ltsummary["created_at"] = time.Now()
	ltsummary["_id"] = uuid.New().String()

	collection := d.client.Database(d.database).Collection(ltsummaryColl)
	_, err := collection.InsertOne(ctx, ltsummary)
	if err != nil {
		return nil, err
	}

	return ltsummary, nil
}

func (d DB) GetLoadTestSummaryByID(loadTestId string) (LoadTestSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var ltsummary LoadTestSummary
	collection := d.client.Database(d.database).Collection(ltsummaryColl)
	err := collection.FindOne(ctx, bson.M{"load_test_id": loadTestId}).Decode(&ltsummary)
	if err != nil {
		return nil, err
	}

	return ltsummary, nil
}
