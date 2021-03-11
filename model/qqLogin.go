package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	AppId       = "101827468"
	AppKey      = "0d2d856e48e0ebf6b98e0d0c879fe74d"
	redirectURI = "http://127.0.0.1:9090/qqLogin"
)

type QQLoginInfo struct {
	ID           uint      `gorm:"primary_key" gqlschema:"query!;querys" description:"id"`
	AccessToken  string    `json:"access_token" gorm:"Type:varchar(255);DEFAULT:'';NOT NULL;" description:""`
	ExpiresIn    string    `json:"expires_in" gorm:"Type:varchar(255);DEFAULT:'';NOT NULL;" description:""`
	RefreshToken string    `json:"refresh_token" gorm:"Type:varchar(255);DEFAULT:'';NOT NULL;" description:""`
	OpenId       string    `json:"openid" gorm:"Type:varchar(255);DEFAULT:'';NOT NULL;" description:""`
	Ret          int       `json:"ret" gorm:"-" description:"登录返回码"`
	Msg          string    `json:"msg" gorm:"-" description:"返回信息"`
	Nickname     string    `json:"nickname" gorm:"Type:varchar(255);DEFAULT:'';" description:"昵称"`
	Figureurl    string    `json:"figureurl" gorm:"Type:varchar(255);DEFAULT:'';" description:"30×30像素的QQ空间头像URL"`
	Figureurl1   string    `json:"figureurl_1" gorm:"Type:varchar(255);DEFAULT:'';" description:"50×50像素的QQ空间头像URL"`
	Figureurl2   string    `json:"figureurl_2" gorm:"Type:varchar(255);DEFAULT:'';" description:"100×100像素的QQ空间头像URL"`
	FigureurlQq1 string    `json:"figureurl_qq_1" gorm:"Type:varchar(255);DEFAULT:'';" description:"40×40像素的QQ头像URL"`
	FigureurlQq2 string    `json:"figureurl_qq_2" gorm:"Type:varchar(255);DEFAULT:'';" description:"100×100像素的QQ头像URL"`
	Gender       string    `json:"gender" gorm:"Type:varchar(10);DEFAULT:'';" description:"性别,获取不到则默认返回'男'"`
	CreatedAt    time.Time `description:"创建时间" gqlschema:"querys"`
	UpdatedAt    time.Time `description:"更新时间" gqlschema:"querys"`
	DeletedAt    *time.Time
	v2           int    `gorm:"-" exclude:"true"`
	Confirm      bool   `gorm:"-" exclude:"true" gqlschema:"requestlogin!" description:"确认1是0否"`
	Code         string `gorm:"-" exclude:"true" gqlschema:"gettoken!" description:"code"`
}

type QQLoginInfos struct {
	TotalCount int
	Edges      []QQLoginInfo
}

func (o QQLoginInfo) Query(params graphql.ResolveParams) (QQLoginInfo, error) {
	p := params.Args
	err := db.Where(p).First(&o).Error
	return o, err
}

func (o QQLoginInfo) Querys(params graphql.ResolveParams) (QQLoginInfos, error) {
	var result QQLoginInfos

	dbselect := GenSelet(db, params)
	dbcount := GenWhere(db.Model(o), params)

	err := dbselect.Find(&result.Edges).Error
	if err != nil {
		return result, err
	}
	err = dbcount.Count(&result.TotalCount).Error
	return result, err
}

func (o QQLoginInfo) Requestlogin(params graphql.ResolveParams) (QQLoginInfo, error) {
	p := params.Args
	o.Confirm = p["confirm"].(bool)
	if o.Confirm != false {
		params := url.Values{}
		params.Add("response_type", "code")
		params.Add("client_id", AppId)
		params.Add("state", "test")
		str := fmt.Sprintf("%s&redirect_uri=%s", params.Encode(), redirectURI)
		loginURL := fmt.Sprintf("%s?%s", "https://graph.qq.com/oauth2.0/authorize", str)
		fmt.Println(loginURL)
		return o, errors.New(loginURL)
	}
	return o, nil
}

func (o QQLoginInfo) Gettoken(params graphql.ResolveParams) (QQLoginInfo, error) {
	p := params.Args
	o.Code = p["code"].(string)
	{
		params := url.Values{}
		params.Add("grant_type", "authorization_code")
		params.Add("client_id", AppId)
		params.Add("client_secret", AppKey)
		params.Add("code", o.Code)
		str := fmt.Sprintf("%s&redirect_uri=%s", params.Encode(), redirectURI)
		loginURL := fmt.Sprintf("%s?%s", "https://graph.qq.com/oauth2.0/token", str)
		response, err := http.Get(loginURL)
		if err != nil {
			return o, err
		}
		defer response.Body.Close()
		bs, _ := ioutil.ReadAll(response.Body)
		body := string(bs)
		resultMap := convertToMap(body)
		o.AccessToken = resultMap["access_token"]
		o.RefreshToken = resultMap["refresh_token"]
		o.ExpiresIn = resultMap["expires_in"]
		resp, err := http.Get(fmt.Sprintf("%s?access_token=%s", "https://graph.qq.com/oauth2.0/me", o.AccessToken))
		if err != nil {
			return o, err
		}
		defer resp.Body.Close()
		bs, _ = ioutil.ReadAll(resp.Body)
		body = string(bs)
		o.OpenId = body[45:77]
	}
	{
		params := url.Values{}
		params.Add("access_token", o.AccessToken)
		params.Add("openid", o.OpenId)
		params.Add("oauth_consumer_key", AppId)

		uri := fmt.Sprintf("https://graph.qq.com/user/get_user_info?%s", params.Encode())
		resp, err := http.Get(uri)
		if err != nil {
			return o, err
		}
		defer resp.Body.Close()
		bs, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(bs, &o); err != nil {
			panic(err)
		}
	}
	if err := db.Where("open_id = ?", o.OpenId).First(&QQLoginInfo{}).Updates(&o).Error; err != nil {
		err = db.Create(&o).Error
		if err != nil {
			fmt.Println(err)
		}
	}
	return o, nil
}

func convertToMap(str string) map[string]string {
	var resultMap = make(map[string]string)
	values := strings.Split(str, "&")
	for _, value := range values {
		vs := strings.Split(value, "=")
		resultMap[vs[0]] = vs[1]
	}
	return resultMap
}
