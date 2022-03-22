package sqldata

import (
	"database/sql"
	"errors"
	"time"
	"txt/pkg/config"
	"txt/pkg/structure"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var (
	Db   *sql.DB
	errs error
)

//处理用户登录
func EuserLogins(emailanduser string) error {
	var Email string
	var username string
	sqlstr := "select email from user where email=?"
	rows := Db.QueryRow(sqlstr, emailanduser)
	rows.Scan(&Email)
	if Email == emailanduser {
		return nil
	}
	rows = Db.QueryRow("select username from user where username=?", emailanduser)
	rows.Scan(&username)
	if username == emailanduser {
		return nil
	}
	return errors.New("用户名不存在！")
}

//展示我收藏的记事
func FindMyCollection(email string) ([]structure.Blogs, error) {
	var blog = make([]structure.Blogs, 0) //读取所有记事
	sqlstr := "select id from collection where people=?"
	rows, err := Db.Query(sqlstr, email)
	if err != nil {
		logrus.Error(err)
		return blog, err
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		sqlstr := "select title,praise,content from txt where id=?"
		rows, err := Db.Query(sqlstr, id)
		if err != nil {
			logrus.Error(err)
			return blog, err
		}
		for rows.Next() {
			var b structure.Blogs //读取所有记事
			b.Id = id
			rows.Scan(&b.Title, &b.Praise, &b.Content)
			blog = append(blog, b)
		}
	}
	return blog, nil
}

//搜索
func FindSearch(key string) ([]structure.Blogs, bool) {
	var Blogs = make([]structure.Blogs, 0) //读取所有记事
	sqlstr := "select id,title,praise,content,label from txt where title like ?;"
	rows, err := Db.Query(sqlstr, "%"+key+"%")
	if err != nil {
		logrus.Error(err)
		return Blogs, false
	}
	for rows.Next() {
		var b structure.Blogs //读取所有记事
		rows.Scan(&b.Id, &b.Title, &b.Praise, &b.Content, &b.Label)
		Blogs = append(Blogs, b)
	}
	return Blogs, true
}

//提交修改邮箱申请
func Subnewemail(email string, newemail string, reason string) bool {
	sqlstr := "select username from user where email=?"
	rows := Db.QueryRow(sqlstr, email)
	var username string
	rows.Scan(&username)
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return false
	}
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")             //上海
	timeStr := time.Now().In(cstSh).Format("2006-01-02 15:04:05") //In更改时区
	sqlstr = "insert into modifyemail values(?,?,?,?,?)"
	ret, err := tx.Exec(sqlstr, email, newemail, reason, timeStr, username)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...修改邮箱申请提交成功")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...,修改邮箱申请提交成功失败")
		return false
	}
	return true
}

//展示后台审核修改邮箱
func Findmodifyemail() ([]structure.Fmodiyemail, bool) {
	var Blogs = make([]structure.Fmodiyemail, 0) //读取所有记事
	sqlstr := "select * from modifyemail "
	rows, err := Db.Query(sqlstr)
	if err != nil {
		logrus.Error(err)
		return Blogs, false
	}
	for rows.Next() {
		var b structure.Fmodiyemail //读取所有记事
		var T string
		loginTime, _ := time.Parse("2006-01-02 15:04:05", T)
		rows.Scan(&b.Email, &b.Newemail, &b.Reason, &loginTime, &b.Username)
		b.Time = time.Time(loginTime).Format("2006-01-02 15:04:05")
		Blogs = append(Blogs, b)
	}
	return Blogs, true
}

//判断是否收藏过记事
func JudgeCollection(email string, id int) bool {
	sqlstr := "select people from collection where people=?; "
	rows := Db.QueryRow(sqlstr, email)
	var c string
	rows.Scan(&c)
	sqlstr = "select id from collection where id=?; "
	rows = Db.QueryRow(sqlstr, id)
	var Id int
	rows.Scan(&Id)
	if c == email && id == Id {
		return false
	} else {
		return true
	}
}

//收藏记事
func InsertCollection(email string, id int) bool {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return false
	}
	sqlstr := "insert into collection values(?,?)"
	ret, err := tx.Exec(sqlstr, email, id)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...,后台用户注册失败")
		return false
	}
	return true
}

