package models

import (
    "fmt"
    modelsUc "go-admin/app/uc/models"
    "go-admin/common/models"
    "go-admin/common/utils"
    "gorm.io/gorm"
    "strconv"
    "strings"
    "time"
)

type PresoLog struct {
    models.Model

    PresoNo string `json:"presoNo" gorm:"type:varchar(30);comment:审批单编号"`
    Type int `json:"type" gorm:"type:tinyint(1);comment:0:初始下单  1 系统审批   2 审批流"` 
    ApproveflowId int `json:"approveflowId" gorm:"type:int unsigned;comment:审批流水主键"` 
    ApproveflowItemId int `json:"approveflowItemId" gorm:"type:int unsigned;comment:审批流条目主键"` 
    UserId int `json:"userId" gorm:"type:int unsigned;comment:审批流条目审批人ID"` 
    OperUser string `json:"operUser" gorm:"type:varchar(20);comment:操作人"` 
    OperContent string `json:"operContent" gorm:"type:varchar(1000);comment:操作内容"` 
    ApproveStatus int `json:"approveStatus" gorm:"type:tinyint(1);comment:审批状态  -1 审批不通过  1 审批通过 0 初始提交审批 -2 超时"` 
    ApproveRemark string `json:"approveRemark" gorm:"type:varchar(400);comment:审批备注"` 
    Step int `json:"step" gorm:"type:tinyint(1);comment:审批流步骤"`
    ApproveRankType int `json:"approveRankType" gorm:"type:tinyint(1);comment:审批层级类型：1-普签 2-会签 3-或签"`
    IsAutoApprove int `json:"isAutoApprove" gorm:"type:tinyint(1);comment:否自动审批 0-否 1-是"`
    TotalPriceLimit float64 `json:"totalPriceLimit" gorm:"type:decimal(10,2);comment:采购限额"`
    LimitType int `json:"limitType" gorm:"type:tinyint(1);comment:采购限额类型(1:订单金额,2:商品金额)"`
    CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
    UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
    models.ControlBy
}

func (PresoLog) TableName() string {
    return "preso_log"
}

func (e *PresoLog) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *PresoLog) GetId() interface{} {
	return e.Id
}

type WorkFlowNodes struct {
    Step int `json:"step" comment:"步骤"`
    Text string `json:"text" comment:"内容"`
    Status int `json:"status" comment:"状态 0-未审核 1-已审核"`
    TextAdditional string `json:"textAdditional" comment:"额外"`
    Remark string `json:"remark" comment:"审批备注"`
    IsAutoApprove int `json:"isAutoApprove" comment:"是否自动审批 0-否 1-是"`
    ApproveUser string `json:"approveUser" comment:"审批人"`
    RejectText []string `json:"rejectText" comment:"驳回商品"`
    Children  []WorkFlowNodesChildren `json:"children"`
}

type WorkFlowNodesChildren struct {
    //Step int `json:"step" comment:"步骤"`
    Text string `json:"text" comment:"内容"`
    ApprovedAt time.Time `json:"approvedAt" comment:"审核时间"`
    Status int `json:"status" comment:"状态 0-待审核 1-已审核 2-未操作"`
    TextAdditional string `json:"textAdditional" comment:"额外"`
    Remark string `json:"remark" comment:"审批备注"`
}

