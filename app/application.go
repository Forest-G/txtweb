package app

import (
	"txt/pkg/sendemail"
	"txt/pkg/structure"
	"txt/redisdata"
	"txt/sqldata"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

//验证登录的用户名或邮箱还有密码是否正确
func SolveLogin(emailanduser string, password string) (bool, interface{}, string, string) {
	flag, email, username := sqldata.UserLogins(emailanduser, password)
	token := uuid.NewV4().String()
	if flag {
		var s structure.Respons
		s.Code = 200
		s.Msg = "登录成功"
		s.Email = email
		s.Token = token
		s.Username = username
		return flag, s, token, username
	} else {
		var s structure.Resp
		s.Code = 202
		s.Msg = "账号或密码错误"
		return flag, s, "", ""
	}
}

//验证功能验证邮箱和验证码是否正确
func SolveValidation(email string, icode string) structure.Resp {
	var s structure.Resp
	err := redisdata.FindVcode_Validation(email, icode)
	if err != nil {
		s.Code = 202
		s.Msg = "验证失败"
	} else {
		s.Code = 200
		s.Msg = "验证成功"
	}
	return s
}

//根据不同类型，发送验证码时需要做不同的判断
func SolveDeliveremail(email string, emailtype string) structure.Resp {
	var s structure.Resp
	if emailtype == "login" {
		err := sendemail.DeliverEmail(email) //发送并获取到验证码
		if err != nil {
			s.Code = 202
			s.Msg = "验证码发送失败！"
			return s
		}
		s.Code = 200
		s.Msg = "验证码已发送！"
		logrus.Info(s.Msg)
		return s
	} else if emailtype == "regist" {
		err := sqldata.Findemail(email, "regist") //判断邮箱是否已经注册
		if err != nil {
			logrus.Error(err)
			s.Code = 203
			s.Msg = "该邮箱已注册！"
			logrus.Info(s.Msg)
		} else {
			err := sendemail.DeliverEmail(email)
			if err != nil {
				s.Code = 202
				s.Msg = "验证码发送失败！"
				return s
			}
			s.Code = 200
			s.Msg = "验证码已发送！"
			logrus.Info(s.Msg)
		}
		return s
	} else if emailtype == "background" {
		err := sqldata.Findemail(email, "background") //判断邮箱与注册邮箱是否一致（验证功能）
		if err != nil {
			logrus.Error(err)
			s.Code = 204
			s.Msg = "该邮箱与注册邮箱不一致！"
			logrus.Info(s.Msg)
		} else {
			err := sendemail.DeliverEmail(email) //发送并获取到验证码
			if err != nil {
				s.Code = 202
				s.Msg = "验证码发送失败！"
				return s
			}
			s.Code = 200
			s.Msg = "验证码已发送！"
			logrus.Info(s.Msg)
		}
		return s
	}
	return s
}

//查询redis验证验证码是否正确
func SolveVcode(email string) string {
	str := redisdata.FindVcode(email)
	return str
}

//解决后台管理员登录
func SolveSuperlogin(email string, password string, icode string) interface{} {
	username, err := sqldata.SuperLogins(email, password)
	var s structure.Resp
	if err != nil {
		s.Code = 202
		s.Msg = err.Error()
		return s
	} else {
		Vcode := SolveVcode(email)
		if Vcode == icode {
			var t structure.Respons
			t.Code = 200
			t.Msg = "登录成功"
			token := uuid.NewV4().String()
			LoginRedis(username, token) //设置token
			t.Email = email
			t.Token = token
			t.Username = username
			return t
		} else {
			s.Code = 202
			s.Msg = "验证码错误"
			return s
		}
	}
}

//注册管理员用户
func SolveSuperRegist(email string, username string, password string, icode string) structure.Resp {
	var s structure.Resp
	err := sqldata.Findemail(email, "regist") //判断注册时邮箱是否重复
	if err != nil {
		s.Code = 201
		s.Msg = err.Error()
		logrus.Error(err.Error())
		return s
	}
	err = sqldata.Findemail(username, "username") //判断注册时用户名是否重复
	if err != nil {
		s.Code = 204
		s.Msg = err.Error()
		logrus.Error(err.Error())
		return s
	}
	b := sqldata.SuperUserRegist(email, username, password) //去注册新用户
	if !b {
		logrus.Error("后台管理员注册失败")
		s.Code = 202
		s.Msg = "后台管理员注册失败"
		return s
	}
	Vcode := SolveVcode(email)
	if Vcode == icode {
		s.Code = 200
		s.Msg = "后台管理员注册成功"
		logrus.Info("后台管理员注册成功")
	} else {
		s.Code = 203
		s.Msg = "验证码错误"
	}
	return s
}

//注册新用户
func SolveRegist(email string, username string, password string, nickname string, icode string) structure.Resp {
	var s structure.Resp
	err := sqldata.Findemail(email, "regist") //判断注册时邮箱是否重复
	if err != nil {
		s.Code = 202
		s.Msg = err.Error()
		logrus.Error(err.Error())
		return s
	}
	err = sqldata.Findemail(username, "username") //
	if err != nil {
		s.Code = 204
		s.Msg = err.Error()
		logrus.Error(err.Error())
		return s
	}
	err = sqldata.UserRegists(email, username, password, nickname) //去注册新用户
	if err != nil {
		logrus.Error(err)
		s.Code = 202
		s.Msg = "注册失败"
		return s
	}
	Vcode := SolveVcode(email)
	if Vcode == icode {
		s.Code = 200
		s.Msg = "注册成功"
	} else {
		s.Code = 202
		s.Msg = "验证码错误"
	}

	return s
}

//用户注销申请
func SolveUserLogout(email string, reason string, icode string) structure.Resp {
	var s structure.Resp
	b := sqldata.SubmitLogout(email, reason)
	if b {
		bo := sqldata.Userlogout(email)
		if bo {
			Vcode := SolveVcode(email)
			if Vcode == icode {
				s.Code = 200
				s.Msg = "用户注销提交成功"
				logrus.Info("用户注销提交成功")
				return s
			} else {
				s.Code = 203
				s.Msg = "验证码错误"
				return s
			}
		} else {
			s.Code = 202
			s.Msg = "用户注销提交失败"
			logrus.Info("用户注销提交失败")
			return s
		}
	} else {
		s.Code = 202
		s.Msg = "用户注销提交失败"
		logrus.Info("用户注销提交失败")
		return s
	}
}

//同意该用户注册申请
func SolveUserRegistS(email string) structure.Resp {
	err := sqldata.SucessRegist(email)
	var s structure.Resp
	if err != nil {
		s.Code = 202
		s.Msg = "用户审核失败！"
	} else {
		s.Code = 200
		s.Msg = "用户审核成功！"
		sendemail.DeliverSucessEmail(email)
	}
	return s
}

//驳回用户注销申请
func SolveCellationfail(email string, reason string) structure.Resp {
	var s structure.Resp
	err := sendemail.DeliverFailcEmail(email, reason)
	if err != nil {
		s.Code = 202
		s.Msg = "驳回用户注销失败！"
		logrus.Error("驳回用户注销失败！")
	} else {
		s.Code = 200
		s.Msg = "驳回用户注销成功！"
		logrus.Info("驳回用户注销成功！")
	}
	return s
}

//驳回该用户注册申请
func SolveUserRegistF(email string, reason string) structure.Resp {
	err := sqldata.FailRegist(email)
	var s structure.Resp
	if err != nil {
		s.Code = 202
		s.Msg = "驳回用户审核失败！"
		logrus.Error("驳回用户审核失败！")
	} else {
		s.Code = 200
		s.Msg = "驳回用户审核成功！"
		logrus.Info("驳回用户审核成功！")
		sendemail.DeliverFailEmail(email, reason)
	}
	return s
}

//忘记密码验证
func SolveForgetPassword(email string, password string, newpassword string) structure.Resp {
	err := sqldata.ForgetPassword(email, password)
	var s structure.Resp
	if err != nil {
		s.Code = 202
		s.Msg = "密码设置失败！"
	} else {
		s.Code = 200
		s.Msg = "密码设置成功！"
	}
	return s
}

//展示用户所有记事
func SolveShowTxt(email string) structure.Res {
	datas, data, err := sqldata.Selecttxt(email)
	if err != nil {
		logrus.Error(err)
	}
	var s structure.Res
	s.Code = 200
	s.Username = data.Username
	s.Avatar = data.Avatar
	s.Msg = datas
	return s
}

//展示用户头像和用户名
func SolveShowua(email string) structure.Resp {
	datas := sqldata.FinduserAndavatar(email)
	var s structure.Resp
	s.Code = 200
	s.Msg = datas
	return s
}

//增加记事
func SolveAddTxt(owner string, title string, content string, label int) structure.Resp {
	b := sqldata.AddTxts(owner, title, content, label)
	var s structure.Resp
	if b {
		s.Code = 200
		s.Msg = "记事增加成功"
		logrus.Info("记事增加成功")
		return s
	} else {
		s.Code = 202
		s.Msg = "记事增加失败！"
		logrus.Info("记事增加失败")
		return s
	}
}

//判断用户登录时用户名是否存在
func SolveExit(emailanduser string) error {
	err := sqldata.EuserLogins(emailanduser)
	return err
}

//显示所有记事
func SolveShowAllTxt() structure.Resp {
	data, err := sqldata.FindAlltxt()
	var s structure.Resp
	if err != nil {
		logrus.Error(err)
		s.Code = 202
		s.Msg = "获取信息失败"
	}
	s.Code = 200
	s.Msg = data
	return s
}

//后天展示待审核修改邮箱
func SolveCheckModifyEmail() structure.Resp {
	data, b := sqldata.Findmodifyemail()
	var s structure.Resp
	if b {
		s.Code = 200
		s.Msg = data
		return s
	}
	s.Code = 202
	s.Msg = "审核信息获取失败"
	logrus.Error(s.Msg)
	return s
}

//提交修改邮箱审核
func SolveSubModifyEmail(email string, newemail string, reason string) structure.Resp {
	b := sqldata.Subnewemail(email, newemail, reason)
	var s structure.Resp
	if b {
		s.Code = 200
		s.Msg = "提交成功"
		return s
	}
	s.Code = 202
	s.Msg = "提交失败"
	logrus.Error(s.Msg)
	return s
}

//搜索
func SolveSearch(key string) structure.Resp {
	data, b := sqldata.FindSearch(key)
	var s structure.Resp
	if !b {
		s.Code = 202
		s.Msg = "获取信息失败"
		logrus.Error(s.Msg)
	}
	s.Code = 200
	s.Msg = data
	return s
}

//收藏
func SolveCollection(email string, id int) structure.Resp {
	b := sqldata.InsertCollection(email, id)
	var s structure.Resp
	if b {
		s.Code = 200
		s.Msg = "文章收藏成功"
		logrus.Info(s.Msg)
	} else {
		s.Code = 202
		s.Msg = "文章收藏失败"
		logrus.Error("文章收藏失败")
	}
	return s
}

//判断是否收藏过记事
func SolveJudgeCollection(email string, id int) structure.Resp {
	b := sqldata.JudgeCollection(email, id)
	var s structure.Resp
	if b {
		s.Code = 200
		s.Msg = "文章未收藏"
		logrus.Info(s.Msg)
	} else {
		s.Code = 202
		s.Msg = "文章已收藏"
		logrus.Error("文章收藏失败")
	}
	return s
}

//展示我收藏的记事
func SolveMyCollection(email string) structure.Resp {
	data, err := sqldata.FindMyCollection(email)
	var s structure.Resp
	if err != nil {
		logrus.Error(err)
		s.Code = 202
		s.Msg = "获取信息失败"
	}
	s.Code = 200
	s.Msg = data
	return s
}

//取消点赞
func SolveCancelLike(id int) structure.Resp {
	var s structure.Resp
	b := sqldata.Txtcancellike(id)
	if b {
		s.Code = 200
		s.Msg = "取消点赞成功"
		logrus.Info(s.Msg)
	} else {
		s.Code = 202
		s.Msg = "取消点赞失败"
		logrus.Error(s.Msg)
	}
	return s
}

//点赞
func SolveLike(id int) structure.Resp {
	var s structure.Resp
	b := sqldata.Txtlike(id)
	if b {
		s.Code = 200
		s.Msg = "点赞成功"
		logrus.Info(s.Msg)
	} else {
		s.Code = 202
		s.Msg = "点赞失败"
		logrus.Error(s.Msg)
	}
	return s
}

//取消收藏
func SolveCancelCollection(people string, id int) structure.Resp {
	var s structure.Resp
	b := sqldata.Scancelcollection(people, id)
	if b {
		s.Code = 200
		s.Msg = "取消收藏成功"
	} else {
		s.Code = 202
		s.Msg = "取消收藏失败"
		logrus.Error(s.Msg)
	}
	return s
}

//后台展示提交用户注销的用户
func SolveLogout() structure.Resp {
	var s structure.Resp
	b, data := sqldata.Findlogoutuser()
	if b {
		s.Code = 200
		s.Msg = data
	} else {
		s.Code = 202
		s.Msg = "注销用户信息获取失败"
		logrus.Error(s.Msg)
	}
	return s
}

//把上传的头像添加到数据库
func SolveaddAvatar(email string, avatar string) structure.Resp {
	err := sqldata.AddAvatar(email, avatar)
	var s structure.Resp
	if err != nil {
		s.Code = 202
		s.Msg = "头像上传失败"
		return s
	} else {
		s.Code = 200
		s.Msg = "头像上传成功"
		return s
	}

}

//删除记事
func SolveDeleteTxt(id int) structure.Resp {
	flag := sqldata.DeleteTxts(id)
	var s structure.Resp
	if flag {
		s.Code = 200
		s.Msg = "记事删除成功"
		return s
	} else {
		s.Code = 202
		s.Msg = "记事删除失败"
		return s
	}
}

//修改记事
func SolveModifyTxt(title string, content string, label int, complete int, id int) structure.Resp {
	err := sqldata.ModifyTxts(title, content, label, complete, id)
	var s structure.Resp
	if err != nil {
		s.Code = 202
		s.Msg = "记事修改失败"
		return s
	} else {
		s.Code = 200
		s.Msg = "记事修改成功"
		return s
	}
}

//修改邮箱
func SolveModifyEmail(email string, newemail string) structure.Resp {
	err := sqldata.ModifyEmails(email, newemail)
	var s structure.Resp
	if err != nil {
		s.Code = 202
		s.Msg = "邮箱改失败"
		return s
	} else {
		s.Code = 200
		s.Msg = "邮箱修改成功"
		return s
	}
}

//登录的总逻辑处理
func SolveLoginall(emailanduser string, password string) interface{} {
	sflag := SolveExit(emailanduser)
	if sflag == nil {
		flag, failstr := SolveFindgroupid(emailanduser) //判断是否在审核中
		if flag {
			b, str, token, username := SolveLogin(emailanduser, password)
			if b {
				logrus.Info("登录成功")
				LoginRedis(username, token) //设置token
			} else {
				logrus.Info("登录失败")
			}
			return str
		} else {
			var s structure.Resp
			s.Code = 203
			s.Msg = failstr
			logrus.Error(s.Msg)
			return s
		}
	} else {
		var s structure.Resp
		s.Code = 201
		s.Msg = "用户名不存在"
		logrus.Error(s.Msg)
		return s
	}
}

//查看记事
func SolveFindtxt(id int) structure.Resp {
	b, data := sqldata.Findtxt(id)
	var s structure.Resp
	if b {
		s.Code = 200
		s.Msg = data
		logrus.Info("记事查看成功")
	} else {
		s.Code = 202
		s.Msg = "该篇内容不存在"
		logrus.Info("该篇内容不存在")
	}
	return s
}

//展示正在审核的用户
func SolveShowExamregist() structure.Resp {
	data := sqldata.FindExamUser()
	var s structure.Resp
	s.Code = 200
	s.Msg = data
	return s
}

//展示个人基本信息
func SolveUserinf(email string) structure.Resp {
	data := sqldata.ShowUserinf(email)
	var s structure.Resp
	s.Code = 200
	s.Msg = data
	return s
}

//用户登录成功后把token 存到redis里面
func LoginRedis(username string, token string) {
	err := redisdata.Ruserlogin(username, token)
	if err != nil {
		logrus.Error(err)
	}
}

//同意用户注销申请
func SolveSucessCellaiton(email string) structure.Resp {
	var s structure.Resp
	b := sqldata.AgreeCellation(email)
	if b {
		s.Code = 200
		s.Msg = "操作成功"
		logrus.Info("同意用户注销申请操作成功")
		sendemail.DeliverSucesscEmail(email)
	} else {
		s.Code = 202
		s.Msg = "操作失败"
		logrus.Info("同意用户注销申请操作失败")
	}
	return s
}

//用户登录成功后把token 存到redis里面
func SolveLoginRedis(username string, token string) structure.Respp {
	err := redisdata.FindRuserlogin(username, token)
	avatar := sqldata.FindAvatar(username)
	var s structure.Respp
	if err != nil {
		s.Code = 202
		s.Msg = "用户身份已过期"
		return s
	} else {
		s.Code = 200
		s.Msg = "用户身份通过"
		s.Avatar = avatar
	}
	return s
}

//展示个人基本信息
func SolveUserModify(email string, sex string, avatar string, name string) structure.Resp {
	err := sqldata.Modifyuserinf(email, sex, avatar, name)
	var s structure.Resp
	if err != nil {
		logrus.Error(err)
		s.Code = 202
		s.Msg = "个人信息更新失败"
	} else {
		s.Code = 200
		s.Msg = "个人信息更新成功"
	}
	return s
}

//后台展示所有注册过的用户
func SolveregistUser() structure.Resp {
	data := sqldata.FindRegistUser()
	var s structure.Resp
	s.Code = 200
	s.Msg = data
	return s
}

//驳回用户修改邮箱申请
func SolveModifyEmailfail(email string, reason string) structure.Resp {
	err := sqldata.ModifyEmailfail(email)
	var s structure.Resp
	if err != nil {
		logrus.Error(err)
		s.Code = 202
		s.Msg = "驳回邮箱修改失败"
	} else {
		s.Code = 200
		s.Msg = "驳回邮箱修改成功"
		sendemail.DeliverModifyEmailfail(email, reason)
	}
	return s
}

//查询该用户是否已经审核过
func SolveFindgroupid(emailanduser string) (bool, string) {
	err := sqldata.Findgroupid(emailanduser)
	if err != nil {
		return false, err.Error()
	} else {
		return true, ""
	}
}
