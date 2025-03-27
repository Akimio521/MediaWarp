package static

import (
	"embed"
)

//go:embed embyExternalUrl/embyWebAddExternalUrl/embyLaunchPotplayer.js
//go:embed dd-danmaku/ede.js
//go:embed jellyfin-danmaku/ede.js
//go:embed emby-web-mod/actorPlus/actorPlus.js
//go:embed emby-web-mod/emby-swiper/emby-swiper.js
//go:embed emby-web-mod/emby-tab/emby-tab.js
//go:embed emby-web-mod/fanart_show/fanart_show.js
//go:embed emby-web-mod/itemSortForNewDevice/itemSortForNewDevice.js
//go:embed emby-web-mod/playbackRate/playbackRate.js
//go:embed emby-web-mod/trailer/trailer.js
//go:embed emby-crx/static/css/style.css
//go:embed emby-crx/static/js/common-utils.js
//go:embed emby-crx/static/js/jquery-3.6.0.min.js
//go:embed emby-crx/static/js/md5.min.js
//go:embed emby-crx/content/main.js
//go:embed jellyfin-crx/static/css/style.css
//go:embed jellyfin-crx/static/js/common-utils.js
//go:embed jellyfin-crx/static/js/jquery-3.6.0.min.js
//go:embed jellyfin-crx/static/js/md5.min.js
//go:embed jellyfin-crx/content/main.js
var EmbeddedStaticAssets embed.FS
