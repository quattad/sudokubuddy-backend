package caching

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// MEMCACHE SETUP
// Cache is an instance of the
// var Cache = memcache.New("127.0.0.1:11211")

// Cache is an in-memory cache with default expiration time and purge time
// GO-CACHE
var Cache = cache.New(5*time.Minute, 10*time.Minute)

// GROUPCACHE
// Experiment with groupcache
// - In-code Distributed Cache
// - Every instance of server is a node in distributed cache

// 1. App asks the GroupCache for the data via a key.
// 2. GroupCache checks the in memory hot cache for the data, if no data, continue.
// 3. GroupCache performs a consistent hash on the key to determine which GroupCache instance has the data.
// 4. GroupCache makes a network request to the GroupCache instance that has the data
// 5. GroupCache returns the data if it exists in memory, if not it asks the App to render or fetch data.
// 6. GroupCache returns data to GroupCache instance that initiated the request.

// Exports groups

// Store simulates a cache
// var Store = map[string]string{
// 	"testkey1": "testvalue1",
// 	"testkey2": "testvalue2",
// }

// CacheUserGroup is a cache for users
// var CacheUserGroup = groupcache.NewGroup("users", 64<<20, groupcache.GetterFunc(
// 	func(ctx context.Context, key string, dest groupcache.Sink) error {
// 		log.Println("Looking up: ", key)

// 		v, ok := Store[key] // replace with key

// 		if !ok {
// 			fmt.Println("Cache miss")
// 			return errors.New("User not found")
// 		}

// 		fmt.Println("Cache hit")
// 		dest.SetString(v)
// 		return nil

// 	},
// ))

// InitializeCache initializes cache
// func InitializeCache() {
// 	// Initialize peer in cluster and add current instance into the pool
// 	peers := flag.String("pool", config.BASEURL, "server pool list")
// 	flag.Parse()

// 	p := strings.Split(*peers, ",")
// 	pool := groupcache.NewHTTPPool(p[0])
// 	pool.Set(p...)
// 	fmt.Println("Running cache...")
// }
