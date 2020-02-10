package backend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/xiusin/iriscms/src/application/controllers"

	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/xiusin/iriscms/src/application/models"
	"github.com/xiusin/iriscms/src/common/helper"
)

type ContentController struct {
	Ctx     iris.Context
	Orm     *xorm.Engine
	Session *sessions.Session
}

func (c *ContentController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("ANY", "/content/index", "Index")
	b.Handle("ANY", "/content/right", "Right")
	b.Handle("ANY", "/content/public-welcome", "Welcome")
	b.Handle("ANY", "/content/news-list", "NewsList")
	b.Handle("ANY", "/content/page", "Page")
	b.Handle("ANY", "/content/add", "AddContent")
	b.Handle("ANY", "/content/edit", "EditContent")
	b.Handle("ANY", "/content/delete", "DeleteContent")
	b.Handle("ANY", "/content/order", "OrderContent")
}

func (c *ContentController) Index() {
	menuid, _ := c.Ctx.URLParamInt64("menuid")
	c.Ctx.ViewData("currentPos", models.NewMenuModel(c.Orm).CurrentPos(menuid))
	c.Ctx.View("backend/content_index.html")
}

func (c *ContentController) Welcome() {
	c.Ctx.View("backend/content_welcome.html")
}

func (c *ContentController) Right() {
	if c.Ctx.Method() == "POST" {
		cats := models.NewCategoryModel(c.Orm).GetContentRightCategoryTree(models.NewCategoryModel(c.Orm).GetAll(), 0)
		c.Ctx.JSON(cats)
		return
	}
	c.Ctx.View("backend/content_right.html")
}

func (c *ContentController) NewsList() {
	catid, _ := c.Ctx.URLParamInt64("catid")
	page, _ := c.Ctx.URLParamInt64("page")
	//rows, _ := c.Ctx.URLParamInt64("rows")
	catogoryModel, err := models.NewCategoryModel(c.Orm).GetCategory(catid)
	if err != nil {
		panic(err)
	}
	if catogoryModel.ModelId < 1 {
		helper.Ajax("找不到关联模型", 1, c.Ctx)
		return
	}
	relationDocumentModel := models.NewDocumentModel(c.Orm).GetByID(catogoryModel.ModelId)
	if relationDocumentModel.Id == 0 {
		helper.Ajax("找不到关联模型", 1, c.Ctx)
		return
	}
	if page > 0 {
		sql := []interface{}{fmt.Sprintf("SELECT * FROM `iriscms_%s` WHERE catid=? and deleted_time IS NULL ORDER BY listorder DESC, id DESC", relationDocumentModel.Table), catid}
		contents, err := c.Orm.QueryString(sql...)
		if err != nil {
			glog.Error(helper.GetCallerFuncName(), err.Error())
			helper.Ajax("获取文档列表错误", 1, c.Ctx)
			return
		}

		sql = []interface{}{fmt.Sprintf("SELECT COUNT(*) total FROM `iriscms_%s` WHERE catid=? and deleted_time IS NULL", relationDocumentModel.Table), catid}
		totals, _ := c.Orm.QueryString(sql...)
		var total = "0"
		if len(totals) > 0 {
			total = totals[0]["total"]
		}
		c.Ctx.JSON(map[string]interface{}{"rows": contents, "total": total})
		return
	}

	fields := helper.EasyuiGridfields{
		"排序": {"field": "listorder", "width": "15", "formatter": "contentNewsListOrderFormatter", "index": "0"},
	}

	// 获取所有字段
	dslFields := models.NewDocumentFieldDslModel(c.Orm).GetList(catogoryModel.ModelId)
	var tMapF = map[string]string{}
	var ff []string
	for _, dsl := range dslFields {
		tMapF[dsl.TableField] = dsl.FormName
		ff = append(ff, dsl.TableField)
	}
	var showInPage = map[string]controllers.FieldShowInPageList{}
	_ = json.Unmarshal([]byte(relationDocumentModel.FieldShowInList), &showInPage)
	if len(showInPage) == 0 {
		helper.Ajax("请配置模型字段显隐属性", 1, c.Ctx)
		return
	}
	var index = 1
	// 系统模型需要固定追加标题, 描述等字段
	if relationDocumentModel.IsSystem == models.SYSTEM_TYPE {
		fields["标题"] = map[string]string{"field": "title", "width": "50", "formatter": "contentNewsListOperateFormatter", "index": strconv.Itoa(index)}
		index++
	}
	for _, field := range ff {
		conf := showInPage[field]
		if conf.Show {
			f := map[string]string{"field": field, "width": "50", "index": strconv.Itoa(index)}
			if conf.Formatter != "" {
				f["formatter"] = conf.Formatter
			}
			fields[tMapF[field]] = f
			index++
		}
	}
	fields["管理操作"] = map[string]string{"field": "id", "width": "50", "formatter": "contentNewsListOperateFormatter", "index": strconv.Itoa(index)}
	table := helper.Datagrid("category_categorylist_datagrid", "/b/content/news-list?grid=datagrid&catid="+strconv.Itoa(int(catid)), helper.EasyuiOptions{
		"toolbar":      "content_newslist_datagrid_toolbar",
		"singleSelect": "true",
	}, fields)

	c.Ctx.ViewData("DataGrid", template.HTML(table))
	c.Ctx.ViewData("catid", catid)
	c.Ctx.ViewData("formatters", template.JS(relationDocumentModel.Formatters))
	c.Ctx.View("backend/content_newslist.html")
}

