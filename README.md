# gostudy
golang学习 golang Mysql数据库操作类

### 1、打开 config.go 设置对应的数据库连接信息


### 2、使用示例
```
func main(){
	DB := cylib.DbHandle

	//1、获取单条数据
	books := DB.Table("cy_books").Where(map[string]interface{}{"status":"1"}).Find()

    //2、其他查询
    books := DB.Table("cy_books").WhereRaw("status=1").Select()
    
	fmt.Println(books)
}
```

### 开启你的DB操作之旅 ^_^