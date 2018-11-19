package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hagelund96/paragliding/struct"
	"github.com/Hagelund96/paragliding/db"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//checks if the given id exists
func checkId(id string) bool {
	idExists := false
	for i := 0; i < len(_struct.IDs); i++ {
		if _struct.IDs[i] == strings.ToUpper(id) {
			idExists = true
			break
		}
	}
	return idExists
}

//checks if url given is valid
func checkURL(u string) bool {
	check, _ := regexp.MatchString("^(http://skypolaris.org/wp-content/uploads/IGS%20Files/)(.*?)(%20)(.*?)(.igc)$", u)
	if check == true {
		return true
	}
	return false
}

//parses ids into json, and encodes and outputs whole array of ids
func replyWithAllTracksId(w http.ResponseWriter, db _struct.TrackDB) {
	http.Header.Set(w.Header(), "content-type", "application/json")
	if len(_struct.IDs) == 0 {
		_struct.IDs = make([]string, 0)
	}
	json.NewEncoder(w).Encode(_struct.IDs)
	return
}

//parses id into json, and encodes and outputs whole track mapped to id
/*func replyWithTracksId(w http.ResponseWriter, db _struct.TrackDB, id string) {
	http.Header.Set(w.Header(), "content-type", "application/json")
	t, _ := db.Get(strings.ToUpper(id))
	api := _struct.Track{t.UniqueId, t.HeaderDate, t.Pilot, t.Glider, t.GliderId, t.TrackLength, }
	json.NewEncoder(w).Encode(api)
}

//parses field into json, and encodes and outputs it
func replyWithField(w http.ResponseWriter, db _struct.TrackDB, id string, field string) {
	http.Header.Set(w.Header(), "content-type", "application/json")
	t, _ := db.Get(strings.ToUpper(id))

	api := _struct.Track{t.UniqueId, t.HeaderDate, t.Pilot, t.Glider, t.GliderId, t.TrackLength}

	switch strings.ToUpper(field) {
	case "PILOT":
		json.NewEncoder(w).Encode(api.Pilot)
	case "GLIDER":
		json.NewEncoder(w).Encode(api.Glider)
	case "GLIDER_ID":
		json.NewEncoder(w).Encode(api.GliderId)
	case "TRACK_LENGTH":
		json.NewEncoder(w).Encode(api.TrackLength)
	case "H_DATE":
		json.NewEncoder(w).Encode(api.HeaderDate)
	default:
		http.Error(w, "Not a valid option", http.StatusNotFound)
		return
	}
}*/
	//redirect /igcinfo to /igcinfo/api
func Handler(w http.ResponseWriter, r *http.Request) {

	//Handling for /igcinfo and for /<rubbish>
	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	// Redirect to /paragliding/api
	http.Redirect(w, r, "/paragliding/api", 302)
}

/*func HandlerIgc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//handling POST /api/igc/   **NEED TO HAVE SLASH, DOES NOT WORK WIHTOUT**
	case "POST":
		//checks that input is not empty
		if r.Body == nil {
			http.Error(w, "Missing body", http.StatusBadRequest)
			return
		}
		var u _struct.URL
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//checks if url is valid
		if checkURL(u.URL) == false {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		track, err := igc.ParseLocation(u.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//calculates total distance
		totalDistance := _struct.CalculatedDistance(track)
		var i _struct.ID
		i.ID = ("ID" + strconv.Itoa(_struct.LastUsed))
		t := _struct.Track{track.UniqueID, track.Header.Date,
			track.Pilot,
			track.GliderType,
			track.GliderID,
			totalDistance}
		//counts up last used
		_struct.LastUsed++
		_struct.Db.Add(t, i)
		return
		//Handling all GETs after /api/
	case "GET":
		parts := strings.Split(r.URL.Path, "/")

		switch len(parts) {
		//handling /api/igc/ and /api/igc/id
		case 5:
			if parts[4] == "" {
				replyWithAllTracksId(w, _struct.Db)
			} else if checkId(parts[4]) {
				replyWithTracksId(w, _struct.Db, parts[4])
			} else {
				http.Error(w, http.StatusText(404), 404)
			}
		case 6:
			//handling /api/igc/id/ and /api/igc/id/field

			if parts[5] == "" {
				if !checkId(parts[4]) /*!idExists*/ //{
				/*	http.Error(w, "ID out of range.", http.StatusNotFound)
					return
				} else {
					replyWithTracksId(w, _struct.Db, parts[4])
				}
			} else {
				if checkId(parts[4]) {
					replyWithField(w, _struct.Db, parts[4], parts[5])
				} else {
					http.Error(w, "Not a valid request", http.StatusBadRequest)
				}
			}
			//handling /api/igc/id/field/
		case 7:
			if parts[6] == "" {
				if checkId(parts[4]) {
					replyWithField(w, _struct.Db, parts[4], parts[5])
				} else {
					http.Error(w, "Not a valid request", http.StatusBadRequest)
				}
			} else {
				http.Error(w, "Not a valid request", http.StatusBadRequest)
			}
		}

		//if instead of case, left behind to show working progress
		/*if len(parts) == 5 {
				if parts[4] == ""{
					replyWithAllTracksId(w, _struct.Db)
				} else {
					http.Error(w, http.StatusText(404), 404)
				}
			} else if (parts[5] == "" && len(parts) == 6) || len(parts) == 5 {
				if !checkId(parts[4]) {
					http.Error(w, "ID out of range.", http.StatusNotFound)
					return
				} else {
					replyWithTracksId(w, _struct.Db, parts[4])
				}
			} else if (parts[6] == "" && len(parts) == 7) || len(parts) == 6 {
				if checkId(parts[4]) {
					replyWithField(w, _struct.Db, parts[4], parts[5])
				} else {
					http.Error(w, "Not a valid request", http.StatusBadRequest)
				}
			} else {
				http.Error(w, "Not a valid request", http.StatusBadRequest)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return*/
	//}
