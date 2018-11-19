package db

import (
	"context"
	"github.com/Hagelund96/paragliding/struct"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"net/http"
)
	//connects to the DB
func MongoConnect() *mongo.Client {
	// Connect to MongoDB
	conn, err := mongo.Connect(context.Background(), "mongodb://hagelund:passord123@ds145053.mlab.com:45053/paragliding", nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return conn
}
	//get all track IDs from the database
func GetTrackID(client *mongo.Client) string {

	db := client.Database("paragliding")     // `paragliding` Database
	collection := db.Collection("tracks") // `track` Collection

	cursor, err := collection.Find(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	resTrack := _struct.Track{}
	length, err := collection.Count(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	ids := "["
	i := int64(0)
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resTrack)
		if err != nil {
			log.Fatal(err)
		}
		ids += resTrack.UniqueId
		if i == length-1 {
			break
		}
		ids += ","
		i++
	}
	ids += "]"
	return ids
}
	//checks if the url already exists in the db
func UrlInDB(url string, trackColl *mongo.Collection) bool {

	cursor, err := trackColl.Find(context.Background(),
		bson.NewDocument(bson.EC.String("url", url)))
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	track := _struct.Track{}

	// Point the cursor at whatever is found
	for cursor.Next(context.Background()) {
		err = cursor.Decode(&track)
		if err != nil {
			log.Fatal(err)
		}
	}

	if track.URL == "" { // If there is an empty field, in this case, `url`, it means the track is not on the database
		return false
	}
	return true
}
	//
func GetTrack(client *mongo.Client, url string) _struct.Track {
	db := client.Database("paragliding")     // `paragliding` Database
	collection := db.Collection("tracks") // `track` Collection

	cursor, err := collection.Find(context.Background(), bson.NewDocument(bson.EC.String("url", url)))

	if err != nil {
		log.Fatal(err)
	}

	resTrack := _struct.Track{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resTrack)
		if err != nil {
			log.Fatal(err)
		}
	}

	return resTrack
}
	//return one track from the db
func GetOneTrack(client *mongo.Client, id string, w http.ResponseWriter) _struct.Track {
	db := client.Database("paragliding")     // `paragliding` Database
	collection := db.Collection("tracks") // `track` Collection
	filter := bson.NewDocument(bson.EC.String("uniqueid", ""+id+""))
	resTrack := _struct.Track{}
	err := collection.FindOne(context.Background(), filter).Decode(&resTrack)
	if err != nil {
		http.Error(w, "File not found!", 404)
	}
	return resTrack

}
	//counts all tracks in db for admin
func CountAllTracks(client *mongo.Client) int64 {
	db := client.Database("paragliding")
	collection := db.Collection("tracks")

	// Count the tracks
	count, _ := collection.Count(context.Background(), nil, nil)

	return count
}
	//admin can delete tracks in db
func DeleteAllTracks(client *mongo.Client) {
	db := client.Database("paragliding")
	collection := db.Collection("tracks")

	// Delete the tracks
	collection.DeleteMany(context.Background(), bson.NewDocument())
}