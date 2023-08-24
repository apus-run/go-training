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

type UserHandler struct {
	svc service.UserService

	log *log.Logger
}

func NewUserHandler(svc service.UserService, logger *log.Logger) *UserHandler {
	return &UserHandler{
		svc: svc,
		log: logger,
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
		Birthday: user.Birthday,
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
	if !isValidDate(req.Birthday) {
		ctx.JSONE(http.StatusBadRequest, "生日格式错误", nil)
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
		Birthday: req.Birthday,
		Profile:  req.Profile,
	})

	if err != nil {
		ctx.JSONE(http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSONOK("用户信息更新成功", nil)
}

func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
