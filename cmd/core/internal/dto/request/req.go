package request

type IdReq struct {
	ID string `json:"id" binding:"required" form:"id" uri:"id"` // ID
}
