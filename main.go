package main

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/asim/go-micro/plugins/wrapper/select/roundrobin/v4"
	opentracing4 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/opentracing/opentracing-go"
	cli2 "github.com/urfave/cli/v2"
	cartPb "github.com/wsjcko/shopcart/protobuf/pb"
	"github.com/wsjcko/shopcartApi/common"
	"github.com/wsjcko/shopcartApi/handler"
	pb "github.com/wsjcko/shopcartApi/protobuf/pb"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	log "go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"net"
	"net/http"
)

var (
	MICRO_API_NAME       = "go.micro.api.shopCartApi" //决定路由： shopCartApi/findAll?user_id=3
	MICRO_SERVICE_NAME   = "go.micro.service.shop.cart"
	MICRO_VERSION        = "latest"
	MICRO_ADDRESS        = "0.0.0.0:8088"
	MICRO_HYSTRIX_HOST   = "0.0.0.0"
	MICRO_HYSTRIX_PORT   = "9096"
	MICRO_CONSUL_ADDRESS = "127.0.0.1:8500"
	MICRO_JAEGER_ADDRESS = "127.0.0.1:6831"
	DOCKER_HOST          = "127.0.0.1"
)

func SetDockerHost(host string) {
	DOCKER_HOST = host
	MICRO_CONSUL_ADDRESS = host + ":8500"
	MICRO_JAEGER_ADDRESS = host + ":6831"
}

func main() {

	function := micro.NewFunction(
		micro.Flags(
			&cli2.StringFlag{ //micro 多个选项 --ip
				Name:  "ip",
				Usage: "docker Host IP(ubuntu)",
				Value: "0.0.0.0",
			},
		),
	)

	function.Init(
		micro.Action(func(c *cli2.Context) error {
			ipstr := c.Value("ip").(string)
			if len(ipstr) > 0 { //后续校验IP
				fmt.Println("docker Host IP(ubuntu)1111", ipstr)
			}
			SetDockerHost(ipstr)
			return nil
		}),
	)

	fmt.Println("DOCKER_HOST ", DOCKER_HOST)

	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			MICRO_CONSUL_ADDRESS,
		}
	})

	//链路追踪
	t, io, err := common.NewTracer(MICRO_API_NAME, MICRO_JAEGER_ADDRESS)
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	// 启动端口 启动监听 上报熔断状态
	go func() {
		err = http.ListenAndServe(net.JoinHostPort(MICRO_HYSTRIX_HOST, MICRO_HYSTRIX_PORT), hystrixStreamHandler)
		if err == http.ErrServerClosed {
			log.Info("httpserver shutdown cased: ", err)
		} else {
			log.Error(err)
		}
	}()

	// New Service
	srv := micro.NewService(
		micro.Name(MICRO_API_NAME),
		micro.Version(MICRO_VERSION),
		micro.Address(MICRO_ADDRESS),
		//添加 consul 注册中心
		micro.Registry(consulRegistry),
		//添加链路追踪 服务端绑定handle 客户端绑定client
		micro.WrapClient(opentracing4.NewClientWrapper(opentracing.GlobalTracer())),
		//添加熔断
		micro.WrapClient(NewClientHystrixWrapper()),
		//添加负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	// Initialise service
	srv.Init()

	// 调用后端服务
	shopCartService := cartPb.NewShopCartService(MICRO_SERVICE_NAME, srv.Client())

	shopCartService.AddCart(context.TODO(), &cartPb.CartInfo{

		UserId:    3,
		ProductId: 6,
		SizeId:    7,
		Num:       7,
	})

	// Register Handler
	if err := pb.RegisterShopCartApiHandler(srv.Server(), &handler.ShopCartApi{ShopCartService: shopCartService}); err != nil {
		log.Error(err)
	}

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		//run 正常执行
		fmt.Println(req.Service() + "." + req.Endpoint())
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(err error) error {
		fmt.Println(err)
		return err
	})
}

func NewClientHystrixWrapper() client.Wrapper {
	return func(i client.Client) client.Client {
		return &clientWrapper{i}
	}
}
