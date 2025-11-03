package main

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Account struct {
	Id      uint
	Balance float64
}

type Transaction struct {
	Id              uint
	From_account_id uint
	To_accout_id    uint
	Amount          float64
}

// 数据表初始化
func InitData(dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Transaction{})
	a := Account{Id: 1, Balance: 300}
	b := Account{Id: 2}
	db.Create([]Account{
		a,
		b,
	})
}

// 交易
func (t *Transaction) TransAToB1(db *gorm.DB, a, b *Account, amount float64) bool {
	ctx := context.Background()
	db.Transaction(func(tx *gorm.DB) error {
		account, err := gorm.G[Account](tx).Select("balance").Where("id = ?", a.Id).Take(ctx)
		// 校验A的余额
		if err != nil || account.Balance < amount {
			return errors.New("余额不足!!")
		}
		// 更新 A和B 的balance
		tx.Model(&a).Update("balance", gorm.Expr("balance - ?", amount))
		tx.Model(&b).Update("balance", gorm.Expr("balance + ?", "abc"))
		// 添加记录到Transaction中
		trans := Transaction{From_account_id: a.Id, To_accout_id: b.Id, Amount: amount}
		tx.Create(&trans)
		return nil
	})
	fmt.Println("处理完成...")
	return false
}

func main() {
	var dsn = "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	// InitData(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// a := Account{Id: 1}
	b := Account{Id: 2}
	result := db.Debug().Model(&b).Update("balance", gorm.Expr("balance + ?", "#@"))
	fmt.Println(result.Error)
	// trans := Transaction{}
	// trans.TransAToB1(db, &a, &b, 200.00)

}
