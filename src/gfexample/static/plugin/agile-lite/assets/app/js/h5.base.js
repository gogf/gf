/*
*	ExMobi5.0 混合 html5 JS 框架之 控件对象类h5.base.js
*	Version	:	1.0.0
*	Author	:	WangMingzhu
*	Email	:	wangmingzhu@nj.fiberhome.com.cn
*   Copyright 2012 (c) 南京烽火星空通信发展有限公司
*/
//扩展jquery方法
(function($){
	/*邮箱自动填充，示例：$("#emailid").addEmailHelp([{domain: '@nj.fiberhome.com.cn', label: '烽火邮箱'},{domain: '@baidu.com', label: '百度邮箱'}])*/
   $.fn.addEmailHelp=function(otherAddress){
		var emails=[{domain: '@qq.com', label: 'qq邮箱'},{domain: '@163.com', label: '163邮箱'},{domain: '@126.com', label: '126邮箱'},{domain: '@hotmail.com', label: 'hotmail邮箱'},{domain: '@sina.com', label:'sina邮箱'},{domain: '@gmail.com',label:'gmail邮箱'},{domain: '@139.com', label: '139邮箱'},{domain:'@yahoo.com.cn', label: 'yahoo中国邮箱'}];
		if(this.length == 0){ 
		    return; 
		}
		this.keyup(function (ev) { 
  			var val = $(this).val(); 
  			var lastInputKey = val.charAt(val.length - 1); 
  			if (lastInputKey == '@') { 
  				var indexOfAt = val.indexOf('@'); 
  				var username = val.substring(0, indexOfAt); 
  				
  				if ($('datalist#emailList').length > 0) { 
  					$('datalist#emailList').remove(); 
  				} 
  		
  				$(this).parent().append('<datalist id="emailList"></datalist>'); 
  				for (var i in emails) { 
  					$('datalist#emailList').append('<option value="' + username + emails[i].domain + '" label="' + emails[i].label + '" />'); 
  				}
  		
  				if(otherAddress != null && typeof otherAddress != 'undefined'){ 
  					for (var i in otherAddress) { 
  						$('datalist#emailList').append('<option value="' + username + otherAddress[i].domain + '" label="' + otherAddress[i].label + '" />'); 
  					}
  				}
  				$(this).attr('list', 'emailList'); 
  	 		}
  		});
    },
	//返回移动运营商的名字，示例$("#phone").mobileISP();
   $.fn.mobileISP=function(){
    	if(!$(this)){
    		return;
    	}
	   var phonenum=$(this).val();
	   if(isNaN(phonenum)==true || phonenum.length<=0){
		   alert("请输入一个数字");
		   return;
		}
	   // 匹配移动手机号
		var PATTERN_CHINAMOBILE = /^1(3[4-9]|5[012789]|8[23478]|4[7]|7[8])\d{8}$/;
		// 匹配联通手机号
		var PATTERN_CHINAUNICOM =/^1(3[0-2]|5[56]|8[56]|4[5]|7[6])\d{8}$/;
		// 匹配电信手机号
		var PATTERN_CHINATELECOM =/^1(3[3])|(8[019])\d{8}$/;
		
		if(PATTERN_CHINATELECOM.test(phonenum)){
			return "chinatelecom";//电信号
		}else if(PATTERN_CHINAUNICOM.test(phonenum)){
			return "chinaunicom";//联通号
		}else if(PATTERN_CHINAMOBILE.test(phonenum)){
			return "chinamobile";//移动号
		}else{
			return "unPhoneNumber";
			alert("请输入一个正确的手机号码");
		}
   }
    //一般用于点击触发动画，执行完动画后再销毁此动画
    $.fn.animateObj=function(anim,thisdisplay){
    	$(this).addClass(anim+' animated-long').one('webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend', function(){
    		$(this).removeClass(anim+' animated-long');
    		if(thisdisplay){
    			$(this).css("display",thisdisplay);
    		}
		});
	}
    //替换class
    $.fn.replaceClassTo=function(reclass){
    	$(this).removeClass();
    	$(this).addClass(reclass);
    }
    //滑动事件
    $.fn.swipeOrientation=function(opts){
    	var id=opts.id;
    	var fun=opts.fun;//执行方法
    	var xRange=opts.xRange||0;//x轴滑动距离
    	var yRange=opts.yRange||0;//y轴滑动距离
    	var orientation=opts.orientation||"left";
    	
    	var xDown=null;
    	var yDown=null;
    	document.getElementById(id).addEventListener('touchstart',function(evt){xDown = evt.touches[0].clientX;yDown=evt.touches[0].clientY;},false);
    	document.getElementById(id).addEventListener('touchmove',function(evt){
    		if(!xDown||!yDown){
    			return;
    		}
    		var xUp=evt.touches[0].clientX;
    		var yUp=evt.touches[0].clientY;
    		var xDiff=xDown-xUp;
    		var yDiff=yDown-yUp;
    		if ( Math.abs( xDiff ) > Math.abs( yDiff ) ) {
    	        if (xDiff>xRange && opts.orientation=="left") {
    	        	fun();
    	        } else if(xDiff<parseInt('-'+xRange) && opts.orientation=="right"){
    	        	fun();
    	        }
    	    } else {
    	        if(yDiff>parseInt(yRange) && opts.orientation=="up") {
    	        	fun();
    	        }else if(yDiff<parseInt('-'+yRange) && opts.orientation=="down"){
    	        	fun();
    	        }
    	    }
    	    xDown = null;
    	    yDown = null;
    	}, false);
    }
    $.fn.swipeToUp=function(opts){
    	var objId=$(this).attr('id').toString();
    	opts.id=objId;
    	opts.orientation="up";
    	$(this).swipeOrientation(opts);
    }
    $.fn.swipeToDown=function(opts){
    	var objId=$(this).attr('id').toString();
    	opts.id=objId;
    	opts.orientation="down";
    	$(this).swipeOrientation(opts);
    }
    $.fn.swipeToLeft=function(opts){
    	var objId=$(this).attr('id').toString();
    	opts.id=objId;
    	opts.orientation="left";
    	$(this).swipeOrientation(opts);
    }
    $.fn.swipeToRight=function(opts){
    	var objId=$(this).attr('id').toString();
    	opts.id=objId;
    	opts.orientation="right";
    	$(this).swipeOrientation(opts);
    }
    //鼠标滑动事件
    $.fn.mouseOrientation=function(opts){
        var id=opts.id;
        var fun=opts.fun;//执行方法
        var xRange=opts.xRange||0;//x轴滑动距离
        var yRange=opts.yRange||0;//y轴滑动距离
        var orientation=opts.orientation||"left";
        
        var xDown=null;
        var yDown=null;
        $('#'+id).mousedown(function(){xDown=event.clientX;yDown=event.clientY;});
        $('#'+id).mouseup(function(){
            if(!xDown||!yDown){
                return;
            }
            var xUp=event.clientX;
            var yUp=event.clientY;
            var xDiff=xDown-xUp;
            var yDiff=yDown-yUp;
            if ( Math.abs( xDiff ) > Math.abs( yDiff ) ) {
                if (xDiff>xRange && opts.orientation=="left") {
                    fun();
                } else if(xDiff<parseInt('-'+xRange) && opts.orientation=="right"){
                    fun();
                }
            } else {
                if(yDiff>parseInt(yRange) && opts.orientation=="up") {
                    fun();
                }else if(yDiff<parseInt('-'+yRange) && opts.orientation=="down"){
                    fun();
                }
            }
            xDown=null;
            yDown=null;
        });
    }
    $.fn.mouseToUp=function(opts){
        var objId=$(this).attr('id').toString();
        opts.id=objId;
        opts.orientation="up";
        $(this).mouseOrientation(opts);
    }
    $.fn.mouseToDown=function(opts){
        var objId=$(this).attr('id').toString();
        opts.id=objId;
        opts.orientation="down";
        $(this).mouseOrientation(opts);
    }
    $.fn.mouseToLeft=function(opts){
        var objId=$(this).attr('id').toString();
        opts.id=objId;
        opts.orientation="left";
        $(this).mouseOrientation(opts);
    }
    $.fn.mouseToRight=function(opts){
        var objId=$(this).attr('id').toString();
        opts.id=objId;
        opts.orientation="right";
        $(this).mouseOrientation(opts);
    }
})(jQuery)

