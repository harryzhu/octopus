<?php
ini_set('mongo.native_long', 1);

require_once("send_request.php");

$r2_data = array();

$r2_data["Bucket"] = "webui";
?>

<h2>SAVE</h2>
<?php

$r2_data["Action"] = "SAVE";
$r2_data["Data"] = array(
	 "_id" => "test/r2-php-test-01.png",
	 "path" => "/Users/harry/5.png",
	 "mime" => "image/png",
);



$res_save = PostMsgPackURL($url_r2_action,$r2_data);

print_r($res_save["info"]["http_code"]);
echo "<br>";
print_r($res_save["output"]);

?>

<h2>DELETE</h2>
<?php

$r2_data["Action"] = "DELETE";
$r2_data["Data"] = array(
	 "_id" => "test/r2-php-test-01.png",
);



$res_delete = PostMsgPackURL($url_r2_action,$r2_data);

print_r($res_delete["info"]["http_code"]);
echo "<br>";
print_r($res_delete["output"]);

?>