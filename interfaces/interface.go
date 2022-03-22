package interfaces

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"txt/app"
	"txt/pkg/config"
	"txt/pkg/sendemail"
	"txt/pkg/structure"

	"github.com/sirupsen/logrus"
)

var (
	jsonerr structure.JsonError
)

//接收前端传来的数据
func ResiveRequest(r io.Reader) ([]byte, error) {
	data, err := ioutil.ReadAll(r) //此处的r是http请求得到的json格式数据-->然后转化为[]byte格式数据.
	if err != nil {
		logrus.Error(err)
		return data, errors.New("接受数据失败")

	}
	return data, nil
}

//返回响应给前端
func SetResponse(w http.ResponseWriter, str interface{}) {
	var STR, err = json.Marshal(str)
	if err != nil {
		logrus.Error(err)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(STR)))
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	w.Write(STR)
}

// var (
// 	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
// 	key   = []byte("super-secret-key")
// 	store = sessions.NewCookieStore(key)
// )

//处理用户登录
/*先判断用户存不存在，然后判断是否在审核，再判断登录成功与否*/
func Login(w http.ResponseWriter, r *http.Request) {
	var userlogin structure.UserLogin
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &userlogin)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	fflag := app.SolveLoginall(userlogin.Emailanduser, userlogin.Password)
	SetResponse(w, fflag)
}

//用户登录身份识别（单点登录）
func IdentiFication(w http.ResponseWriter, r *http.Request) {
	var e structure.Userun //获取前端发来的用户名
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	c1 := e.Token                            //获取到token
	s := app.SolveLoginRedis(e.Username, c1) //跟redis里面token进行比较
	//s, b := app.SolveLoginRedis(e.Username, c1) //跟redis里面token进行比较
	// if !b {
	// 	SetResponse(w, s)
	// 	return
	// }
	SetResponse(w, s)
}

//忘记密码进行验证后修改新的密码
func PasswordSucess(w http.ResponseWriter, r *http.Request) {
	var e structure.ModifyPassword
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	s := app.SolveForgetPassword(e.Email, e.Password, e.Newpassword)
	SetResponse(w, s)
}

//显示所有记事
func Showalltxt(w http.ResponseWriter, r *http.Request) {
	datas := app.SolveShowAllTxt()
	SetResponse(w, datas)
}

//收藏记事
func Collection(w http.ResponseWriter, r *http.Request) {
	var e structure.Collection
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	s := app.SolveCollection(e.Email, e.Id)
	SetResponse(w, s)
}

//展示我收藏的记事
func ShowCollection(w http.ResponseWriter, r *http.Request) {
	var e structure.UserEmail
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	s := app.SolveMyCollection(e.Email)
	SetResponse(w, s)
}

//搜索
func Search(w http.ResponseWriter, r *http.Request) {
	var e structure.Keyword
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	s := app.SolveSearch(e.Keryword)
	SetResponse(w, s)
}

//提交邮箱修改申请
func SubModifyEmail(w http.ResponseWriter, r *http.Request) {
	var e structure.SubModiEmail
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	s := app.SolveSubModifyEmail(e.Email, e.Newemail, e.Reason)
	SetResponse(w, s)
}

//后台展示待审核邮箱
func CheckModifyEmail(w http.ResponseWriter, r *http.Request) {
	s := app.SolveCheckModifyEmail()
	SetResponse(w, s)
}

//取消记事收藏
func CancelCollection(w http.ResponseWriter, r *http.Request) {
	var e structure.Collection
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	s := app.SolveCancelCollection(e.Email, e.Id)
	SetResponse(w, s)
}

//点赞记事本
func Like(w http.ResponseWriter, r *http.Request) {
	var e structure.TxtId
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	s := app.SolveLike(e.Id)
	SetResponse(w, s)
}

//取消点赞
func CancelLike(w http.ResponseWriter, r *http.Request) {
	var e structure.TxtId
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	s := app.SolveCancelLike(e.Id)
	SetResponse(w, s)
}