// GetWorkFlowNodes 获取审批单的审批流程
func (e *PresoLog) GetWorkFlowNodes(tx *gorm.DB, presoNoList []string) (resMap map[string][]WorkFlowNodes) {
    var presoList []Preso
    err := tx.Where("preso_no in ?", presoNoList).Preload("PresoDetails").Order("preso_no asc").Find(&presoList).Error
    if err != nil {
        return
    }
    // 审批单
    presoMap := make(map[string]Preso)
    for _, preso := range presoList {
        presoMap[preso.PresoNo] = preso
    }
    var presoLogList []PresoLog
    err = tx.Table(e.TableName()).Where("preso_no in ?", presoNoList).Order("preso_no asc, step asc").Find(&presoLogList).Error
    if err != nil {
        return
    }
    var userIds []int
    // 审批单logs
    presoLogMap := make(map[string][]PresoLog)
    for _, log := range presoLogList {
        presoLogMap[log.PresoNo] = append(presoLogMap[log.PresoNo], log)
        userIds = append(userIds, log.UserId)
    }

    // 用户名称map
    var modelUserInfo modelsUc.UserInfo
    users, _ := modelUserInfo.GetUsersByIds(tx, userIds)
    userMap := make(map[int]string)
    for _, user := range *users {
        userMap[user.Id] = user.UserName
    }

    resMap = make(map[string][]WorkFlowNodes)
    for presoNo, presoLogs := range presoLogMap {
        tmpMap := make(map[int][]PresoLog)
        for _, log := range presoLogs {
            tmpMap[log.Step] = append(tmpMap[log.Step], log)
        }

        step := 1
        for {
            logs, ok := tmpMap[step]
            if !ok {
                break
            }
            status := 0
            text := ""
            textAdditional := ""
            remark := ""
            isAutoApprove := 0
            approveUser := ""
            children := []WorkFlowNodesChildren{}
            if len(logs) == 1 {
                isAutoApprove = logs[0].IsAutoApprove
                if preso, ok := presoMap[presoNo]; ok && preso.Step >= step && logs[0].ApproveStatus != 0 && logs[0].ApproveStatus != 10 {
                    status = 1
                }

                userName := ""
                if _, ok := userMap[logs[0].UserId]; ok {
                    userName = userMap[logs[0].UserId]
                    approveUser = userName
                }
                if logs[0].LimitType == 1 {
                    textAdditional = "(订单金额>"+ strconv.FormatFloat(logs[0].TotalPriceLimit, 'f', 2, 64) +"元起批)"
                } else if logs[0].LimitType == 2 {
                    textAdditional = "(商品金额>"+ strconv.FormatFloat(logs[0].TotalPriceLimit, 'f', 2, 64) +"元起批)"
                }
                tmpTextAdditional := ""
                if logs[0].ApproveStatus == 1 || logs[0].ApproveStatus == -1 {
                    tmpTextAdditional = utils.TimeFormat(logs[0].UpdatedAt)
                }
                if logs[0].ApproveRankType == 1 { // 普签
                    text = userName
                    textAdditional = textAdditional + " " + tmpTextAdditional
                    remark = logs[0].ApproveRemark
                } else if logs[0].ApproveRankType == 3 { // 或签
                    text = "或签"
                    children = append(children, WorkFlowNodesChildren{
                        Text: userName,
                        ApprovedAt: logs[0].UpdatedAt,
                        Status: status,
                        TextAdditional: tmpTextAdditional,
                        Remark: logs[0].ApproveRemark,
                    })
                }
            } else if len(logs) > 1 {
                text = "或签"
                if logs[0].LimitType == 1 {
                    textAdditional = "(订单金额>"+ strconv.FormatFloat(logs[0].TotalPriceLimit, 'f', 2, 64) +"元起批)"
                } else if logs[0].LimitType == 2 {
                    textAdditional = "(商品金额>"+ strconv.FormatFloat(logs[0].TotalPriceLimit, 'f', 2, 64) +"元起批)"
                }
                for _, log := range logs {
                    if log.IsAutoApprove == 1 {
                        isAutoApprove = 1
                    }
                    childrenTextAdditional := ""
                    userName := ""
                    if _, ok := userMap[log.UserId]; ok {
                        userName = userMap[log.UserId]
                        approveUser = approveUser + "/" + userName
                    }
                    if log.ApproveStatus == 1 || log.ApproveStatus == -1 {
                        childrenTextAdditional = utils.TimeFormat(log.UpdatedAt)
                    }
                    childrenStatus := 0
                    preso, ok := presoMap[presoNo]
                    if ok && preso.Step >= step {
                        childrenStatus = 2
                        if log.ApproveStatus != 0 && log.ApproveStatus != 10 {
                            childrenStatus = 1
                            status = 1
                        }
                    }
                    children = append(children, WorkFlowNodesChildren{
                        Text: userName,
                        ApprovedAt: log.UpdatedAt,
                        Status: childrenStatus,
                        TextAdditional: childrenTextAdditional,
                        Remark: log.ApproveRemark,
                    })
                }
            }
            var rejectText []string
            for _, detail := range presoMap[presoNo].PresoDetails {
                if detail.ApproveStatus == -1 && detail.Step == step {
                    rejectText = append(rejectText, fmt.Sprintf("驳回：%s ￥%.2fx%v", detail.ProductName, detail.MarketPrice, detail.Quantity))
                }
            }

            resMap[presoNo] = append(resMap[presoNo], WorkFlowNodes{
                Step: step,
                Text: text,
                Status: status,
                TextAdditional: textAdditional,
                Remark: remark,
                IsAutoApprove: isAutoApprove,
                ApproveUser: strings.Trim(approveUser, "/"),
                RejectText: rejectText,
                Children: children,
            })
            step++
        }

    }

    return
}


// GetPresoLogsMap 获取以step为key的 审批记录map
func (e *PresoLog) GetPresoLogsMap(tx *gorm.DB, presoNo string) (presoLogsMap map[int][]PresoLog) {
    var presoLogs []PresoLog
    err := tx.Model(PresoLog{}).Where("preso_no = ?", presoNo).Order("step asc, id asc").Find(&presoLogs).Error
    if err != nil {
        return
    }
    presoLogsMap = make(map[int][]PresoLog)
    for _, log := range presoLogs {
        presoLogsMap[log.Step] = append(presoLogsMap[log.Step], log)
    }

    return
}
