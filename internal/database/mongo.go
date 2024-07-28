package database

import "go.mongodb.org/mongo-driver/mongo"

func GetOne(db *mongo.Client, collection string, filter interface{}, result interface{}) error {
	err := db.Database(DatabaseName).Collection(collection).FindOne(nil, filter).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func InsertOne(db *mongo.Client, collection string, document interface{}) error {
	_, err := db.Database(DatabaseName).Collection(collection).InsertOne(nil, document)
	if err != nil {
		return err
	}
	return nil
}

func UpdateOne(db *mongo.Client, collection string, filter interface{}, update interface{}) error {
	_, err := db.Database(DatabaseName).Collection(collection).UpdateOne(nil, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func DeleteOne(db *mongo.Client, collection string, filter interface{}) error {
	_, err := db.Database(DatabaseName).Collection(collection).DeleteOne(nil, filter)
	if err != nil {
		return err
	}
	return nil
}

func GetMany(db *mongo.Client, collection string, filter interface{}, result interface{}) error {
	cur, err := db.Database(DatabaseName).Collection(collection).Find(nil, filter)
	if err != nil {
		return err
	}
	defer cur.Close(nil)
	err = cur.All(nil, result)
	if err != nil {
		return err
	}
	return nil
}
