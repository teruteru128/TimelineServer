package db

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/TinyKitten/Timeline/config"
	"github.com/TinyKitten/Timeline/models"

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

func TestFindUser(t *testing.T) {
	dummy := models.NewUser("id", "name", "password", "test@example.com")
	err := ins.Create("users", dummy)
	if err != nil {
		t.Errorf(err.Error())
	}
	exist, err := ins.FindUser("id")
	if err != nil {
		t.Errorf(err.Error())
	}
	if exist == nil {
		t.Errorf("User not found")
	}

	_, err = ins.FindUser("ugly_betty")
	if err == nil {
		t.Errorf("already registered")
	}
}

func TestEmptyInstance(t *testing.T) {
	badIns := MongoInstance{}
	_, err := badIns.FindUser("hoge")
	if err == nil {
		t.Errorf("bad instance")
	}
	err = badIns.Create("bad", empty{})
	if err == nil {
		t.Errorf("bad instance")
	}
}
