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

# 本地测试
## 构建所有版本
```BASH
goreleaser release --snapshot --clean --skip-publish
```
## 流程测试
```BASH
goreleaser release --snapshot --clean --skip-publish --rm-dist
```

## 参数含义
1. **snapshot**: 这将跳过标签验证，并生成一个快照版本
2. **clean**: 清除上一次生成的 dist
3. **skip-publish**: 跳过推送镜像到 Docker Registry
4. **rm-dist**：删除 dist

# 媒体服务器部分 API 具体响应
## EmbyServer
### playbackInfo
**EmbyServer version 4.8**
AlistStrm

```
{
    "MediaSources": [
        {
            "Protocol": "File",
            "Id": "37bcb3827aeb2095f1b6293e9a189072",
            "Path": "/资源盘/ANIOpen/2024-10/[ANi] 魔王陛下，RETRY！R - 04 [1080P][Baha][WEB-DL][AAC AVC][CHT].mp4",
            "Type": "Default",
            "Container": "strm",
            "Size": 0,
            "Name": "重来吧，魔王大人！ S02E04  -1080p -AVC -AAC -ANi -Baha",
            "IsRemote": false,
            "HasMixedProtocols": false,
            "SupportsTranscoding": true,
            "SupportsDirectStream": false,
            "SupportsDirectPlay": false,
            "IsInfiniteStream": false,
            "RequiresOpening": false,
            "RequiresClosing": false,
            "RequiresLooping": false,
            "SupportsProbing": false,
            "MediaStreams": [],
            "Formats": [],
            "RequiredHttpHeaders": {},
            "DirectStreamUrl": "/videos/61162/master.m3u8?DeviceId=e12ab84a-2948-463a-9cb0-ecebcac2949e&MediaSourceId=37bcb3827aeb2095f1b6293e9a189072&PlaySessionId=e643532a23ca4a60b2c70d4210c018a9&api_key=e824c422e02047dfa820dbcab5490bb9&VideoCodec=h264,h265,hevc,av1&VideoBitrate=1400000&TranscodingMaxAudioChannels=2&SegmentContainer=ts&MinSegments=1&BreakOnNonKeyFrames=True&SubtitleStreamIndexes=-1&ManifestSubtitles=vtt&h264-profile=high,main,baseline,constrainedbaseline,high10&h264-level=62&hevc-codectag=hvc1,hev1,hevc,hdmv&TranscodeReasons=ContainerBitrateExceedsLimit",
            "AddApiKeyToDirectStreamUrl": false,
            "TranscodingUrl": "/videos/61162/master.m3u8?DeviceId=e12ab84a-2948-463a-9cb0-ecebcac2949e&MediaSourceId=37bcb3827aeb2095f1b6293e9a189072&PlaySessionId=e643532a23ca4a60b2c70d4210c018a9&api_key=e824c422e02047dfa820dbcab5490bb9&VideoCodec=h264,h265,hevc,av1&VideoBitrate=1400000&TranscodingMaxAudioChannels=2&SegmentContainer=ts&MinSegments=1&BreakOnNonKeyFrames=True&SubtitleStreamIndexes=-1&ManifestSubtitles=vtt&h264-profile=high,main,baseline,constrainedbaseline,high10&h264-level=62&hevc-codectag=hvc1,hev1,hevc,hdmv&TranscodeReasons=ContainerBitrateExceedsLimit",
            "TranscodingSubProtocol": "hls",
            "TranscodingContainer": "ts",
            "ReadAtNativeFramerate": false,
            "ItemId": "61162"
        }
    ],
    "PlaySessionId": "e643532a23ca4a60b2c70d4210c018a9"
}
```
普通视频
```
{
    "MediaSources": [
        {
            "Protocol": "File",
            "Id": "56f692be38b2c4ae3c8f155e77006d46",
            "Path": "/media/媒体库/动漫追番/青之箱 (2024)/Season 1/青之箱 S01E04  -HEVC -AAC -KitaujiSub.mkv",
            "Type": "Default",
            "Container": "mkv",
            "Size": 470388402,
            "Name": "青之箱 S01E04  -HEVC -AAC -KitaujiSub",
            "IsRemote": false,
            "HasMixedProtocols": false,
            "RunTimeTicks": 14172790000,
            "SupportsTranscoding": true,
            "SupportsDirectStream": false,
            "SupportsDirectPlay": false,
            "IsInfiniteStream": false,
            "RequiresOpening": false,
            "RequiresClosing": false,
            "RequiresLooping": false,
            "SupportsProbing": false,
            "MediaStreams": [
                {
                    "Codec": "hevc",
                    "ColorSpace": "bt709",
                    "TimeBase": "1/1000",
                    "VideoRange": "SDR",
                    "DisplayTitle": "1080p HEVC",
                    "IsInterlaced": false,
                    "BitRate": 2655163,
                    "BitDepth": 10,
                    "RefFrames": 1,
                    "IsDefault": true,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Height": 1080,
                    "Width": 1920,
                    "AverageFrameRate": 23.976025,
                    "RealFrameRate": 23.976025,
                    "Profile": "Main 10",
                    "Type": "Video",
                    "AspectRatio": "16:9",
                    "Index": 0,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Protocol": "File",
                    "PixelFormat": "yuv420p10le",
                    "Level": 120,
                    "IsAnamorphic": false,
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 0
                },
                {
                    "Codec": "aac",
                    "Language": "jpn",
                    "TimeBase": "1/1000",
                    "DisplayTitle": "Japanese AAC 5.1 (默认)",
                    "DisplayLanguage": "Japanese",
                    "IsInterlaced": false,
                    "ChannelLayout": "5.1",
                    "BitRate": 320000,
                    "Channels": 6,
                    "SampleRate": 48000,
                    "IsDefault": true,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Profile": "LC",
                    "Type": "Audio",
                    "Index": 1,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 0
                },
                {
                    "Codec": "ass",
                    "Language": "chi",
                    "TimeBase": "1/1000",
                    "Title": "chs_jp",
                    "DisplayTitle": "Chinese Simplified (默认 ASS)",
                    "DisplayLanguage": "Chinese Simplified",
                    "IsInterlaced": false,
                    "IsDefault": true,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Subtitle",
                    "Index": 2,
                    "IsExternal": false,
                    "DeliveryMethod": "External",
                    "DeliveryUrl": "/Videos/61256/56f692be38b2c4ae3c8f155e77006d46/Subtitles/2/0/Stream.ass?api_key=e824c422e02047dfa820dbcab5490bb9",
                    "IsExternalUrl": false,
                    "IsTextSubtitleStream": true,
                    "SupportsExternalStream": true,
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 0,
                    "SubtitleLocationType": "InternalStream"
                },
                {
                    "Codec": "ass",
                    "Language": "chi",
                    "TimeBase": "1/1000",
                    "Title": "cht_jp",
                    "DisplayTitle": "Chinese Simplified (ASS)",
                    "DisplayLanguage": "Chinese Simplified",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Subtitle",
                    "Index": 3,
                    "IsExternal": false,
                    "DeliveryMethod": "External",
                    "DeliveryUrl": "/Videos/61256/56f692be38b2c4ae3c8f155e77006d46/Subtitles/3/0/Stream.ass?api_key=e824c422e02047dfa820dbcab5490bb9",
                    "IsExternalUrl": false,
                    "IsTextSubtitleStream": true,
                    "SupportsExternalStream": true,
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 0,
                    "SubtitleLocationType": "InternalStream"
                },
                {
                    "Codec": "otf",
                    "TimeBase": "1/90000",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Attachment",
                    "Index": 4,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Path": "FOT-ModeMinBLargeStd-B.0.E8KEG6LU.otf",
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 49752,
                    "MimeType": "application/vnd.ms-opentype"
                },
                {
                    "Codec": "otf",
                    "TimeBase": "1/90000",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Attachment",
                    "Index": 5,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Path": "FOT-SeuratProN-DB.0.JIIL4FRM.otf",
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 196388,
                    "MimeType": "application/vnd.ms-opentype"
                },
                {
                    "Codec": "otf",
                    "TimeBase": "1/90000",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Attachment",
                    "Index": 6,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Path": "FOT-TsukuMinPr5N-B.0.8X1XX3RE.otf",
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 82360,
                    "MimeType": "application/vnd.ms-opentype"
                },
                {
                    "Codec": "ttf",
                    "TimeBase": "1/90000",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Attachment",
                    "Index": 7,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Path": "FZLANTINGHEI-DB-GBK.0.WOAP9IVQ.ttf",
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 15908,
                    "MimeType": "application/x-truetype-font"
                },
                {
                    "Codec": "ttf",
                    "TimeBase": "1/90000",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Attachment",
                    "Index": 8,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Path": "方正粗雅宋_GBK.0.ELFFB6DY.ttf",
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 89012,
                    "MimeType": "application/x-truetype-font"
                },
                {
                    "Codec": "ttf",
                    "TimeBase": "1/90000",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Attachment",
                    "Index": 9,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Path": "方正兰亭大黑_GBK.0.5PPVFM7L.ttf",
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 13900,
                    "MimeType": "application/x-truetype-font"
                },
                {
                    "Codec": "ttf",
                    "TimeBase": "1/90000",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Attachment",
                    "Index": 10,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Path": "方正兰亭黑_GBK.0.OAU1SQQO.ttf",
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 33180,
                    "MimeType": "application/x-truetype-font"
                },
                {
                    "Codec": "ttf",
                    "TimeBase": "1/90000",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Attachment",
                    "Index": 11,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Path": "方正兰亭圆_GBK_准.0.D3ABRYMS.ttf",
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 354176,
                    "MimeType": "application/x-truetype-font"
                },
                {
                    "Codec": "ttf",
                    "TimeBase": "1/90000",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Attachment",
                    "Index": 12,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Path": "方正准雅宋_GBK.0.T77HXJCC.ttf",
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 69468,
                    "MimeType": "application/x-truetype-font"
                },
                {
                    "Codec": "ttf",
                    "TimeBase": "1/90000",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Attachment",
                    "Index": 13,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Path": "华康翩翩体W3-A.0.B7WYAJQK.ttf",
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 29816,
                    "MimeType": "application/x-truetype-font"
                },
                {
                    "Codec": "ttf",
                    "TimeBase": "1/90000",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Attachment",
                    "Index": 14,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Path": "华康翩翩体W5-A.0.CVH5DGCU.ttf",
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 58720,
                    "MimeType": "application/x-truetype-font"
                },
                {
                    "Codec": "ttf",
                    "TimeBase": "1/90000",
                    "IsInterlaced": false,
                    "IsDefault": false,
                    "IsForced": false,
                    "IsHearingImpaired": false,
                    "Type": "Attachment",
                    "Index": 15,
                    "IsExternal": false,
                    "IsTextSubtitleStream": false,
                    "SupportsExternalStream": false,
                    "Path": "华康手札体W5-A.0.B8DWXF5W.ttf",
                    "Protocol": "File",
                    "ExtendedVideoType": "None",
                    "ExtendedVideoSubType": "None",
                    "ExtendedVideoSubTypeDescription": "None",
                    "AttachmentSize": 56332,
                    "MimeType": "application/x-truetype-font"
                }
            ],
            "Formats": [],
            "Bitrate": 2655163,
            "RequiredHttpHeaders": {},
            "DirectStreamUrl": "/videos/61256/master.m3u8?DeviceId=e12ab84a-2948-463a-9cb0-ecebcac2949e&MediaSourceId=56f692be38b2c4ae3c8f155e77006d46&PlaySessionId=592ccdf484644f3db25993c423f392ad&api_key=e824c422e02047dfa820dbcab5490bb9&VideoCodec=h264,h265,hevc,av1&AudioCodec=ac3,mp3,aac&VideoBitrate=1430000&AudioBitrate=320000&AudioStreamIndex=1&TranscodingMaxAudioChannels=2&SegmentContainer=ts&MinSegments=1&BreakOnNonKeyFrames=True&SubtitleStreamIndexes=-1&ManifestSubtitles=vtt&h264-profile=high,main,baseline,constrainedbaseline,high10&h264-level=62&hevc-codectag=hvc1,hev1,hevc,hdmv&TranscodeReasons=ContainerBitrateExceedsLimit",
            "AddApiKeyToDirectStreamUrl": false,
            "TranscodingUrl": "/videos/61256/master.m3u8?DeviceId=e12ab84a-2948-463a-9cb0-ecebcac2949e&MediaSourceId=56f692be38b2c4ae3c8f155e77006d46&PlaySessionId=592ccdf484644f3db25993c423f392ad&api_key=e824c422e02047dfa820dbcab5490bb9&VideoCodec=h264,h265,hevc,av1&AudioCodec=ac3,mp3,aac&VideoBitrate=1430000&AudioBitrate=320000&AudioStreamIndex=1&TranscodingMaxAudioChannels=2&SegmentContainer=ts&MinSegments=1&BreakOnNonKeyFrames=True&SubtitleStreamIndexes=-1&ManifestSubtitles=vtt&h264-profile=high,main,baseline,constrainedbaseline,high10&h264-level=62&hevc-codectag=hvc1,hev1,hevc,hdmv&TranscodeReasons=ContainerBitrateExceedsLimit",
            "TranscodingSubProtocol": "hls",
            "TranscodingContainer": "ts",
            "ReadAtNativeFramerate": false,
            "DefaultAudioStreamIndex": 1,
            "DefaultSubtitleStreamIndex": 2,
            "ItemId": "61256"
        }
    ],
    "PlaySessionId": "592ccdf484644f3db25993c423f392ad"
}
```