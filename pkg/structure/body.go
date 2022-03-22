package structure

//注册用户的基本信息
type UserRegist struct {
	Email    string `json:"email"`    //邮箱
	Username string `json:"username"` //用户名
	Password string `json:"password"` //密码
	Nickname string `json:"nickname"` //昵称
	Icode    string `json:"icode"`    //验证码
}

//登录用户的基本信息
type UserLogin struct {
	Emailanduser string `json:"emailanduser"` //用户名或邮箱
	Password     string `json:"password"`     //密码
}

//个人中心
type User struct {
	Email    string `json:"email"`    //邮箱
	Username string `json:"username"` //用户名
	Nickname string `json:"nickname"` //昵称
	Sex      string `json:"sex"`      //性别
	Avatar   string `json:"avatar"`   //头像
	Name     string `json:"name"`     //姓名
	Time     string `json:"time"`     //注册时间
}

//用户注销
type UserCellation struct {
	Email  string `json:"email"`  //邮箱
	Reason string `json:"reason"` //理由
	Icode  string `json:"icode"`  //验证码
}

//用户注销
type Cancellation struct {
	Email    string `json:"email"`    //邮箱
	Username string `json:"username"` //用户名
	Reason   string `json:"reason"`   //理由
}

//个人中心
type Usermodify struct {
	Email  string `json:"email"`  //邮箱
	Sex    string `json:"sex"`    //性别
	Avatar string `json:"avatar"` //头像
	Name   string `json:"name"`   //姓名
}

//后台管理员登录
type SuperUserLogin struct {
	Email    string `json:"email"`    //用户名或邮箱
	Password string `json:"password"` //密码
	Icode    string `json:"icode"`    //验证码
}

//后台管理员注册
type SuperUserRegist struct {
	Email    string `json:"email"`    //用户名或邮箱
	Username string `json:"username"` //用户名
	Password string `json:"password"` //密码
	Icode    string `json:"icode"`    //验证码
}

//获取用户邮箱
type UserEmail struct {
	Email string `json:"email"` //邮箱
}

//获取用户邮箱
type EmailandReason struct {
	Email  string `json:"email"` //邮箱
	Reason string `json:"reason"`
}

//获取用户邮箱
type Userun struct {
	Username string `json:"username"` //用户名
	Token    string `json:"token"`
}

//获取用户头像和用户名
type UserUAvatar struct {
	Username string `json:"username"` //用户名
	Avatar   string `json:"avatar"`   //头像
}

//上传用户头像
type UserAvatar struct {
	Email  string `json:"email"`  //邮箱
	Avatar string `json:"avatar"` //头像
}

//前端发来的json格式错误
type JsonError struct {
	Information string `json:"information"`
}

//博客中的每一个记事
type UserLabel struct {
	Username string `json:"username"` //用户名
	Avatar   string `json:"avatar"`   //头像
	Title    string `json:"title"`    //文章标题
	Content  string `json:"content"`  //内容
	Label    int    `json:"label"`    //标签
	Time     string `json:"time"`     //发布时间
	Praise   string `json:"praise"`   //点赞数
	Id       int    `json:"id"`       //记事id
}

//获取记事id
type TxtId struct {
	Id int `json:"id"`
}

//获取用户所有记事
type UserTxt struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Label    int    `json:"label"`
	Time     string `json:"time"`
	Complete int    `json:"complete"`
}

//修改记事
type UserModifyTxt struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Label    int    `json:"label"`
	Complete int    `json:"complete"`
}

//展示审核的用户
type UserExamine struct {
	Email    string `json:"email"`    //邮箱
	Username string `json:"username"` //用户名
	Time     string `json:"time"`
}

//增加记事(时间后端在加上)
type UserAddTxt struct {
	Email   string `json:"email"` //邮箱
	Title   string `json:"title"`
	Content string `json:"content"`
	Label   int    `json:"label"`
}

//删除记事（时间是主键）
type UserDeleteTxt struct {
	Time string `json:"time"` //记事发布时间
}

//修改邮箱
type ModiEmail struct {
	Email    string `json:"email"`    //邮箱
	Newemail string `json:"newemail"` //新邮箱
}

//提交修改邮箱
type SubModiEmail struct {
	Email    string `json:"email"`    //邮箱
	Newemail string `json:"newemail"` //新邮箱
	Reason   string `json:"reason"`
}

//个人中心
type AllregistUser struct {
	Email    string `json:"email"`    //邮箱
	Username string `json:"username"` //用户名
	Time     string `json:"time"`     //注册时间
	Avatar   string `json:"avatar"`   //头像
}

//展示后台审核修改邮箱
type Fmodiyemail struct {
	Email    string `json:"email"`    //邮箱
	Username string `json:"username"` //用户名
	Newemail string `json:"newemail"` //新邮箱
	Reason   string `json:"reason"`
	Time     string `json:"time"` //注册时间
}

//邮箱和验证码
type EmailandIcode struct {
	Email string `json:"email"` //邮箱
	Icode string `json:"icode"` //验证码
}

//忘记密码
type ModifyPassword struct {
	Email       string `json:"email"`    //邮箱
	Password    string `json:"password"` //密码
	Newpassword string `json:"newpassword"`
}

//博客
type Blog struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Praise int    `json:"praise"`
	Label  int    `json:"label"`
}

//收藏记事
type Blogs struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Praise  int    `json:"praise"`
	Content string `json:"content"`
	Label  int    `json:"label"`
}

//收藏
type Collection struct {
	Email string `json:"email"` //邮箱
	Id    int    `json:"id"`
}

//搜索关键字
type Keyword struct {
	Keryword string `json:"keyword"`
}

//发送响应状态码和响应信息
type Resp struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

//发送响应状态码和响应信息（显示信息）
type Respo struct {
	Code   int         `json:"code"`
	Msg    interface{} `json:"msg"`
	Avatar string      `json:"avatar"`
	Email  string      `json:"email"`
}

//发送响应状态码和响应信息（显示我的记事本）
type Res struct {
	Code     int         `json:"code"`
	Username string      `json:"username"`
	Avatar   string      `json:"avatar"`
	Msg      interface{} `json:"msg"`
}

//发送响应状态码和响应信息和邮箱（登录）
type Respon struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Email string      `json:"email"`
	Token string      `json:"token"`
}

//发送响应状态码和响应信息和邮箱（登录）
type Respons struct {
	Code     int         `json:"code"`
	Msg      interface{} `json:"msg"`
	Email    string      `json:"email"`
	Username string      `json:"username"`
	Token    string      `json:"token"`
}

//免密登录
type Respp struct {
	Code     int         `json:"code"`
	Msg      interface{} `json:"msg"`
	Email    string      `json:"email"`
	Username string      `json:"username"`
	Token    string      `json:"token"`
	Avatar   string      `json:"avatar"`
}
