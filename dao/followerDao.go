package dao

import (
	"log"
	"sync"
)

// Follow 对应数据库follower表结构体
type Follow struct {
	id          int64
	user_id     int64 //被关注者的id
	follower_id int64 //粉丝id
	cancel      int
}

var (
	followDao  *FollowerDao //操作该dao层crud的结构体变量。
	followOnce sync.Once    //单例限定，去限定申请一个followDao结构体变量。
)

type FollowerDao struct {
}

// NewFollowDaoInstance 生成并返回followDao的单例对象。
func NewFollowDaoInstance() *FollowerDao {
	followOnce.Do(
		func() {
			followDao = &FollowerDao{}
		})
	return followDao
}

// 获取粉丝
func (*FollowerDao) GetFollowerCnt(userId int64) ([]int64, error) {
	var ids []int64
	if err := Db.Debug().Model(Follow{}).
		Where("user_id = ?", userId).
		Pluck("follower_id", &ids).Error; err != nil {
		//没有关注任何人
		if "record not found" == err.Error() {
			return nil, nil
		}
		// 查询出错。
		log.Println(err.Error())
		return nil, err
	}
	//查询成功
	return ids, nil
}

// 获取关注
func (*FollowerDao) GetFollowingCnt(userId int64) ([]int64, error) {
	var ids []int64
	if err := Db.Debug().Model(Follow{}).
		Where("follower_id = ?", userId).
		Pluck("user_id", &ids).Error; err != nil {
		//没有关注任何人
		if "record not found" == err.Error() {
			return nil, nil
		}
		// 查询出错。
		log.Println(err.Error())
		return nil, err
	}
	//查询成功
	return ids, nil
}
