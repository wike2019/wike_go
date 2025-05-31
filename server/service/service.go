package service

import (
	model2 "github.com/wike2019/wike_go/common"
	"github.com/wike2019/wike_go/pkg/core"
	"github.com/wike2019/wike_go/pkg/func/ctl"
	model "github.com/wike2019/wike_go/server/model"
	"gorm.io/gorm"
)

type Service struct {
	DB *core.CoreDb
}

func (this *Service) GetApiService(query *ctl.Page, body *model.APISearch, header *model2.Empty, context *ctl.Controller) (*ctl.PageList[model.API], error) {
	data, total, err := ctl.ListItem[model.API](this.DB.DB, query.Offset, query.Count, body, nil)
	if err != nil {
		return nil, err
	}
	context.FormatTotal(total)

	return ctl.List[model.API]("获取api列表成功", data, context), nil
}

func (this *Service) GetApiFrontService(query *model2.Empty, body *model.APISearch, header *model2.Empty, context *ctl.Controller) (*ctl.DataList[[]model.APIFront], error) {
	data, err := ctl.ListItemAll[model.API](this.DB.DB, body, nil)
	if err != nil {
		return nil, err
	}
	res := this.convertToAPIFront(data)

	return ctl.Item[[]model.APIFront]("获取api列表成功", res, context), nil
}
func CoreService(DB *core.CoreDb) *Service {
	return &Service{DB: DB}
}
func (this *Service) convertToAPIFront(apiList []model.API) []model.APIFront {
	groupedAPIs := make(map[string][]model.ChildAPi)

	// 按 Group 分类
	for _, api := range apiList {
		child := model.ChildAPi{
			Value:  api.ID,
			Label:  api.Name,
			Path:   api.Path,
			Method: api.Method,
		}
		groupedAPIs[api.Group] = append(groupedAPIs[api.Group], child)
	}

	// 组装 APIFront 结构
	var apiFrontList []model.APIFront
	for group, children := range groupedAPIs {
		apiFrontList = append(apiFrontList, model.APIFront{
			Value:    0,
			Label:    group,
			Children: children,
		})
	}

	return apiFrontList
}

func (this *Service) DictionaryListService(query *ctl.Page, body *model.SysDictionarySearch, header *model2.Empty, context *ctl.Controller) (*ctl.PageList[model.SysDictionary], error) {
	data, total, err := ctl.ListItem[model.SysDictionary](this.DB.DB, query.Offset, query.Count, body, nil)
	if err != nil {
		return nil, err
	}
	context.FormatTotal(total)

	return ctl.List[model.SysDictionary]("获取字典列表成功", data, context), nil
}
func (this *Service) DictionaryItemService(query *model2.IDSearch, body *model2.Empty, header *model2.Empty, context *ctl.Controller) (*ctl.DataList[model.SysDictionary], error) {
	res, err := ctl.GetItem[model.SysDictionary](this.DB.DB, query.ID, func(db *gorm.DB) *gorm.DB {
		return db.Preload("SysDictionaryDetails", func(db *gorm.DB) *gorm.DB {
			return db.Order("`sort` ASC") // 按 Author 的 name 字段排序
		})
	})
	if !ctl.IsExist(err) {
		return ctl.Item[model.SysDictionary]("获取字典列表成功", res, context), nil
	}
	return ctl.Item[model.SysDictionary]("获取字典列表成功", res, context), err
}
func (this *Service) DictionarySearchService(query *model2.KeyWordsSearch, body *model2.Empty, header *model2.Empty, context *ctl.Controller) (*ctl.DataList[model.SysDictionary], error) {
	res, err := ctl.GetItemByFunc[model.SysDictionary](this.DB.DB, func(db *gorm.DB) *gorm.DB {
		return db.Where("type=?", query.KeyWords).Preload("SysDictionaryDetails", func(db *gorm.DB) *gorm.DB {
			return db.Order("`sort` ASC").Where("status=2") // 按 Author 的 name 字段排序
		})
	})
	return ctl.Item[model.SysDictionary]("获取字典列表成功", res, context), err
}

func (this Service) RoleListService(query *model.Role, body *model2.Empty, header *model2.Empty, context *ctl.Controller) (*ctl.DataList[model.RoleService], error) {
	list, err := ctl.ListItemAll[model.Role](this.DB.DB, query, nil)
	res := make([]string, 0)
	for _, item := range list {
		res = append(res, item.Child)
	}
	return ctl.Item[model.RoleService]("获取api列表成功", model.RoleService{
		List:  list,
		Names: res,
	}, context), err
}

func (this Service) RuleCreateService(query *model2.Empty, body *model.RuleInput, header *model2.Empty, context *ctl.Controller) (*ctl.DataList[model2.Empty], error) {
	err := this.DB.DB.Transaction(func(tx *gorm.DB) error {
		delObj := model.Rule{}
		err := tx.Where(" role =?", body.Role).Delete(&delObj).Error
		for _, item := range body.APIId {
			if item == 0 {
				continue
			}
			old := model.Rule{}
			err = tx.Where("api_id=? and  role =?", item, body.Role).First(&old).Error
			if ctl.IsExist(err) {
				continue
			}
			api := model.API{}
			err = tx.Where("id=?", item).First(&api).Error
			if err != nil {
				return err
			}
			data := &model.Rule{}
			data.APIId = item
			data.Role = body.Role
			data.Path = api.Path
			data.Method = api.Method
			err = tx.Create(data).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	return ctl.Item[model2.Empty]("同步接口成功", model2.Empty{}, context), err
}

func (this Service) RuleInfoService(query *model.RoleInput, body *model2.Empty, header *model2.Empty, context *ctl.Controller) (*ctl.DataList[model.APIDataForRole], error) {
	data, err := ctl.ListItemAll[model.Rule](this.DB.DB, query, nil)
	if err != nil {
		return nil, err
	}
	ids := make([]uint, 0)
	for _, item := range data {
		ids = append(ids, item.APIId)
	}
	res := model.APIDataForRole{}
	res.ApiIDs = ids
	res.APIList = data
	var roleList []model.Role
	err = this.DB.DB.Raw(`WITH RECURSIVE parent_chain AS (
   -- 初始查询，找到child = 'admin'的记录
   SELECT id, parent, child
   FROM sys_role
   WHERE child = ?

   UNION ALL

   -- 递归查询，找到每个parent的parent
   SELECT sr.id, sr.parent, sr.child
   FROM sys_role sr
   INNER JOIN parent_chain pc ON sr.child = pc.parent
)
SELECT *
FROM parent_chain;`, query.Role).Debug().Find(&roleList).Error
	if err != nil {
		return nil, err
	}
	roleStr := make([]string, 0)
	for _, item := range roleList {
		if item.Parent != "" {
			roleStr = append(roleStr, item.Parent)
		}
	}
	var apiList []model.Rule
	err = this.DB.DB.Raw(`SELECT *
FROM (
   SELECT *, ROW_NUMBER() OVER (PARTITION BY api_id ORDER BY created_at ASC) AS rn
   FROM sys_rule
   WHERE role IN (?) and sys_rule.deleted_at IS NULL

) t
WHERE t.rn = 1;`, roleStr).Debug().Find(&apiList).Error
	if err != nil {
		return nil, err
	}
	res.APIParentList = apiList
	return ctl.Item[model.APIDataForRole]("同步接口成功", res, context), err
}