//博客
func FindAlltxt() ([]structure.Blog, error) {
	var blog = make([]structure.Blog, 0) //读取所有记事
	sqlstr := "select id,title,praise,label from txt"
	rows, err := Db.Query(sqlstr)
	if err != nil {
		logrus.Error(err)
		return blog, err
	}
	for rows.Next() {
		var b structure.Blog
		rows.Scan(&b.Id, &b.Title, &b.Praise, &b.Label)
		blog = append(blog, b)
	}
	return blog, nil
}

//处理用户登录
func UserLogins(emailanduser string, password string) (bool, string, string) {
	sqlstr := "select email,username,password from user"
	rows, err := Db.Query(sqlstr)
	if err != nil {
		logrus.Error(err)
		return false, "", ""
	}
	for rows.Next() {
		var email string
		var username string
		var Password string
		rows.Scan(&email, &username, &Password)
		if (email == emailanduser || username == emailanduser) && password == Password {
			return true, email, username
		}
	}
	return false, "", ""
}

//判断邮箱和用户名是否重复注册和验证功能时判断邮箱是否正确
func Findemail(str string, Type string) error {
	var EMA string
	if Type == "regist" {
		rows := Db.QueryRow("select email from user where email=?", str)
		rows.Scan(&EMA)
		if EMA != "" {
			return errors.New("该邮箱已注册！")
		}
	} else if Type == "username" {
		rows := Db.QueryRow("select username from user where username=?", str)
		rows.Scan(&EMA)
		if EMA != "" {
			return errors.New("用户名已存在！")
		}
	} else if Type == "backgroud" {
		rows := Db.QueryRow("select email from user where email=?", str)
		rows.Scan(&EMA)
		if EMA != str {
			return errors.New("该邮箱与注册邮箱不一致！")
		}
	}
	return nil
}

//后台管理员注册
func SuperUserRegist(email string, username string, password string) bool {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return false
	}
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")             //上海
	timeStr := time.Now().In(cstSh).Format("2006-01-02 15:04:05") //In更改时区
	sqlstr := "insert into user(email,username,password,nickname,groupid,time) values(?,?,?,?,?,?)"
	ret, err := tx.Exec(sqlstr, email, username, password, "管理员", 2, timeStr)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...,后台用户注册失败")
		return false
	}
	return true
}

//把groupid 改为3
func Userlogout(email string) bool {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return false
	}
	sqlstr := "update user set groupid=3 where email=?"
	ret, err := tx.Exec(sqlstr, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...用户注销申请成功")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...,用户注销申请失败")
		return false
	}
	return true
}

//提交用户注销
func SubmitLogout(email string, reason string) bool {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return false
	}
	var username string
	rows := Db.QueryRow("select username from user where email=?", email)
	rows.Scan(&username)
	sqlstr := "insert into logout values(?,?,?)"
	ret, err := tx.Exec(sqlstr, email, reason, username)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...用户注销申请成功")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...,用户注销申请失败")
		return false
	}
	return true
}

//判断邮箱（后台）
func Findemailb(email string) error {
	var EMA string
	rows := Db.QueryRow("select email from user where email=?", email)
	rows.Scan(&EMA)
	if EMA != email {
		return errors.New("该邮箱与注册邮箱不一致！")
	}
	return nil
}

//判断后台管理员登录
func SuperLogins(email string, password string) (string, error) {
	var EMA string
	rows := Db.QueryRow("select email from user where email=?", email)
	rows.Scan(&EMA)
	if EMA == "" {
		return "", errors.New("该邮箱未注册！")
	}
	var username string
	rows = Db.QueryRow("select username from user where email=?", email)
	rows.Scan(&username)
	var PAS string
	var ID int
	rows = Db.QueryRow("select email,password,groupid from user where email=?", email)
	rows.Scan(&EMA, &PAS, &ID)
	if ID == 2 {
		if EMA == email && password == PAS {
			return username, nil
		} else {
			return "", errors.New("账号或密码错误")
		}
	} else {
		return "", errors.New("该用户没有管理员权限")
	}
}

