package cache

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)
 
type Movie struct {
   Id          string `json:"id"`
   Title       string `json:"title"`
   Description string `json:"description"`
}

/*
{
  "id": "1",
  "title": "Avenger",
  "description": "Best movie of year"
}
*/

func (cache redisCache) CreateMovie(ctx *gin.Context, movie *Movie) (*Movie, error) {
	c := cache.getClient()
	
	movie.Id = uuid.New().String()

	json, err := json.Marshal(movie)

	if err != nil {
		return nil, err
	}

	c.HSet(ctx, "movies", json)

	if err != nil {
		return nil, err
	}

	return movie, nil
 }
  
 func (cache redisCache) GetMovie(ctx *gin.Context, id string) (*Movie, error) {
	c := cache.getClient()
	val, err := c.HGet(ctx, "movies", id).Result()
  
	if err != nil {
		return nil, err
	}
	movie := &Movie{}
	err = json.Unmarshal([]byte(val), movie)
  
	if err != nil {
		return nil, err
	}
	return movie, nil
 }
  
 func (cache redisCache) GetMovies(ctx *gin.Context) ([]*Movie, error) {
	c := cache.getClient()

	movies := []*Movie{}

	val, err := c.HGetAll(ctx, "movies").Result()

	if err != nil {
		return nil, err
	}

	for _, item := range val {
		movie := &Movie{}
		err := json.Unmarshal([]byte(item), movie)

		if err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}
  
	return movies, nil
 }
  
 func (cache redisCache) UpdateMovie(ctx *gin.Context, movie *Movie) (*Movie, error) {
	c := cache.getClient()

	json, err := json.Marshal(&movie)

	if err != nil {
		return nil, err
	}

	c.HSet(ctx, "movies", movie.Id, json)


	if err != nil {
		return nil, err
	}

	return movie, nil
 }
 func (cache redisCache) DeleteMovie(ctx *gin.Context, id string) error {
	c := cache.getClient()

	numDeleted, err := c.HDel(ctx, "movies", id).Result()

	if numDeleted == 0 {
		return errors.New("movie to delete not found")
	}
	if err != nil {
		return err
	}

	return nil
 }
 
type MovieService interface {
   GetMovie(ctx *gin.Context, id string) (*Movie, error)
   GetMovies(ctx *gin.Context) ([]*Movie, error)
   CreateMovie(ctx *gin.Context, movie *Movie) (*Movie, error)
   UpdateMovie(ctx *gin.Context, movie *Movie) (*Movie, error)
   DeleteMovie(ctx *gin.Context, id string) error
}

type redisCache struct {
	host string
	db   int
	exp  time.Duration
	ctx *gin.Context
 }
  
 func NewRedisCache(ctx *gin.Context, host string, db int, exp time.Duration) MovieService {
	var newCache redisCache;

	newCache.host = host;
	newCache.db = db;
	newCache.exp = exp;
	newCache.ctx = ctx;

	return &newCache
 }
  
 func (cache redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
 }