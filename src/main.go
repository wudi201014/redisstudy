package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"math/rand"
	"strconv"
	"time"
)


func main(){
	conn,err := redis.Dial("tcp","127.0.0.1:6379")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

/*
	//程序还有问题，最后的总量还是扣了2w，只是中间显示会重复，那么可不可能是只是同时读取？
	//因为decr和get不是同时执行的，不是原子操作
	for i :=0;i<10000;i++{
		conn.Do("decr","stock")
		conn.Flush()
		stock,err := conn.Do("get","stock")
		if err != nil {
			continue
		}
		log.Println("stock : ", string(stock.([]byte)))
	}*/


	//arttimeset   ID - 创建时间
	//artscoreset  ID - Score
	//artuserset   ID - User
	//vote := 432

	//随机生成10个文章并加入时间戳
	arts := []string{}
	for i:=0;i<10;i++{
		arts = append(arts,fmt.Sprintf("art%d",i))
		c,b := conn.Do("zadd","arttimeset",time.Now().Unix(),arts[i])
		log.Println(c,b)
		a := rand.Int31n(10)
		log.Println("sleep" ,a)
		time.Sleep(time.Duration(a)*time.Second)
	}


	//模拟100次用户投票,随机给1-10投票
	for i:=0;i<100;i++{
		user := fmt.Sprintf("user%d",rand.Int31n(30))
		art2v := rand.Int31n(10)
		art2vote :=fmt.Sprintf("artuser%d",art2v)

		a,b := conn.Do("sadd",art2vote,user)
		log.Println(a,b)
		if a.(int64)==1{
			//成功,加分
			a,b = conn.Do("zscore","arttimeset",fmt.Sprintf("art%d",art2v))
			if a==nil{
				//不存在
				continue
			}
			t,_ := strconv.ParseInt(string(a.([]uint8)),10,64)
			if t<time.Now().Add(time.Second*-20).Unix(){
				a,b = conn.Do("zincrby","artscoreset",vote,fmt.Sprintf("art%d",art2v))
				a,b = conn.Do("hincrby","artnumvote",fmt.Sprintf("art%d",art2v),1)
				log.Println(user,"加分",art2v)
			}
		}else {
			log.Printf("%s 已经投过",user)

		}
	}


	//按排名获取文章
	a,_ := conn.Do("zrevrange","artscoreset",0,-1,"withscores")
	for i,name := range a.([]interface{}){
		log.Println(i,string(name.([]uint8)))
	}

}