//}

//handler for /api shows uptime description and version in json
func HandlerApi(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) == 4 && parts[3] == "" {
		api := _struct.Information{_struct.Uptime(), _struct.Description, _struct.Version}
		json.NewEncoder(w).Encode(api)
	} else {
		http.Error(w, http.StatusText(404), 404)
	}
}

//Handling for /paragliding/api/track
func HandlerTrack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	//Handling GET /paragliding/api/track for all ids  in database
	case http.MethodGet:

		client := db.MongoConnect()

		ids := db.GetTrackID(client)

		fmt.Fprint(w, ids)

	case http.MethodPost:

		//handling post /igcinfo/api/igc for sending a url and returning an id for that url
		pattern := ".*.igc"

		URL := &_struct.URL{}

		var error = json.NewDecoder(r.Body).Decode(URL)
		if error != nil {
			fmt.Fprintln(w, "Error!! ", error)
			return
		}
		res, err := regexp.MatchString(pattern, URL.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if res {
			ID := rand.Intn(1000)

			track, err := igc.ParseLocation(URL.URL)
			if err != nil {
				fmt.Fprintln(w, "Error made: ", err)
				return
			}

			track.UniqueID = strconv.Itoa(ID)

			trackFile := _struct.Track{}

			timestamp := time.Now().Second()
			timestamp = timestamp * 1000

			totalDistance := _struct.CalculatedDistance(track)

			client := db.MongoConnect()

			collection := client.Database("paragliding").Collection("tracks")

			// Checking for duplicates so that the user doesn't add into the database igc files with the same URL
			duplicate := db.UrlInDB(URL.URL, collection)

			if !duplicate {

				trackFile = _struct.Track{
				track.UniqueID, track.Header.Date, track.Pilot, track.GliderType, track.GliderID, totalDistance, URL.URL}

				res, err := collection.InsertOne(context.Background(), trackFile)
				if err != nil {
					log.Fatal(err)
				}

				id := res.InsertedID

				if id == nil {
					http.Error(w, "", 300)
				}

				// Encoding the ID of the track that was just added to DB
				fmt.Fprint(w, "{\n\"id\":\""+track.UniqueID+"\"\n}")

			} else {

				trackInDB := db.GetTrack(client, URL.URL)
				// If there is another file in igcFilesDB with that URL return and tell the user that that IGC FILE is already in the database
				http.Error(w, "409 already exists!", http.StatusConflict)
				fmt.Fprintln(w, "\nID in db: ", trackInDB.UniqueId)
				return

			}

		}
	default:
		http.Error(w, "Not implemented", http.StatusNotImplemented)
		return
	}

}

