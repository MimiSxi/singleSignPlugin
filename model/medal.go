package model

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 勋章
type Medal struct {
	ID        uint      `gorm:"primary_key" gqlschema:"query!;querys;checkcertrecord!" description:"勋章id"`
	Title     string    `gorm:"Type:varchar(255);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:"勋章名称"`
	CreatedAt time.Time `description:"创建时间" gqlschema:"querys"`
	UpdatedAt time.Time `description:"更新时间" gqlschema:"querys"`
	DeletedAt *time.Time
	v2        int    `gorm:"-" exclude:"true"`
	Confirm   bool   `gorm:"-" exclude:"true" gqlschema:"getmedallist!" description:"确认1是0否"`
	Userid    string `gorm:"-" exclude:"true" gqlschema:"checkcertrecord!" description:"用户ID"`
}

type Medals struct {
	TotalCount int
	Edges      []Medal
}

func (o Medal) Query(params graphql.ResolveParams) (Medal, error) {
	p := params.Args
	err := db.Where(p).First(&o).Error
	return o, err
}

func (o Medal) Querys(params graphql.ResolveParams) (Medals, error) {
	var result Medals

	dbselect := GenSelet(db, params)
	dbcount := GenWhere(db.Model(o), params)

	err := dbselect.Find(&result.Edges).Error
	if err != nil {
		return result, err
	}
	err = dbcount.Count(&result.TotalCount).Error
	return result, err
}

// 获取勋章列表
func (o Medal) Getmedallist(params graphql.ResolveParams) (Medal, error) {
	p := params.Args
	o.Confirm = p["confirm"].(bool)
	if o.Confirm != false {
		resp, err := http.Post("https://www.zmlxj.com/app.php/Medal/ajax_get_medal_list",
			"application/x-www-form-urlencoded",
			strings.NewReader("token=b5afc7b7a1d16e58a0d1983154c58e4c"))
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return o, err
		}
		fmt.Println("===================================")
		fmt.Println(string(body))
		//{"id":0,"msg":"\u6210\u529f","data":{"list":[{"id":"2515","title":"\u65b0\u589e\u52cb\u7ae0-\u4fee\u6539"},
		//{"id":"269","title":"\u4e2d\u56fd\u5730\u8d28\u5927\u5b66\u5316\u77f3\u6797\u52cb\u7ae0"},
		//{"id":"239","title":"\u6700\u65c5\u884c\u5bb6\u6b22\u8fce\u52cb\u7ae0"}],"time_":"0.0035s","memory_":"96kb"}}
		fmt.Println("===================================")

	}
	return o, nil
}

//检查勋章获取记录
func (o Medal) Checkcertrecord(params graphql.ResolveParams) (Medal, error) {
	p := params.Args
	o.Userid = p["userid"].(string)
	resp, err := http.Post("https://www.zmlxj.com/app.php/Medal/ajax_get_check_cert_record",
		"application/x-www-form-urlencoded",
		strings.NewReader("id="+p["id"].(string)+"&token=b5afc7b7a1d16e58a0d1983154c58e4c&user_id="+o.Userid))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return o, err
	}
	fmt.Println("===================================")
	fmt.Println(string(body))
	//{"id":0,"msg":"\u6210\u529f","data":{"list":[{"id":"2515","title":"\u65b0\u589e\u52cb\u7ae0-\u4fee\u6539"},
	//{"id":"269","title":"\u4e2d\u56fd\u5730\u8d28\u5927\u5b66\u5316\u77f3\u6797\u52cb\u7ae0"},
	//{"id":"239","title":"\u6700\u65c5\u884c\u5bb6\u6b22\u8fce\u52cb\u7ae0"}],"time_":"0.0035s","memory_":"96kb"}}
	fmt.Println("===================================")

	return o, nil
}
