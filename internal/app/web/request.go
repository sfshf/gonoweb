package web

type Pagination struct {
	Page     int `json:"page" form:"page" binding:""`         // <=0 获取所有；>0 进行分页
	PageSize int `json:"pageSize" form:"pageSize" binding:""` // <=0 获取所有；>0 进行分页
}