func HandlerID(w http.ResponseWriter, r *http.Request) {
	//Handling /igcinfo/api/igc/<id>
	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	idURL := mux.Vars(r)

	rNum, _ := regexp.Compile(`[0-9]+`)
	if !rNum.MatchString(idURL["id"]) {
		http.Error(w, "400 - Bad Request", http.StatusBadRequest)
		return
	}

	client := db.MongoConnect()

	collection := client.Database("paragliding").Collection("tracks")

	cursor, err := collection.Find(context.Background(), nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	track := _struct.Track{}


	for cursor.Next(context.Background()) {
		err = cursor.Decode(&track)
		if err != nil {
			log.Fatal(err)
		}

		if track.UniqueId == idURL["id"] {
			http.Header.Set(w.Header(), "content-type", "application/json")
			api := _struct.Track{track.UniqueId, track.HeaderDate, track.Pilot, track.Glider, track.GliderId, track.TrackLength, track.URL}
			json.NewEncoder(w).Encode(api)

		} else {
			//Handling if user type different id from ids stored
			http.Error(w, "404 - ID not in database ", http.StatusNotFound)

		}

	}

}

func HandlerField(w http.ResponseWriter, r *http.Request) {

	//Handling for GET /api/igc/<id>/<field>
	w.Header().Set("Content-Type", "application/json")

	urlFields := mux.Vars(r)

	var rNum, _ = regexp.Compile(`[a-zA-Z_]+`)

	//attributes := &Attributes{}

	// Regular Expression for IDs

	regExID, _ := regexp.Compile("[0-9]+")

	if !regExID.MatchString(urlFields["id"]) {
		http.Error(w, "400 - Bad Request, invalid ID.", http.StatusBadRequest)
		return
	}

	if !rNum.MatchString(urlFields["field"]) {
		http.Error(w, "400 - Bad Request, wrong parameters", http.StatusBadRequest)
		return
	}
	client := db.MongoConnect()

	trackDB := _struct.Track{}

	trackDB = db.GetOneTrack(client, urlFields["id"], w)
	// Taking the field variable from the URL path and converting it to lower case to skip some potential errors
	field := urlFields["field"]

	http.Header.Set(w.Header(), "content-type", "application/json")

	api := _struct.Track{trackDB.UniqueId, trackDB.HeaderDate, trackDB.Pilot, trackDB.Glider, trackDB.GliderId, trackDB.TrackLength, trackDB.URL}

	switch strings.ToUpper(field) {
	case "PILOT":
		json.NewEncoder(w).Encode(api.Pilot)
	case "GLIDER":
		json.NewEncoder(w).Encode(api.Glider)
	case "GLIDER_ID":
		json.NewEncoder(w).Encode(api.GliderId)
	case "TRACK_LENGTH":
		json.NewEncoder(w).Encode(api.TrackLength)
	case "H_DATE":
		json.NewEncoder(w).Encode(api.HeaderDate)
	case "URL":
		json.NewEncoder(w).Encode(api.URL)
	default:
		http.Error(w, "Not a valid option", http.StatusNotFound)
		return
	}
}
	//counts all tracks in DB for admin
func AdminTracksCount(w http.ResponseWriter, r *http.Request) {

	//w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemennted", http.StatusNotImplemented)
		return
	}

	client := db.MongoConnect()

	fmt.Fprintf(w, "tracks in DB : %d", db.CountAllTracks(client))
}
	//deletes tracks in db
func AdminTracks(w http.ResponseWriter, r *http.Request) {

	//w.Header().Set("Content-Type", "application/json")

	if r.Method != "DELETE" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	client := db.MongoConnect()

	// Notifying the admin first for the current count of the track
	fmt.Fprintf(w, "tracks removed from DB : %d", db.CountAllTracks(client))

	// Deleting all the track in DB
	db.DeleteAllTracks(client)

}

//Old functions that we went away from, left behind to show working progress
/*func handlerIdAndField(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	idExists := false
	for i := 0; i < len(IDs); i++ {
		if IDs[i] == strings.ToUpper(parts[4]) {
			idExists = true
			break
		}
	}
	if !idExists {
		http.Error(w, "ID out of range.", http.StatusNotFound)
		return
	}
	t, _ := db.Get(strings.ToUpper(parts[4]))
	if len(parts) == 5 {
		http.Header.Set(w.Header(), "content.type", "application/json")
		json.NewEncoder(w).Encode(t)
	}
	if len(parts) == 6 {
		switch strings.ToUpper(parts[5]) {
		case "PILOT":
			fmt.Fprint(w, t.Pilot)
		case "GLIDER":
			fmt.Fprint(w, t.Glider)
		case "GLIDER_ID":
			fmt.Fprint(w, t.GliderId)
		case "TRACK_LENGTH":
			fmt.Fprint(w, t.TrackLength)
		case "H_DATE":
			fmt.Fprint(w, t.HeaderDate)
		default:
			http.Error(w, "Not a valid option", http.StatusNotFound)
			return
		}
	}
	if len(parts) > 6 {
		http.Error(w, "Too many /.", http.StatusNotFound)
	}
}*/