# 开发相关文档

## 项目目录结构
```
MEDIAWARP
├─.github
│  └─workflows          # 工作流相关
├─api                   # Alist、EmbyServer等服务器的开放接口
├─config                # MediaWarp配置文件
├─constants             # 存放一些常量
├─core                  # MediaWarp核心组织，包括配置读取、日志等
├─docs                  # MediaWarp相关文档
├─handlers              # MediaWarp请求处理函数
│  └─handlers_emby      # 处理服务器为EmbySever的一些函数
├─img                   # MediaWarp演示图片
├─logs                  # MediaWarp输出日志
├─middleware            # Gin框架用到的中间件，包括访问日志、客户端过滤
├─utils                   # 一些工具集合
├─resources             # 内嵌的一些JavaScript脚本、CSS样式
│  ├─css
│  └─js
├─router                # MediaWarp的API路由部分
├─schemas               # API响应结构体
│  ├─schemas_alist
│  └─schemas_emby
└─static                # MediaWarp自定义静态文件存储目录
```