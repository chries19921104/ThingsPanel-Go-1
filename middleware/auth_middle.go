package middleware

import (
	response "ThingsPanel-Go/utils"
	"fmt"
	"strings"

	"ThingsPanel-Go/initialize/redis"
	jwt "ThingsPanel-Go/utils"

	adapter "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/adapter/context"
	context2 "github.com/beego/beego/v2/server/web/context"
)

// AuthMiddle 中间件
func AuthMiddle() {
	// 不 需要验证的url
	noLogin := map[string]interface{}{
		"api/plugin/device/sub-device-detail": 0,
		"api/plugin/register":                 0,
		"api/plugin/device/config":            0,
		"api/system/logo/index":               0,
		"api/open/data":                       0,
		"api/auth/login":                      0,
		"api/auth/refresh":                    0,
		"api/auth/register":                   1,
		"/ws":                                 2,
	}
	var filterLogin = func(ctx *context.Context) {
		url := strings.TrimLeft(ctx.Input.URL(), "/")
		if !isAuthExceptUrl(strings.ToLower(url), noLogin) {
			//获取TOKEN
			fmt.Println(ctx.Request.Header["Authorization"])
			if len(ctx.Request.Header["Authorization"]) == 0 {
				response.SuccessWithMessage(401, "Unauthorized", (*context2.Context)(ctx))
				return
			}
			authorization := ctx.Request.Header["Authorization"][0]
			userToken := authorization[7:]
			_, err := jwt.ParseCliamsToken(userToken)
			if err != nil || redis.GetStr(userToken) != "1" {
				// 异常
				response.SuccessWithMessage(401, "Unauthorized", (*context2.Context)(ctx))
				return
			}
			// 验证用户是否存在
			// 双层校验
			// if psql.Mydb.Where("id = ? and name = ? ", user.ID, user.Name).First(&models.Users{}).Error == gorm.ErrRecordNotFound {
			// 	response.SuccessWithMessage(401, "Unauthorized", (*context2.Context)(ctx))
			// 	return
			// }
			// s, _ := cache.Bm.IsExist(c.TODO(), userToken)
			// if !s {
			// 	response.SuccessWithMessage(401, "Unauthorized", (*context2.Context)(ctx))
			// 	return
			// }
		}
	}
	adapter.InsertFilter("/api/*", adapter.BeforeRouter, filterLogin)
}

// 不需要授权的url返回true
func isAuthExceptUrl(url string, m map[string]interface{}) bool {
	urlArr := strings.Split(url, "/")
	// url大于4个长度只判断前四个是否在不需授权map中
	if len(urlArr) > 4 {
		url = fmt.Sprintf("%s/%s/%s/%s", urlArr[0], urlArr[1], urlArr[2], urlArr[3])
	}
	_, ok := m[url]
	return ok
}
