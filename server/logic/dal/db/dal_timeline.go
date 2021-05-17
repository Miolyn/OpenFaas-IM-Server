package db

import (
	"OpenFaaS-Logic/model"
	"OpenFaaS-Logic/pkg/logger"
	"fmt"

	"github.com/jinzhu/gorm"
)

type DalTimeLine struct {
	a AffairDB
}

func NewDalTimeLine() *DalTimeLine {

	return &DalTimeLine{}
}

func (d *DalTimeLine) InitAffair(tx *gorm.DB) {
	d.a.InitAffairsWithTx(tx)
}

func (d *DalTimeLine) GetDBHandler() *gorm.DB {
	if d.a.IsInit() {
		return d.a.tx
	}
	return DBConn.FindersIMDB
}

//device 作为唯一键 存在则修改conn id 不存在则创建
//同一个timeline对应的seq递增
func (d *DalTimeLine) Insert(m *model.Timeline) (timeline *model.Timeline, err error) {
	var (
		seq int64
	)
	db := d.GetDBHandler()
	db = db.Table(m.TableName())

	seq, err = d.GetTimelineSeq(m.TimelineID)
	if err != nil {
		logger.Logger.Errorf("GetTimelineSeq error:%v", err)
		return timeline, fmt.Errorf("GetTimelineSeq error:%v", err)
	}
	m.Seq = seq
	m.Status = model.TIMELINE_STATUS_NOT_READ
	logger.Logger.Debug("ready to insert ", *m)
	err = db.Create(m).Error
	if err != nil {
		logger.Logger.Errorf("Create error:%v", err)
		return nil, fmt.Errorf("Create error:%v", err)
	}
	return m, nil
}

func (d *DalTimeLine) GetTimelineSeq(timelineId int64) (int64, error) {
	db := d.GetDBHandler()
	m := new(model.Timeline)
	db = db.Table(m.TableName())
	err := db.Where("timeline_id=?", timelineId).Find(m).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("GetTimelineSeq error :%v", err)
	}
	return m.Seq + 1, nil
}

func (d *DalTimeLine) AckMessage(ackId int64) (err error) {
	db := d.GetDBHandler()
	m := new(model.Timeline)
	db = db.Table(m.TableName())
	err = db.Where("id = ?", ackId).Update("status", model.TIMELINE_STATUS_READ).Error
	return
}

func (d *DalTimeLine) GetUnReadTimelineTotalNum(timelineID int64) (cnt int64) {
	db := d.GetDBHandler().Debug()
	m := &model.Timeline{}
	db = db.Table(m.TableName())
	db.Where("timeline_id = ? and status = ?", timelineID, model.TIMELINE_STATUS_NOT_READ).Count(&cnt)
	return
}

func (d *DalTimeLine) GetUnReadTimeline(timelineID, page, pageSize int64) (timelines []*model.Timeline, err error) {
	db := d.GetDBHandler()
	m := &model.Timeline{}
	db = db.Table(m.TableName())
	pageNum := (page - 1) * pageSize
	err = db.Where("timeline_id = ? and status = ?", timelineID, model.TIMELINE_STATUS_NOT_READ).
		Offset(pageNum).Limit(pageSize).
		Find(&timelines).Error
	if gorm.IsRecordNotFoundError(err) {
		err = nil
		return
	}
	return
}

func (d *DalTimeLine) GetTimelineBySeqTotalNum(timelineID, seq int64) (cnt int64) {
	db := d.GetDBHandler().Debug()
	m := &model.Timeline{}
	db = db.Table(m.TableName())
	db.Where("timeline_id = ? and seq >= ?", timelineID, seq).Count(&cnt)
	return
}

func (d *DalTimeLine) GetTimelineBySeqWithPaginator(timelineId, seq, page, pageSize int64) (timelines []*model.Timeline, err error) {
	db := d.GetDBHandler()
	var timeline model.Timeline
	db = db.Table(timeline.TableName())
	pageNum := (page - 1) * pageSize
	err = db.Where("timeline_id = ? and seq >= ?", timelineId, seq).
		Offset(pageNum).Limit(pageSize).
		Find(&timelines).Error
	if gorm.IsRecordNotFoundError(err) {
		err = nil
	}
	return
}
