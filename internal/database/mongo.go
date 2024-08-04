package database

func GetOne(db *Service, collection string, filter interface{}, result interface{}) error {
	err := db.Db.Collection(collection).FindOne(nil, filter).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func InsertOne(db *Service, collection string, document interface{}) error {
	_, err := db.Db.Collection(collection).InsertOne(nil, document)
	if err != nil {
		return err
	}
	return nil
}

func UpdateOne(db *Service, collection string, filter interface{}, update interface{}) error {
	_, err := db.Db.Collection(collection).UpdateOne(nil, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func DeleteOne(db *Service, collection string, filter interface{}) error {
	_, err := db.Db.Collection(collection).DeleteOne(nil, filter)
	if err != nil {
		return err
	}
	return nil
}

func GetMany(db *Service, collection string, filter interface{}, result interface{}) error {
	cur, err := db.Db.Collection(collection).Find(nil, filter)
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
