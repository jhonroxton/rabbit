package routers

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	gl "rabbit/globals"
	"time"
)

var (
	g errgroup.Group
)

func RunServive() {
	//初始化路由组，
	router := Routers()

	//api 路由服务
	s := &http.Server{
		Addr:         ":" + gl.V.GetString("myGin.port"),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	fmt.Printf(`
         路由监听端口为:%s
`, gl.V.GetString("myGin.port"))

	// api 监听
	g.Go(func() error {
		err := s.ListenAndServe()
		if err != nil {
			return err
		}
		return s.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		panic("路由监听出错：" + err.Error())
	}

}
