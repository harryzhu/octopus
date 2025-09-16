<?php
ini_set('mongo.native_long', 1);

require_once("send_request.php");



$mongo_action_url = "http://localhost:9090/mongo/action";

$mongo_data = array();
$mongo_data["Bucket"] = "image";

?>

<h2>SAVE:</h2>
<?php

$mongo_data["Action"] = "SAVE";
$mongo_data["Data"] = array(
	"_id" => "test-01",
	 "username" => "__testman__",
	 "image_width" => 512,
	 "image_height" => 512,
);

$res_save = PostMsgPackURL($mongo_action_url,$mongo_data);

print_r($res_save["info"]["http_code"]);
echo "<br>";
print_r($res_save["output"]);

?>

<h2>GET:</h2>
<?php

$mongo_data["Action"] = "GET";

$filter = array("image_width"=>array("\$gt"=>500));
$option = array("sort" =>array("username"=>1),
	"skip" =>0,
	"limit" => 0,
	"projection"=>array("_id"=>1,"image_width"=>1,"image_height"=>1),
);

$mongo_data["Data"] = array(
	 "filter" => $filter,
	 "option" => $option,
);

$res_get = PostMsgPackURL($mongo_action_url,$mongo_data);

print_r($res_get["info"]["http_code"]);
echo "<br><br>";
$res_get_output = $res_get["output"];
if(!empty($res_get_output)){
	$out = json_decode($res_get_output, true);
	$rows = $out["Data"];
	foreach ($rows as $row) {
		foreach ($row as $k=>$v) {
			echo $k." ====> ".$v."<br>";
		}
		
	}
}
?>

<h2>UPDATE:</h2>
<?php

$mongo_data["Action"] = "UPDATE";
$mongo_data["Data"] = array(
	"_id" => "test-01",
	 "image_height" => 1024,
);

$res_update = PostMsgPackURL($mongo_action_url,$mongo_data);

print_r($res_update["info"]["http_code"]);
echo "<br><br>";
$res_update_output = $res_update["output"];
print_r($res_update_output);
$res_update_output_data = json_decode($res_update["output"], true)["Data"];
echo "<br><br>";
print_r($res_update_output_data);

?>

<h2>DELETE:</h2>
<?php
$mongo_data["Action"] = "DELETE";
$mongo_data["Data"] = array(
	"_id" => "test-01",
);

$res_delete = PostMsgPackURL($mongo_action_url,$mongo_data);

print_r($res_delete["info"]["http_code"]);
echo "<br><br>";
$res_delete_output = $res_delete["output"];
print_r($res_delete_output);




?>
