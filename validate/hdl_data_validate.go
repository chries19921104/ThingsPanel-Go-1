package valid

type SoupDataPaginationValidate struct {
	CurrentPage int    `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int    `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	ShopName    string `json:"shop_name,omitempty" alias:"店铺名称" valid:"MaxSize(99)"`
	Id          string `json:"id,omitempty" alias:"Id" valid:"MaxSize(36)"`
	Limit       int    `json:"limit,omitempty" alias:"导出限制"`
	TableNumber string `json:"table_number,omitempty" alias:"桌号" valid:"MaxSize(99)"`
	StartTime   string `json:"start_time,omitempty" alias:"订单检索开始时间" valid:"MaxSize(99)"`
	EndTime     string `json:"end_time,omitempty" alias:"订单检索结束时间" valid:"MaxSize(99)"`
}

type RspSoupDataPaginationValidate struct {
	CurrentPage int                      `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int                      `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	Data        []map[string]interface{} `json:"data" alias:"返回数据"`
	Total       int64                    `json:"total" alias:"总数" valid:"Max(10000)"`
}
