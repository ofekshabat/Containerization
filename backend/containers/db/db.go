package db

import (
	"context"
	"strings"

	. ".."
	mechola "../../db"
	"go.mongodb.org/mongo-driver/bson"
)

// Returns a single container's info
func GetInfo(name string) (ContainerInfo, error) {
	var containerInfo ContainerInfo
	filter := bson.M{"containername": name}
	err := mechola.ContainersCollection.FindOne(context.Background(), filter).Decode(&containerInfo)
	if err != nil {
		return ContainerInfo{}, err
	}

	return containerInfo, nil
}

// Returns a list of all containers in the system
func List() ([]ContainerInfo, error) {
	ctx := context.Background()
	cur, err := mechola.ContainersCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var containers []ContainerInfo
	for cur.Next(ctx) {
		var containerInfo ContainerInfo
		err := cur.Decode(&containerInfo)
		if err != nil {
			panic(err)
		}
		containers = append(containers, containerInfo)
	}

	return containers, nil
}

// Adds a container's info to the database
func Create(containerInfo ContainerInfo) error {
	data, err := bson.Marshal(&containerInfo)
	_, err = mechola.ContainersCollection.InsertOne(context.Background(), data)
	return err
}

// Deletes a container from the database
func Delete(containerName string) (error, int) {
	data, err := mechola.ContainersCollection.DeleteOne(context.Background(), bson.M{"containername": containerName})
	if err != nil {
		return err, 0
	}
	return nil, int(data.DeletedCount)
}

// Updates a container's info in the database
func Update(containerName string, updateInfo map[string]interface{}) (error, int) {
	updateInfoLowerCase := make(map[string]interface{})
	for k, v := range updateInfo {
		if k != "state" {
			updateInfoLowerCase[strings.ToLower(k)] = v
		}
	}
	updateResult, err := mechola.ContainersCollection.UpdateOne(context.Background(), bson.M{"containername": containerName}, bson.M{"$set": updateInfoLowerCase})
	if err != nil {
		return err, 0
	}
	return nil, int(updateResult.MatchedCount)
}
