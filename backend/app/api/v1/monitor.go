package v1

import (
	"time"

	"github.com/1Panel-dev/1Panel/backend/app/api/v1/helper"
	"github.com/1Panel-dev/1Panel/backend/app/dto"
	"github.com/1Panel-dev/1Panel/backend/app/model"
	"github.com/1Panel-dev/1Panel/backend/constant"
	"github.com/1Panel-dev/1Panel/backend/global"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/net"
)

func (b *BaseApi) LoadMonitor(c *gin.Context) {
	var req dto.MonitorSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	if err := global.VALID.Struct(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	req.StartTime = req.StartTime.Add(8 * time.Hour)
	req.EndTime = req.EndTime.Add(8 * time.Hour)

	var backdatas []dto.MonitorData
	if req.Param == "all" || req.Param == "cpu" || req.Param == "memory" || req.Param == "load" {
		var bases []model.MonitorBase
		if err := global.DB.
			Where("created_at > ? AND created_at < ?", req.StartTime, req.EndTime).
			Find(&bases).Error; err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}

		var itemData dto.MonitorData
		itemData.Param = "base"
		for _, base := range bases {
			itemData.Date = append(itemData.Date, base.CreatedAt)
			itemData.Value = append(itemData.Value, base)
		}
		backdatas = append(backdatas, itemData)
	}
	if req.Param == "all" || req.Param == "io" {
		var bases []model.MonitorIO
		if err := global.DB.
			Where("created_at > ? AND created_at < ?", req.StartTime, req.EndTime).
			Find(&bases).Error; err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}

		var itemData dto.MonitorData
		itemData.Param = "io"
		for _, base := range bases {
			itemData.Date = append(itemData.Date, base.CreatedAt)
			itemData.Value = append(itemData.Value, base)
		}
		backdatas = append(backdatas, itemData)
	}
	if req.Param == "all" || req.Param == "network" {
		var bases []model.MonitorNetwork
		if err := global.DB.
			Where("name = ? AND created_at > ? AND created_at < ?", req.Info, req.StartTime, req.EndTime).
			Find(&bases).Error; err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}

		var itemData dto.MonitorData
		itemData.Param = "network"
		for _, base := range bases {
			itemData.Date = append(itemData.Date, base.CreatedAt)
			itemData.Value = append(itemData.Value, base)
		}
		backdatas = append(backdatas, itemData)
	}
	helper.SuccessWithData(c, backdatas)
}

func (b *BaseApi) GetNetworkOptions(c *gin.Context) {
	netStat, _ := net.IOCounters(true)
	var options []string
	options = append(options, "all")
	for _, net := range netStat {
		options = append(options, net.Name)
	}
	helper.SuccessWithData(c, options)
}

func (b *BaseApi) GetIOOptions(c *gin.Context) {
	diskStat, _ := disk.IOCounters()
	var options []string
	options = append(options, "all")
	for _, net := range diskStat {
		options = append(options, net.Name)
	}
	helper.SuccessWithData(c, options)
}
