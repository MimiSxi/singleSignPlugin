package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type UserInfo struct {
	ID           uint      `gorm:"primary_key" gqlschema:";query!;querys" description:"id"`
	GeoShow      string    `json:"geo_show" gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	Etime        string    `json:"etime" gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	UserId       string    `json:"user_id" gorm:"primary_key;" gqlschema:"query" description:"用户id"`
	IsVip        string    `json:"is_vip" gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:"是否会员"`
	UserNick     string    `json:"user_nick" gorm:"Type:varchar(255);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:"用户昵称"`
	UserIcon     string    `json:"user_icon" gorm:"Type:varchar(255);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:"用户头像"`
	UserIel      string    `json:"user_tel" gorm:"Type:varchar(128);DEFAULT:'';NOT NULL;" gqlschema:"querys;sendletter!;register!;login!" description:"联系电话"`
	ParkId       string    `json:"park_id" gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	VolunteerSex string    `json:"volunteer_sex" gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	RedPacket    string    `json:"red_packet" gorm:"Type:varchar(1000);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:"红包"`
	Pid          string    `json:"pid" gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	IsShowMoney  string    `json:"is_show_money" gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	Total        string    `json:"total_" gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	Wxacode      string    `json:"wxacode" gorm:"Type:varchar(128);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:""`
	IsWork       string    `json:"is_work" gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	Stime        string    `json:"stime" gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	Stype        string    `json:"stype" gorm:"Type:varchar(1000);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:""`
	Ltime        string    `json:"ltime" gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	LabelId      string    `json:"label_id" gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	UserPwd      string    `json:"user_pwd" gorm:"Type:varchar(1280);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:""`
	CreatedAt    time.Time `description:"创建时间" gqlschema:"querys"`
	UpdatedAt    time.Time `description:"更新时间" gqlschema:"querys"`
	DeletedAt    *time.Time
	v2           int    `gorm:"-" exclude:"true"`
	Code         string `gorm:"-" exclude:"true" gqlschema:"register!;login!" description:"验证码"`
}

type User struct {
	GeoShow      string `json:"geo_show"`
	Etime        string `json:"etime"`
	UserID       string `json:"user_id"`
	IsVip        string `json:"is_vip"`
	UserNick     string `json:"user_nick"`
	UserIcon     string `json:"user_icon"`
	UserTel      string `json:"user_tel"`
	ParkID       string `json:"park_id"`
	VolunteerSex string `json:"volunteer_sex"`
	RedPacket    string `json:"red_packet"`
	Pid          string `json:"pid"`
	IsShowMoney  string `json:"is_show_money"`
	Total        string `json:"total_"`
	Wxacode      string `json:"wxacode"`
	IsWork       string `json:"is_work"`
	Stime        string `json:"stime"`
	Stype        string `json:"stype"`
	Ltime        string `json:"ltime"`
	LabelID      string `json:"label_id"`
	UserPwd      string `json:"user_pwd"`
}

type UserInfos struct {
	TotalCount int
	Edges      []UserInfo
}

func (o UserInfo) Query(params graphql.ResolveParams) (UserInfo, error) {
	p := params.Args
	err := db.Where(p).First(&o).Error
	return o, err
}

func (o UserInfo) Querys(params graphql.ResolveParams) (UserInfos, error) {
	var result UserInfos

	dbselect := GenSelet(db, params)
	dbcount := GenWhere(db.Model(o), params)

	err := dbselect.Find(&result.Edges).Error
	if err != nil {
		return result, err
	}
	err = dbcount.Count(&result.TotalCount).Error
	return result, err
}

func (o UserInfo) Sendletter(params graphql.ResolveParams) (UserInfo, error) {
	p := params.Args
	o.UserIel = p["userIel"].(string)
	resp, err := http.Post("https://www.zmlxj.com/app.php/Login/send_letter",
		"application/x-www-form-urlencoded",
		strings.NewReader("tel_="+o.UserIel+"&token=b5afc7b7a1d16e58a0d1983154c58e4c&country=86"))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return o, err
	}
	v, _ := UnescapeUnicode(body)
	id := gjson.Get(string(v), "id")
	msg := gjson.Get(string(v), "msg")
	if id.Int() != 0 {
		return o, errors.New(msg.String())
	}
	return o, nil
}

func (o UserInfo) Register(params graphql.ResolveParams) (UserInfo, error) {
	p := params.Args
	o.UserIel = p["userIel"].(string)
	o.Code = p["code"].(string)
	resp, err := http.Post("https://www.zmlxj.com/app.php/Login/register",
		"application/x-www-form-urlencoded",
		strings.NewReader("tel_="+o.UserIel+"&token=b5afc7b7a1d16e58a0d1983154c58e4c&code="+o.Code))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return o, err
	}
	v, _ := UnescapeUnicode(body)
	id := gjson.Get(string(v), "id")
	msg := gjson.Get(string(v), "msg")
	if id.Int() != 0 {
		return o, errors.New(msg.String())
	}
	data := gjson.Get(string(v), "data.user")
	fmt.Println(data)
	if err := json.Unmarshal([]byte(data.String()), &o); err != nil {
		panic(err)
	}
	// todo 重复入库
	err = db.Create(&o).Error
	if err != nil {
		fmt.Println(err)
	}
	return o, nil
}

func (o UserInfo) Login(params graphql.ResolveParams) (UserInfo, error) {
	p := params.Args
	o.UserIel = p["userIel"].(string)
	o.Code = p["code"].(string)
	resp, err := http.Post("https://www.zmlxj.com/app.php/Login/login",
		"application/x-www-form-urlencoded",
		strings.NewReader("tel_="+o.UserIel+"&token=b5afc7b7a1d16e58a0d1983154c58e4c&code="+o.Code))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return o, err
	}
	v, _ := UnescapeUnicode(body)
	id := gjson.Get(string(v), "id")
	msg := gjson.Get(string(v), "msg")
	if id.Int() != 0 {
		return o, errors.New(msg.String())
	}
	data := gjson.Get(string(v), "data.user")
	fmt.Println(data)
	if err := json.Unmarshal([]byte(data.String()), &o); err != nil {
		panic(err)
	}
	err = db.Create(&o).Error
	if err != nil {
		fmt.Println(err)
	}
	return o, nil
}
