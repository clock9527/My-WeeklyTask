package main

import (
	"database/sql/driver"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type Employee struct {
	Id         uint    `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func (emp Employee) Value() (driver.Value, error) {
	return []interface{}{emp.Name, emp.Department, emp.Salary}, nil
}

func InitDB() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True"
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("Connect failed")
		panic(err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	ddl := `create table if not exists employees(
	id bigint NOT NULL AUTO_INCREMENT,
	name varchar(200),
	department varchar(400),
	salary Decimal(17,2),
	Primary Key (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

	r, err := db.Exec(ddl)
	if err != nil {
		panic(err)
	}
	fmt.Println(r.RowsAffected())

	return

}

func InsertData() {
	sqlStr := "Insert Into employees(name,department,salary) values(?,?,?)"
	// r, err := db.Exec(sqlStr, "张三", "营销部", 7000.00)
	// if err != nil {
	// 	panic(err)
	// }
	// i, err := r.LastInsertId()
	// fmt.Println(i)

	// 使用db.NamedExec + 结构体  批量插入数据
	// sqlStr = "Insert Into employees(name,department,salary) values(:name,:department,:salary)"
	// emps := []*Employee{{Name: "李四", Department: "生产部", Salary: 6000.00},
	// 	{Name: "王二", Department: "技术部", Salary: 9000.00}}
	// db.NamedExec(sqlStr, emps)

	// 使用sqlX.In 拼接SQL
	sqlStr = "Insert Into employees(name,department,salary) values(?),(?)"
	emp1 := Employee{Name: "麻子", Department: "人资部", Salary: 5600}
	emp2 := Employee{Name: "张傻", Department: "财资部", Salary: 5600}
	query, args, _ := sqlx.In(sqlStr, []interface{}{emp1, emp2}...)
	fmt.Println(query)
	fmt.Println(db.Rebind(query))
	fmt.Println(args)
	_, err := db.Exec(query, args...)
	if err != nil {
		panic(err)
	}
}

// 使用select 基于department进行查询
func getEmployeesByDeptForSelect(dept string) []Employee {
	emps := []Employee{}
	queryStr := "Select * From employees t Where t.department = ?"
	err := db.Select(&emps, queryStr, dept)
	if err != nil {
		panic(err)
	}
	return emps
}

// 使用NamedQuery+Rows+map 基于department进行查询
func getEmployeesByDeptForRows1(dept string) []Employee {
	args := map[string]interface{}{"dept": dept}
	rows, err := db.NamedQuery("Select * From employees t Where t.department = :dept", args)
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
		return nil
	}
	defer rows.Close()
	emps := RowsDeal[Employee](Employee{}, rows)

	return emps
}

// 使用NamedQuery+Rows+结构体 基于department进行查询
func getEmployeesByDeptForRows2(dept string) []Employee {
	emp0 := Employee{
		Department: dept,
	}
	rows, err := db.NamedQuery("Select * From employees t Where t.department = :department", emp0)
	if err != nil {
		fmt.Printf("db.NamedQuery failed %v\n", err)
	}
	defer rows.Close()
	emps := RowsDeal[Employee](Employee{}, rows)
	return emps
}

// rows处理
func RowsDeal[T Employee](t T, rows *sqlx.Rows) []T {
	var ts []T
	for rows.Next() {
		err := rows.StructScan(&t)
		if err != nil {
			fmt.Printf("scan failed %v\n", err)
			continue
		}
		ts = append(ts, t)
	}
	return ts
}

func MaxSalaryEmployee1() Employee {
	var emp = Employee{}
	queryStr := "Select * From employees t Where t.salary = (Select Max(e.salary) max_salary From employees e)"
	err := db.Get(&emp, queryStr)
	if err != nil {
		fmt.Printf("db.NamedQuery failed %v\n", err)
	}
	return emp
}

// 返回最高薪资对应的员工,返回切片
// Select * From employees t Where t.salary = (Select Max(e.salary) max_salary From employees e)
func MaxSalaryEmployee2() []Employee {
	queryStr := "Select Max(salary) max_salary From employees"
	var maxSalary float64
	err := db.Get(&maxSalary, queryStr)
	if err != nil {
		fmt.Printf("select failed %v\n", err)
		return nil
	}
	emps := []Employee{}
	err = db.Select(&emps, "Select * From employees t Where t.salary = ?", maxSalary)
	if err != nil {
		fmt.Printf("select failed %v \n", err)
	}
	return emps
}

func main() {
	InitDB()
	// InsertData()
	fmt.Println("--查询技术部的员工信息--")
	fmt.Println("方式1: ", getEmployeesByDeptForSelect("技术部"))
	fmt.Println("方式2: ", getEmployeesByDeptForRows1("技术部"))
	fmt.Println("方式3: ", getEmployeesByDeptForRows2("技术部"))

	fmt.Println("--薪资最高的员工信息--")
	fmt.Println("方式1: ", MaxSalaryEmployee1())
	fmt.Println("方式2: ", MaxSalaryEmployee2())

}
