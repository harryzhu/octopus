<?php
require_once("send_request.php");

$mem_data = array();
$mem_data["Bucket"] = "images";
?>

<h2>SAVE</h2>
<?php

$mem_data["Action"] = "SAVE";

$mem_data["Data"] = array(
	"_id" => "user_profile_nickname",
	"mc_value" => "harry",
);


$t1=time();
$res = PostMsgPackURL($url_mem_action,$mem_data);
echo time()-$t1 ."s <br>";
echo "status_code: ". $res["info"]["http_code"]."<br>";

print_r($res["output"]);

?>

<h2>GET</h2>
<?php

$mem_data["Action"] = "GET";

$mem_data["Data"] = array(
	"_id" => "user_profile_nickname",
);



$res = PostMsgPackURL($url_mem_action,$mem_data);
echo "status_code: ". $res["info"]["http_code"]."<br>";
if($res["info"]["http_code"] < 400){
	$out = json_decode($res["output"],true);
	$rows = $out["Data"];
	//print_r($rows);
	foreach ($rows as $row) {
		echo $row["_id"]." ==> ".$row["mc_value"]."<br>";
	}
}

?>


<h2>DELETE</h2>
<?php

$mem_data["Action"] = "DELETE";

$mem_data["Data"] = array(
	"_id" => "user_profile_nickname",
);



$res = PostMsgPackURL($url_mem_action,$mem_data);
echo "status_code: ". $res["info"]["http_code"]."<br>";
if($res["info"]["http_code"] < 400){
	$out = json_decode($res["output"],true);
	$rows = $out["Data"];
	print_r($rows);
}

?>


<h2>GET after DELETE: </h2>
<p>should be error because of delete before, [mc_value] should be empty.</p>
<?php

$mem_data["Action"] = "GET";

$mem_data["Data"] = array(
	"_id" => "user_profile_nickname",
);



$res = PostMsgPackURL($url_mem_action,$mem_data);
echo "status_code: ". $res["info"]["http_code"]."<br>";
if($res["info"]["http_code"] < 400){
	$out = json_decode($res["output"],true);
	$rows = $out["Data"];
	print_r($rows);
	foreach ($rows as $row) {
		echo $row["_id"]." ==> ".$row["mc_value"]."<br>";
	}
}else{
	$out = json_decode($res["output"],true);
	print_r($out);
}

?>

