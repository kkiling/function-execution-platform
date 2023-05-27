package mongodb

import "go.mongodb.org/mongo-driver/bson"

func objectToBsonMap(obj interface{}) (bson.M, error) {
	data, err := bson.Marshal(obj)
	if err != nil {
		return nil, err
	}

	result := bson.M{}
	err = bson.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
