package utils

import (
	"sort"

	"github.com/TinyKitten/TimelineServer/models"
)

type Posts []models.Post

func (p Posts) Len() int {
	return len(p)
}

func (p Posts) Less(i, j int) bool {
	return p[i].CreatedAt.After(p[j].CreatedAt)
}

func (p Posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func SortByPostDates(posts []models.Post) []models.Post {
	sort.Sort(Posts(posts))
	return posts
}
