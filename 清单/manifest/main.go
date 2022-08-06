package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

var (
	DB *gorm.DB
)
type Todo struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Status bool `json:"status"`
}
func initMysql()(err error){
	dsn:="root:123456@tcp(127.0.0.1:3306)/manifest?charset-utf8mb4&parseTime=True&loc=Local"
	DB,err=gorm.Open("mysql",dsn)
	if err!=nil{
		return
	}
	return DB.DB().Ping()
}
func main() {
	//数据库
	err:=initMysql()
	if err!=nil{
		panic(err)
	}
	defer DB.Close()
	DB.AutoMigrate(&Todo{})

	r:=gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static","static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK,"index.html",nil)
	})

	v1Group:=r.Group("v1")
	{
		//待办事项
		//添加
		v1Group.POST("/todo", func(c *gin.Context) {
			//请求数据
			var todo Todo
			c.BindJSON(&todo)
			//	存入数据库 响应
			err =DB.Create(&todo).Error
			if err!=nil{
				c.JSON(http.StatusOK,gin.H{
					"error":err.Error(),
				})
			}else{
				c.JSON(http.StatusOK,todo)
			}
		})
		//查看所有
		v1Group.GET("/todo", func(c *gin.Context) {
			var todoList []Todo
			err:=DB.Find(&todoList).Error
			if err!=nil{
				c.JSON(http.StatusOK,err.Error())
			}else {
				c.JSON(http.StatusOK,todoList)
			}
		})
		//修改
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id,ok:=c.Params.Get("id")
			if !ok{
				c.JSON(http.StatusOK,gin.H{"error":"id不存在"})
				return
			}
			var todo Todo
			err:=DB.Where("id=?",id).First(&todo).Error
			if err!=nil{
				c.JSON(http.StatusOK,err)
				return
			}
			c.BindJSON(&todo)
			if err=DB.Save(&todo).Error;err!=nil{
				c.JSON(http.StatusOK,err)
			}else {
				c.JSON(http.StatusOK,todo)
			}

		})
		//删除
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id,ok:=c.Params.Get("id")
			if !ok{
				c.JSON(http.StatusOK,gin.H{"error":"id不存在"})
				return
			}
			if err=DB.Where("id=?",id).Delete(Todo{}).Error;err!=nil{
				c.JSON(http.StatusOK,err)
			}else{
				c.JSON(http.StatusOK,gin.H{id:"deleted"})
			}
		})
	}
	r.Run()
}
