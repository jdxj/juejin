package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/jdxj/juejin/db"

	"github.com/astaxie/beego/logs"
	"github.com/jdxj/juejin/client"
	"github.com/jdxj/juejin/module"
)

func NewCollector() *Collector {
	coll := new(Collector)
	coll.stop = make(chan int)
	coll.wg = &sync.WaitGroup{}
	return coll
}

type Collector struct {
	stop chan int
	wg   *sync.WaitGroup
}

func (coll *Collector) Start() {
	wg := coll.wg
	wg.Add(1)

	go func() {
		err := coll.GetIDFromDB()
		if err != nil {
			logs.Error("%s", err)
		}
		wg.Done()
	}()
}

func (coll *Collector) Stop() {
	close(coll.stop)
	coll.wg.Wait()
	logs.Info("collector stopped")
}

func (coll *Collector) GetIDFromDB() error {
	mysql := db.MySQL

	for offset, amount := 0, 100; ; {
		select {
		case <-coll.stop:
			logs.Info("stop GetIDFromDB")
			return nil
		default:
		}

		var usersID []string
		query := fmt.Sprintf("SELECT id FROM user limit %d,%d", offset, amount)
		rows, err := mysql.Query(query)
		if err != nil {
			return err
		}

		var count int
		for id := ""; rows.Next(); count++ {
			if err := rows.Scan(&id); err != nil {
				return err
			}
			usersID = append(usersID, id)
		}
		rows.Close()

		// 没数据了
		if count <= 0 {
			break
		}
		offset += count

		for _, v := range usersID {
			err := coll.CollectFollowee(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (coll *Collector) CollectFollowee(target string) error {
	// 初始路径
	path := fmt.Sprintf("https://follow-api-ms.juejin.im/v1/getUserFolloweeList?uid=%s&currentUid=5891f7048d6d81006c412f30&src=web",
		target)
	followeesID := make([]string, 0)

	var infos []*module.FollowInfo
	var err error

lab:
	for infos, err = coll.followee(path); len(infos) != 0 && err == nil; {
		select {
		case <-coll.stop:
			logs.Info("stop CollectFollowee, followeesIDCount: %d", len(followeesID))
			break lab
		default:
		}

		var before string
		for i, v := range infos {
			// 需要记录最后一个人的关注时间
			if i == len(infos)-1 {
				before = v.FollowDatetime
			}
			followeesID = append(followeesID, v.Followee.FolloweeId)
		}

		// 继续下次寻找
		queryPath := path + "&before=" + url.QueryEscape(before)
		infos, err = coll.followee(queryPath)
	}

	if err != nil {
		return fmt.Errorf("for: %s", err)
	}

	return coll.CollectUser(followeesID)
}

func (coll *Collector) followee(path string) ([]*module.FollowInfo, error) {
	resp, err := client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("client: %s", err)
	}
	defer resp.Body.Close()

	// 解析通用数据包装
	decoder := json.NewDecoder(resp.Body)
	respData := &module.RespData{}
	err = decoder.Decode(respData)
	if err != nil {
		return nil, fmt.Errorf("decode respData: %s", err)
	}

	// 解析关注数据
	followInfo := make([]*module.FollowInfo, 0)
	reader := bytes.NewReader(respData.Data)
	decoder = json.NewDecoder(reader)
	err = decoder.Decode(&followInfo)
	if err != nil {
		return nil, fmt.Errorf("decode folowInfo: %s", err)
	}
	return followInfo, nil
}

func (coll *Collector) CollectUser(target []string) error {
	var userInfos []*module.UserInfo

lab:
	for _, id := range target {
		select {
		case <-coll.stop:
			logs.Info("stop CollectUser, userInfosCount: %d", len(userInfos))
			break lab
		default:
		}

		info, err := coll.UserInfo(id)
		if err != nil {
			return err
		}
		userInfos = append(userInfos, info)
	}

	return coll.InsertUserInfo(userInfos)
}

func (coll *Collector) UserInfo(target string) (*module.UserInfo, error) {
	pathPrefix := "https://lccro-api-ms.juejin.im/v1/get_multi_user?"
	query := "uid=5891f7048d6d81006c412f30&device_id=1587714064961&token=eyJhY2Nlc3NfdG9rZW4iOiJ5SEVWRWZ4bFI1SktLNmcwIiwicmVmcmVzaF90b2tlbiI6ImJCWnZUZUYzdVJKeHhGR2wiLCJ0b2tlbl90eXBlIjoibWFjIiwiZXhwaXJlX2luIjoyNTkyMDAwfQ==&src=web&ids=%s&cols=viewedEntriesCount|role|totalCollectionsCount|allowNotification|subscribedTagsCount|appliedEditorAt|email|followersCount|postedEntriesCount|latestCollectionUserNotification|commentedEntriesCount|weeklyEmail|collectedEntriesCount|postedPostsCount|username|latestLoginedInAt|totalHotIndex|blogAddress|selfDescription|latestCheckedNotificationAt|emailVerified|totalCommentsCount|installation|blacklist|weiboId|mobilePhoneNumber|apply|followeesCount|deviceType|editorType|jobTitle|company|latestVoteLikeUserNotification|authData|avatarLarge|mobilePhoneVerified|objectId|createdAt|updatedAt"

	query = fmt.Sprintf(query, target)
	path := pathPrefix + query

	resp, err := client.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析通用数据
	decoder := json.NewDecoder(resp.Body)
	respData := &module.RespData{}
	err = decoder.Decode(respData)
	if err != nil {
		return nil, err
	}

	var userInfo *module.UserInfo
	reader := bytes.NewReader(respData.Data)
	decoder = json.NewDecoder(reader)
	infoMap := make(map[string]*module.UserInfo)

	err = decoder.Decode(&infoMap)
	if err != nil {
		return nil, err
	}

	// 只有一条数据
	for _, v := range infoMap {
		userInfo = v
		break
	}
	return userInfo, nil
}

func (coll *Collector) InsertUserInfo(userInfos []*module.UserInfo) error {
	mysql := db.MySQL
	query := "INSERT INTO user (id,name,followee_count,follower_count,rank_index,views,levels,power,collection) VALUE (?,?,?,?,?,?,?,?,?)"
	stmt, err := mysql.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, info := range userInfos {
		_, err := stmt.Exec(info.ObjectId, info.Username, info.FolloweesCount, info.FollowersCount, info.RankIndex,
			info.TotalViewsCount, info.Level, info.JuejinPower, info.TotalCollectionsCount)
		if err != nil {
			if strings.Index(err.Error(), "Duplicate entry") >= 0 {
				logs.Warn("%s", err)
				continue
			}
			return err
		}
	}
	return nil
}
