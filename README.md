[license]: /LICENSE
[license-badge]: https://img.shields.io/github/license/Akimio521/MediaWarp?style=flat-square&a=1
[prs]: https://github.com/Akimio521/MediaWarp
[prs-badge]: https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square
[issues]: https://github.com/Akimio521/MediaWarp/issues/new
[issues-badge]: https://img.shields.io/badge/Issues-welcome-brightgreen.svg?style=flat-square
[release]: https://github.com/Akimio521/MediaWarp/releases/latest
[release-badge]: https://img.shields.io/github/v/release/Akimio521/MediaWarp?style=flat-square
[docker]: https://hub.docker.com/r/akimio/mediawarp
[docker-badge]: https://img.shields.io/docker/pulls/akimio/mediawarp?color=%2348BB78&logo=docker&label=pulls

<div align="center">

# MediaWarp

MediaWarp 是**前置于EmbyServer的API服务器**，修改了原版EmbyServer的API返回内容以实现特殊功能  

[![license][license-badge]][license]
[![prs][prs-badge]][prs]
[![issues][issues-badge]][issues]
[![release][release-badge]][release]
[![docker][docker-badge]][docker]



[功能](#功能) •
[TODO LIST](#todo-list) •
[更新日志](#更新日志) •
[Star History](#star-history)

</div>

# 功能
- 对于协议为Http的Strm文件可以实现302直链播放，流量不经过EmbyServer
  - **推荐配合[AutoFilm](https://github.com/Akimio521/AutoFilm)使用**
  - 已通过测试客户端（Web、iOS Emby、Infuse、Conflux、Fileball、Vidhub）
- 屏蔽特定客户端访问
  
  <img src="./img/client_filter.png" alt="" width=500px /> 
- 自定义Web前端样式（HTML、CSS、JavaScript）
  - 效果演示：

    <img src="./img/index.jpg" alt="首页" width=310px /> 
    <img src="./img/movie.jpg" alt="电影" width=310px />
    <img src="./img/series.jpg" alt="电视剧" width=310px />

  - 感谢以下作者提供相关脚本、前端样式（排名不分先后）：
    - [Nolovenodie/emby-crx](https://github.com/Nolovenodie/emby-crx)
    - [newday-life/emby-front-end-mod](https://github.com/newday-life/emby-front-end-mod)
    - [chen3861229/embyExternalUrl](https://github.com/chen3861229/embyExternalUrl)

# TODO LIST
- [x] Strm文件内部为标准http协议链接实现302重定向
- [x] 屏蔽特定客户端访问
- [x] 提供多种Web前端样式
- [ ] Strm文件内部为本地软连接实现特定解析方式（访问Alist获取raw_url后返回给客户端实现302重定向）
- [ ] 利用Mysql/PostgreSQL/Redis优化Infuse媒体库模式下扫库体验
- [ ] 多服务器负载均衡
- [ ] 多服务器转码推流

# 更新日志
[更新日志](./docs/UpdateLog.md)

# Star History
<a href="https://github.com/Akimio521/MediaWarp/stargazers">
    <img width="500" alt="Star History Chart" src="https://api.star-history.com/svg?repos=Akimio521/MediaWarp&type=Date">
</a> 