func (c *ContentController) Page() {
	catid, _ := c.Ctx.URLParamInt64("catid")
	if catid == 0 {
		helper.Ajax("页面错误", 1, c.Ctx)
		return
	}
	pageModel := models.NewPageModel(c.Orm)
	page := pageModel.GetPage(catid)
	var hasPage bool = false
	if page.Title != "" {
		hasPage = true
	}
	var res bool
	if c.Ctx.Method() == "POST" {
		page.Title = c.Ctx.FormValue("title")
		page.Content = c.Ctx.FormValue("content")
		page.Keywords = c.Ctx.FormValue("keywords")
		page.Updatetime = int64(helper.GetTimeStamp())
		if !hasPage {
			page.Catid = catid
			res = pageModel.AddPage(page)
		} else {
			res = pageModel.UpdatePage(page)
		}
		if res {
			helper.Ajax("发布成功", 0, c.Ctx)
		} else {
			helper.Ajax("发布失败", 1, c.Ctx)
		}
		return
	}

	c.Ctx.ViewData("catid", catid)
	c.Ctx.ViewData("info", page)
	c.Ctx.View("backend/content_page.html")

}

type customForm map[string]string

func (c customForm) MustCheck() bool {
	var ok bool
	if _, ok = c["catid"]; !ok {
		return false
	}
	if _, ok = c["mid"]; !ok {
		return false
	}
	if _, ok = c["table_name"]; !ok {
		return false
	}
	return true
}

