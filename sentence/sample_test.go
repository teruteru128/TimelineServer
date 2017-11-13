package sentence

import "testing"

func TestGetSampleText(t *testing.T) {
	n := 50
	mc := 10 // 衝突最大許容回数
	c := 0   // 累計衝突回数
	prev := ""
	for n > 0 {
		n--
		text := getSampleText()
		if text == "" {
			t.Errorf("%d: Text is empty", n)
		}
		if prev == text {
			c++
			if c == mc {
				// 衝突許容回数超過
				t.Errorf("Text collision count exceeded: %d %s", c, text)
			}
		}
		prev = text
	}
}
func TestGetRandomPost(t *testing.T) {
	sender, post := GetRandomPost()
	if sender == nil && post == nil {
		t.Errorf("NIL")
	}
	if sender.UserID == "" {
		t.Errorf("Invalid userID: %s", sender.UserID)
	}

	if sender.DisplayName == "New User" {
		t.Errorf("sender display name not changed")
	}

	if sender.EMail == "" {
		t.Errorf("Empty email")
	}

	if sender.Password == "" {
		t.Errorf("Empty password")
	}

	if post.UserID == "" {
		t.Errorf("Empty sender")
	}

	if post.Text == "" {
		t.Errorf("Empty post body")
	}
}