//同意用户注销申请
func AgreeCellation(email string) bool {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return false
	}
	sqlstr := "delete from user where email=?"
	ret, err := tx.Exec(sqlstr, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	logrus.Infof("sucess:%d", affRow)
	sqlstr = "delete from logout where email=?"
	ret, err = tx.Exec(sqlstr, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err = ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	logrus.Infof("sucess:%d", affRow)
	sqlstr = "delete from txt where owner=?"
	ret, err = tx.Exec(sqlstr, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err = ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	logrus.Infof("sucess:%d", affRow)
	tx.Commit()
	return true
}

//处理用户注册把信息插入数据库
func UserRegists(email string, username string, password string, nickname string) error {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return err
	}
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")             //上海
	timeStr := time.Now().In(cstSh).Format("2006-01-02 15:04:05") //In更改时区
	sqlstr := "insert into user values(?,?,?,?,?,?,?,?,?)"
	ret, err := tx.Exec(sqlstr, email, username, password, nickname, "", 0, "", timeStr, "")
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...")
		return errors.New("注册失败")
	}
	return nil
}

//显示用户基本信息
func ShowUserinf(email string) structure.User {
	var s structure.User
	sqlstr := "select email,username,nickname,sex,avatar,time,name from user where email=?"
	rows := Db.QueryRow(sqlstr, email)
	var b []byte
	var t time.Time
	rows.Scan(&s.Email, &s.Username, &s.Nickname, &s.Sex, &b, &t, &s.Name)
	s.Time = time.Time(t).Format("2006-01-02 15:04:05")
	s.Avatar = string(b)
	return s
}

//修改个人中心基本信息
func Modifyuserinf(email string, sex string, avatar string, name string) error {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return err
	}
	sqlstr := "update user set sex=?,avatar=?,name=? where email=?"
	b := []byte(avatar)
	ret, err := tx.Exec(sqlstr, sex, b, name, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...")
		tx.Commit() // 提交事务
		return nil
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...")
		return errors.New("更新失败")
	}

}

//展示提交注销的用户
func Findlogoutuser() (bool, []structure.Cancellation) {
	var usr = make([]structure.Cancellation, 0)
	sqlstr := "select * from logout"
	rows, err := Db.Query(sqlstr)
	if err != nil {
		logrus.Error(err)
		return false, usr
	}
	for rows.Next() {
		var u structure.Cancellation
		rows.Scan(&u.Email, &u.Reason, &u.Username)
		usr = append(usr, u)
	}
	return true, usr
}

//点赞
func Txtlike(id int) bool {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return false
	}
	sqlstr := "update txt set praise=praise+1 where  id=?"
	ret, err := tx.Exec(sqlstr, id)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...点赞成功")
		tx.Commit() // 提交事务
		return true
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...点赞失败")
		return false
	}
}

//取消点赞
func Txtcancellike(id int) bool {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return false
	}
	sqlstr := "update txt set praise=praise-1 where  id=?"
	ret, err := tx.Exec(sqlstr, id)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...取消点赞成功")
		tx.Commit() // 提交事务
		return true
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...取消点赞失败")
		return false
	}
}

//取消收藏
func Scancelcollection(people string, id int) bool {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return false
	}
	sqlstr := "delete from collection where people=? and id=?"
	ret, err := tx.Exec(sqlstr, people, id)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...取消收藏成功")
		tx.Commit() // 提交事务
		return true
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...取消收藏失败")
		return false
	}

}

//查询所有记事
func Selecttxt(email string) ([]structure.UserTxt, structure.UserUAvatar, error) {
	var txt = make([]structure.UserTxt, 0) //读取所有记事
	var u structure.UserUAvatar            //读取用户民和头像
	sqlstr := "select id,title,content,label,time,complete from txt where owner=?"
	rows, err := Db.Query(sqlstr, email)
	if err != nil {
		logrus.Error(err)
		return txt, u, err
	}
	var ti time.Time
	for rows.Next() {
		var t structure.UserTxt
		rows.Scan(&t.Id, &t.Title, &t.Content, &t.Label, &ti, &t.Complete)
		t.Time = time.Time(ti).Format("2006-01-02 15:04:05")
		txt = append(txt, t)
	}
	sqlstr = "select username,avatar from user where email=?"
	row := Db.QueryRow(sqlstr, email)
	var b []byte //先用字节数组读取数据库内容在传给结构体
	row.Scan(&u.Username, &b)
	u.Avatar = string(b)
	return txt, u, nil
}

