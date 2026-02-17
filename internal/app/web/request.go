package web

type Pagination struct {
	Page     int `json:"page" form:"page" binding:"gt=0"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"gte=2,lte=20"`
}
