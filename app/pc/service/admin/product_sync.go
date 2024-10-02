package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"go-admin/app/pc/models"
	"go-admin/app/pc/service/admin/dto"
	wc "go-admin/app/wc/models"
	"go-admin/common/client"
	"go-admin/common/client/sp"
	spDto "go-admin/common/dto/sp"
	"go-admin/common/global"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ProductSync struct {
	service.Service
}

const VendorName = "西域智慧供应链（上海）股份公司"

var skuImageMap map[string]spDto.ProductImageResult
var spClient *sp.AccessToken

func (e *ProductSync) SupplierProductSync(c *gin.Context) error {

	spApiUrl := client.ReadConfig("apiUrl.sp")
	clientId := client.ReadConfig("accessToken.client_id")
	clientSecret := client.ReadConfig("accessToken.client_secret")
	username := client.ReadConfig("accessToken.username")
	password := client.ReadConfig("accessToken.password")
	// sp初始化
	spClient = sp.New(clientId, clientSecret, username, password, spApiUrl)
	// 取消息池里的商品
	result, err := spClient.GetPushMsg(&spDto.GetMsgReq{Type: 22})
	if err != nil || !result.Success {
		return err
	}

	if len(result.Result) == 0 {
		return errors.New("当前消息池无商品 ")
	}

	// 获取到所有的SKU
	skuCodes := getSkuCodes(result)
	// 批量获取图片
	productImages, err := spClient.GetProductImage(&spDto.GetProductImageReq{
		Sku: skuCodes,
	})
	if err != nil || !result.Success {
		return err
	}

	skuImageMap = make(map[string]spDto.ProductImageResult)
	if err = utils.StructColumn(&skuImageMap, productImages.Result, "", "Sku"); err != nil {
		return err
	}

	// 遍历商品 获取商品详情及商品图片
	for _, msgResult := range result.Result {
		err = e.makeOrm(c, msgResult.Result.TenantId)
		if err != nil {
			// 数据库连接失败
			e.syncFailedPush(&spDto.SyncFailedPush{
				SkuId:           msgResult.Result.SkuId,
				RejectionReason: "数据库连接失败, 数据库未初始化或状态不可用",
				TenantName:      msgResult.Result.TenantName,
				MessageId:       msgResult.Id,
			})
			continue
		}
		// 数据校验
		productSyncReq := dto.ProductSyncReq{}

		e.generateImages(&productSyncReq, &msgResult)
		// 查商品详情
		detail, err := spClient.GetProductDetail(&spDto.GetProductDetailReq{
			Sku: msgResult.Result.SkuId,
		})
		if err != nil || !result.Success {
			e.syncFailedPush(&spDto.SyncFailedPush{
				SkuId:           msgResult.Result.SkuId,
				RejectionReason: err.Error(),
				TenantName:      msgResult.Result.TenantName,
				MessageId:       msgResult.Id,
			})
			continue
		}
		// 数据校验 末级产线是否存在
		// 主产线
		categoryId := getLastCategoriesId(detail.Result.Category)
		productSyncReq.ProductCategoryPrimary = categoryId
		productSyncReq.ProductCategory = append(productSyncReq.ProductCategory, dto.ProductCategory{
			CategoryId:   categoryId,
			CategoryName: "",
		})
		e.generate(&productSyncReq, &detail.Result)
		// 货主信息
		err = e.generateVendorId(&productSyncReq)
		if err != nil {
			// 插入错误
			e.syncFailedPush(&spDto.SyncFailedPush{
				SkuId:           msgResult.Result.SkuId,
				RejectionReason: fmt.Sprintf("货主信息(%s)不存在", VendorName),
				TenantName:      msgResult.Result.TenantName,
				MessageId:       msgResult.Id,
			})
			continue
		}
		// 品牌信息
		err = e.generateBrand(&productSyncReq, detail.Result.BrandName)
		if err != nil {
			// 插入错误
			e.syncFailedPush(&spDto.SyncFailedPush{
				SkuId:           msgResult.Result.SkuId,
				RejectionReason: err.Error(),
				TenantName:      msgResult.Result.TenantName,
				MessageId:       msgResult.Id,
			})
			continue
		}
		// 判断商品货主+货主SKU是否已经存在
		productService := Product{e.Service}
		var product models.Product
		if productService.vendorIsBind(&dto.ProductInsertReq{VendorId: productSyncReq.VendorId, SupplierSkuCode: productSyncReq.SupplierSkuCode}, &product) {
			// 存在 更新
			productUpdateReq := dto.ProductUpdateReq{}
			_ = copier.Copy(&productUpdateReq, &product)
			_ = copier.Copy(&productUpdateReq, &productSyncReq)
			err = e.Update(&productUpdateReq)
			if err != nil {
				e.syncFailedPush(&spDto.SyncFailedPush{
					SkuId:           msgResult.Result.SkuId,
					RejectionReason: err.Error(),
					TenantName:      msgResult.Result.TenantName,
					MessageId:       msgResult.Id,
				})
				continue
			}
		} else {
			// 不存在
			productInsertReq := dto.ProductInsertReq{}
			_ = copier.Copy(&productInsertReq, &productSyncReq)
			err = e.Insert(&productInsertReq)
			if err != nil {
				e.syncFailedPush(&spDto.SyncFailedPush{
					SkuId:           msgResult.Result.SkuId,
					RejectionReason: err.Error(),
					TenantName:      msgResult.Result.TenantName,
					MessageId:       msgResult.Id,
				})
				continue
			}
		}
		// 成功就直接消费掉消息
		spClient.MessageDelete(&spDto.MessageDelete{Id: msgResult.Id})
	}
	return nil
}

