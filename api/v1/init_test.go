package v1

import (
	"fmt"
	"os"
	"testing"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/db"
	"github.com/TinyKitten/TimelineServer/logger"
	"github.com/labstack/gommon/log"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

var th *APIHandler

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "latest", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	redisResource, err := pool.Run("redis", "latest", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		conf := config.DBConfig{
			Server:   fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp")),
			Database: "testing",
		}
		cacheConf := config.CacheConfig{
			Server: "redis://localhost:" + redisResource.GetPort("6379/tcp"),
		}
		ins, err := db.NewMongoInstance(conf, cacheConf)
		if err != nil {
			log.Fatalf("Could not getting mongo instance: %s", err)
		}
		logger := logger.NewLogger()

		th = &APIHandler{
			db:     ins,
			logger: logger,
		}

		return ins.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(redisResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}
