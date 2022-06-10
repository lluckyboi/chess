package dao

import "MyChess/server/model"

func SelectUserByUserName(username string)(model.User,error){
	User := model.User{}
	sqlstr := "select id,username,User_mail from user where username=?"
	//单行查询
	errs := Db.QueryRow(sqlstr, username)

	//错误处理
	if errs.Err() != nil {
		return User, errs.Err()
	}
	err := errs.Scan(&User.Id, &User.UserName, &User.UserMail)
	if err != nil {
		return User, err
	}
	return User, nil
}

func InsertUser(user model.User) error {
	sqlstr := "insert into user(username, user_mail)values(?,?);"
	_, errs := Db.Exec(sqlstr, user.UserName,user.UserMail)
	if errs != nil {
		return errs
	}
	return nil
}

