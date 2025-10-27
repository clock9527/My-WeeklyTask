package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
	1 指针
	1-1 整数指针
	1-2 整数切片的指针
*/
// 1-1整数指针
func intPointer(intP *int) {
	if intP != nil {
		*intP += 10
	}
}

// 1-2整数切片的指针
func sliceIntPointer(intSP *[]int) {
	if len(*intSP) == 0 {
		fmt.Println("切片为空...")
		return
	}

	for i := range *intSP {
		(*intSP)[i] *= 2
	}
}

/*
   2 Goroutine
   2-1 奇数偶数
   2-2 任务调度器
*/
// 2-1 奇数偶数
func numOddEven(n int) {
	// 使用 WaitGroup 等待两个协程完成
	// var wg sync.WaitGroup
	// wg.Add(2) // 等待两个协程
	go func() {
		// defer wg.Done() // 协程结束时通知 WaitGroup
		for i := 1; i <= n; i += 2 {
			fmt.Println("goroutine1", i)
		}
	}()
	go func() {
		// defer wg.Done() // 协程结束时通知 WaitGroup
		for i := 2; i <= n; i += 2 {
			fmt.Println("goroutine2", i)
		}
	}()

	// 等待所有协程完成
	// wg.Wait()
	time.Sleep(time.Second * 10) // 主协程等待一会，让子协程打印完成  或使用 wg.Wait()
	fmt.Println("所有数字打印完成!")
}

// 2-2 任务调度器
type Task func(int) // 将任务定义为一个无参无返回值的函数类型

type TaskScheduler struct { // 定义调度器
	tasks []Task
}

func (tSch TaskScheduler) exec() { // 任务并行执行方法
	for i, task := range tSch.tasks {
		go func() {
			begin := time.Now()
			task(i)
			end := time.Now()
			fmt.Println(i, ":", end.Sub(begin))
		}()
	}
	time.Sleep(time.Second * 2) //主线程等待一会
}

/*
	3 面向对象
	3-1 Shape 接口
	3-2 组合
*/
// 3-1 Shape 接口
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	width  float64
	height float64
}

func (r *Rectangle) Area() float64 {
	return r.width * r.height
}
func (r *Rectangle) Perimeter() float64 {
	return (r.width + r.height) * 2
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}
func (c Circle) Perimeter() float64 {
	return 3.14 * c.Radius * 2
}

// 3-2 组合
type Person struct {
	Name string
	Age  uint
}

type Employee struct {
	Person
	EmployeeID uint
}

func (emp Employee) PrintInfo() {
	fmt.Printf("员工信息：id=%d; name=%s; age=%d\n", emp.EmployeeID, emp.Person.Name, emp.Age)
}

/*
	4 Channel
	4-1 协程之间的通信
	4-2 通道的缓冲机制
*/
// 4-1 协程之间的通信
func goroutineBaseTest(msg int) {
	var wg sync.WaitGroup
	wg.Add(2) // 等待两个协程
	ch := make(chan int)
	go func() {
		defer wg.Done()
		ch <- msg
	}()
	go func() {
		defer wg.Done()
		fmt.Println("接收到了信息：", <-ch)
	}()
	wg.Wait()
	fmt.Println("完成")
}

// 4-2 通道的缓冲机制
func chanProductConsumer(n int) {
	var wg sync.WaitGroup
	wg.Add(2)
	ch := make(chan int, 100)
	go func() {
		fmt.Println("生产...")
		defer wg.Done()
		for i := 0; i < n; i++ {
			ch <- i
		}
		close(ch)
	}()
	go func() {
		// time.Sleep(time.Second * 3)
		fmt.Println("消费...")
		defer wg.Done()
		for c := range ch {
			fmt.Println("接收到消息:", c)
		}
	}()

	wg.Wait()
	fmt.Println("处理完成...")
}

/*
	5 锁机制
	5-1 共享的计数器
	5-2 原子操作
*/
// 5-1 共享的计数器
func lockSafeCounter(n int) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	sum := 0
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			for j := 0; j < n; j++ {
				sum++
			}
			fmt.Printf("协程%d, 完成...\n", gid)
		}(i)
	}
	// time.Sleep(time.Second)
	wg.Wait()
	fmt.Println(sum)
}

// 5-2 原子操作
func automicOpt(n int) {
	var num int32
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < n; j++ {
				atomic.AddInt32(&num, 1)
			}
		}()
	}
	wg.Wait()
	value := atomic.LoadInt32(&num)
	fmt.Println(value)
}
func main() {

	// 1-1
	// n := 10
	// intPointer(&n)
	// fmt.Println("整数指针: ", n)

	// 1-2
	// sliceP := []int{2, 3, 4}
	// sliceIntPointer(&sliceP)
	// fmt.Println("整数切片的指针: ", sliceP)

	// 2-1
	// numOddEven(10)

	// 2-2
	// taskSch := TaskScheduler{[]Task{func(a int) { fmt.Println("任务1-task-ID: ", a) }, func(a int) { fmt.Println("任务2-task-ID: ", a) }}}
	// taskSch.exec()

	// 3-1
	// var s Shape
	// s = &Rectangle{width: 3.5, height: 3.3}
	// fmt.Printf("面积：%.3f;  周长：%.3f\n", s.Area(), s.Perimeter())

	// s = Circle{Radius: 3.0}
	// fmt.Println(s.Area(), "  ", s.Perimeter())

	// 3-2
	// emp := Employee{Person{Name: "Tom", Age: 25}, 1}
	// emp.PrintInfo()

	// 4-1
	// goroutineBaseTest(3)

	// 4-2
	// chanProductConsumer(10)

	// 5-1
	// lockSafeCounter(1000)

	// 5-2
	automicOpt(1000)
}
