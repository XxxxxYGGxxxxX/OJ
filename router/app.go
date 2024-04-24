package router

import (
	_ "OJ/docs"
	"OJ/middlewares"
	"OJ/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 路由服务 注册路由
func Router() *gin.Engine {
	r := gin.Default()

	// swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 路由规则
	// 测试
	// 公用方法
	r.GET("/ping", service.Ping)
	// 问题
	r.GET("/problem-list", service.GetProblemList)
	r.GET("/problem-detail", service.GetProblemDetail)
	// 用户
	r.GET("/user-detail", service.GetUserDetail)
	r.POST("/login", service.Login)
	r.POST("/send-code", service.SendCode)
	r.POST("/register", service.Register)

	// 排行榜
	r.GET("/rank-list", service.GetRankList)
	// 提交记录
	r.GET("/submit-list", service.GetSubmitList)
	// 分类列表
	r.GET("/category-list", service.GetCategoryList)

	// 管理员私有方法
	authAdmin := r.Group("/admin", middlewares.AuthAdminCheck())
	// 问题创建
	authAdmin.POST("problem-create", service.ProblemCreate)
	// 问题修改
	authAdmin.PUT("problem-modify", service.ProblemModify)
	// 分类创建
	authAdmin.GET("/category-create", service.CategoryCreat)
	// 分类修改
	authAdmin.PUT("/category-modify", service.CategoryModify)
	// 分类删除
	authAdmin.DELETE("/category-delete", service.CategoryDelete)

	// 用户私有方法
	authUser := r.Group("/user", middlewares.AuthUserCheck())
	// 代码提交
	authUser.POST("/submit", service.Submit)

	return r
}
