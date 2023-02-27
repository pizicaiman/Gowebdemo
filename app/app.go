package app

import (
	"awesomeProject/biz/service"
	"awesomeProject/common"
	"awesomeProject/controller"
	"awesomeProject/dao"
	"awesomeProject/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

package app

import (
	"errors"
	"fmt"
	// 引入gin web插件
	"github.com/gin-gonic/gin"
	//引入mysql插件
	_ "github.com/go-sql-driver/mysql"
	//引入 http插件
	"net/http"

	//项目 自定义包
	"awesomeProject/biz/dao"
	"awesomeProject/biz/service"
	"ibc/common"
	"awesomeProject/controller"
	"awesomeProject/util"
)

func Run() error {
	//初始化数据库连接
	dao.InitDB()
	defer dao.CloseDB()
	//创建service
	accountService := service.NewAccountService()
	activitiesService := service.NewOfflineActivitiesService()

	//定义路由
	router := gin.Default()
	router.GET("/", homeEndpoint)
	accountGroup := router.Group("/account")
	{
		accountCtl := controller.NewAccountController(accountService)
		accountGroup.POST("/", wrapHandler(accountCtl.InsertAccount))
		accountGroup.GET("/:userId", wrapHandler(accountCtl.FindAccount))
		accountGroup.GET("/account-flow/:accountId",                      wrapHandler(accountCtl.FindAccountFlow))
	}

	return router.Run(":8080")
}


type EndpointFunc func(*gin.Context) (any, error)

//定义返回
func homeEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS"})
}


// 定义统一处理器
func wrapHandler(handler EndpointFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理跨域请求
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		Cors(c)


		//从header 获取token,
		var token = c.GetHeader("token")
		//解析token,如果token不存在，或者 解析错误,则返回 没权限
		userId, err := util.TokenHandle(token)
		if err != nil {
			fmt.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"code": "ERROR", "message": "token is invalid"})
			return
		}
		//解析后 userId 放入上下文中,后面的service服务，可以获取该用户信息
		c.Set("userId", userId)
		data, err := handler(c)
		if err != nil {
			var systemErr *common.Error
			if errors.As(err, &systemErr) {
				c.JSON(systemErr.Status, common.NewErrorResult(systemErr.Code, systemErr.Msg))
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"code": "ERROR", "message": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "data": data})
	}
}

//跨域请求处理
func Cors(context *gin.Context) {
	method := context.Request.Method
	// 必须，接受指定域的请求，可以使用*不加以限制，但不安全
	//context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Origin", context.GetHeader("Origin"))
	fmt.Println(context.GetHeader("Origin"))
	// 必须，设置服务器支持的所有跨域请求的方法
	context.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	// 服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段
	context.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
	// 可选，设置XMLHttpRequest的响应对象能拿到的额外字段
	context.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
	// 可选，是否允许后续请求携带认证信息Cookir，该值只能是true，不需要则不设置
	context.Header("Access-Control-Allow-Credentials", "true")
	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		context.AbortWithStatus(http.StatusNoContent)
		return
	}
	context.Next()
}