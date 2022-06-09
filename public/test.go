package main

import (
	"github.com/RaymondCode/simple-demo/controller"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//users := database.FindAllUser()
	//for i := 0; i < len(users); i++ {
	//	if users[i].Username == "将进酒" {
	//		fmt.Print("查找成功")
	//		break
	//	}
	//}
	print(controller.GetLastId())
}
