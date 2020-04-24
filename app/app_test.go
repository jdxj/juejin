package app

import (
	"fmt"
	"testing"

	"github.com/jdxj/juejin/module"

	"github.com/jdxj/juejin/db"
)

func TestCollector_UserInfo(t *testing.T) {
	coll := new(Collector)
	userInfo, err := coll.UserInfo("")
	if err != nil {
		t.Fatalf("%s", err)
	}
	fmt.Printf("%#v", *userInfo)
}

func TestCollector_GetIDFromDB(t *testing.T) {
	coll := new(Collector)
	err := coll.GetIDFromDB()
	if err != nil {
		t.Fatalf("%s", err)
	}

	db.MySQL.Close()
}

func TestCollector_InsertUserInfo(t *testing.T) {
	userInfo := &module.UserInfo{
		ObjectId: "5dedf940f265da33d361f96d",
		Username: "fadsfadsfa",
	}

	coll := new(Collector)
	err := coll.InsertUserInfo([]*module.UserInfo{userInfo})
	if err != nil {
		t.Fatalf("%s", err)
	}

	db.MySQL.Close()
}
