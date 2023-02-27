package controller

import (
	"github.com/gin-gonic/gin"
	//项目
	"awesomeProject/biz/model/dto"
	"awesomeProject/biz/service"
)

func NewAccountController(accountService service.UserAccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

// InsertAccount 新增账户信息
func (me *AccountController) InsertAccount(c *gin.Context) (any, error) {
	var req dto.CreateAccountRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, err
	}

	account, err := me.accountService.InsertAccount(req)
	if err != nil {
		return nil, err
	}
	return account, err
}

// FindAccount 查找账户信息
func (me *AccountController) FindAccount(c *gin.Context) (any, error) {
	userId := c.Param("userId")
	account, err := me.accountService.FindAccount(userId)
	if err != nil {
		return nil, err
	}
	return account, err

}

// FindAccountFlow 根据账户ID查找账户
func (me *AccountController) FindAccountFlow(c *gin.Context) (any, error) {
	accountId := c.Param("accountId")
	account, err := me.accountService.FindAccountFlow(accountId)
	if err != nil {
		return nil, err
	}
	return account, err
}

type AccountController struct {
	accountService service.UserAccountService
}