//同意用户注销申请
func CellationSucess(w http.ResponseWriter, r *http.Request) {
	var e structure.UserEmail //获取前端发来的邮箱
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	datas := app.SolveSucessCellaiton(e.Email)
	SetResponse(w, datas)
}

//驳回用户注销申请
func CellationFail(w http.ResponseWriter, r *http.Request) {
	var e structure.EmailandReason //获取前端发来的邮箱
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error("json格式错误")
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	datas := app.SolveCellationfail(e.Email, e.Reason)
	SetResponse(w, datas)
}

//后台审核同意用户审核
func RegistSucess(w http.ResponseWriter, r *http.Request) {
	var e structure.UserEmail //获取前端发来的邮箱
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	datas := app.SolveUserRegistS(e.Email)

	SetResponse(w, datas)
}

//用户注销
func UserCellation(w http.ResponseWriter, r *http.Request) {
	var e structure.UserCellation
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		return
	}
	datas := app.SolveUserLogout(e.Email, e.Reason, e.Icode)
	SetResponse(w, datas)
}

//后台展示所有提交注销申请的用户
func Cellation(w http.ResponseWriter, r *http.Request) {
	datas := app.SolveLogout()
	SetResponse(w, datas)
}

//后台驳回用户审核
func RegistFail(w http.ResponseWriter, r *http.Request) {
	var e structure.EmailandReason //获取前端发来的邮箱
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		return
	}
	datas := app.SolveUserRegistF(e.Email, e.Reason)
	SetResponse(w, datas)
}

//处理用户注册
func Regist(w http.ResponseWriter, r *http.Request) {
	var userregist structure.UserRegist
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &userregist)
	if err != nil {
		logrus.Error(err)
		return
	}
	s := app.SolveRegist(userregist.Email, userregist.Username, userregist.Password, userregist.Nickname, userregist.Icode)
	SetResponse(w, s)
}

//展示审核注册的用户
func ShowExamregist(w http.ResponseWriter, r *http.Request) {
	datas := app.SolveShowExamregist()
	SetResponse(w, datas)
}

//个人中心展示基本信息
func Userinf(w http.ResponseWriter, r *http.Request) {
	var e structure.UserEmail
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		return
	}
	datas := app.SolveUserinf(e.Email)
	SetResponse(w, datas)
}

//修改个人中心基本信息
func UserinfModify(w http.ResponseWriter, r *http.Request) {
	var u structure.Usermodify
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &u)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	datas := app.SolveUserModify(u.Email, u.Sex, u.Avatar, u.Name)
	SetResponse(w, datas)
}

//后台显示所有已经注册过的用户信息
func AllregistUser(w http.ResponseWriter, r *http.Request) {
	datas := app.SolveregistUser()
	SetResponse(w, datas)
}

//修改邮箱
func ModifyEmail(w http.ResponseWriter, r *http.Request) {
	var e structure.ModiEmail
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	datas := app.SolveModifyEmail(e.Email, e.Newemail)
	SetResponse(w, datas)
}

//后台管理员注册
func Superregist(w http.ResponseWriter, r *http.Request) {
	var e structure.SuperUserRegist
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		return
	}
	s := app.SolveSuperRegist(e.Email, e.Username, e.Password, e.Icode)
	SetResponse(w, s)
}

//后台管理员登录
func SuperLogin(w http.ResponseWriter, r *http.Request) {
	var e structure.SuperUserLogin
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	datas := app.SolveSuperlogin(e.Email, e.Password, e.Icode)
	SetResponse(w, datas)
}

//驳回用户的邮箱申请
func CheckModifyEmailFail(w http.ResponseWriter, r *http.Request) {
	var e structure.EmailandReason
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	datas := app.SolveModifyEmailfail(e.Email, e.Reason)
	SetResponse(w, datas)
}

//同意用户的邮箱申请
func CheckModifyEmailSucess(w http.ResponseWriter, r *http.Request) {
	var e structure.ModiEmail
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	datas := app.SolveModifyEmail(e.Email, e.Newemail)
	SetResponse(w, datas)
}

