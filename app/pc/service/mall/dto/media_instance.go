package dto

import (

	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type MediaInstanceGetPageReq struct {
	dto.Pagination     `search:"-"`
    MediaInstanceOrder
}

type MediaInstanceOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:media_instance"`
    MediaDir string `form:"mediaDirOrder"  search:"type:order;column:media_dir;table:media_instance"`
    MediaName string `form:"mediaNameOrder"  search:"type:order;column:media_name;table:media_instance"`
    
}

func (m *MediaInstanceGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type MediaInstanceInsertReq struct {
    Id int `json:"-" comment:""` // 
    MediaDir string `json:"mediaDir" comment:"文件路径"`
    MediaName string `json:"mediaName" comment:"文件名"`
    common.ControlBy
}

func (s *MediaInstanceInsertReq) Generate(model *models.MediaInstance)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.MediaDir = s.MediaDir
    model.MediaName = s.MediaName
}

func (s *MediaInstanceInsertReq) GetId() interface{} {
	return s.Id
}

type MediaInstanceUpdateReq struct {
    Id int `uri:"id" comment:""` // 
    MediaDir string `json:"mediaDir" comment:"文件路径"`
    MediaName string `json:"mediaName" comment:"文件名"`
    common.ControlBy
}

func (s *MediaInstanceUpdateReq) Generate(model *models.MediaInstance)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.MediaDir = s.MediaDir
    model.MediaName = s.MediaName
}

func (s *MediaInstanceUpdateReq) GetId() interface{} {
	return s.Id
}

// MediaInstanceGetReq 功能获取请求参数
type MediaInstanceGetReq struct {
     Id int `uri:"id"`
}
func (s *MediaInstanceGetReq) GetId() interface{} {
	return s.Id
}

// MediaInstanceDeleteReq 功能删除请求参数
type MediaInstanceDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *MediaInstanceDeleteReq) GetId() interface{} {
	return s.Ids
}
