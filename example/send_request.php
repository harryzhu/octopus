<?php
$url_mem_action = "http://localhost:9090/memcache/action";
$url_mongo_action = "http://localhost:9090/mongo/action";
$url_r2_action = "http://localhost:9090/r2/action";

function PostMsgPackURL($url,$data){ 
	$packer = new \MessagePack(false);
	$packed = $packer->pack($data);
	$headerArray =array("Content-Type:application/octet-stream","Accept:application/json");
	$curl = curl_init();
	curl_setopt($curl, CURLOPT_URL, $url);
	curl_setopt($curl, CURLOPT_USERPWD, "admin:123");  
	curl_setopt($curl, CURLOPT_SSL_VERIFYPEER, FALSE);
	curl_setopt($curl, CURLOPT_SSL_VERIFYHOST,FALSE);
	curl_setopt($curl, CURLOPT_POST, 1);
	curl_setopt($curl, CURLOPT_POSTFIELDS, $packed);
	curl_setopt($curl, CURLOPT_HTTPHEADER,$headerArray);
	curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
	curl_setopt($curl, CURLOPT_TIMEOUT, 5);
	$output = curl_exec($curl);
	$resp = array(
		"info" =>curl_getinfo($curl),
		"output" =>$output,
	);

	curl_close($curl);
	return $resp;
}

?>
