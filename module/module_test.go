package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/jdxj/juejin/client"
)

func TestFollowInfo(t *testing.T) {
	respData := new(RespData)

	resp, err := client.Get("https://follow-api-ms.juejin.im/v1/getUserFolloweeList?uid=5891f7048d6d81006c412f30&currentUid=5891f7048d6d81006c412f30&src=web")
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(respData)
	if err != nil {
		t.Fatalf("%s", err)
	}

	followInfo := make([]*FollowInfo, 0)

	reader := bytes.NewReader(respData.Data)
	decoder = json.NewDecoder(reader)
	err = decoder.Decode(&followInfo)
	if err != nil {
		t.Fatalf("%s", err)
	}
	for _, v := range followInfo {
		fmt.Printf("followee: %v, followTime: %s\n", *v.Followee, v.FollowDatetime)
	}
}

func TestUserInfo(t *testing.T) {
	respData := &RespData{}
	resp, err := client.Get("https://lccro-api-ms.juejin.im/v1/get_multi_user?uid=5891f7048d6d81006c412f30&device_id=1587714064961&token=eyJhY2Nlc3NfdG9rZW4iOiJ5SEVWRWZ4bFI1SktLNmcwIiwicmVmcmVzaF90b2tlbiI6ImJCWnZUZUYzdVJKeHhGR2wiLCJ0b2tlbl90eXBlIjoibWFjIiwiZXhwaXJlX2luIjoyNTkyMDAwfQ%3D%3D&src=web&ids=58d7e6e2a22b9d00646882b5&cols=viewedEntriesCount%7Crole%7CtotalCollectionsCount%7CallowNotification%7CsubscribedTagsCount%7CappliedEditorAt%7Cemail%7CfollowersCount%7CpostedEntriesCount%7ClatestCollectionUserNotification%7CcommentedEntriesCount%7CweeklyEmail%7CcollectedEntriesCount%7CpostedPostsCount%7Cusername%7ClatestLoginedInAt%7CtotalHotIndex%7CblogAddress%7CselfDescription%7ClatestCheckedNotificationAt%7CemailVerified%7CtotalCommentsCount%7Cinstallation%7Cblacklist%7CweiboId%7CmobilePhoneNumber%7Capply%7CfolloweesCount%7CdeviceType%7CeditorType%7CjobTitle%7Ccompany%7ClatestVoteLikeUserNotification%7CauthData%7CavatarLarge%7CmobilePhoneVerified%7CobjectId%7CcreatedAt%7CupdatedAt")
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(respData)
	if err != nil {
		t.Fatalf("%s", err)
	}
	fmt.Printf("data: %s\n", respData.Data)

	reader := bytes.NewReader(respData.Data)
	decoder = json.NewDecoder(reader)

	user := make(map[string]*UserInfo)
	err = decoder.Decode(&user)
	if err != nil {
		t.Fatalf("%s", err)
	}
	for k, v := range user {
		fmt.Printf("id: %s, userInfo: %v\n", k, *v)
	}
}

func TestURLESC(t *testing.T) {
	//str := "https://follow-api-ms.juejin.im/v1/getUserFolloweeList?uid=5891f7048d6d81006c412f30&currentUid=5891f7048d6d81006c412f30&src=web"
	res := url.QueryEscape("2019-08-26T08:41:19.239Z")
	fmt.Printf("%s\n", res)
}
