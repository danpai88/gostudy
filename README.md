# gostudy
golang学习 Mysql数据库操作类

### 1、打开 cylib/cydb.go 文件中的 Connect 方法，设置对应的数据库连接信息


### 2、使用示例
```
	DB := cylib.DbHandle
	books := DB.Table("cy_books").Where(map[string]string{"status":"1"}).Find()

```

### 开启你的DB操作之旅 ^_^