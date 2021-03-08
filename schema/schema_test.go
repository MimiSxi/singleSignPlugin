package schema

func initDB() {
	// if err := model.OpenDatabase("sqlite3", ":memory:", 10, 100, 1000); err != nil {
	// 	panic("failed to connect database")
	// }
	// model.Run(func(db *gorm.DB) {
	// 	password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	// 	account := model.Account{
	// 		UserName:  "cadmin",
	// 		PassWord:  string(password),
	// 		ReferID:   1,
	// 		ReferType: "refertype",
	// 	}
	// 	db.Save(&account)
	// })
}
