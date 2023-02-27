package dao

import (
	"awesomeProject/biz/model/entity"
	"database/sql"
	"fmt"
	"strings"
)

type UserAccountRepository interface {
	Insert(account entity.UserAccount) error
	Update(account entity.UserAccount) error
	DeleteById(id string) (bool, error)
	FindById(userId string) (entity.UserAccount, error)
	ListAll() ([]entity.UserAccount, error)
	ExistsByUserId(userId string) (error, bool)
}
type UserAccountFlowRepository interface {
	Insert(accountFlow entity.UserAccountFlow) (int, error)
	DeleteById(id string) (bool, error)
	FindById(id string) (entity.UserAccountFlow, error)
	ListAll() ([]entity.UserAccountFlow, error)
	ListByAccountId(accountId string) (*[]entity.UserAccountFlow, error)
	FindByActivities(id string, id2 int64) (*[]entity.UserAccountFlow, error)
}

type UserNewOfflineActivitiesRepo interface {
	Insert(accountFlow entity.UserAccountFlow) (int, error)
	DeleteById(id string) (bool, error)
	FindById(id string) (entity.UserAccountFlow, error)
	ListAll() ([]entity.UserAccountFlow, error)
	ListByAccountId(accountId string) (*[]entity.UserAccountFlow, error)
	FindByActivities(id string, id2 int64) (*[]entity.UserAccountFlow, error)
}

type UserAccountFlowRepository interface {
	Insert(accountFlow entity.UserAccountFlow) (int, error)
	DeleteById(id string) (bool, error)
	FindById(id string) (entity.UserAccountFlow, error)
	ListAll() ([]entity.UserAccountFlow, error)
	ListByAccountId(accountId string) (*[]entity.UserAccountFlow, error)
	FindByActivities(id string, id2 int64) (*[]entity.UserAccountFlow, error)
}

type accountRepo struct {
}

type accountFlowRepo struct {
}

// ===========================新增===================================================
type newactivitiesRepo struct {
}

type offlineactivitiesRepo struct {
}

// ===========================隔离===================================================

// NewAccountRepository 申明实现了 UserAccountRepository 接口
func NewAccountRepository() UserAccountRepository {
	return &accountRepo{}
}

// Insert  < 具体实现 NewAccountRepository
func (u accountRepo) Insert(account entity.UserAccount) error {
	sql := `insert into user_account(id,user_id,blockchain_address,balance)values(:id,:user_id,:blockchain_address,:balance)`
	result, err := BaseDb.NamedExec(sql, &account)
	if err != nil {
		return err
	}
	fmt.Println(result)

	return nil

}
func (u accountRepo) ExistsByUserId(userId string) (error, bool) {
	sqlStr := `select * from user_account where user_id=?`
	var account entity.UserAccount
	err := BaseDb.Get(&account, sqlStr, userId)
	if err == sql.ErrNoRows {
		return nil, false
	}
	if err != nil {
		return err, false
	}
	if account == (entity.UserAccount{}) {
		return nil, true
	}

	return nil, false

}

// Update  < 具体实现 NewAccountRepository
func (u accountRepo) Update(account entity.UserAccount) error {

	sqlColum := `update user_account set `
	if account.Balance >= -1 {
		sqlColum += ` balance =:balance ,`
	}
	if account.BlockchainAddress != "" {
		sqlColum += ` blockchain_address =:blockchain_address ,`
	}
	sqlColum = strings.TrimSuffix(sqlColum, ",")
	sqlColum += " where id =:id "

	_, err := BaseDb.NamedExec(sqlColum, &account)
	return err

}

// DeleteById  < 具体实现 NewAccountRepository
func (u accountRepo) DeleteById(id string) (bool, error) {
	sql := "update user_account set is_delete=1 where id =?"
	_, err := BaseDb.Exec(sql, id)
	if err != nil {
		return false, err
	}
	return true, err
}

