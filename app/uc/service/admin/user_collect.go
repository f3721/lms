package admin

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/jinzhu/copier"
	modelsPc "go-admin/app/pc/models"
	pcClient "go-admin/common/client/pc"
	"gorm.io/gorm"

	"go-admin/app/uc/models"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type UserCollect struct {
	service.Service
}

// GetPage 获取UserCollect列表
func (e *UserCollect) GetPage(c *dto.UserCollectGetPageReq, p *actions.DataPermission, list *[]models.UserCollect, count *int64) error {
	var err error
	var data models.UserCollect

	err = e.Orm.Model(&data).Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("UserCollectService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// GetListPage  获取UserCollect列表(包含商品名 商品图片)
func (e *UserCollect) GetListPage(c *dto.UserCollectGetPageReq, p *actions.DataPermission, resList *[]dto.UserCollectGetListPageRes, count *int64) error {
	var err error
	var list []models.UserCollect

	err = e.GetPage(c, p, &list, count)
	if err != nil {
		return err
	}
	var goodIds []int
	for _, collect := range list {
		goodIds = append(goodIds, collect.GoodsId)
	}

	// 获取商品信息
	goodsResult := pcClient.ApiByDbContext(e.Orm).GetGoodsById(goodIds)

	goodsResultInfo := &struct {
		response.Response
		Data []modelsPc.Goods
	}{}
	goodsResult.Scan(goodsResultInfo)
	goodsMap := make(map[int]modelsPc.Goods)
	for _, goods := range goodsResultInfo.Data {
		goodsMap[goods.Id] = goods
	}

	var resListCopy []dto.UserCollectGetListPageRes
	for _, collect := range list {
		goods, ok := goodsMap[collect.GoodsId]
		productName := "商品不存在"
		productImage := ""
		if ok {
			productName = goods.Product.NameZh
			if goods.Product.MediaRelationship != nil && len(*goods.Product.MediaRelationship) > 0 {
				mediaRelationship := (*goods.Product.MediaRelationship)[0]
				productImage = mediaRelationship.MediaInstant.MediaDir
			}
		}
		resListCopy = append(resListCopy, dto.UserCollectGetListPageRes{UserCollect: collect, ProductName: productName, ProductImage: productImage})
	}
	copier.Copy(&resList, resListCopy)
	return nil
}

// Get 获取UserCollect对象
func (e *UserCollect) Get(d *dto.UserCollectGetReq, p *actions.DataPermission, model *models.UserCollect) error {
	var data models.UserCollect

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUserCollect error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建UserCollect对象
func (e *UserCollect) Insert(c *dto.UserCollectInsertReq) error {
	var err error
	var data models.UserCollect
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("UserCollectService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改UserCollect对象
func (e *UserCollect) Update(c *dto.UserCollectUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.UserCollect{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("UserCollectService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除UserCollect
func (e *UserCollect) Remove(d *dto.UserCollectDeleteReq, p *actions.DataPermission) error {
	var data models.UserCollect

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveUserCollect error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
