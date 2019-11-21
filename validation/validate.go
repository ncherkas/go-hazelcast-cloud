package validation

import (
	"fmt"
	"math"
	"time"

	"github.com/hazelcast/hazelcast-cloud-go-demo/common"
	"github.com/hazelcast/hazelcast-go-client/core"
)

// Request model for validation input
type Request struct {
	UserID               int       `json:"userId"`
	AirportCode          string    `json:"airportCode"`
	TransactionTimestamp time.Time `json:"transactionTimestamp"`
}

// Response model for validation output
type Response struct {
	Valid   bool
	Message string
}

const radiusOfEarthM = 6371000

// Apply a fraud detection validation
func Apply(req *Request, usersMap core.Map, airportsMap core.Map) (*Response, error) {
	fmt.Println("Handling Validate request", req)

	usrID := req.UserID

	usr, err := getUser(usersMap, usrID)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		newUsr := &common.User{
			UserID:               req.UserID,
			LastCardUsePlace:     req.AirportCode,
			LastCardUseTimestamp: req.TransactionTimestamp,
		}
		if err := usersMap.Set(usrID, newUsr); err != nil {
			return nil, fmt.Errorf("Failed to set user: %w", err)
		}
		return &Response{Valid: true, Message: "User data saved for future validations"}, nil
	}

	lastAirport, err := getAirport(airportsMap, usr.LastCardUsePlace)
	if err != nil {
		return nil, err
	}

	nextAirport, err := getAirport(airportsMap, req.AirportCode)
	if err != nil {
		return nil, err
	}

	mins := req.TransactionTimestamp.Sub(usr.LastCardUseTimestamp).Minutes()
	meters := haversine(nextAirport.Latitude, nextAirport.Longitude, lastAirport.Latitude, lastAirport.Longitude)
	speed := meters / mins
	valid := !(speed > 13000)

	var msg string
	if valid {
		msg = "Transaction is OK"
	} else {
		msg = "Transaction is suspicious"
	}

	usr.LastCardUsePlace = req.AirportCode
	usr.LastCardUseTimestamp = req.TransactionTimestamp

	if err := usersMap.Set(usrID, usr); err != nil {
		return nil, fmt.Errorf("Failed to set user: %w", err)
	}

	return &Response{valid, msg}, nil
}

func getUser(usersMap core.Map, id int) (*common.User, error) {
	val, err := usersMap.Get(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user: %w", err)
	}
	if val != nil {
		return val.(*common.User), nil
	}
	return nil, nil
}

func getAirport(airportsMap core.Map, code string) (*common.Airport, error) {
	val, err := airportsMap.Get(code)
	if err != nil {
		return nil, fmt.Errorf("Failed to get airport: %w", err)
	}
	if val == nil {
		return nil, fmt.Errorf("Airport by code '%v' not found", code)
	}
	return val.(*common.Airport), nil
}

func haversine(lat1 float64, longit1 float64, lat2 float64, longit2 float64) float64 {
	rlat1 := toRadians(lat1)
	rlongit1 := toRadians(longit1)
	rlat2 := toRadians(lat2)
	rlongit2 := toRadians(longit2)

	rlatDiff := rlat1 - rlat2
	rlongitDiff := rlongit1 - rlongit2

	hav := math.Pow(math.Sin(rlatDiff/2), 2) + math.Pow(math.Sin(rlongitDiff/2), 2)*math.Cos(rlat1)*math.Cos(rlat2)

	return 2 * radiusOfEarthM * math.Asin(math.Sqrt(hav))
}

func toRadians(angdeg float64) float64 {
	return float64(angdeg) / float64(180) * math.Pi
}