// FindById  < 具体实现 NewAccountRepository
func (u accountRepo) FindById(userId string) (entity.UserAccount, error) {
	sql := "select * from user_account where user_id=?"
	var userAccount entity.UserAccount
	err := BaseDb.Get(&userAccount, sql, userId)
	if err != nil {
		return userAccount, err
	}
	return userAccount, err

}

// ListAll  < 具体实现 NewAccountRepository
func (u accountRepo) ListAll() ([]entity.UserAccount, error) {
	sql := "select * from user_account"
	var accountFlow []entity.UserAccount
	err := BaseDb.Get(&accountFlow, sql)

	return accountFlow, err
}

func NewAccountFlowRepository() UserAccountFlowRepository {
	return &accountFlowRepo{}
}

// Insert  < 具体实现 NewAccountFlowRepository
func (accountFlowRepo) Insert(account entity.UserAccountFlow) (int, error) {
	sql := `insert into user_account_flow(user_id,account_id,type,amount,blockchain_address,before_balance,balance,activities_id,transaction_info)
     values(:user_id,:account_id,:type,:amount,:blockchain_address,:before_balance,:balance,:activities_id,:transaction_info)`
	result, err := BaseDb.NamedExec(sql, account)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	account.Id = id
	return int(id), nil

}

// DeleteById  < 具体实现 NewAccountFlowRepository
func (accountFlowRepo) DeleteById(id string) (bool, error) {
	sql := "update user_account_flow set is_delete=1 where id =?"
	_, err := BaseDb.Exec(sql, id)
	if err != nil {
		return false, err
	}
	return true, err
}

// FindById  < 具体实现 NewAccountFlowRepository
func (accountFlowRepo) FindById(id string) (entity.UserAccountFlow, error) {
	sql := "select * from user_account_flow where id=?"
	var userAccount entity.UserAccountFlow
	err := BaseDb.Get(&userAccount, sql, id)
	if err != nil {
		return userAccount, err
	}
	return userAccount, err

}

// ListAll  < 具体实现 NewAccountFlowRepository
func (accountFlowRepo) ListAll() ([]entity.UserAccountFlow, error) {
	sql := "select * from user_account_flow where is_delete =0"
	var accountFlow []entity.UserAccountFlow
	err := BaseDb.Get(&accountFlow, sql)

	return accountFlow, err
}

// ListByAccountId  < 具体实现 NewAccountFlowRepository
func (accountFlowRepo) ListByAccountId(accountId string) (*[]entity.UserAccountFlow, error) {
	sqlStr := "select * from user_account_flow where is_delete =0 and account_id =?"
	var accountFlow []entity.UserAccountFlow
	err := BaseDb.Select(&accountFlow, sqlStr, accountId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &accountFlow, err
}
func (accountFlowRepo) FindByActivities(accountId string, activitiesId int64) (*[]entity.UserAccountFlow, error) {
	sqlStr := "select * from user_account_flow where is_delete =0 and account_id =? and activities_id=?"
	var accountFlow []entity.UserAccountFlow
	err := BaseDb.Select(&accountFlow, sqlStr, accountId, activitiesId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &accountFlow, err
}

// NewOfflineActivitiesRepo 申明实现了 NewOfflineActivitiesRepo接口
func NewOfflineActivitiesRepo() UserNewOfflineActivitiesRepo {
	return &newactivitiesRepo{}
}

// Insert  < 具体实现 NewOfflineActivitiesRepo
func (newactivitiesRepo) Insert(account entity.UserNewOfflineActivities) (int, error) {
	sql := `insert into user_account_flow(user_id,account_id,type,amount,blockchain_address,before_balance,balance,activities_id,transaction_info)
     values(:user_id,:account_id,:type,:amount,:blockchain_address,:before_balance,:balance,:activities_id,:transaction_info)`
	result, err := BaseDb.NamedExec(sql, account)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil

}
