package main

// Pakage: go get github.com/redis/go-redis/v8

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// API: https://nominatim.openstreetmap.org/search?city=Ho Chi Minh&format=json
// TEST: http://localhost:5000/api?q=Ho Chi Minh

type API struct {
	cache *redis.Client
}

func NewAPI() *API {
	redisAddress := fmt.Sprintf("%s:6379", os.Getenv("REDIS_URL")); 

	rdb := redis.NewClient(&redis.Options{
        Addr:     redisAddress,
        Password: "", // no password set
        DB:       0,  // use default DB
    })

	return &API{
		cache: rdb,
	}
}

type NominatimResponse struct {
	PlaceID     int      `json:"place_id"`
	Licence     string   `json:"licence"`
	OsmType     string   `json:"osm_type"`
	OsmID       int64    `json:"osm_id"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	PlaceRank   int      `json:"place_rank"`
	Importance  float64  `json:"importance"`
	Addresstype string   `json:"addresstype"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Boundingbox []string `json:"boundingbox"`
}

type APIResponse struct {
	Cache bool `json:"cache"`
	Data []NominatimResponse `json:"data"`
}

func (a* API) getData(ctx context.Context, query string) ([]NominatimResponse, bool, error) {
	// Check data cached by Redis
	value, err := a.cache.Get(ctx, query).Result();

    if err == redis.Nil {
        fmt.Println("Key does not exist, call source api");

		escapedQ := url.PathEscape(query);

		city := fmt.Sprintf("https://nominatim.openstreetmap.org/search?city=%s&format=json", escapedQ)

		res, err := http.Get(city);

		if(err != nil) {
			return nil, false, err;
		}

		data := make([]NominatimResponse, 0);

		err = json.NewDecoder(res.Body).Decode(&data);

		if(err != nil) {
			return nil, false, err;
		}

		b, err := json.Marshal(data);

		if(err != nil) {
			return nil, false, err;
		}

		// Set Redis value
		err = a.cache.Set(ctx, query, bytes.NewBuffer(b).Bytes(), time.Second*15).Err();

		if(err != nil) {
			return nil, false, err;
		}

		// Return response
		return data, false, nil;

    } else if err != nil {
		fmt.Printf("error calling redis %v\n", err);
		return nil, false, err;
    } else {
		// Cache hit
		data := make([]NominatimResponse, 0);

		err := json.Unmarshal(bytes.NewBufferString(value).Bytes(), &data);

		if(err != nil) {
			return nil, false, err;
		}

		return data, true, nil;
    }
	
}

func (a* API) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In the handler...");

	query := r.URL.Query().Get("q");

	// GET METHOD
	data, cacheHit, err := a.getData(r.Context(), query);

	if(err != nil) {
		fmt.Printf("error calling data source %v\n", err);
		w.WriteHeader(http.StatusInternalServerError);
		return;
	}

	res := APIResponse{
		Cache: cacheHit,
		Data: data,
	};

	err = json.NewEncoder(w).Encode(res);

	if(err != nil) {
		fmt.Printf("error encoding response %v\n", err);
		w.WriteHeader(http.StatusInternalServerError);
		return;
	}

}

func main() {
	api := NewAPI();

	fmt.Println("Starting server...");

	http.HandleFunc("/api", api.Handler);

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil);
}