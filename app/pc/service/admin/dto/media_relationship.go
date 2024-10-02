package dto

import (
    "go-admin/app/pc/models"
    "go-admin/common/dto"
    common "go-admin/common/models"
)

type MediaRelationshipGetPageReq struct {
    dto.Pagination `search:"-"`
    MediaRelationshipOrder
}

type MediaRelationshipOrder struct {
    Id             string `form:"idOrder"  search:"type:order;column:id;table:media_relationship"`
    MediaTypeId    string `form:"mediaTypeIdOrder"  search:"type:order;column:media_type_id;table:media_relationship"`
    BuszId         string `form:"buszIdOrder"  search:"type:order;column:busz_id;table:media_relationship"`
    MediaInstantId string `form:"mediaInstantIdOrder"  search:"type:order;column:media_instant_id;table:media_relationship"`
    Seq            string `form:"seqOrder"  search:"type:order;column:seq;table:media_relationship"`
    ImgFeatureId   string `form:"imgFeatureIdOrder"  search:"type:order;column:img_feature_id;table:media_relationship"`
    Watermark      string `form:"watermarkOrder"  search:"type:order;column:watermark;table:media_relationship"`
    CreatedAt      string `form:"createdAtOrder"  search:"type:order;column:created_at;table:media_relationship"`
    UpdatedAt      string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:media_relationship"`
    DeletedAt      string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:media_relationship"`
    CreateBy       string `form:"createByOrder"  search:"type:order;column:create_by;table:media_relationship"`
    CreateByName   string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:media_relationship"`
    UpdateBy       string `form:"updateByOrder"  search:"type:order;column:update_by;table:media_relationship"`
    UpdateByName   string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:media_relationship"`
}

func (m *MediaRelationshipGetPageReq) GetNeedSearch() interface{} {
    return *m
}

type MediaRelationshipInsertReq struct {
    Id                     int                    `json:"-" comment:"id"` // id
    MediaTypeId            int                    `json:"mediaTypeId" comment:"媒体类型id"`
    BuszId                 string                 `json:"buszId" comment:"业务id"`
    MediaInstantId         int                    `json:"mediaInstantId" comment:""`
    Seq                    int                    `json:"seq" comment:"排序号"`
    ImgFeatureId           string                 `json:"imgFeatureId" comment:"图片的特征id"`
    Watermark              int                    `json:"watermark" comment:"图片是否有水印,0无水印,1,图片仅供参考,2,产品仅为附件"`
    MediaInstanceInsertReq MediaInstanceInsertReq `json:"mediaInstanceInsertReq"`
    common.ControlBy
}

func (s *MediaRelationshipInsertReq) Generate(model *models.MediaRelationship) {
    if s.Id == 0 {
        model.Model = common.Model{Id: s.Id}
    }
    model.MediaTypeId = s.MediaTypeId
    model.BuszId = s.BuszId
    model.MediaInstantId = s.MediaInstantId
    model.Seq = s.Seq
    model.ImgFeatureId = s.ImgFeatureId
    model.Watermark = s.Watermark
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.CreateByName = s.CreateByName
}

func (s *MediaRelationshipInsertReq) GetId() interface{} {
    return s.Id
}

type MediaRelationshipUpdateReq struct {
    Id                     int                    `uri:"id" comment:"id"` // id
    MediaTypeId            int                    `json:"mediaTypeId" comment:"媒体类型id"`
    BuszId                 string                 `json:"buszId" comment:"业务id"`
    MediaInstantId         int                    `json:"mediaInstantId" comment:""`
    Seq                    int                    `json:"seq" comment:"排序号"`
    ImgFeatureId           string                 `json:"imgFeatureId" comment:"图片的特征id"`
    Watermark              int                    `json:"watermark" comment:"图片是否有水印,0无水印,1,图片仅供参考,2,产品仅为附件"`
    MediaInstanceInsertReq MediaInstanceInsertReq `json:"mediaInstanceInsertReq"`
    common.ControlBy
}

func (s *MediaRelationshipUpdateReq) Generate(model *models.MediaRelationship) {
    if s.Id == 0 {
        model.Model = common.Model{Id: s.Id}
    }
    model.MediaTypeId = s.MediaTypeId
    model.BuszId = s.BuszId
    model.MediaInstantId = s.MediaInstantId
    model.Seq = s.Seq
    model.ImgFeatureId = s.ImgFeatureId
    model.Watermark = s.Watermark
    model.UpdateBy = s.UpdateBy
    model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *MediaRelationshipUpdateReq) GetId() interface{} {
    return s.Id
}

// MediaRelationshipGetReq 功能获取请求参数
type MediaRelationshipGetReq struct {
    Id int `uri:"id"`
}

func (s *MediaRelationshipGetReq) GetId() interface{} {
    return s.Id
}

// MediaRelationshipDeleteReq 功能删除请求参数
type MediaRelationshipDeleteReq struct {
    Ids []int `json:"ids"`
}

func (s *MediaRelationshipDeleteReq) GetId() interface{} {
    return s.Ids
}
