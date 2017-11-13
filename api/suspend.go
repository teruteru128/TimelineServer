package api

func (h *handler) checkSuspended(oid string) bool {
	u, err := h.db.FindUserByOID(oid)
	if err != nil {
		return true
	}

	return u.Suspended
}
