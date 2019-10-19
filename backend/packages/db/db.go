package db

import (
	"context"
	"errors"
	"fmt"

	. ".."
	mechola "../../db"
	"go.mongodb.org/mongo-driver/bson"
)

// Returns a list of all packages in the system
func ListPackages() ([]Package, error) {
	ctx := context.Background()
	cur, err := mechola.PackagesCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var packages []Package
	for cur.Next(ctx) {
		var _package Package
		err := cur.Decode(&_package)
		if err != nil {
			panic(err)
		}
		packages = append(packages, _package)
	}
	return packages, nil
}

// Returns a single package in the system
func GetPackageInfo(name string) (Package, error) {
	var packageInfo Package
	filter := bson.M{"name": name}
	err := mechola.PackagesCollection.FindOne(context.Background(), filter).Decode(&packageInfo)
	if err != nil {
		return Package{}, err
	}
	return packageInfo, nil
}

// Adds a package to the database
func CreatePackage(packageInfo Package) error {
	data, err := bson.Marshal(&packageInfo)
	_, err = mechola.PackagesCollection.InsertOne(context.Background(), data)
	return err
}

// Deletes a package from the database
func DeletePackage(packageName string) error {
	data, err := mechola.PackagesCollection.DeleteOne(context.Background(), bson.M{"name": packageName})
	if err != nil {
		return err
	}
	if data.DeletedCount == 0 {
		fmt.Println("Deleted 0")
		return errors.New("Couldn't Delete")
	}
	return nil
}
