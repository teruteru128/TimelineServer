package utils

import (
	"testing"
	"time"

	"github.com/TinyKitten/TimelineServer/models"
	"gopkg.in/mgo.v2/bson"
)

func TestSortByPostDates(t *testing.T) {
	olderPost := models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "old")         // 2年前
	newerPost := models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "new")         // ついさっき
	middlePost := models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "1 year ago") // 1年前
	now := time.Now()
	old := now.AddDate(-2, 0, 0)
	olderPost.CreatedAt = old
	newerPost.CreatedAt = now
	whileAgo := now.AddDate(-1, 0, 0)
	middlePost.CreatedAt = whileAgo

	posts := []models.Post{
		*middlePost,
		*olderPost,
		*newerPost,
	}
	expected := []models.Post{
		*newerPost,
		*middlePost,
		*olderPost,
	}

	sorted := SortByPostDates(posts)
	for i := range sorted {
		if sorted[i].CreatedAt != expected[i].CreatedAt {
			t.Fatalf("not sorted\n actual:\n %v,\n expected:\n %v\n", sorted, expected)
		}
	}
}
