<p>v1.0 完成路由(基于map和tree)</p>
<p>v1.1 完成中间件(middleware,builder设计模式)，用中间件hook响应发送</p>
<p>v1.2 修改了json会多次发送header的bug，完成缓慢退出的功能(gracefulshutdown)</p>
<p>v1.3 增加了基于压缩前缀树的路由(radixtrie，但只实现了静态匹配)，增加了路由组功能(server/routegroup)，增加了静态文件路由(file)</p>