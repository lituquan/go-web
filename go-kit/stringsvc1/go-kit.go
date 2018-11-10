1.handler
	底层：			
		//Handler接口
		type Handler interface {
			ServeHTTP(ResponseWriter, *Request)  // 路由实现器
		}
		//HandlerFunc 实现了Handler接口
		type HandlerFunc func(ResponseWriter, *Request)
		func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
			f(w, r)
		}
	gokit实现的handler,主要接口实现ServeHTTP：
		type Server struct {
			//...
		}
		func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if len(s.finalizer) > 0 {
			//6.defer 语句顺序执行所有finalizer
			}
			for _, f := range s.before {
			//1.顺序执行所有before
			}
			//2.执行dec  解析请求
			request, err := s.dec(ctx, r)
			//3.执行e  
			response, err := s.e(ctx, request)
			for _, f := range s.after {
			//4.执行after
			}
			//5.执行enc
			err := s.enc(ctx, w, response)
		}
2.demo分析：
	//业务对象：服务
	svc := stringService{}
	//Server.handler对象构建
	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc), //3-e:业务解析实现
		decodeUppercaseRequest,	//2-dec：解析请求
		encodeResponse,		//5-enc:编码响应
	)
	//路由注册：默认DefaultServeMux维护一个map
	http.Handle("/uppercase", uppercaseHandler)
	//开服务
	log.Fatal(http.ListenAndServe(":8080", nil))
	
	/*
		分析：
	*/
	//请求体json反序列化为对象：uppercaseRequest
	//请求报文-->request.body-->uppercaseRequest
	func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
		var request uppercaseRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, err
		}
		return request, nil
	}
	//uppercaseRequest对象的s值传给:Uppercase(s)
	//uppercaseRequest.S-->uppercaseResponse
	func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
		return func(_ context.Context, request interface{}) (interface{}, error) {
			req := request.(uppercaseRequest)
			v, err := svc.Uppercase(req.S)
			if err != nil {
				return uppercaseResponse{v, err.Error()}, nil
			}
			return uppercaseResponse{v, ""}, nil
		}
	}
	//uppercaseResponse对象json序列化:响应体
	//uppercaseResponse-->response.body
	func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
		return json.NewEncoder(w).Encode(response)
	}

3.过程复现
	/*
	合起来写：
	*/
	type Myhandler struct{
	}
	func (*Myhandler)ServeHTTP(w http.ResponseWriter, r *http.Request){
		//dec  解析请求体
		var request uppercaseRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return
		}
		//e：调用业务逻辑
		svc := stringService{}
		v, _ := svc.Uppercase(request.S)
		response:=uppercaseResponse{v, ""}
		//enc 编码返回对象
		json.NewEncoder(w).Encode(response)
		return 
	}
	
4.go底层包：
	(1)注册路由：Handle方法
		http.Handle("/uppercase", uppercaseHandler)	
		func Handle(pattern string, handler Handler) {
			DefaultServeMux.Handle(pattern, handler) 
		}	
		// Handle registers the handler for the given pattern. 
		func (mux *ServeMux) Handle(pattern string, handler Handler) {
			//注册路由对象,map存储路径和处理方法
			//路由表：map[string]muxEntry
			mux.m[pattern] = muxEntry{explicit: true, h: handler, pattern: pattern}	
		}

		type ServeMux struct {
			mu    sync.RWMutex
			m     map[string]muxEntry
			hosts bool // whether any patterns contain hostnames
		}		
		type muxEntry struct {
			explicit bool
			h        Handler
			pattern  string
		}
	(2)服务监听、服务执行		
		http.Server-->tcp.Serve-->handler.severHttp
		1)服务监听
		// handler to invoke, http.DefaultServeMux if nil
		func ListenAndServe(addr string, handler Handler) error {
			server := &Server{Addr: addr, Handler: handler}
			return server.ListenAndServe()
		}
		//Server
		type Server struct{
		}
		func (srv *Server) ListenAndServe() error {
			addr := srv.Addr
			if addr == "" {
				addr = ":http"
			}
			ln, err := net.Listen("tcp", addr)//tcp协议
			if err != nil {
				return err
			}
			return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
		}
		func (srv *Server) Serve(l net.Listener) error {
			//一直等待客户端：持续服务
			for {
				//堵塞等待
				rw, e := l.Accept()
				//建立连接
				c, err := srv.newConn(rw)
				//执行一次服务
				//一个连接分配一个协程
				go c.serve()
			}
		}
		2)服务执行:路由匹配和执行
		//http1.1调用
		c.serve(){
			serverHandler{c.server}.ServeHTTP(w, w.req)
		}
		func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
			if r.RequestURI == "*" {
				if r.ProtoAtLeast(1, 1) {
					w.Header().Set("Connection", "close")
				}
				w.WriteHeader(StatusBadRequest)
				return
			}
			h, _ := mux.Handler(r)//路由匹配
			h.ServeHTTP(w, r)//执行控制器
		}

	参考：
		https://github.com/go-kit/kit/tree/master/examples/stringsvc1
		http://gokit.io/examples/stringsvc.html
		https://www.jianshu.com/p/be3d9cdc680b