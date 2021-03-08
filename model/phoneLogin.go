package model

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type UserInfo struct {
	ID           uint      `gorm:"primary_key" gqlschema:";query!;querys" description:"id"`
	GeoShow      uint      `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	Etime        uint      `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	UserId       string    `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"query" description:"用户id"`
	IsVip        uint      `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:"是否会员"`
	UserNick     string    `gorm:"Type:varchar(255);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:"用户昵称"`
	UserIcon     string    `gorm:"Type:varchar(255);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:"用户头像"`
	UserIel      string    `gorm:"Type:varchar(128);DEFAULT:'';NOT NULL;" gqlschema:"querys;sendletter!;register!;login!" description:"联系电话"`
	ParkId       uint      `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	VolunteerSex uint      `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	RedPacket    string    `gorm:"Type:varchar(1000);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:"红包"`
	Pid          uint      `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	IsShowMoney  uint      `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	Total        uint      `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	Wxacode      string    `gorm:"Type:varchar(128);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:""`
	IsWork       uint      `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	Stime        uint      `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	Stype        string    `gorm:"Type:varchar(1000);DEFAULT:'';NOT NULL;" gqlschema:"querys" description:""`
	Ltime        uint      `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	LabelId      uint      `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"querys" description:""`
	CreatedAt    time.Time `description:"创建时间" gqlschema:"querys"`
	UpdatedAt    time.Time `description:"更新时间" gqlschema:"querys"`
	DeletedAt    *time.Time
	v2           int    `gorm:"-" exclude:"true"`
	Code         string `gorm:"-" exclude:"true" gqlschema:"register!;login!" description:"验证码"`
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
	fmt.Println(string(body)) //{"id":0,"msg":"\u6210\u529f"}

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
	fmt.Println(string(body))
	//{"id":0,"msg":"\u6210\u529f","data":
	//{"user":{"geo_show":"1","etime":"0","user_id":"74f74a6b6799106cf66796eb009bcbb8","is_vip":"0","user_nick":"\u624b\u673a\u7528\u62373915",
	//"user_icon":"http:\/\/www.zmlxj.com\/zmlxjshare.jpg","user_tel":"18862213915","park_id":"0","volunteer_sex":"0","red_packet":"0.00",
	//"pid":"0","is_show_money":"0","total_":"5","wxacode":"","is_work":"0","stime":"0","stype":"news","ltime":"0","label_id":"0","user_pwd":""}}}
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
	fmt.Println(string(body))
	//{"id":0,"msg":"\u6210\u529f","data":
	//{"user":{"geo_show":"1","etime":"0","user_id":"74f74a6b6799106cf66796eb009bcbb8","is_vip":"0","user_nick":"\u624b\u673a\u7528\u62373915",
	//"user_icon":"http:\/\/www.zmlxj.com\/zmlxjshare.jpg","user_tel":"18862213915","park_id":"0","volunteer_sex":"0","red_packet":"0.00",
	//"pid":"0","is_show_money":"0","total_":"5","wxacode":"","is_work":"0","stime":"0","stype":"news","ltime":"0","label_id":"0","user_pwd":""}}}
	return o, nil
}
