function oneanimateBackward(){
	$('#topic').replaceClassTo("topic-backward");
	$('#subtitle').replaceClassTo("subtitle-backward");
	$('#rdohtml').replaceClassTo("rdohtml-backward");
	$('#rdocss').replaceClassTo("rdocss-backward");
	$('#rdojs').replaceClassTo("rdojs-backward");
	
	$('#txt-up').val('oneanimate');
	$('#txt-down').val('oneanimateBackward');
}
function oneanimate(){
	//从2跳回来，或者从1跳过来
	if($('#boxplug').attr("class")=="boxplug-forward"||$('#boxplug').attr("class")=="boxplug-down"){
		$('#boxplug').replaceClassTo("boxplug-backward");
		$('#topic2').replaceClassTo("topic2-backward");
		$('#topic').replaceClassTo("topic-down");
		initplugs();
	}else{
		$('#topic').replaceClassTo("topic-forward");
		$('#subtitle').replaceClassTo("subtitle-forward");
	}
	
	//加载圆
	$('#rdohtml').replaceClassTo("rdohtml-forward");
	$('#rdocss').replaceClassTo("rdocss-forward");
	$('#rdojs').replaceClassTo("rdojs-forward");
	
	$('#txt-up').val('twoanimate');
	$('#txt-down').val('oneanimateBackward');
	
}
function twoanimate(){
	//从3跳回来
	if($('#phone').attr("class")=="phone-forward"){
		$('#down').replaceClassTo("animate-show");
		$('#phone').replaceClassTo("phone-backward");
		$('#pull').replaceClassTo("pull-backward");
		$('#topic3').replaceClassTo("topic3-backward");
		
		$('#topic2').replaceClassTo("topic-down");
		$('#boxplug').replaceClassTo("boxplug-down");
		plugsload();
	}
	//从2跳过来
	if($('#rdohtml').attr("class")=="rdohtml-forward"){
		$('#rdohtml').replaceClassTo("rdohtml-backward");
		$('#rdocss').replaceClassTo("rdocss-backward");
		$('#rdojs').replaceClassTo("rdojs-backward");
		$('#topic').replaceClassTo("topic-up");
		$('#topic2').replaceClassTo("topic2-forward");
		$('#boxplug').replaceClassTo("boxplug-forward");
		plugsload();
	}
	$('#txt-up').val('threeanimate');
	$('#txt-down').val('oneanimate');
}
function plugsload(){
	$('#plugchat').replaceClassTo("plugchat-forward");
	$('#plughome').replaceClassTo("plughome-forward");
	$('#plugfile').replaceClassTo("plugfile-forward");
	$('#plugreadmark').replaceClassTo("plugreadmark-forward");
	$('#plugtel').replaceClassTo("plugtel-forward");
	$('#plugfriends').replaceClassTo("plugfriends-forward");
	$('#plugalbum').replaceClassTo("plugalbum-forward");
	$('#plugposition').replaceClassTo("plugposition-forward");
	$('#plugcomments').replaceClassTo("plugcomments-forward");
	$('#plugearth').replaceClassTo("plugearth-forward");
	$('#plugmessage').replaceClassTo("plugmessage-forward");
	$('#plugimg').replaceClassTo("plugimg-forward");
	$('#plugshare').replaceClassTo("plugshare-forward");
	$('#plugsearch').replaceClassTo("plugsearch-forward");
}
function initplugs(){
	$('#plugchat').replaceClassTo("plugchat");
	$('#plughome').replaceClassTo("plughome");
	$('#plugfile').replaceClassTo("plugfile");
	$('#plugreadmark').replaceClassTo("plugreadmark");
	$('#plugtel').replaceClassTo("plugtel");
	$('#plugfriends').replaceClassTo("plugfriends");
	$('#plugalbum').replaceClassTo("plugalbum");
	$('#plugposition').replaceClassTo("plugposition");
	$('#plugcomments').replaceClassTo("plugcomments");
	$('#plugearth').replaceClassTo("plugearth");
	$('#plugmessage').replaceClassTo("plugmessage");
	$('#plugimg').replaceClassTo("plugimg");
	$('#plugshare').replaceClassTo("plugshare");
	$('#plugsearch').replaceClassTo("plugsearch");
}
function threeanimate(){
	initplugs();
	$('#topic2').replaceClassTo("topic-up");
	$('#boxplug').replaceClassTo("boxplug-up");
	$('#down').replaceClassTo("animate-hide");
	
	$('#phone').replaceClassTo("phone-forward");
	$('#pull').replaceClassTo("pull-forward");
	$('#topic3').replaceClassTo("topic3-forward");
	
	$('#txt-up').val('closewelcome');
	$('#txt-down').val('twoanimate');
}

$("#welcome").swipeToUp({fun:swipeToUp});
$("#welcome").swipeToDown({fun:swipeToDown});
function swipeToUp(){
	var upStr=$('#txt-up').val().toString();
	$phone.runFunction(eval(upStr));
}
function swipeToDown(){
	var downStr=$('#txt-down').val().toString();
	$phone.runFunction(eval(downStr));
}
function closewelcome(){
	A.Controller.section('#follow_section');
	//A.Controller.close&&A.Controller.close();
}
$("#welcome").mouseToUp({fun:swipeToUp});
$("#welcome").mouseToDown({fun:swipeToDown});
function mouseToUp(){
    var upStr=$('#txt-up').val().toString();
    $phone.runFunction(eval(upStr));
}
function mouseToDown(){
    var downStr=$('#txt-down').val().toString();
    $phone.runFunction(eval(downStr));
}
$(document).ready(function(){
	$("#content").css("height",$phone.bodyHeight()+"px");
});
