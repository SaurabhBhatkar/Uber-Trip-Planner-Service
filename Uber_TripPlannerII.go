package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MyJsonRequest struct {
}

//JsonCoordinates     struct to hold coordinates of two endpoints
type JsonCoordinates struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	StartLat  float64 `json:"start_latitude" bson:"start_latitude"`
	StartLng  float64 `json:"start_longitude" bson:"start_longitude"`
	EndLat    float64 `json:"end_latitude" bson:"end_latitude"`
	EndLng    float64 `json:"end_longitude" bson:"end_longitude"`
}

type MyJsonResult struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}



type MyJsonNameReq2 struct {
	Name       string        `json:"name"`
	Address    string        `json:"address"`
	City       string        `json:"city"`
	State      string        `json:"state"`
	Zip        string        `json:"zip"`
	Id         bson.ObjectId `json:"id" bson:"_id"`
	Coordinate struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"coordinate"`
}

type MyJsonName struct {
	Prices []struct {
		//CurrencyCode         string  `json:"currency_code"`
		//DisplayName          string  `json:"display_name"`
		Distance float64 `json:"distance"`
		Duration int     `json:"duration"`
		//Estimate             string  `json:"estimate"`
		//HighEstimate         int     `json:"high_estimate"`
		//LocalizedDisplayName string  `json:"localized_display_name"`
		LowEstimate int `json:"low_estimate"`
		//Minimum              int     `json:"minimum"`
		//		ProductID            string  `json:"product_id"`
		ProductID string `json:"product_id" bson:"product_id"`
		//	SurgeMultiplier int    `json:"surge_multiplier"`
	} `json:"prices"`
}

type MyJsonNameInput struct {
	StartLocationID string   `json:"starting_from_location_id"`
	LocationIDs     []string `json:"location_ids"`
}


type OnRideComplete struct {
	BestRouteLocationIDs      []bson.ObjectId `json:"Best_route_location_ids"`
	ID                        bson.ObjectId   `json:"_id" bson:"_id"`
	TotalUberCost             int             `json:"total_uber_costs"`
	TotalUberDuration         int             `json:"total_uber_duration"`
	TotalDistance             float64         `json:"total_distance"`
	NextDestinationLocationID bson.ObjectId   `json:"next_destination_location_id" bson:"next_destination_location_id"`
	StartingFromLocationID    bson.ObjectId   `json:"starting_from_location_id" bson:"starting_from_location_id"`
	Status                    string          `json:"status" bson:"status"`
	UberWaitTimeEta           int             `json:"uber_wait_time_eta" bson:"uber_wait_time_eta"`
	CurrentIndex              int             `json:"CurrentIndex"`

}

type OnRideCompleted struct {
	BestRouteLocationIDs []bson.ObjectId `json:"Best_route_location_ids"`
	ID                   bson.ObjectId   `json:"_id" bson:"_id"`
	TotalUberCost        int             `json:"total_uber_costs"`
	TotalUberDuration    int             `json:"total_uber_duration"`
	TotalDistance        float64         `json:"total_distance"`
	StartingFromLocationID bson.ObjectId `json:"starting_from_location_id" bson:"starting_from_location_id"`
	Status                 string        `json:"status" bson:"status"`

		}
}

type EstimateTime struct {
	Eta int `json:"eta"`
}



//ReadRideRequest Starts Here
func ReadRideRequest(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	session, err := mgo.Dial("mongodb://tr!l)l)lanner:^^&&@ds045064.mongolab.com:45064/robots1")

	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	RideSchema := session.DB("robots1").C("Rides")
	if err != nil {
		log.Fatal(err)
	}

	id := p.ByName("name")
	Oid := bson.ObjectIdHex(id)
	var RideGetResult OnRideComplete
	RideSchema.FindId(Oid).One(&RideGetResult)
	if err != nil {
		log.Fatal(err)
	}
	b2, err := json.Marshal(RideGetResult)
	if err != nil {
	}
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, string(b2))

}

//GEt ends

//Create REquest starts
func CreateRideRequest(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	var myjson3 MyJsonNameInput
	d1 := json.NewDecoder(req.Body)
	err := d1.Decode(&myjson3)
	session, err := mgo.Dial("mongodb://tr!l)l)lanner:^^&&@ds045064.mongolab.com:45064/robots1")

	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	RideSchema := session.DB("robots1").C("Rides")

	var ArrayBestRouteID []bson.ObjectId
	ArrayBestRouteID = make([]bson.ObjectId, len(myjson3.LocationIDs))
	fmt.Println("Zone-------AAAAAAAAAAAAAAAAAAAAAA")

	fmt.Println("Lengthhhhhhhhhh", len(myjson3.LocationIDs))
	var StartLocation MyJsonNameReq2
	LocResultArr := make([]MyJsonNameReq2, len(myjson3.LocationIDs))
	UberResultArr := make([]MyJsonName, len(myjson3.LocationIDs))

	TotalUberCost := 0
	TotalUberDuration := 0
	TotalDistance := 0.0

	MinCost := 100
	//	CostIndex := 0
	fmt.Println("ZoneBBBBBBBBBBBBBB", len(myjson3.LocationIDs))

	Totalleft := len(myjson3.LocationIDs)
	//outerfor
	for OutLoop := 0; OutLoop < len(myjson3.LocationIDs); OutLoop++ {
		fmt.Println("OutLoop :-----------------------------------")
		fmt.Println(OutLoop)

		if OutLoop == 0 {
			id := myjson3.StartLocationID
			oid := bson.ObjectIdHex(id)
			fmt.Println("Testinggggg", oid)
			RideSchema.FindId(oid).One(&StartLocation)
		} else {
			oid := ArrayBestRouteID[OutLoop-1]
			RideSchema.FindId(oid).One(&StartLocation)
		}
		fmt.Println("Journey-=-=-=-=--=-=---=--=-=-=--Initiated")
		fmt.Println(OutLoop)
		fmt.Println(StartLocation)

		CostIndex := 0

		//InnerLoop Initiated

		for InLoop := 0; InLoop < Totalleft; InLoop++ {
			fmt.Println("PointA")
			id := myjson3.LocationIDs[InLoop]
			fmt.Println("PointB", id)
			oid := bson.ObjectIdHex(id)
			fmt.Println("PointC", oid)
			RideSchema.FindId(oid).One(&LocResultArr[InLoop])
			fmt.Println(LocResultArr[InLoop])

			var ServerToken = "aMjYkKq3URmOqlZ6t72kPfHNjZRBlqBPk4ANTKo-"
			Url := "https://sandbox-api.uber.com/v1/estimates/price?"

			AppendUrl := "start_latitude=" + strconv.FormatFloat(StartLocation.Coordinate.Lat, 'f', 6, 64) + "&start_longitude=" + strconv.FormatFloat(StartLocation.Coordinate.Lng, 'f', 6, 64) + "&end_latitude=" + strconv.FormatFloat(LocResultArr[InLoop].Coordinate.Lat, 'f', 6, 64) + "&end_longitude=" + strconv.FormatFloat(LocResultArr[InLoop].Coordinate.Lng, 'f', 6, 64) + "&server_token=" + ServerToken
			Url += AppendUrl
			fmt.Println("URL =====" + Url)
			res, err := http.Get(Url)
			if err != nil {
				log.Fatal(err)
			}
			droids, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)

			}

			err = json.Unmarshal(droids, &UberResultArr[InLoop])

			fmt.Println("UberResultArr	 :")
			fmt.Println(InLoop)
			fmt.Println(UberResultArr[InLoop])
			if err != nil {
				log.Fatal(err)
			}


			if MinCost > UberResultArr[InLoop].Prices[0].LowEstimate {
				MinCost = UberResultArr[InLoop].Prices[0].LowEstimate
				CostIndex = InLoop
				fmt.Println("MinCost :")
				fmt.Println(MinCost)
				fmt.Println("CostIndex :")
				fmt.Println(CostIndex)

			}
		} //Inloop ends
		TotalUberCost += UberResultArr[CostIndex].Prices[0].LowEstimate
		TotalUberDuration += UberResultArr[CostIndex].Prices[0].Duration
		TotalDistance += UberResultArr[CostIndex].Prices[0].Distance

		fmt.Println("LocationResArr_+_+_+_+_+_+_+_+_+_+_+_+_+-")
		fmt.Println(OutLoop)
		fmt.Println(LocResultArr)

		ArrayBestRouteID[OutLoop] = LocResultArr[CostIndex].Id
		LocResultArr[CostIndex] = LocResultArr[len(LocResultArr)-(OutLoop+1)]
		UberResultArr[CostIndex] = UberResultArr[len(LocResultArr)-(OutLoop+1)]
		myjson3.LocationIDs[CostIndex] = myjson3.LocationIDs[len(LocResultArr)-(OutLoop+1)]
		fmt.Println("Totalleft:")
		fmt.Println(Totalleft)
		fmt.Println("CostIndex:")
		fmt.Println(CostIndex)

		Totalleft = Totalleft - 1
		CostIndex = 0
		fmt.Println("LocationResultarray-------------------------------------------------")
		fmt.Println(OutLoop)
		fmt.Println(LocResultArr)
		fmt.Println("Totalleft:")
		fmt.Println(Totalleft)

		fmt.Println("Intermediate_uber_costs:")
		fmt.Println(TotalUberCost)
		fmt.Println("Intermediate_uber_duration:")
		fmt.Println(TotalUberDuration)
		fmt.Println("Intermediate_distance:")
		fmt.Println(TotalDistance)

		fmt.Println("End of outer loop")
		fmt.Println(OutLoop)

	} //OutLoop ends

	//fmt.Println("Start of round trip")

	var ReturnPoint MyJsonNameReq2
	RideSchema.FindId(ArrayBestRouteID[len(ArrayBestRouteID)-1]).One(&StartLocation)
	RideSchema.FindId(bson.ObjectIdHex(myjson3.StartLocationID)).One(&ReturnPoint)

	ServerToken := "aMjYkKq3URmOqlZ6t72kPfHNjZRBlqBPk4ANTKo-"
	Url := "https://sandbox-api.uber.com/v1/estimates/price?"
	AppendUrl := "start_latitude=" + strconv.FormatFloat(StartLocation.Coordinate.Lat, 'f', 6, 64) + "&start_longitude=" + strconv.FormatFloat(StartLocation.Coordinate.Lng, 'f', 6, 64) + "&end_latitude=" + strconv.FormatFloat(ReturnPoint.Coordinate.Lat, 'f', 6, 64) + "&end_longitude=" + strconv.FormatFloat(ReturnPoint.Coordinate.Lng, 'f', 6, 64) + "&server_token=" + ServerToken
	Url += AppendUrl
	fmt.Println("URL =====" + Url)
	res, err := http.Get(Url)
	if err != nil {
		log.Fatal(err)
	}
	droids, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)

	}

	var MyNewUberResultFirst MyJsonName
	err = json.Unmarshal(droids, &MyNewUberResultFirst)

	fmt.Println("Check beginning results", MyNewUberResultFirst)

	fmt.Println(MyNewUberResultFirst.Prices[0].Distance)
	fmt.Println(MyNewUberResultFirst.Prices[0].Duration)
	fmt.Println(MyNewUberResultFirst.Prices[0].LowEstimate)
	fmt.Println("Check end results")
	fmt.Println("End of round trip")

	fmt.Println("Total_uber_costs:")
	fmt.Println(TotalUberCost)
	fmt.Println("Total_uber_duration:")
	fmt.Println(TotalUberDuration)
	fmt.Println("Total_distance:")
	fmt.Println(TotalDistance)
	TotalUberCost = TotalUberCost + MyNewUberResultFirst.Prices[0].LowEstimate
	TotalUberDuration += MyNewUberResultFirst.Prices[0].Duration
	TotalDistance += MyNewUberResultFirst.Prices[0].Distance
	fmt.Println(" New Total_uber_costs:")
	fmt.Println(TotalUberCost)
	fmt.Println(" New Total_uber_duration:")
	fmt.Println(TotalUberDuration)
	fmt.Println(" New Total_distance:")
	fmt.Println(TotalDistance)

	var MyUberResult OnRideComplete

	MyUberResult.TotalUberCost = TotalUberCost

	fmt.Println(MyUberResult.TotalUberCost)
	MyUberResult.TotalUberDuration = TotalUberDuration
	MyUberResult.TotalUberDuration = TotalUberDuration
	MyUberResult.TotalDistance = TotalDistance

	MyUberResult.BestRouteLocationIDs = ArrayBestRouteID
	MyUberResult.Status = "Planning"
	MyUberResult.StartingFromLocationID = bson.ObjectIdHex(myjson3.StartLocationID)
	MyUberResult.ID = bson.NewObjectId()
	MyUberResult.NextDestinationLocationID = MyUberResult.BestRouteLocationIDs[0]

	var MyUberResults OnRideCompleted
	MyUberResults.BestRouteLocationIDs = make([]bson.ObjectId, len(MyUberResult.BestRouteLocationIDs))

	for a, _ := range MyUberResult.BestRouteLocationIDs {
		MyUberResults.BestRouteLocationIDs[a] = MyUberResult.BestRouteLocationIDs[a]
	}
	MyUberResults.ID = MyUberResult.ID
	MyUberResults.StartingFromLocationID = MyUberResult.StartingFromLocationID
	MyUberResults.Status = MyUberResult.Status
	MyUberResults.TotalUberCost = MyUberResult.TotalUberCost
	MyUberResults.TotalDistance = MyUberResult.TotalDistance
	MyUberResults.TotalUberDuration = MyUberResult.TotalUberDuration

	err = RideSchema.Insert(MyUberResults)

	if err != nil {
		fmt.Println("Error while inserting record")
		fmt.Println(err)
	} else {
		fmt.Println("Inserted Successfully")

	}

	Oid2 := MyUberResult.ID
	fmt.Println("Oid2")
	fmt.Println(Oid2)
	var RideGetResult OnRideComplete
	RideSchema.FindId(Oid2).One(&RideGetResult)

	b2, err := json.Marshal(MyUberResults)
	if err != nil {
		fmt.Println(err)

	}

	fmt.Println(ArrayBestRouteID)

	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, string(b2))

}

//ends

//UpdateRideRequest begins
func UpdateRideRequest(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	session, err := mgo.Dial("mongodb://tr!l)l)lanner:^^&&@ds045064.mongolab.com:45064/robots1")
	//  session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	RideSchema := session.DB("robots1").C("Rides")
	//	RideSession := session.DB("robots1").C("people")


	id := p.ByName("name")
	fmt.Println("id")
	fmt.Println(id)
	Oidz := bson.ObjectIdHex(id)
	fmt.Println("Oid")
	fmt.Println(Oidz)
	var RideGetResult OnRideComplete
	RideSchema.FindId(Oidz).One(&RideGetResult)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UberGetResult.Id:", RideGetResult.ID)
	fmt.Println("UberGetResult.Best_route_location_ids:", RideGetResult.BestRouteLocationIDs)
	fmt.Println("UberGetResult.Distance:", RideGetResult.TotalDistance)
	fmt.Println("UberGetResult.Duration", RideGetResult.TotalUberDuration)
	fmt.Println("UberGetResult.Low_Estimate_Price", RideGetResult.UberWaitTimeEta)

	fmt.Println("Id2:", RideGetResult.NextDestinationLocationID.String())
	fmt.Println("Checking whether true or false:", RideGetResult.NextDestinationLocationID == bson.ObjectId(""))

	if RideGetResult.NextDestinationLocationID == bson.ObjectId("") {
		RideGetResult.NextDestinationLocationID = RideGetResult.BestRouteLocationIDs[0]
		fmt.Println("Setting: ", RideGetResult.NextDestinationLocationID.String())

		fmt.Println("First Location")

	} else {
		for i := 0; i < len(RideGetResult.BestRouteLocationIDs)-1; i++ {
			if RideGetResult.BestRouteLocationIDs[0] == RideGetResult.NextDestinationLocationID {
				fmt.Println("Need to go to next location")
				RideGetResult.NextDestinationLocationID = RideGetResult.BestRouteLocationIDs[i]
			} //if
		} //for
		//		UberGetResult.Next_destination_location_id = UberGetResult.Best_route_location_ids[UberGetResult.CurrentIndex]
	} //else

	var Starting MyJsonNameReq2
	var NextPoint1 MyJsonNameReq2
	fmt.Println("Printing id")
	if RideGetResult.CurrentIndex == 0 {
		fmt.Println(RideGetResult.StartingFromLocationID)
		RideSchema.FindId(RideGetResult.StartingFromLocationID).One(&Starting)
		RideSchema.FindId(RideGetResult.BestRouteLocationIDs[0]).One(&NextPoint1)

		//		RideSession.FindId(RideGetResult.StartingFromLocationID).One(&Starting)
		//		RideSession.FindId(RideGetResult.BestRouteLocationIDs[0]).One(&NextPoint1)

		RideGetResult.Status = "Planning"

		RideGetResult.NextDestinationLocationID = RideGetResult.BestRouteLocationIDs[0]
	} else if RideGetResult.CurrentIndex == len(RideGetResult.BestRouteLocationIDs) {
		RideSchema.FindId(RideGetResult.BestRouteLocationIDs[RideGetResult.CurrentIndex-1]).One(&Starting)
		RideSchema.FindId(RideGetResult.StartingFromLocationID).One(&NextPoint1)
		RideGetResult.NextDestinationLocationID = RideGetResult.StartingFromLocationID

		RideGetResult.Status = "Finished"

	} else {
		RideSchema.FindId(RideGetResult.BestRouteLocationIDs[RideGetResult.CurrentIndex-1]).One(&Starting)
		RideSchema.FindId(RideGetResult.BestRouteLocationIDs[RideGetResult.CurrentIndex]).One(&NextPoint1)
		RideGetResult.NextDestinationLocationID = RideGetResult.BestRouteLocationIDs[RideGetResult.CurrentIndex]
		RideGetResult.Status = "Planning"
	}
	fmt.Println("Starting.Coordinate.Lat")
	fmt.Println("Starting.Coordinate.Lng")
	fmt.Println(Starting.Coordinate.Lat)
	fmt.Println(Starting.Coordinate.Lng)

	var ServerToken = "ZfFnTv1EKoy6SwKeHMuecmxy2IL8coZe-n5zC6No"
	Url := "https://sandbox-api.uber.com/v1/estimates/price?"
	AppendUrl := "start_latitude=" + strconv.FormatFloat(Starting.Coordinate.Lat, 'f', 6, 64) + "&start_longitude=" + strconv.FormatFloat(Starting.Coordinate.Lng, 'f', 6, 64) + "&end_latitude=" + strconv.FormatFloat(NextPoint1.Coordinate.Lat, 'f', 6, 64) + "&end_longitude=" + strconv.FormatFloat(NextPoint1.Coordinate.Lng, 'f', 6, 64) + "&server_token=" + ServerToken
	Url += AppendUrl
	fmt.Println("URL: " + Url)
	res, err := http.Get(Url)
	if err != nil {
		log.Fatal(err)
	}
	droids, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var rideresFirst MyJsonName
	err = json.Unmarshal(droids, &rideresFirst)

	fmt.Println("rideresFirst.Prices[0].Low_Estimate")
	fmt.Println(rideresFirst.Prices[0].LowEstimate)

	fmt.Println("rideresFirst.Prices[0].ProductID")
	fmt.Println(rideresFirst.Prices[0].ProductID)

	apiUrl := "https://sandbox-api.uber.com/v1/requests"

	JsonParam := JsonCoordinates{}
	JsonParam.ProductID = rideresFirst.Prices[0].ProductID
	JsonParam.StartLat = Starting.Coordinate.Lat
	JsonParam.StartLng = Starting.Coordinate.Lng
	JsonParam.EndLat = NextPoint1.Coordinate.Lat
	JsonParam.EndLng = NextPoint1.Coordinate.Lng

	JsonStr, err := json.Marshal(JsonParam)
	if err != nil {
		fmt.Println("UBER Error")
		//	log.Fatal(err)
	}

	req, err = http.NewRequest("POST", apiUrl, bytes.NewBuffer(JsonStr))
	if err != nil {
		fmt.Println("UBER Error")
		//	log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicmVxdWVzdCJdLCJzdWIiOiJlNjRiZmYwOC04ZDQ0LTQ5MzQtODRmNy0zZWRmNDYyZDM5YjgiLCJpc3MiOiJ1YmVyLXVzMSIsImp0aSI6ImQ1YzdjZjg3LTI3OTktNGQ5Ni04MDdiLTI3Yzk2ZjExYTlmOCIsImV4cCI6MTQ1MDgzMTIzNiwiaWF0IjoxNDQ4MjM5MjM1LCJ1YWN0IjoiMm5lcjZqeTFrdDR6UTRzS2ZpNWRpdEtPdmFWVUMzIiwibmJmIjoxNDQ4MjM5MTQ1LCJhdWQiOiJtcklieTFCeWdGSmpENkF1ZGpMNDhGVTFIdHkyUkZLUiJ9.ZYnJQqOaRDVrh9lB_yovyPO1GkvvcjSoUtlRLskXxqxvJdcE0ibY-c6xJrMU5ZpuKiYbTGhWR2y-HBsVtsfSp5QlAY4Zz4HsNiVuL78sbQdiqcIasQ69DP7iGazr9E4grEc6qZlPrzl39oK8qojeT5K_jaRWwFLW9H-cPVQdv0FG0GR69Cs_N7YpRZ-NR-iAASm-92Zg-LHiyljd49FM2AcswkBTNNQuHoayLxwW8kFRHQSRIkji-mOTTJAbJqeZyg4GgwMw74cwW349mk2DRKjuI2jVrF1SDkIROGtnP7Ifd17ppNZ8JTkXdsVC1_8LOsWbEP3MzN8cayc5YaJDXQ")

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("UBER ErrorStage1")
		//	log.Fatal(err)
	}

	bodys, errs := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if errs != nil {
		fmt.Println("UBER ErrorStage2")
		//	log.Fatal(err)
	}

	var TimeEst EstimateTime
	errs = json.Unmarshal(bodys, &TimeEst)
	if errs != nil {
		fmt.Println("UBER Decoding error Partially unmarshalled")
		//	log.F	atal(err)
	}

	RideGetResult.UberWaitTimeEta = TimeEst.Eta
	fmt.Println("RideGetResult.CurrentIndex1")
	fmt.Println(RideGetResult.CurrentIndex)

	RideGetResult.CurrentIndex = RideGetResult.CurrentIndex + 1
	if RideGetResult.CurrentIndex > len(RideGetResult.BestRouteLocationIDs) {
		RideGetResult.CurrentIndex = 0
	}
	fmt.Println("UberGetResult.CurrentIndex2")
	fmt.Println(RideGetResult.CurrentIndex)
	//
	// if UberGetResult.CurrentIndex <= len(UberGetResult.Best_route_location_ids) {
	// 	UberGetResult.Status = "Planning"
	// } else {
	// 	UberGetResult.Status = "Finished"
	// }
	fmt.Println("UberGetResult.Uber_wait_time_eta")
	fmt.Println(RideGetResult.UberWaitTimeEta)

	fmt.Println("UberGetResult.Best_route_location_ids")
	fmt.Println(RideGetResult.BestRouteLocationIDs)

	fmt.Println("UberGetResult.Total_distance")
	fmt.Println(RideGetResult.TotalDistance)

	fmt.Println("UberGetResult.Total_uber_duration")
	fmt.Println(RideGetResult.TotalUberDuration)

	fmt.Println("UberGetResult.Total_uber_costs")
	fmt.Println(RideGetResult.TotalUberCost)

	fmt.Println("UberGetResult.Uber_wait_time_eta")
	fmt.Println(RideGetResult.Status)
	fmt.Println("UberGetResult.Next_destination_location_id")
	fmt.Println(RideGetResult.NextDestinationLocationID)
	fmt.Println("UberGetResult.Starting_from_location_id")
	fmt.Println(RideGetResult.StartingFromLocationID)

	RideSchema.UpdateId(Oidz, RideGetResult)

	b2, err := json.Marshal(RideGetResult)
	if err != nil {
	}

	rw.WriteHeader(http.StatusOK)

	fmt.Fprintf(rw, string(b2))
	//TripStatus.UberWaitTimeEta = TimeEst.Eta

	//

}


func main() {
	mux := httprouter.New()
	mux.GET("/trips/:name", ReadRideRequest)
	mux.POST("/trips", CreateRideRequest)
	mux.PUT("/trips/:name/request", UpdateRideRequest)
	//	mux.DELETE("/rideplans/:name", DeleteRideRequest)

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()

}
