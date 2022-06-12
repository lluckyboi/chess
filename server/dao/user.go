package dao

import "MyChess/server/model"

func SelectUserByUserName(username string) (model.User, error) {
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
	_, errs := Db.Exec(sqlstr, user.UserName, user.UserMail)
	if errs != nil {
		return errs
	}
	return nil
}

func UpdateWinCount(userid int) error {
	wt := 0
	sqlstr := "select win_count from user where id=?"
	err := Db.QueryRow(sqlstr, userid).Scan(&wt)

	sqlstr = "update user set win_count=? where id=?"
	_, err = Db.Exec(sqlstr, wt+1, userid)
	return err
}

func SelectWinCount(userid int) int {
	var wt int
	sqlstr := "select win_count from user where id=? "
	Db.QueryRow(sqlstr, userid).Scan(&wt)
	return wt
}
