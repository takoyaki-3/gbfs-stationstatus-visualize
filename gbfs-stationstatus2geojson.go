package main

import (
	"log"
	"strconv"

	geojson "github.com/takoyaki-3/go-geojson"
	json "github.com/takoyaki-3/go-json"
	"github.com/takoyaki-3/goc"
)

type StationStatus struct {
	TTL  int `json:"ttl"`
	Data struct {
		Stations []struct {
			IsRenting         bool   `json:"is_renting"`
			StationID         string `json:"station_id"`
			IsInstalled       bool   `json:"is_installed"`
			IsReturning       bool   `json:"is_returning"`
			LastReported      int    `json:"last_reported"`
			NumBikesAvailable int    `json:"num_bikes_available"`
			NumDocksAvailable int    `json:"num_docks_available"`
		} `json:"stations"`
	} `json:"data"`
	Version     string `json:"version"`
	LastUpdated int    `json:"last_updated"`
}

type StationInformaton struct {
	TTL  int `json:"ttl"`
	Data struct {
		Stations []StationStr `json:"stations"`
	} `json:"data"`
	Version     string `json:"version"`
	LastUpdated int    `json:"last_updated"`
}
type StationStr struct {
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Name      string  `json:"name"`
	Capacity  int     `json:"capacity"`
	RegionID  string  `json:"region_id"`
	StationID string  `json:"station_id"`
}

func main(){
	paths,_ := goc.Dirwalk("./gbfs-stationstatus")

	StationInformaton := StationInformaton{}
	if err:=json.LoadFromPath("./station_information.json",&StationInformaton);err!=nil{
		log.Fatalln(err)
	}
	stations := map[string]StationStr{}
	for _,station := range StationInformaton.Data.Stations{
		stations[station.StationID] = station
	}

	stationStatuses := []StationStatus{}
	for _,path := range paths {
		stationStatus := StationStatus{}
		if err:=json.LoadFromPath(path, &stationStatus);err!=nil{
			log.Fatalln(err)
		}
		stationStatuses = append(stationStatuses,stationStatus)
	}

	FeatureCollection := geojson.FeatureCollection{}
	FeatureCollection.Type = "FeatureCollection"

	for _,stationStatus := range stationStatuses{
		for _,station := range stationStatus.Data.Stations{
			s := stations[station.StationID]
			f := geojson.Feature{}
			f.Type = "Feature"
			f.Geometry = geojson.Geometry{
				Type: "Point",
				Coordinates: []float64{s.Lon,s.Lat},
			}
			prop := map[string]string{}
			prop["stop_id"] = station.StationID
			prop["num_bikes_available"] = strconv.Itoa(station.NumBikesAvailable)
			prop["num_docks_available"] = strconv.Itoa(station.NumDocksAvailable)
			if station.NumDocksAvailable+station.NumBikesAvailable == 0{
				continue
			}
			prop["rate"] = strconv.Itoa(station.NumBikesAvailable*10/(station.NumDocksAvailable+station.NumBikesAvailable))
			prop["last_reported"] = strconv.Itoa(station.LastReported+3600*9)
			f.Properties = prop
			FeatureCollection.Features = append(FeatureCollection.Features, f)
		}
	}

	if err:=json.DumpToFile(FeatureCollection,"./gbfs-stationstatus.geojson");err!=nil{
		log.Fatalln(err)
	}
}
