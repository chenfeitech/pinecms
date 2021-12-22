package router

import (
	"github.com/xiusin/pine"
	"github.com/xiusin/pinecms/src/application/controllers/backend/wechat"
)

// InitModuleRouter 函数内不可出现"}", 如需实现自定义解析, 务必定义新函数
// {holder}占位符作为替换位置, 不可删除
func InitModuleRouter(backendRouter *pine.Router, app *pine.Application) {
	// holder
}


func InitSubModuleRouter(app *pine.Application, admin *pine.Router)  {
	wechat.InitRouter(app, admin)
}
