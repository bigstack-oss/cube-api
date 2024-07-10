package tuning

import (
	"context"

	definition "github.com/bigstack-oss/cube-api/internal/definition/v1"
	cubeMongo "github.com/bigstack-oss/cube-api/internal/helpers/mongo"
	log "go-micro.dev/v5/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getTuningRecords() ([]definition.Tuning, error) {
	mdb := cubeMongo.NewHelper(cubeMongo.NewDefaultConf("tuning"))
	defer mdb.Disconnect(context.Background())

	colls, err := mdb.GetAllCollections(definition.TuningDB())
	if err != nil {
		return nil, err
	}

	tunings := []definition.Tuning{}
	for _, coll := range colls {
		cursor, err := mdb.GetQueryCursor(
			definition.TuningDB(),
			coll,
			bson.M{},
		)
		if err != nil {
			log.Errorf("Failed to get cursor for %s (%s)", coll, err.Error())
			continue
		}

		appendTuningRecords(cursor, &tunings)
	}

	return tunings, nil
}

func appendTuningRecords(cursor *mongo.Cursor, tunings *[]definition.Tuning) {
	for cursor.Next(context.Background()) {
		tuning := definition.Tuning{}
		if err := cursor.Decode(&tuning); err != nil {
			log.Errorf("Failed to decode tuning record (%s)", err.Error())
			continue
		}

		*tunings = append(*tunings, tuning)
	}
	if cursor.Err() != nil {
		log.Errorf("Failed to iterate tuning cursor (%s)", cursor.Err().Error())
	}
}

func syncTuningRecord(tuning definition.Tuning) {
	mdb := cubeMongo.NewHelper(cubeMongo.NewDefaultConf("tuning"))
	defer mdb.Disconnect(context.Background())
	filter := bson.M{"node.id": tuning.Node.ID, "name": tuning.Name}
	update := bson.M{"$set": tuning}

	err := mdb.UpdateOne(
		definition.TuningDB(),
		definition.TuningCollection(tuning.Name),
		filter,
		update,
		cubeMongo.CreateRecordIfNotExist,
	)
	if err != nil {
		log.Errorf(
			"Failed to sync tuning record for %s (%s)",
			tuning.Name,
			err.Error(),
		)
	}
}

func updateRecordStatus(tuning *definition.Tuning) error {
	mdb := cubeMongo.NewHelper(cubeMongo.NewDefaultConf("tuning"))
	defer mdb.Disconnect(context.Background())

	return mdb.UpdateOne(
		definition.TuningDB(),
		definition.TuningCollection(tuning.Name),
		bson.M{"node.id": tuning.Node.ID, "name": tuning.Name},
		tuning,
		options.Update().SetUpsert(true),
	)
}
