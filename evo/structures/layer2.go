package structures

type GetDocumentsRequest struct {
	StartAt    *int    `json:"start_at,omitempty"`
	StartAfter *int    `json:"start_after,omitempty"`
	Where      *[]byte `json:"where,omitempty"`
	OrderBy    *[]byte `json:"order_by,omitempty"`
	Limit      *int    `json:"limit,omitempty"`
}
