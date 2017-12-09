package db

import (
	"time"

	"github.com/TinyKitten/TimelineServer/models"
	"gopkg.in/mgo.v2/bson"
)

const (
	// EventCol DB上のEvent用カラム
	EventCol = "event"
)

// GetEvents イベントをDBから取得する
func (m *MongoInstance) GetEvents(userID bson.ObjectId) (*[]models.Event, error) {
	sess := m.session.Clone()
	defer sess.Close()
	var events []models.Event
	if err := sess.DB(m.db()).C(EventCol).
		Find(bson.M{"to_user_id": userID}).
		All(&events); err != nil {
		return nil, err
	}
	return &events, nil
}

// InsertEvent イベントをDBに挿入する
func (m *MongoInstance) InsertEvent(fromID, toID bson.ObjectId, eventType models.EventType) (*models.Event, error) {
	sess := m.session.Clone()
	defer sess.Close()

	event := models.Event{
		FromUserID:  fromID,
		ToUserID:    toID,
		Type:        eventType,
		AlreadyRead: false,
		CreatedAt:   time.Now(),
	}

	err := m.Insert(EventCol, event)
	if err != nil {
		return nil, err
	}
	return &event, err
}

// InsertEvent イベントをDBに挿入する
func (m *MongoInstance) InsertPostEvent(fromID, toID, postID bson.ObjectId, eventType models.EventType) (*models.Event, error) {
	sess := m.session.Clone()
	defer sess.Close()

	event := models.Event{
		FromUserID:   fromID,
		ToUserID:     toID,
		Type:         eventType,
		AlreadyRead:  false,
		CreatedAt:    time.Now(),
		TargetPostID: postID,
	}

	err := m.Insert(EventCol, event)
	if err != nil {
		return nil, err
	}
	return &event, err
}

// DeleteEvent イベントをDBから削除
func (m *MongoInstance) DeleteEvent(id bson.ObjectId) error {
	sess := m.session.Clone()
	defer sess.Close()

	return sess.DB(m.db()).C(PostsCol).Remove(bson.M{"_id": id})
}
