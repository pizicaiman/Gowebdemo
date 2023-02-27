package service

import (
	"awesomeProject/biz/model/dto"
	"awesomeProject/biz/model/entity"
	"awesomeProject/dao"
	"awesomeProject/util"
	"database/sql"
	"math"
	"strconv"
	"strings"
)

type UserAccountService interface {
	// InsertAccount 新增账户信息
	InsertAccount(accountReq dto.CreateAccountRequest) (*entity.UserAccount, error)
	// FindAccount 查找账户信息
	FindAccount(userId string) (any, error)
	// FindAccountFlow 根据账户ID查找账户
	FindAccountFlow(accountId string) ([]entity.UserAccountFlow, error)
	//RewardAccount 扫码奖励
	RewardAccount(scanning dto.Scanning) (any, error)
}

func NewAccountService() UserAccountService {
	return &accountService{
		AccountRepo:           dao.NewAccountRepository(),
		AccountFlowRepo:       dao.NewAccountFlowRepository(),
		OfflineActivitiesRepo: dao.NewOfflineActivitiesRepo(),
	}
}

type accountService struct {
	AccountRepo           dao.UserAccountRepository
	AccountFlowRepo       dao.UserAccountFlowRepository
	OfflineActivitiesRepo dao.OfflineActivitiesRepository
}

// InsertAccount 新增账户信息
func (me *accountService) InsertAccount(accountDto dto.CreateAccountRequest) (*entity.UserAccount, error) {
	err, b := me.AccountRepo.ExistsByUserId(accountDto.UserId)
	if err != nil {
		return nil, err
	}
	if b {
		return nil, err
	}
	node, _ := util.NewWorker(10)
	account := entity.UserAccount{
		UserId:            accountDto.UserId,
		Id:                "AC" + strconv.Itoa(int(node.GetId())),
		BlockchainAddress: accountDto.Address,
	}
	err = me.AccountRepo.Insert(account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// FindAccount 查找账户信息
func (me *accountService) FindAccount(userId string) (any, error) {
	account, err := me.AccountRepo.FindById(userId)
	if err == sql.ErrNoRows {
		accDto := dto.CreateAccountRequest{}
		accDto.UserId = userId
		address, err := indiev_wallet.GetAddress()
		if err != nil {
			return nil, err
		}
		accDto.Address = address
		accountEntity, err := me.InsertAccount(accDto)
		if err != nil {
			return nil, err
		}
		return accountEntity, nil
	}
	if account.BlockchainAddress != "" {
		chainBalance := util.GetBalance(account.BlockchainAddress)
		chainBalance = strings.TrimLeft(chainBalance, "0x")
		f, err := strconv.ParseUint(chainBalance, 16, 32)
		if err != nil {
			return nil, err
		}
		//金额转换
		account.Balance = int64(float64(f) / math.Pow(10, 8))

	}
	return account, err

}

// FindAccountFlow 根据账户ID查找账户
func (me *accountService) FindAccountFlow(accountId string) ([]entity.UserAccountFlow, error) {
	accountFlows, err := me.AccountFlowRepo.ListByAccountId(accountId)
	if err != nil {
		return nil, err
	}
	return *accountFlows, nil
}

type scanningDto struct {
	RewardValue int64  `json:"rewardValue" db:"reward_value"`
	Unit        string `json:"unit" db:"reward_value"`
}
