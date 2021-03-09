package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 勋章
type Medal struct {
	ID        string    `json:"id" gorm:"primary_key" gqlschema:"query!;querys" description:"勋章id"`
	Title     string    `json:"title" gorm:"Type:varchar(255);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:"勋章名称"`
	CreatedAt time.Time `description:"创建时间" gqlschema:"querys"`
	UpdatedAt time.Time `description:"更新时间" gqlschema:"querys"`
	DeletedAt *time.Time
	v2        int    `gorm:"-" exclude:"true"`
	Confirm   bool   `gorm:"-" exclude:"true" gqlschema:"getmedallist!" description:"确认1是0否"`
	Userid    string `gorm:"-" exclude:"true" gqlschema:"checkcertrecord!" description:"用户ID"`
	Medalid   string `gorm:"-" exclude:"true" gqlschema:"checkcertrecord!" description:"勋章ID"`
}

type Medals struct {
	TotalCount int
	Edges      []Medal
}

type ReceiveMedal struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func UnescapeUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
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
			return o, errors.New("http.Post error")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return o, err
		}
		v, _ := UnescapeUnicode(body)
		id := gjson.Get(string(v), "id")
		msg := gjson.Get(string(v), "msg")
		data := gjson.Get(string(v), "data.list")
		if id.Int() != 0 {
			return o, errors.New(msg.String())
		}
		var receiveMedal []Medal
		if err := json.Unmarshal([]byte(data.String()), &receiveMedal); err != nil {
			panic(err)
		}
		for _, val := range receiveMedal {
			err = db.Create(&val).Error
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return o, nil
}

//检查勋章获取记录
func (o Medal) Checkcertrecord(params graphql.ResolveParams) (Medal, error) {
	p := params.Args
	resp, err := http.Post("https://www.zmlxj.com/app.php/Medal/ajax_get_check_cert_record",
		"application/x-www-form-urlencoded",
		strings.NewReader("id="+p["medalid"].(string)+"&token=b5afc7b7a1d16e58a0d1983154c58e4c&user_id="+p["userid"].(string)))
	if err != nil {
		return o, errors.New("http.Post error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return o, err
	}
	v, _ := UnescapeUnicode(body)
	msg := gjson.Get(string(v), "msg")
	data := gjson.Get(string(v), "data")
	fmt.Println(v)
	fmt.Println(msg)
	fmt.Println(data)
	return o, nil
}
