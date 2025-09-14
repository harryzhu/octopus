# Octopus
一个使用 MessagePack 格式数据，封装 MongoDB 和 Cloudflare R2 服务的 保存（Save）、更新（Update）、删除（Delete）、获取（Get）操作，并通过 HTTP 提供服务的工具。

1）主要方便跨语言、跨机器读写 MongoDB 和 R2 服务，API Gateway：可以使用python、php、nodjs、go ... 等语言，简单的将数据包装为 MessagePack 格式，传递给 Octopus 即可完成对 MongoDB 或 R2 的读写操作。

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
# --admin-user=admin --admin-password=123 调用 Octopus 服务设置有认证，HTTP Basic Auth
# --host=localhost Octopus设计为内网基础服务组件，此处应尽量使用内网地址，避免公网直接访问（若必须，则应启用HTTPS）
# --port=9090 Octopus使用的端口

# 其他参数
--with-mongo 是否启用 MongoDB 读写服务，默认启用
--with-r2 是否启用Cloudflare R2 对象存储服务，默认启用
--with-memcache 是否启用内存缓存服务，默认启用
--memcache-size-mb=32 内存缓存的最大内存占用，默认32MB

```

# Example
1）使用 PHP、Python 语言，使用 MessagePack 数据格式，通过 Octopus 操作 Mongo DB 和 R2 的用法，请看 `example` 目录；
2）PHP需要先安装 msgpack 扩展：
```
pecl install msgpack
```
3）Python 需要先安装 msgpack 和 requests ：
```
pip install msgpack
pip install requests
```

# PHP 用法：
1）先创建一个公共函数（一般放到一个文件中（`send_request.php`），其他文件要使用，`requrie_once("send_request.php");` 即可） ，用于将数据发送给 Octopus：
```php

<?php
$url_mongo_action = "http://localhost:9090/mongo/action";
$url_r2_action = "http://localhost:9090/r2/action";
$url_mem_action = "http://localhost:9090/memcache/action";

function PostMsgPackURL($url,$data){ 
	$packer = new \MessagePack(false);
	// 将 PHP 数据包转成 messagepack 格式；
	// 所以需要按上述方法给PHP安装messagepack扩展并启用；
	$packed = $packer->pack($data);
	$headerArray =array("Content-Type:application/octet-stream","Accept:application/json");
	$curl = curl_init();
	curl_setopt($curl, CURLOPT_URL, $url);
	// 此处的 admin:123 即为在启动 octopus 是用参数 --admin-user=admin --admin-password=123 指定的用户名和密码
	curl_setopt($curl, CURLOPT_USERPWD, "admin:123");  
	curl_setopt($curl, CURLOPT_SSL_VERIFYPEER, FALSE);
	curl_setopt($curl, CURLOPT_SSL_VERIFYHOST,FALSE);
	curl_setopt($curl, CURLOPT_POST, 1);
	curl_setopt($curl, CURLOPT_POSTFIELDS, $packed);
	curl_setopt($curl, CURLOPT_HTTPHEADER,$headerArray);
	curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
	curl_setopt($curl, CURLOPT_TIMEOUT, 5);
	$output = curl_exec($curl);
	// 返回的数据，包含 头部信息 info 和 实际数据 output
	$resp = array(
		"info" =>curl_getinfo($curl),
		"output" =>$output,
	);

	curl_close($curl);
	return $resp;
}

?>
```
2）包装PHP中的数据：
```php
<?php
// 定义一个最终发送给 Octopus 的数组；
// Bucket 用于指定：
// R2 中表示你需要保存文件的bucketName，
// mongodb中表示具体要操作的 collection（集合） ，注意，不是数据库
//（数据库在环境变量中指定：比如 export MONGODATABASE="StableDiffusion"，即指明数据库名称为 StableDiffusion）
// 比如下面数组中将 Bucket 字段设置为 "image"， 即表示 这些操作将作用于
// 数据库 StableDiffusion下面的 image 集合。
$mongo_data = array();
$mongo_data["Bucket"] = "image";
// 字段首字母需要大写Action，Bucket，Data，其中Action 必须是 SAVE、GET、UPDATE、DELETE 中的一个（大写），不能是其他字符；
// Action 这里指明了是 SAVE ，即新加，假如数据已经存在，会报错，对于已经存在的数据，修改应该用更新 UPDATE
// SAVE、UPDATE、DELETE 必须在 Data 字段中 指定 _id 字段
// GET 查询无需 _id 字段，需要 filter 和 option 字段
// _id 字段在 R2 服务中表示即将保存在 bucket 中的键；
// _id 在 mongodb 服务中表示即将保存在 mongodb 中的键（mongodb objectid）；
$mongo_data["Action"] = "SAVE";
$mongo_data["Data"] = array(
	"_id" => "test-01",
	 "username" => "__testman__",
	 "image_width" => 512,
	 "image_height" => 512,
);

// 使用 PostMsgPackURL 函数发送
$res_save = PostMsgPackURL($mongo_action_url,$mongo_data);

print_r($res_save["info"]["http_code"]);
echo "<br>";
// 输出结果
print_r($res_save["output"]);

?>
```

# Python 用法：
```python

# Message 的字段 Action、Bucket、Data 首字母大写，
# Action 必须是SAVE、GET、UPDATE、DELETE 中的一个（大写）
# _id 在 R2 的 bucket 里面保存的路径
# path 文件所在的本地路径，如果是远程服务器，需要先上传到服务器上，然后此处填入服务器上的地址
# mime 文件类型

import msgpack
import requests

message = {
    "Action": "SAVE",
    "Bucket": "your-bucket-name",
    "Data": {
        "_id":"r2-test-01.png",
        "path": "/Users/harryzhu/images/1.png",
        "mime": "image/png"
    }
}

# 包装成 MessagePack 格式
m = msgpack.packb(message, use_bin_type=True)


sess = requests.session()
sess.keep_alive = False
url = "http://localhost:9090/r2/action"

r = sess.post(url, data=m, auth=('admin','123'), stream=True)
print(r.text)

```



