package db

import (
	"OpenFaaS-Logic/model"
	"OpenFaaS-Logic/pkg/dto"
	"fmt"

	"github.com/jinzhu/gorm"
)

type DalTimelineMap struct {
}

func NewDalTimelineMap() *DalTimelineMap {

	return &DalTimelineMap{}
}

func (d *DalTimelineMap) GetDBHandler() *gorm.DB {
	return DBConn.FindersIMDB
}

// timeline类型，1.个人 2.单聊 3.群聊
// object_ids timeline对象id: type=1 时为user_id，type=2时为user_id,user_id，type=3时为group_id
func (d *DalTimelineMap) GetTimeLineMap(timeLineType int, args ...interface{}) (timelineMaps *model.TimelineMap, err error) {
	m := new(model.TimelineMap)
	db := d.GetDBHandler().Table(m.TableName()).Debug()
	whereStr := ""
	objectIds := ""
	if timeLineType == dto.TIMELINE_PERSON {
		if len(args) != 1 {
			return nil, fmt.Errorf("请传入 userID%s", "")
		}
		objectIds = args[0].(string)
		m.Type = dto.TIMELINE_PERSON
		whereStr = fmt.Sprintf("type=%d and object_ids='%v'", dto.TIMELINE_PERSON, objectIds)

	} else if timeLineType == dto.TIMELINE_DOUBLE {
		if len(args) != 2 {
			return nil, fmt.Errorf("请传入 双方userID%s", "")
		}
		objectIds1 := fmt.Sprintf("%v,%v", args[0], args[1])
		objectIds2 := fmt.Sprintf("%v,%v", args[1], args[0])
		objectIds = objectIds1
		m.Type = dto.TIMELINE_DOUBLE
		whereStr = fmt.Sprintf("type=%d and (object_ids='%v' or object_ids='%v')", dto.TIMELINE_DOUBLE, objectIds1, objectIds2)

	} else if timeLineType == dto.TIMELINE_GROUP {
		if len(args) != 1 {
			return nil, fmt.Errorf("请传入 group_id%s", "")
		}
		objectIds = fmt.Sprintf("%v", args[0])
		m.Type = dto.TIMELINE_GROUP
		whereStr = fmt.Sprintf("type=%d and object_ids='%v'", dto.TIMELINE_GROUP, objectIds)
	}
	m.ObjectIds = objectIds
	err = db.FirstOrCreate(m, whereStr).Error

	if err != nil {
		return nil, fmt.Errorf("数据库查询失败 ：%v", err)
	}

	return m, nil
}
