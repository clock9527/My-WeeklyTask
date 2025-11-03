package main

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 用户
type User struct {
	ID       uint
	Username string
	Password string
	Email    string
	PostNums uint
	Posts    []Post
}

// 文章
type Post struct {
	gorm.Model
	Title   string
	Content string
	UserID  uint
	// CommNums uint
	CommStat string
	Comments []Comment
}

// 评论
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	PostID  uint   `gorm:"not null;index"`
}

// 查询某个用户发布的所有文章及其对应的评论信息
func ShowPostComments(db *gorm.DB, username string) {
	ctx := context.Background()
	user, err := gorm.G[User](db.Debug()).Where("username=?", "Tom").Preload("Posts", nil).Preload("Posts.Comments", nil).Find(ctx)
	// user := User{}
	// err := db.Preload(clause.Associations).Find(&user).Error
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}

// 查询评论数量最多的文章信息 - 使用 raw sql
func MaxCommentsPosts1(db *gorm.DB) {
	//  最多评论对应的 postid
	var max_post_id []uint
	db.Debug().Raw("Select Max(ac.cnt) max_cnt From (Select post_id,Count(id) as cnt From comments c Group By c.post_id) ac Where 1=?", 1).Scan(&max_post_id)

	// 最多评论对应的 post
	var posts []Post
	db.Debug().Raw("Select p.* From posts p where p.id  In ?", max_post_id).Scan(&posts)
	fmt.Println(posts)
}

// 查询评论数量最多的文章信息 使用 Generics API + Raw SQL
func MaxCommentsPosts2(db *gorm.DB) {
	// count 子查询
	queryCount := gorm.G[Comment](db.Debug()).Select("post_id,count(*) as cnt_comms").Group("post_id").Order("cnt_comms desc").Limit(1)

	// 获取最大评论数
	var max_comm_num uint
	db.Table("(?) a", queryCount).Select("Max(cnt_comms)").Scan(&max_comm_num)

	// 最大评论数对应的文章ID
	var max_post_id []uint
	db.Table("(?) aa Where aa.cnt_comms = ?", queryCount, max_comm_num).Select("post_id").Scan(&max_post_id)

	// 基于最大评论数进行处理,获取相应信息
	type Result struct {
		Username string
		Email    string
		Title    string
		Content  string
		Comms    string
	}
	var results []Result
	db.Debug().Raw(`
		Select u.username,u.email,p.title,p.content,c.content as comms 
		  From users u
		  Left Join posts p
		    On u.id = p.user_id
		 Inner Join comments c
		    On p.id = c.post_id
		 Where p.id In ?`, max_post_id).Scan(&results)
	fmt.Println(results)
}

// hook 文章创建时自动更新用户的文章数量
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	var posts int64
	tx.Debug().Model(p).Where("user_id = ?", p.UserID).Count(&posts)
	fmt.Println(posts)
	user := User{ID: p.UserID}
	tx.Debug().Model(&user).Update("post_nums", posts)
	return nil
}

// hook 评论删除时检查文章的评论数量，如果评论数量为 0,更新文章的评论状态为 "无评论"。
/*
func (comm *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Println("before:", comm)
	var postID uint
	tx.Debug().Model(comm).Where("id = ?", comm.ID).Select("post_id").Scan(&postID)
	fmt.Println("before:", postID)
	comm.Model.ID = comm.ID
	comm.PostID = postID
	return nil
}

func (comm *Comment) AfterDelete(tx *gorm.DB) (err error) {
	var comms int64
	tx.Debug().Model(comm).Where("post_id = ?", comm.PostID).Scan(&comms)
	fmt.Println("after:", comms)
	if comms == 0 {
		tx.Debug().Model(&Post{}).Where("ID = ?", comm.PostID).Update("comm_stat", "无评论")
	}
	return
}
*/
func (comm *Comment) AfterDelete(tx *gorm.DB) (err error) {
	var posts []uint
	tx.Debug().Table("posts p").
		Select("p.id").
		Joins("Left Join comments c On p.id = c.post_id And c.deleted_at Is Null").
		Group("p.id").Having("Count(c.id) = ?", 0).
		Scan(&posts)

	if len(posts) > 0 {
		tx.Debug().Model(&Post{}).Where("id In ?", posts).Update("comm_stat", "无评论")
	}
	return nil
}

func main() {
	var dsn = "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
	// users := []*User{
	// 	{ID: 1, Username: "Jim", Password: "123", Email: "123@q.com"},
	// 	{ID: 2, Username: "Tom", Password: "123", Email: "153@q.com"},
	// 	{ID: 3, Username: "Lily", Password: "123", Email: "133@q.com"},
	// 	{ID: 4, Username: "Jhon", Password: "123", Email: "143@q.com"},
	// }
	// db.Create(users)

	// posts := []*Post{
	// 	{Title: "HA", Content: "abcabccc", UserID: 1},
	// 	{Title: "LA", Content: "bbccddses", UserID: 1},
	// 	{Title: "MA", Content: "abddddsss", UserID: 2},
	// 	{Title: "DA", Content: "hhjjyyrraa", UserID: 3},
	// }
	// db.Create(posts)

	// Comment{Model: gorm.Model{ID: 3}}
	// comments := []*Comment{
	// 	{Content: "Good", PostID: 1},
	// 	{Content: "OK!!", PostID: 1},
	// 	{Content: "Good33", PostID: 3},
	// 	{Content: "Good44", PostID: 4},
	// }
	// db.Create(comments)

	// ShowPostComments(db, "Tom")
	// MaxCommentsPosts2(db)

	db.Debug().Delete(&Comment{}, 3) // &Comment{},3

}
