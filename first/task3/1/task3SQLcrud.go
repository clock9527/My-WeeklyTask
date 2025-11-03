package main

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
	1 SQL语句练习
	1-1 基本CRUD操作
	1-2 事务语句
*/
// 1-1
type Student struct {
	ID    uint
	Name  string
	Age   uint
	Grade string
}

func execStudent(dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	student := Student{Name: "张三", Age: 20, Grade: "三年级"}
	// 创建表
	db.AutoMigrate(&student)

	ctx := context.Background()
	// 插入数据
	// err = gorm.G[Student](db).Create(ctx, &student)

	// 查询年龄大于18岁-1
	students, err := gorm.G[Student](db).Where("age > ?", 18).Find(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("年龄大于18岁-1:", students)
	// 查询年龄大于18岁-2
	students2 := db.Debug().Where("age > ?", 18).Find(&student)
	fmt.Println("年龄大于18岁-2:", students2, student)

	// 查询年龄大于18岁-3
	var students3 []Student
	d := db.Debug().Find(&students3, "age > ?", 18)
	fmt.Println("年龄大于18岁-3:", d, students3)

	// 将名称“张三” 的年级修改为“四年级”-1
	rowsAffected, err := gorm.G[Student](db.Debug()).Where("name = ?", "张三").Update(ctx, "grade", "四年级")
	if err != nil {
		panic(err)
	}
	fmt.Println("改为“四年级”-1:", rowsAffected)

	// 将名称“张三” 的年级修改为“四年级”-2
	// d2 := db.Debug().Model(&Student{}).Where("name = ?", "张三").Update("grade", "五年级")
	d2 := db.Debug().Model(&Student{}).Where("name = ?", "张三").Updates(map[string]interface{}{"grade": "五年级"})
	fmt.Println(d2)

	// 删除年龄小于15岁
	// db.Create([]Student{
	// 	Student{Name: "李四", Age: 15, Grade: "三年级"},
	// 	Student{Name: "麻子", Age: 14, Grade: "三年级"},
	// 	Student{Name: "王二", Age: 15, Grade: "三年级"},
	// })
	stds, err := gorm.G[Student](db.Debug()).Find(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(stds)
	// 删除年龄小于15岁-1
	rowsAffected2, err := gorm.G[Student](db.Debug()).Where("age < ?", 15).Delete(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(rowsAffected2)

	// 删除年龄小于15岁-2
	// db.Debug().Where("age < ?", 16).Delete(&Student{})
	d3 := db.Debug().Delete(&Student{}, "age < ?", 16)
	fmt.Println(d3)
}

func main() {
	var dsn = "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	execStudent(dsn)
}