//生成验证码并发送给前端
func CodeCheck(w http.ResponseWriter, r *http.Request) {
	var e sendemail.Remail
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		return
	}
	s := app.SolveDeliveremail(e.Email, e.Emailtype)
	SetResponse(w, s)
	// Vcode = email.DeliverEmail(e.Email) //发送并获取到验证码
	// b := app.SolveVcode(Vcode)          //处理是否验证成功
	// if b {
	// 	logrus.Info("验证码发送成功")
	// } else {
	// 	logrus.Info("验证码发送失败")
	// }

}

//查询redis去验证邮箱和验证码是否匹配
func Validation(w http.ResponseWriter, r *http.Request) {
	var e structure.EmailandIcode
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		return
	}
	s := app.SolveValidation(e.Email, e.Icode)
	SetResponse(w, s)
}

//登录成功后展示所有记事
func Showtxt(w http.ResponseWriter, r *http.Request) {
	var e structure.UserEmail //获取前端发来的邮箱
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		return
	}
	datas := app.SolveShowTxt(e.Email)
	SetResponse(w, datas)
}

//增加记事
func AddTxt(w http.ResponseWriter, r *http.Request) {
	var e structure.UserAddTxt
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		return
	}
	datas := app.SolveAddTxt(e.Email, e.Title, e.Content, e.Label)
	SetResponse(w, datas)
}

//上传头像（把头像信息存到数据库内）
func AvatarImage(w http.ResponseWriter, r *http.Request) {
	var e structure.UserAvatar
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		return
	}
	datas := app.SolveaddAvatar(e.Email, e.Avatar)
	SetResponse(w, datas)
}

//删除记事
func DeleteTxt(w http.ResponseWriter, r *http.Request) {
	var e structure.TxtId
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		return
	}
	datas := app.SolveDeleteTxt(e.Id)
	SetResponse(w, datas)
}

//修改记事
func ModifyTxt(w http.ResponseWriter, r *http.Request) {
	var e structure.UserModifyTxt
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	datas := app.SolveModifyTxt(e.Title, e.Content, e.Label, e.Complete, e.Id)
	SetResponse(w, datas)
}

//查看记事
func FindTxt(w http.ResponseWriter, r *http.Request) {
	var e structure.TxtId
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		return
	}
	datas := app.SolveFindtxt(e.Id)
	SetResponse(w, datas)
}

//判断是否收藏过记事
func JudgeCollection(w http.ResponseWriter, r *http.Request) {
	var e structure.Collection
	data, err := ResiveRequest(r.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		logrus.Error(err)
		jsonerr.Information = "json格式错误"
		SetResponse(w, jsonerr)
		return
	}
	s := app.SolveJudgeCollection(e.Email, e.Id)
	SetResponse(w, s)
}

//把前端传的图片存到文件夹里
func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	var i structure.Resp
	r.ParseMultipartForm(100)
	mForm := r.MultipartForm
	for k, f := range mForm.File {
		// k is the key of file part
		file, fileHeader, err := r.FormFile(k)
		if err != nil {
			logrus.Error("inovke FormFile error:", err, f)
			i.Code = 202
			SetResponse(w, i)
			return
		}
		defer file.Close()
		// store uploaded file into local path
		s := fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(100000000))
		arr := strings.Split(fileHeader.Filename, ".")
		localFileName := config.C.Imagespath + "/" + s + "." + arr[1]
		m := s + "." + arr[1] //生成图片的地址，随机数加后缀
		out, err := os.Create(localFileName)
		if err != nil {
			logrus.Error("failed to open the file %s for writing", localFileName)
			i.Code = 202
			i.Msg="图片上传失败"
			SetResponse(w, i)
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			logrus.Error("copy file err:%s\n", err)
			i.Code = 202
			i.Msg="图片上传失败"
			SetResponse(w, i)
			return
		}
		logrus.Info("file %s uploaded ok\n", fileHeader.Filename)
		i.Code = 200
		i.Msg = "http://120.76.128.13:8088/uploads/" + m
		SetResponse(w, i)
	}
}