//查询用户的用户名和头像
func FinduserAndavatar(email string) structure.UserUAvatar {
	var u structure.UserUAvatar
	sqlstr := "select username,avatar from user where email=?"
	rows := Db.QueryRow(sqlstr, email)
	var b []byte //先用字节数组读取数据库内容在传给结构体
	rows.Scan(&u.Username, &b)
	u.Avatar = string(b)
	return u
}

//查询正在审核的用户
func FindExamUser() []structure.UserExamine {
	var u = make([]structure.UserExamine, 0)
	sqlstr := "select email,username,time from user where groupid=0"
	rows, err := Db.Query(sqlstr)
	if err != nil {
		logrus.Error(err)
	}
	var ti time.Time
	for rows.Next() {
		var t structure.UserExamine
		rows.Scan(&t.Email, &t.Username, &ti)
		t.Time = time.Time(ti).Format("2006-01-02 15:04:05")
		u = append(u, t)
	}
	return u
}

//查询正在审核的用户
func FindRegistUser() []structure.AllregistUser {
	var u = make([]structure.AllregistUser, 0)
	sqlstr := "select email,username,time,avatar from user where groupid=1"
	rows, err := Db.Query(sqlstr)
	if err != nil {
		logrus.Error(err)
	}
	var ti time.Time
	for rows.Next() {
		var t structure.AllregistUser
		var b []byte
		rows.Scan(&t.Email, &t.Username, &ti, &b)
		t.Time = time.Time(ti).Format("2006-01-02 15:04:05")
		t.Avatar = string(b)
		u = append(u, t)
	}
	return u
}

//增加记事
func AddTxts(owner string, title string, content string, label int) bool {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return false
	}

	var cstSh, _ = time.LoadLocation("Asia/Shanghai")             //上海
	timeStr := time.Now().In(cstSh).Format("2006-01-02 15:04:05") //In更改时区
	b := []byte(content)
	sqlstr := "insert into txt(owner,title,content,label,time,praise,complete) values(?,?,?,?,?,?,?)"
	ret, err := tx.Exec(sqlstr, owner, title, b, label, timeStr, 0, 0)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...增加记事成功")
		tx.Commit() // 提交事务
		return true
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...增加记事失败")
		return false
	}

}

//同意用户申请（groupid 由0改为1）
func SucessRegist(email string) error {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return err
	}
	sqlstr := "update user set groupid=1 where email=?"
	ret, err := tx.Exec(sqlstr, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...用户审核通过")
		tx.Commit() // 提交事务
		return nil
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...用户审核失败")
		return errors.New("审核驳回")
	}
}

//驳回用户邮箱修改申请
func ModifyEmailfail(email string) error {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return err
	}
	sqlstr := "delete from modifyemail where email=?"
	ret, err := tx.Exec(sqlstr, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...用户邮箱修改驳回成功")
		tx.Commit() // 提交事务
		return nil
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...用户邮箱修改驳回失败")
		return errors.New("审核驳回")
	}
}

//不同意该用户申请（删除）
func FailRegist(email string) error {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return err
	}
	sqlstr := "delete from user  where email=?"
	ret, err := tx.Exec(sqlstr, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...用户审核驳回成功")
		tx.Commit() // 提交事务
		return nil
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...用户审核驳回失败")
		return errors.New("审核驳回")
	}
}

//忘记密码要设置新的密码
func ForgetPassword(email string, password string) error {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return err
	}
	sqlstr := "update user set password=? where email=?"
	ret, err := tx.Exec(sqlstr, password, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	if affRow == 1 {
		logrus.Info("事务提交啦...密码修改成功")
		tx.Commit() // 提交事务
		return nil
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦...增加记事失败")
		return errors.New("密码修改失败")
	}

}

//把上传的头像增加到数据库
func AddAvatar(email string, avatar string) error {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return err
	}
	a := []byte(avatar) //把图片信息（字符串）转成字节存到数据库
	sqlstr := "update user set avatar=? where email=?"
	ret, err := tx.Exec(sqlstr, a, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	if affRow == 1 {
		logrus.Info("事务提交啦头像上传成功")
		tx.Commit() // 提交事务
		return nil
	} else {
		tx.Rollback()
		logrus.Info("事务回滚啦头像上传失败")
		return errors.New("头像上传失败")
	}
}

