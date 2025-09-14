# Octopus
一个使用 MessagePack 格式数据，封装 MongoDB 和 Cloudflare R2 服务的 保存（Save）、更新（Update）、删除（Delete）、获取（Get）操作，并通过 HTTP 提供服务的工具。

1）主要方便跨语言、跨机器读写 MongoDB 和 R2 服务，API Gateway：可以使用python、php、nodjs、go ... 语言，简单的将数据包装为 MessagePack 格式，传递给 Octopus 即可完成对 MongoDB 或 R2 的读写操作。

2）调用者仅需关心数据，无需关心不同语言中如何与MongoDB 和 R2 服务驱动交互。

## Usage
1）通过环境变量，设置 MongoDB 的访问信息，Linux可以设置在 /etc/profile， Windows可以直接在环境变量中添加（字段名须与下面一致）
```
# MongoDB
# 如果不启用 MongoDB 则无需设置
export MONGOCONN="mongodb://localhost:27017"
export MONGODATABASE="StableDiffusion"

# Cloudflare R2
# 如果不启用 R2 则无需设置
export CFR2BUCKET="Your-R2-Bucket-Name"
export CFR2ID="Your-AccessID"
export CFR2KEYID="Your-AccessKeyID"
export CFR2KEYSECRET="Your-AccessKeySecret"
```

2）关闭当前终端界面，重新打开终端，环境变量设置就会生效。
3）在 Release 页面下载适合你平台的发行版，解压缩，会得到一个名称为 octopus 的可执行程序，将它复制到 /usr/local/bin/ 目录中（可以放在你喜欢的任意目录中，但后面的启动文件中的路径也需要变更为你设置的路径）；  
4）在终端中启动 octopus：
```
/usr/local/bin/octopus server --debug=false --host=localhost --port=9090 --admin-user=admin --admin-password=123

# 参数说明：
# --debug=true 可以在终端看到运行时到调试信息，一切调试好之后，此参数应该设置为 false
# --admin-user=admin --admin-password=123 Octopus 服务设置有调用认证，HTTP Basic Auth
# --host=localhost Octopus设计为内网基础服务组件，此处应尽量使用内网地址，避免公网直接访问（若必须，则应启用HTTPS）
# --port=9090 Octopus使用的端口

# 其他参数
--with-mongo 是否启用 MongoDB 读写服务，默认启用
--with-r2 是否启用Cloudflare R2 对象存储服务，默认启用
--with-memcache 是否启用内存缓存服务，默认启用
--memcache-size-mb=32 内存缓存的最大内存占用，默认32MB

```

# Example
1) 使用 PHP、Python 语言，使用 MessagePack 数据格式，通过 Octopus 操作 Mongo DB 和 R2 的用法，请看 example 目录；


