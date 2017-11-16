package api

import (
	"gopkg.in/mgo.v2/bson"
)

func (h *handler) checkSuspended(oid bson.ObjectId) (bool, error) {
	u, err := h.db.FindUserByOID(oid)
	if err != nil {
		return false, err
	}

	return u.Suspended, nil
}
