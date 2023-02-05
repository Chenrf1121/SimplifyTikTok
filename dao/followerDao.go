package dao

import "log"

// TableFollower 对应数据库follower表结构体
type TableFollower struct {
	id          int64
	user_id     int64
	follower_id int64
	cancel      int
}

func GetFollowerCnt(userId int64) (int64, error) {
	followers := []TableFollower{}
	log.Println("Call GetFollowerCnt")
	if err := Db.Where("user_id = ?", userId).Find(&followers); err != nil {
		log.Println(err.Error)
		return 0, err.Error
	}

	return int64(len(followers)), nil
}
