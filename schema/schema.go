package schema

import (
	"github.com/Fiber-Man/funplugin"
	"github.com/Fiber-Man/singleSignPlugin/model"
	"github.com/graphql-go/graphql"
)

var phoneLoginSchema *funplugin.ObjectSchema
var medalSchema *funplugin.ObjectSchema
var qqloginSchema *funplugin.ObjectSchema

var load = false

func Init() {
	// InitAccount()
}

func marge(oc *funplugin.ObjectSchema) {
	for k, v := range oc.Query {
		queryFields[k] = v
	}
	for k, v := range oc.Mutation {
		mutationFields[k] = v
	}
}

var queryFields = graphql.Fields{
	// "account":  &queryAccount,
	// "accounts": &queryAccountList,
	// "authority":  &queryAuthority,
	// "authoritys": &queryAuthorityList,
}

var mutationFields = graphql.Fields{
	// "createAccount": &createAccount,
	// "updateAccount": &updateAccount,
}

// NewSchema 用于插件主程序调用
func NewPlugSchema(pls funplugin.PluginManger) funplugin.Schema {
	if load != true {

		phoneLoginSchema, _ = pls.NewSchemaBuilder(model.UserInfo{})
		marge(phoneLoginSchema)

		medalSchema, _ = pls.NewSchemaBuilder(model.Medal{})
		marge(medalSchema)

		qqloginSchema, _ = pls.NewSchemaBuilder(model.QQLoginInfo{})
		marge(qqloginSchema)

		load = true
	}

	// roleSchema, _ := pls.NewSchemaBuilder(model.Role{})
	// marge(roleSchema)

	// roleAccountSchema, _ := pls.NewSchemaBuilder(model.RoleAccount{})
	// marge(roleAccountSchema)

	return funplugin.Schema{
		Object: map[string]*graphql.Object{
			// "account": accountType,

			"phoneLogin": phoneLoginSchema.GraphQLType,
			"medal":      medalSchema.GraphQLType,
			"qqlogin":    qqloginSchema.GraphQLType,

			// "role":        roleSchema.GraphQLType,
			// "roleaccount": roleAccountSchema.GraphQLType,
		},
		Query:    queryFields,
		Mutation: mutationFields,
	}
}