//添加内容
func (c *ContentController) AddContent() {
	if c.Ctx.Method() == "POST" {
		mid, _ := strconv.Atoi(c.Ctx.FormValue("mid"))
		if mid < 1 {
			helper.Ajax("模型参数错误， 无法确定所属模型", 1, c.Ctx)
			return
		}
		var data = customForm{}
		postData := c.Ctx.FormValues()
		for formName, values := range postData {
			data[formName] = values[0] //	todo 多项值根据字段类型合并 strings.Join(values, ",")
		}
		if !data.MustCheck() {
			helper.Ajax("缺少必要参数", 1, c.Ctx)
			return
		}
		data["status"] = "0"
		data["created_time"] = time.Now().In(helper.GetLocation()).Format(helper.TIME_FORMAT)
		model := models.NewDocumentModel(c.Orm).GetByID(int64(mid))
		var fields []string
		var values []interface{}
		for k, v := range data {
			if k == "table_name" {
				continue
			}
			fields = append(fields, "`"+k+"`")
			values = append(values, v)
		}
		// 通用字段更新

		if model.IsSystem == models.SYSTEM_TYPE {

		}
		params := append([]interface{}{fmt.Sprintf("INSERT INTO `iriscms_%s` (%s) VALUES (%s)", model.Table, strings.Join(fields, ","), strings.TrimRight(strings.Repeat("?,", len(values)), ","))}, values...)
		// 先直接入库对应表内
		insertID, err := c.Orm.Exec(params...)
		if err != nil {
			glog.Error(err)
			helper.Ajax("添加失败:"+err.Error(), 1, c.Ctx)
			return
		}
		id, _ := insertID.LastInsertId()
		if id > 0 {
			helper.Ajax(id, 0, c.Ctx)
		} else {
			helper.Ajax("添加失败", 1, c.Ctx)
		}
		return
	}
	//根据catid读取出相应的添加模板
	catid, _ := c.Ctx.URLParamInt64("catid")
	if catid == 0 {
		helper.Ajax("参数错误", 1, c.Ctx)
		return
	}
	cat, err := models.NewCategoryModel(c.Orm).GetCategory(catid)
	if err != nil {
		helper.Ajax("读取数据错误:"+err.Error(), 1, c.Ctx)
		return
	}
	if cat.Catid == 0 {
		helper.Ajax("不存在的分类", 1, c.Ctx)
		return
	}
	c.Ctx.ViewData("category", cat)
	c.Ctx.ViewData("form", template.HTML(buildModelForm(c.Orm, cat.ModelId, nil)))
	c.Ctx.ViewData("submitURL", template.HTML("/b/content/add"))
	c.Ctx.View("backend/model_publish.html")
}

//修改内容
func (c *ContentController) EditContent() {
	//根据catid读取出相应的添加模板
	catid, _ := c.Ctx.URLParamInt64("catid")
	id, _ := c.Ctx.URLParamInt64("id")
	if catid < 1 || id < 1 {
		helper.Ajax("参数错误", 1, c.Ctx)
		return
	}
	catogoryModel, err := models.NewCategoryModel(c.Orm).GetCategory(catid)
	if err != nil {
		panic(err)
	}
	if catogoryModel.ModelId < 1 {
		helper.Ajax("找不到关联模型", 1, c.Ctx)
		return
	}
	relationDocumentModel := models.NewDocumentModel(c.Orm).GetByID(catogoryModel.ModelId)
	if relationDocumentModel.Id == 0 {
		helper.Ajax("找不到关联模型", 1, c.Ctx)
		return
	}
	sql := []interface{}{fmt.Sprintf("SELECT * FROM `iriscms_%s` WHERE catid=? and deleted_time IS NULL AND id = ? LIMIT 1", relationDocumentModel.Table), catid, id}
	contents, err := c.Orm.QueryString(sql...)
	if err != nil {
		glog.Error(helper.GetCallerFuncName(), err.Error())
		helper.Ajax("获取文章内容错误", 1, c.Ctx)
		return
	}

	if len(contents) == 0 {
		helper.Ajax("文章不存在或已删除", 1, c.Ctx)
		return
	}
	if c.Ctx.Method() == "POST" {
		var data = customForm{}
		postData := c.Ctx.FormValues()
		for formName, values := range postData {
			data[formName] = values[0] //	todo 多项值根据字段类型合并 strings.Join(values, ",")
		}
		if !data.MustCheck() {
			helper.Ajax("缺少必要参数", 1, c.Ctx)
			return
		}
		delete(data, "id")
		data["status"] = "0"
		data["updated_time"] = time.Now().In(helper.GetLocation()).Format(helper.TIME_FORMAT)
		var sets []string
		var values []interface{}
		for k, v := range data {
			if k == "table_name" {
				continue
			}
			sets = append(sets, "`"+k+"`= ?")
			values = append(values, v)
		}
		// 通用字段更新
		if relationDocumentModel.IsSystem == models.SYSTEM_TYPE {

		}
		values = append(values, id, catid)
		params := append([]interface{}{fmt.Sprintf("UPDATE `iriscms_%s` SET %s WHERE id=? and catid=?", relationDocumentModel.Table, strings.Join(sets, ", "))}, values...)
		insertID, err := c.Orm.Exec(params...)
		if err != nil {
			glog.Error(err)
			helper.Ajax("修改失败:"+err.Error(), 1, c.Ctx)
			return
		}
		res, _ := insertID.RowsAffected()
		if res > 0 {
			helper.Ajax("修改成功", 0, c.Ctx)
		} else {
			helper.Ajax("修改失败", 1, c.Ctx)
		}
		return
	}
	c.Ctx.ViewData("form", template.HTML(buildModelForm(c.Orm, catogoryModel.ModelId, contents[0])))
	c.Ctx.ViewData("category", catogoryModel)
	c.Ctx.ViewData("submitURL", template.HTML("/b/content/edit"))
	c.Ctx.View("backend/model_publish.html")
}