//删除记事（根据id）
func DeleteTxts(id int) bool {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return false
	}
	sqlstr := "delete from txt where id=?"
	ret, err := tx.Exec(sqlstr, id)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return false
	}
	if affRow == 1 {
		logrus.Info("记事删除成功")
		tx.Commit() // 提交事务
		return true
	} else {
		tx.Rollback()
		logrus.Info("记事删除失败")
		return false
	}
}

//修改记事
func ModifyTxts(title string, content string, label int, complete int, id int) error {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return err
	}
	b := []byte(content)
	sqlstr := "update txt set title=?,content=?,label=?,complete=? where id=?"
	ret, err := tx.Exec(sqlstr, title, b, label, complete, id)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	if affRow == 1 {
		logrus.Info("记事修改成功")
		tx.Commit() // 提交事务
		return nil
	} else {
		tx.Rollback()
		logrus.Info("记事修改失败")
		return errors.New("记事修改失败")
	}
}

//修改邮箱
func ModifyEmails(email string, newemail string) error {
	tx, err := Db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		logrus.Error(err)
		return err
	}
	sqlstr := "update user set email=? where email=?;"
	ret, err := tx.Exec(sqlstr, newemail, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	affRow, err := ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	logrus.Infof("sucess:%d", affRow)
	sqlstr = "update collection set people=? where people=?"
	ret, err = tx.Exec(sqlstr, newemail, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	affRow, err = ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	logrus.Infof("sucess:%d", affRow)
	sqlstr = "update txt set owner=? where owner=?"
	ret, err = tx.Exec(sqlstr, newemail, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	affRow, err = ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	logrus.Infof("sucess:%d", affRow)
	sqlstr = "delete from modifyemail  where email=?"
	ret, err = tx.Exec(sqlstr, newemail, email)
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	affRow, err = ret.RowsAffected() //获取操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		logrus.Error(err)
		return err
	}
	logrus.Infof("sucess:%d", affRow)
	tx.Commit()
	return nil
}

//查看记事
func Findtxt(id int) (bool, structure.UserLabel) {
	var data structure.UserLabel
	sqlstr := "select owner from txt where id=? "
	rows := Db.QueryRow(sqlstr, id)
	var owner string
	rows.Scan(&owner)
	sqlstr = "select username,avatar from user where email=? "
	rows = Db.QueryRow(sqlstr, owner)
	var b []byte
	var username string
	rows.Scan(&username, &b)
	avatar := string(b)
	data.Username = username
	data.Avatar = avatar
	sqlstr = "select title,content,label,time,praise from txt where id=? "
	rows = Db.QueryRow(sqlstr, id)
	var bt []byte
	var ti time.Time
	rows.Scan(&data.Title, &bt, &data.Label, &ti, &data.Praise)
	data.Time = time.Time(ti).Format("2006-01-02 15:04:05")
	data.Id = id
	data.Content = string(bt)
	if data.Time == "" {
		return false, data
	}
	return true, data
}

//查询头像
func FindAvatar(username string) string {
	var b []byte
	sqlstr := "select avatar from user where username=? "
	rows := Db.QueryRow(sqlstr, username)
	rows.Scan(&b)
	avatar := string(b)
	return avatar
}

//查询该用户是否已经审核过
func Findgroupid(emailanduser string) error {
	sqlstr := "select groupid from user where email=?"
	rows := Db.QueryRow(sqlstr, emailanduser)
	var b = 0
	rows.Scan(&b)
	if b == 3 {
		return errors.New("该用户已提交注销申请，暂时无法登录")
	}
	if b == 1 {
		return nil //代表审核过了
	} else {
		sqlstr := "select groupid from user where username=? "
		rows := Db.QueryRow(sqlstr, emailanduser)
		rows.Scan(&b)
		if b == 3 {
			return errors.New("该用户已提交注销申请，暂时无法登录")
		}
		if b == 0 {
			return errors.New("该用户正在审核")
		} else {
			return nil
		}
	}
}

//连接数据库
func Opengowebsql() {
	go func() {
		Db, errs = sql.Open("mysql", config.C.Host)
		if errs != nil {
			logrus.Info(errs)
		}
		//尝试与数据库建立连接（校验open的第二个参数是否正确）
		errs = Db.Ping()
		if errs != nil {
			return
		}
		Db.SetMaxOpenConns(1000)
	}()
}
