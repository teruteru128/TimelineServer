package db

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"

	"gopkg.in/ory-am/dockertest.v3"
)

var ins *MongoInstance

type empty struct{}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "3.0", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		conf := config.DBConfig{
			Server:   fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp")),
			Database: "testing",
		}
		ins = &MongoInstance{Conf: conf}

		return ins.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	dummy := models.NewUser("id", "name", "password", "test@example.com")
	err := ins.Create("users", dummy)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestDuplicated(t *testing.T) {
	dummy := models.NewUser("dup", "name", "password", "dup@example.com")
	dummy2 := models.NewUser("dup", "name", "password", "dup@example.com")
	err := ins.Create("users", dummy)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = ins.Create("users", dummy2)
	if err == nil {
		t.Errorf("Duplicated")
	}
}
