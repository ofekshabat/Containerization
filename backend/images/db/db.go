package db

import (
	"context"

	. ".."
	mechola "../../db"
	"go.mongodb.org/mongo-driver/bson"
)

// Returns the list of all images in the system
func List() ([]Image, error) {
	ctx := context.Background()
	cur, err := mechola.ImagesCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var images []Image
	for cur.Next(ctx) {
		var image Image
		err := cur.Decode(&image)
		if err != nil {
			panic(err)
		}
		images = append(images, image)
	}
	return images, nil
}

// Returns a single image in the system
func GetInfo(name string) (Image, error) {
	var image Image
	filter := bson.M{"name": name}
	err := mechola.ImagesCollection.FindOne(context.Background(), filter).Decode(&image)
	if err != nil {
		return image, err
	}
	return image, nil
}

// Adds an image's info to the database
func Create(image Image) error {
	data, err := bson.Marshal(&image)
	_, err = mechola.ImagesCollection.InsertOne(context.Background(), data)
	return err
}

// Deletes an image from the database
func Delete(imageName string) (error, int) {
	data, err := mechola.ImagesCollection.DeleteOne(context.Background(), bson.M{"name": imageName})
	if err != nil {
		return err, 0
	}
	return nil, int(data.DeletedCount)
}
