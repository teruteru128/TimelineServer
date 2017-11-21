package db

import "testing"
import "github.com/TinyKitten/TimelineServer/models"

func TestGetPosts(t *testing.T) {
	p0 := models.NewPost("siketyan", "hoge")
	p1 := models.NewPost("homo", "fuga")

	err := ins.Insert("posts", p0)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = ins.Insert("posts", p1)
	if err != nil {
		t.Errorf(err.Error())
	}

	res, err := ins.GetPosts(2)
	if err != nil {
		t.Errorf(err.Error())
	}
	if res == nil {
		t.Errorf("EMPTY")
	}

	if len(*res) != 2 {
		t.Errorf("Less than 2: %d", len(*res))
	}

	if (*res)[0].UserID != p0.UserID && (*res)[0].Text != p0.Text {
		t.Errorf("Not matched: %s", (*res)[0].Text)
	}

	if (*res)[1].UserID != p1.UserID && (*res)[1].Text != p1.Text {
		t.Errorf("Not matched: %s", (*res)[1].Text)
	}
}

func TestLimitedGetPosts(t *testing.T) {
	p0 := models.NewPost("siketyan", "hoge")
	p1 := models.NewPost("homo", "fuga")

	err := ins.Insert("posts", p0)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = ins.Insert("posts", p1)
	if err != nil {
		t.Errorf(err.Error())
	}

	res, err := ins.GetPosts(1)
	if err != nil {
		t.Errorf(err.Error())
	}
	if res == nil {
		t.Errorf("EMPTY")
	}

	if len(*res) != 1 {
		t.Errorf("Larger than 1: %d", len(*res))
	}

	if (*res)[0].UserID != p0.UserID && (*res)[0].Text != p0.Text {
		t.Errorf("Not matched: %s", (*res)[0].Text)
	}
}