var $phone={
		//打开文件
		openFile:function(){
			$native.openFileSelector();
		},
		//打开日期
		openDateTime:function(){
			$native.openDateTimeSelector();
		},
		//打开二维码扫描
		openBarCode:function(){
			$native.decode();
		},
		//退出应用
		exit:function(){
			$native.exit('是否退出程序？');
		},
		//关闭窗口
		close:function(){
			$native.close();
		},
		//打开exmobi页面
		openExmobiPage:function(href){
			$native.openExMobiPage(href);
		},
		//判断是否是第一次加载
		hasLoad:function(){
			//一般html处理方法
	    	if($.cookie('hasLoad') == null){
				$.cookie('hasLoad','y');
				return false;
	    	}else{
	    		return true;
	    	}
	    },
	    //运行函数
	    runFunction:function(funname){
	    	funname();
	    },
		//自动布局
		initLayout:function() {
			var footerHeight=$("footer").height()==null?"0":$("footer").height();
			var headerHeight=$("header").height()==null?"0":$("header").height();
			$("article").css({"top":headerHeight,"bottom":footerHeight});
		},
		//屏幕高度
		bodyHeight:function(){
			var footerHeight=$("footer").height()==null?"0":$("footer").height();
			var headerHeight=$("header").height()==null?"0":$("header").height();
			var bodyHeight=$(window).height()-footerHeight-headerHeight;
			return bodyHeight;
		},
		//是否第一次加载
		hasLoad:function(){
			//一般html处理方法
			/*var strcookie=document.cookie;
			if(strcookie.indexOf("hasLoad=y")>=0){
				return true;
			}else{
				document.cookie="hasLoad=y";
				return false;
			}*/
			//exmobi处理方法
			if(CacheUtil.getCache("hasLoad")){
				return true;
			}else{
				CacheUtil.setCache("hasLoad","y");
				return false;
			}
		}
};