//删除内容
func (c *ContentController) DeleteContent() {
	catid, _ := c.Ctx.URLParamInt64("catid")
	id, _ := c.Ctx.URLParamInt64("id")
	if catid < 1 || id < 1 {
		helper.Ajax("参数错误", 1, c.Ctx)
		return
	}
	catogoryModel, err := models.NewCategoryModel(c.Orm).GetCategory(catid)
	if err != nil {
		panic(err)
	}
	if catogoryModel.ModelId < 1 {
		helper.Ajax("找不到关联模型", 1, c.Ctx)
		return
	}
	relationDocumentModel := models.NewDocumentModel(c.Orm).GetByID(catogoryModel.ModelId)
	if relationDocumentModel.Id == 0 {
		helper.Ajax("找不到关联模型", 1, c.Ctx)
		return
	}
	sqlOrArgs := []interface{}{fmt.Sprintf("UPDATE `iriscms_%s` SET `deleted_time`='"+time.Now().In(helper.GetLocation()).Format(helper.TIME_FORMAT)+"' WHERE id = ? and catid=?", relationDocumentModel.Table), id, catid}
	res, err := c.Orm.Exec(sqlOrArgs...)
	if err != nil {
		glog.Error(helper.GetCallerFuncName(), err.Error())
		helper.Ajax("删除失败", 1, c.Ctx)
		return
	}
	if ret, _ := res.RowsAffected(); ret > 0 {
		helper.Ajax("删除成功", 0, c.Ctx)
	} else {
		helper.Ajax("删除失败", 1, c.Ctx)
	}
}

//排序内容
func (c *ContentController) OrderContent() {
	data := c.Ctx.FormValues()
	var order = map[string]string{}
	for k, v := range data {
		order[strings.ReplaceAll(strings.ReplaceAll(k, "order[", ""), "]", "")] = v[0]
	}
	id, _ := c.Ctx.URLParamInt64("catid")
	if id < 1 {
		helper.Ajax("参数错误", 1, c.Ctx)
		return
	}
	catogoryModel, err := models.NewCategoryModel(c.Orm).GetCategory(id)
	if err != nil {
		panic(err)
	}
	if catogoryModel.ModelId < 1 {
		helper.Ajax("找不到关联模型", 1, c.Ctx)
		return
	}
	relationDocumentModel := models.NewDocumentModel(c.Orm).GetByID(catogoryModel.ModelId)
	if relationDocumentModel.Id == 0 {
		helper.Ajax("找不到关联模型", 1, c.Ctx)
		return
	}
	for artID, orderNum := range order {
		sqlOrArgs := []interface{}{fmt.Sprintf("UPDATE `iriscms_%s` SET `listorder`=? , updated_time = '"+time.Now().In(helper.GetLocation()).Format(helper.TIME_FORMAT)+"' WHERE id = ? and catid=?", relationDocumentModel.Table), orderNum, artID, id}
		if _, err := c.Orm.Exec(sqlOrArgs...); err != nil {
			glog.Error(helper.GetCallerFuncName(), err.Error())
			helper.Ajax("更新文档排序失败", 1, c.Ctx)
			return
		}
	}
	helper.Ajax("更新排序成功", 0, c.Ctx)
}