func getSkuCodes(result *spDto.GetMsgResp) []string {
	skuCodes := make([]string, 0)
	for _, msgResult := range result.Result {
		skuCodes = append(skuCodes, msgResult.Result.SkuId)
	}
	return lo.Uniq[string](skuCodes)
}

// 品牌信息
func (e *ProductSync) generateBrand(product *dto.ProductSyncReq, brandName string) error {
	brandService := Brand{e.Service}
	var brand models.Brand
	err := brandService.FindBrandInfoByName(brandName, "", &brand)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 新增
		err = brandService.Insert(&dto.BrandInsertReq{
			BrandZh: brandName,
			Status:  1,
			Confirm: 1,
		})
		if err != nil {
			return err
		}
		err = brandService.FindBrandInfoByName(brandName, "", &brand)
		if err != nil {
			return err
		}
		product.BrandId = brand.Id
		return nil
	}
	if err != nil {
		return err
	}
	product.BrandId = brand.Id
	return nil
}

// 基础信息
func (e *ProductSync) generate(product *dto.ProductSyncReq, req *spDto.ProductDetail) {
	product.SupplierSkuCode = req.Sku
	product.NameZh = req.Name
	product.SalesUom = req.SaleUnit
	product.SalesMoq = req.Moq
	product.MfgModel = req.MfgSku
	product.PackWeight = req.Weight
	if req.IsReturn == 1 {
		product.RefundFlag = 0
	} else {
		product.RefundFlag = 1
	}
	product.BriefDesc = findStringSubmatch(req.Introduction) + findStringSubmatch(req.WareQD)
	//product.SalesPhysicalFactor = req.WareNum
	product.Tax = "0.13"
	product.Status = 1 // 默认待审核
}

// 图片信息
func (e *ProductSync) generateImages(productSyncReq *dto.ProductSyncReq, msg *spDto.GetMsgResult) {
	if images, ok := skuImageMap[msg.Result.SkuId]; ok {
		for _, image := range images.SkuPic {
			imagePath := strings.Split(image.Path, ",")
			productSyncReq.ProductImage = append(productSyncReq.ProductImage, dto.MediaInstanceInsertReq{
				MediaDir:  imagePath[0],
				MediaName: filepath.Base(imagePath[0]),
				WaterMark: 0,
				Seq:       image.OrderSort,
			})
		}
	}
}

