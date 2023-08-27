package handler

import (
	"net/http"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/pkg/errors"

	"project-layout/internal/domain/entity"
	"project-layout/internal/service"
	"project-layout/internal/web/dto"
	"project-layout/pkg/ginx"
	"project-layout/pkg/jwtx"
	"project-layout/pkg/log"
)

var (
	bizLoginType = "user:login"
)

type UserHandler struct {
	svc     service.UserService
	codeSvc service.CodeService

	log *log.Logger
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, logger *log.Logger) *UserHandler {
	return &UserHandler{
		svc:     svc,
		codeSvc: codeSvc,
		log:     logger,
	}
}

func (h *UserHandler) Login(ctx *ginx.Context) {
	var req dto.LoginReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
		return
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.JSONE(http.StatusOK, "账号或者密码不正确，请重试", nil)
		return
	}

	// 测试使用1分钟
	// expireAt := time.Now().Add(time.Minute)
	// 正常设置为30分钟，要将过期时间设置更长一些
	expireAt := time.Now().Add(time.Minute * 30)
	token, err := jwtx.GenerateToken(
		jwtx.WithUserAgent(ctx.Request.UserAgent()),
		jwtx.WithSecretKey(jwtx.SecretKey),
		jwtx.WithUserId(user.ID),
		jwtx.WithExpireAt(expireAt),
	)

	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
		return
	}

	// 将token放入响应头中
	ctx.Header("x-jwt-token", token)

	ctx.JSONOK("登录成功", dto.LoginResp{
		ID:   user.ID,
		Name: user.Name,
	})
}

func (h *UserHandler) LoginSMS(ctx *ginx.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := h.codeSvc.Verify(ctx, bizLoginType, req.Phone, req.Code)
	if err != nil {
		ctx.JSONE(5, "系统异常", nil)
		return
	}
	if !ok {
		ctx.JSONE(4, "验证码错误", nil)
		return
	}

	// 验证码是对的
	// 登录或者注册用户
	user, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSONE(4, "系统错误", nil)
		return
	}
	if err == service.ErrInvalidUserOrPassword {
		ctx.JSONE(http.StatusOK, "账号或者密码不正确，请重试", nil)
		return
	}

	// 测试使用1分钟
	// expireAt := time.Now().Add(time.Minute)
	// 正常设置为30分钟，要将过期时间设置更长一些
	expireAt := time.Now().Add(time.Minute * 30)
	token, err := jwtx.GenerateToken(
		jwtx.WithUserAgent(ctx.GetHeader("User-Agent")),
		jwtx.WithSecretKey(jwtx.SecretKey),
		jwtx.WithUserId(user.ID),
		jwtx.WithExpireAt(expireAt),
	)

	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
		return
	}
	// 将token放入响应头中
	ctx.Header("x-jwt-token", token)

	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSONOK("登录成功", dto.LoginResp{
		ID:   user.ID,
		Name: user.Name,
	})
}

// SendSMSLoginCode 发送短信验证码
func (h *UserHandler) SendSMSLoginCode(ctx *ginx.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 你也可以用正则表达式校验是不是合法的手机号
	if req.Phone == "" {
		ctx.JSONE(http.StatusBadRequest, "请输入手机号码", nil)
		return
	}
	err := h.codeSvc.Send(ctx, bizLoginType, req.Phone)
	switch err {
	case nil:
		ctx.JSONOK("发送成功", nil)
	case service.ErrCodeSendTooMany:
		ctx.JSONE(http.StatusBadRequest, "短信发送太频繁，请稍后再试", nil)
	default:
		ctx.JSONE(http.StatusBadRequest, "系统错误", nil)

		// 要打印日志
		h.log.Errorf("发送短信验证码失败: %v", err)

		return
	}
}

func (h *UserHandler) Register(ctx *ginx.Context) {
	var req dto.RegisterReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(http.StatusOK, err.Error(), nil)
	}

	if req.Password != req.ConfirmPassword {
		ctx.JSONE(http.StatusOK, "两次输入的密码不一致", nil)
		return
	}

	isPassword, err := regexp.
		MustCompile(`^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`, regexp.None).
		MatchString(req.Password)
	if !isPassword || err != nil {
		ctx.JSONE(http.StatusOK, "密码必须包含数字、特殊字符，并且长度不能小于 8 位", nil)
		return
	}

	isEmail, err := regexp.
		MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, regexp.None).
		MatchString(req.Email)

	if !isEmail || err != nil {
		ctx.JSONE(http.StatusOK, "邮箱不合法", nil)
		return
	}

	_, err = h.svc.Register(
		ctx.Request.Context(),
		entity.User{
			Name:     req.Name,
			Phone:    req.Phone,
			Email:    req.Email,
			Password: req.ConfirmPassword,
		},
	)

	if errors.As(err, service.ErrUserDuplicateEmailOrPhone) {
		ctx.JSONE(http.StatusOK, "邮箱或者手机号已经存在", nil)
		return
	}

	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSONOK("注册成功", nil)
}

func (h *UserHandler) Profile(ctx *ginx.Context) {
	cs := ctx.MustGet("claims").(*jwtx.CustomClaims)

	// h.log.Infof("claims: %v", cs)

	user, err := h.svc.Profile(ctx, cs.UserID)

	if err != nil {
		ctx.JSONE(http.StatusBadRequest, "用户信息没找到", nil)
		return
	}

	ctx.JSONOK("OK", dto.UserInfoResp{
		Email:    user.Email,
		Phone:    user.Phone,
		Gender:   user.Gender,
		NickName: user.NickName,
		RealName: user.RealName,
		Birthday: user.Birthday.Format(time.DateOnly),
		Profile:  user.Profile,
	})
}

func (h *UserHandler) UpdateProfile(ctx *ginx.Context) {
	var req dto.UpdateProfileReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
		return
	}

	if req.Gender != 0 && req.Gender != 1 && req.Gender != 2 {
		ctx.JSONE(http.StatusBadRequest, "性别参数错误", nil)
		return
	}

	if len(req.NickName) < 2 || len(req.NickName) > 20 {
		ctx.JSONE(http.StatusBadRequest, "昵称为2 ~ 20个字符", nil)
		return
	}

	if req.RealName == "" {
		ctx.JSONE(http.StatusBadRequest, "真实姓名不能为空", nil)
		return
	}
	if len(req.RealName) < 2 || len(req.RealName) > 20 {
		ctx.JSONE(http.StatusBadRequest, "真实姓名为2 ~ 20个字符", nil)
		return
	}

	if req.Birthday == "" {
		ctx.JSONE(http.StatusBadRequest, "生日不能为空", nil)
		return
	}

	// TODO: 其实没有必要直接校验具体的格式 而是应该校验日期的有效性
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.JSONE(http.StatusBadRequest, "日期格式错误", nil)
		return
	}

	if len(req.Profile) > 200 {
		ctx.JSONE(http.StatusBadRequest, "个人简介不能超过200个字符", nil)
		return
	}

	cs := ctx.MustGet("claims").(*jwtx.CustomClaims)

	err = h.svc.UpdateProfile(ctx, entity.User{
		ID: cs.UserID,

		Gender:   req.Gender,
		NickName: req.NickName,
		RealName: req.RealName,
		Birthday: birthday,
		Profile:  req.Profile,
	})

	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSONOK("用户信息更新成功", nil)
}