func showhandleUploadFile(w http.ResponseWriter, r *http.Request) {
	remPartOfURL := r.URL.Path[len("/uploads/"):] // get everything after the /hello/ part of the URL

	file, err := os.Open("images/" + remPartOfURL)

	if err != nil {
		w.Write([]byte(err.Error()))
	}
	defer file.Close()
	buff, _ := ioutil.ReadAll(file)
	w.Write(buff)
}

func Init() {
	http.HandleFunc("/user/login", Login)              //用户登录
	http.HandleFunc("/IdentiFication", IdentiFication) //用户身份验证（获取cookie）
	http.HandleFunc("/user/regist", Regist)            //用户注册
	http.HandleFunc("/email", CodeCheck)               //发送验证码 成功就存到redis
	http.HandleFunc("/user/alltxt", Showtxt)           //登录成功后展示我的所有记事
	http.HandleFunc("/validation", Validation)         //验证功能（修改邮箱密码之前先验证一次）
	http.HandleFunc("/imag", AvatarImage)              //上传头像
	http.HandleFunc("/fogetpassword", PasswordSucess)  //忘记密码验证后改变密码

	http.HandleFunc("/modifyemail", ModifyEmail)          //自认证修改邮箱
	http.HandleFunc("/submodifyemail", SubModifyEmail)    //提交修改邮箱申请
	http.HandleFunc("/user/userinf", Userinf)             //个人中心
	http.HandleFunc("/user/userinfmodify", UserinfModify) //修改个人中心基本信息

	http.HandleFunc("/txt/add", AddTxt)                  //增加记事
	http.HandleFunc("/txt/delete", DeleteTxt)            //删除记事
	http.HandleFunc("/txt/modify", ModifyTxt)            //修改记事
	http.HandleFunc("/txt/find", FindTxt)                //查看记事
	http.HandleFunc("/alltxt", Showalltxt)               //博客
	http.HandleFunc("/showcollection", ShowCollection)   //展示我收藏的记事
	http.HandleFunc("/collection", Collection)           //收藏记事
	http.HandleFunc("/judgecollection", JudgeCollection) //判断是否收藏过记事

	http.HandleFunc("/cancelcollection", CancelCollection) //取消收藏记事
	http.HandleFunc("/like", Like)                         //点赞
	http.HandleFunc("/cancellike", CancelLike)             //取消点赞
	http.HandleFunc("/search", Search)                     //搜索

	http.HandleFunc("/examineuserregist", ShowExamregist)              //展示后台正在审核用户
	http.HandleFunc("/allregistuser", AllregistUser)                   //后台显示所有已经注册过的用户信息
	http.HandleFunc("/user/regist/sucess", RegistSucess)               //同意用户注册申请
	http.HandleFunc("/user/regist/fail", RegistFail)                   //驳回用户注册申请
	http.HandleFunc("/superlogin", SuperLogin)                         //后台管理员登录
	http.HandleFunc("/superregist", Superregist)                       //后台管理员注册
	http.HandleFunc("/user/cancellation", UserCellation)               //用户注销
	http.HandleFunc("/cancellation", Cellation)                        //后台展示用户注销
	http.HandleFunc("/cancellationsucess", CellationSucess)            //同意用户注销
	http.HandleFunc("/cancellationfail", CellationFail)                //驳回用户注销
	http.HandleFunc("/supermodifyemail", ModifyEmail)                  //管理员审核修改邮箱
	http.HandleFunc("/checkmodifyemail", CheckModifyEmail)             //后台展示待审核修改邮箱
	http.HandleFunc("/checkmodifyemailsucess", CheckModifyEmailSucess) //后台同意待审核修改邮箱
	http.HandleFunc("/checkmodifyemailfail", CheckModifyEmailFail)     //后台驳回待审核修改邮箱

	http.HandleFunc("/upload", handleUploadFile)
	http.HandleFunc("/uploads/", showhandleUploadFile)

}