// 货主信息
func (e *ProductSync) generateVendorId(productSyncReq *dto.ProductSyncReq) error {
	vendors := wc.Vendors{}
	res, err := vendors.FindOneByName(e.Orm, VendorName)
	if err != nil {
		return err
	}
	productSyncReq.VendorId = res.Id
	return nil
}

// Insert 创建Product对象
func (e *ProductSync) Insert(c *dto.ProductInsertReq) error {
	productService := Product{e.Service}
	// 自动生成SKU
	skuCode := productService.createSku()
	c.SkuCode = skuCode

	var err error
	var data models.Product
	c.Generate(&data)
	//数据插入校验
	err = productService.validateForm(c, "sync")
	if err != nil {
		return err
	}

	tx := productService.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = tx.Create(&data).Error
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("ProductService Insert error:%s \r\n", err)
		return err
	}

	// 判断是否上传了图片
	if len(c.ProductImage) > 0 {
		err = productService.createProductImage(tx, c.ProductImage, &data)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 创建商品分类
	err = productService.createProductCategory(tx, c.ProductCategory, data.SkuCode, c.ProductCategoryPrimary)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 保存商品扩展属性
	_ = productService.SaveProductExtAttribute(c)

	// 生成日志
	dataLog, _ := json.Marshal(&c)
	productLog := models.ProductLog{
		DataId:     data.Id,
		Type:       global.LogTypeCreate,
		Data:       string(dataLog),
		BeforeData: "",
		AfterData:  string(dataLog),
		ControlBy:  c.ControlBy,
	}
	_ = productLog.CreateLog("product", e.Orm)

	return nil
}

// Update 修改Product对象
func (e *ProductSync) Update(c *dto.ProductUpdateReq) error {
	productService := Product{e.Service}
	var data = models.Product{}
	e.Orm.First(&data, c.GetId())

	var form dto.ProductInsertReq
	copier.Copy(&form, &c)
	err := productService.validateForm(&form, "sync")
	if err != nil {
		return err
	}

	err = productService._update(c, data)
	if err != nil {
		return err
	}
	return nil
}

func (e *ProductSync) makeOrm(c *gin.Context, tId string) error {
	// 获取租户ID拿到数据库前缀
	tenantId, _ := strconv.Atoi(tId)
	tenantDBPrefix := global.GetSyncTenantDBNameWithDB(e.Orm, tenantId)
	tenantDb := sdk.Runtime.GetDbByKey(tenantDBPrefix)
	if tenantDb == nil {
		// 数据库连接失败
		return errors.New("数据库连接失败！")
	}
	// 将tenant-id放到db.statements.context中，用于接口或方法中 获取数据库前缀
	c.Request.Header.Set("tenant-id", global.EncryptTenantId(tenantId))
	tenantDb = tenantDb.Session(&gorm.Session{
		Context: c,
	})
	e.Orm = tenantDb
	return nil
}

// 错误信息同步并消费消息
func (e *ProductSync) syncFailedPush(in *spDto.SyncFailedPush) error {
	syncFailed := spDto.SyncFailed{
		RejectionState: "4",
	}
	syncFailedSkus := make([]spDto.SyncFailedSkus, 0)
	syncFailedSkus = append(syncFailedSkus, spDto.SyncFailedSkus{
		SkuId:           in.SkuId,
		RejectionReason: in.RejectionReason,
		TenantName:      in.TenantName,
	})
	syncFailed.Skus = syncFailedSkus
	result, err := spClient.SyncFailedPush(&syncFailed)
	if err != nil {
		return err
	}
	if result.Success {
		spClient.MessageDelete(&spDto.MessageDelete{
			Id: in.MessageId,
		})
	}
	return nil
}

func getLastCategoriesId(categories []string) (categoryId int) {
	// 去除末尾可能的空值
	categories = lo.Compact[string](categories)
	categoryIdStr := categories[len(categories)-1]
	categoryId, _ = strconv.Atoi(categoryIdStr)
	return
}

func findStringSubmatch(htmlStr string) string {
	if htmlStr == "" {
		return ""
	}
	pattern, _ := regexp.Compile("<body.*?>(.*?)</body>")
	result := pattern.FindStringSubmatch(htmlStr)
	return result[1]
}
