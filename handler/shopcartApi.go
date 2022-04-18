package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
	cartPb "github.com/wsjcko/shopcart/protobuf/pb"
	pb "github.com/wsjcko/shopcartApi/protobuf/pb"
	"strconv"
)

type ShopCartApi struct {
	ShopCartService cartPb.ShopCartService
}

// FindAll pb.Call 通过API向外暴露为/cartApi/findAll，接收http请求
// 即：/cartApi/call请求会调用go.micro.api.shop.cartApi 服务的pb.Call方法
func (e *ShopCartApi) FindAll(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	log.Info("接受到 /cartApi/findAll 访问请求")
	fmt.Println("接受到 /cartApi/findAll 访问请求")
	if _, ok := req.Get["user_id"]; !ok {
		//rsp.StatusCode= 500
		return errors.New("参数异常")
	}
	userIdString := req.Get["user_id"].Values[0]
	fmt.Println(userIdString)
	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		return err
	}
	//获取购物车所有商品
	cartAll, err := e.ShopCartService.GetAll(context.TODO(), &cartPb.CartFindAll{UserId: userId})
	//数据类型转化
	b, err := json.Marshal(cartAll)
	if err != nil {
		return err
	}
	rsp.StatusCode = 200
	rsp.Body = string(b)
	return nil
}
