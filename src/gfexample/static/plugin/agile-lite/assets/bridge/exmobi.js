var $native = (function(){
	var _native = {};

	var progressbar;
	_native.showMask = function(cb){
		progressbar = progressbar||new ProgressBar(); 
		progressbar.setMessage("加载中");
    	progressbar.show();
		cb&&cb();
	};
	_native.hideMask = function(cb){
		progressbar&&progressbar.cancel();
		cb&&cb();
	};
	
	_native.alert = function(str, cb){
		ExMobiWindow.alert(str, function(){
			cb&&cb();
		});
	};
	
	_native.confirm = function(str, okcb, cancelcb){
		ExMobiWindow.confirm(str, function(){
			okcb&&okcb();
		},function(){
			cancelcb&&cancelcb();
		});
	};
	
	_native.toast = function(str){
		var toast = new Toast();
		var duration = 2000;
		toast.setDuration(duration);
		toast.setText(str);
		toast.show();
		setTimeout(function(){
			toast = null;
		}, duration);
	};
	
	//关闭uixml页面，有id则为某个id的uixml窗口，没id则为当前窗口
	_native.close = function(id){
		if(id){
			nativePage.executeScript('PageUtil.getWindowById("'+id+'").close()');
		}else{
			ExMobiWindow.close();
		}
	};
	
	//要在新的uixml打开webview
	_native.openWebView = _native.openWebview = function(hash, isBlank, transition){
	
		var urlObj = A.util.parseURL(hash);
		
		var url = urlObj.getURL();
		if(urlObj.getProtocol().indexOf('http')==-1){
			var curUrl = location.href;
			var curPaths = curUrl.split('/');
			curPaths[curPaths.length-1] = url;
			url = curPaths.join('/');
		}
		
		var id = urlObj.getFragment()||urlObj.getFilename();
		transition = transition||'slideleft';
		var transitionObj = {
			none:true, slideright:"slideleft", slideleft:"slideright", slidedown:"slideup", slideup:"slidedown", zoom:true, fade:true,curlup:true
		};
	
		var style = (transitionObj[transition]?"openanimation:"+transition+";closeanimation:"+(transitionObj[transition]==true?transition:transitionObj[transition])+";":"");
	
		var html = [];
		html.push('<html id="'+id+'" isbridge="true" style="'+style+'">');
		html.push('<head>');
		html.push('<meta charset="UTF-8"/>');
		html.push('<title show="false">agile</title>');
		
		html.push('<script>');
		html.push('<![CDATA[');
		html.push(']]>');
		html.push('</script>');
		html.push('</head>');
		html.push('<body style="margin:0px;padding:0px;">');
		html.push('<webview id="browser" url="'+url.replace(/\&/g,'&amp;')+'" backmonitor="true"/>');// backMonitor="true"
		html.push('</body>');
		html.push('</html>');
		
		ExMobiWindow.openData(html.join(''), isBlank, false, '', urlObj.getQuerystring());
	};
	
	//要打开uixml页面
	_native.openNativePage = function(hash, isBlank){
		ExMobiWindow.open(hash, isBlank);
	};
	
	//打开日期时间控件
	_native.openDateTimeSelector = function(opts){
		var options = {
			mode : 'date',
			val : '',
			callback : null
		};
		$.extend(options, opts);
		var  timewindow = new TimeWindow();
		timewindow.mode = options.mode;
		timewindow.initialvalue = options.val;
		timewindow.onCallback = function(){
			var str = timewindow.isSuccess()?timewindow.result:null;
			options.callback&&options.callback(str);
			timewindow = null;
		};
		timewindow.startWindow();
	};
	
	//打开扫码框
	_native.openDecodeScan = function(callback){
		var decode = new Decode();
		decode.onCallback = function(){
			var str = decode.isSuccess()?decode.result:null;
			callback&&callback(str);
			decode = null;
		};
		decode.startDecode();
	};
	
	//文件选择
	_native.openFileSelector = function(opts){
		var options = {
			filter : null,
			path : null,
			callback : null
		};
		$.extend(options, opts);
		var filechoice = new FileChoice();
		if(options.filter) filechoice.filter = options.filter;
		if(options.path) filechoice.defaultPath = options.path;
		filechoice.onCallback = function(path){
			options.callback&&options.callback(path);
			filechoice = null;
		};
		filechoice.start();
	
	};
	
	//图片选择
	_native.openImgSelector = function(callback){
		var imageChoice = new ImageChoice();
		imageChoice.onCallback = function(path){
			callback&&callback(imageChoice.getFilePaths());
			imageChoice = null;
		};
		imageChoice.start();
	};
	
	//拍照
	_native.openCameraSelector = function(opts){
		var options = {
			mode : 'still',
			callback : null
		};
		$.extend(options, opts);
		var camerawindow = new CameraWindow();
		camerawindow.mode = options.mode;
	    camerawindow.onCallback = function(){
	    	var str = camerawindow.isSuccess()?camerawindow.value:null;
	    	options.callback&&options.callback(str);
	    	camerawindow = null;
	    };
	    camerawindow.startCamera();
	
	};
	
	
	//打开文件选择
	_native.openFileGroupSelector = function(callback){
		A.actionsheet(
			[
				{
					text : '拍照',
				    handler : function(){
				    	$native.openCameraSelector({
				    		callback : function(path){
				    			callback&&callback(path?[path]:null);
				    		}
				    	});
					}
				},{
					text : '照片',
				    handler : function(){
				    	$native.openImgSelector(function(paths){
				    		callback&&callback(paths);
				    	});
					}
				},{
				    text : '文件',
				    handler : function(){
				    	$native.openFileSelector({
				    		callback : function(path){
				    			callback&&callback(path?[path]:null);
				    		}
				    	});
				    }
				}
			], 'alizarin');
	};
	_native._appinfo;
	_native.getAppSetting = function(){
		if($native._appinfo) return $native._appinfo;
		try{
			var obj = {};	
			if(typeof ClientUtil=='undefined'){
				var protocol = location.protocol;			
				var port = location.port;
				var ip = location.host.replace(':'+port,'');
				var domain = protocol+'//'+ip+':'+port;
				obj.ip = ip;
				obj.port = port;
				obj.domain = domain;
			}else{
				
				obj = ClientUtil.getSetting();
				
				var domain = 'http://'+obj.ip+':'+obj.port;
				obj.domain = domain;
			}
			$native._appinfo = obj;
			return obj;
		}catch(e){
			return {};
		}
		
	};
	
	_native.getParameter = function(k){
		return ExMobiWindow.getParameter(k);
	};
	
	_native.getParameters = function(){
		return ExMobiWindow.getParameters();
	};
	
	_native.session = function(){
		if(arguments.length==1){
			try{
				return JSON.parse(ExMobiWindow.getStringSession(arguments[0]));
			}catch(e){
				return ExMobiWindow.getStringSession(arguments[0]);
			}
			
		}else if(arguments.length==2){
			var v = arguments[1]||'';
			try{
				return ExMobiWindow.setStringSession(arguments[0], JSON.stringify(v));
			}catch(e){
				return ExMobiWindow.setStringSession(arguments[0], v);
			}
			
		}else{
			return null;
		}
	
	};
	
	//初始化cache
	_native.cache = function(){
		var k,v;
		k = arguments[0];
		if(arguments.length==2){
			v = arguments[1]||'';
			try{
				CacheUtil.setCache(k, A.JSON.stringify(v));
				//return true;
			}catch(e){
				CacheUtil.setCache(k, v.toString());
				//return false;
			}
		}else{
			try{
				return A.JSON.parse(CacheUtil.getCache(k));
			}catch(e){
				return CacheUtil.getCache(k);
			}
		}
	};
	
	_native.exit = function(msg){
		if(msg){
			ClientUtil.exit(msg);
		}else{
			ClientUtil.exitNoAsk();
		}
	
	};
	
	return A.util.readyAlarm(_native, '$native', 'plusready');
})();

var $util = (function(){
	var _util = {};
	
	_util._cacheMap = {
		index : 0
	};
	
	_util.queryJSON = function(query){
		if(!query){
			return {};
		}
		try{
			var properties = query.replace(/&/g, "',").replace(/=/g, ":'") + "'";  
		    var obj = null;  
		    var template = "var obj = {p}";  
		    eval(template.replace(/p/g, properties));
		}catch(e){
			return {};
		}
	    
	    return obj==null?{}:obj;
	};
	
	_util._showPageLoading = false;
	_util.go = function(opts, handler){
		if(!opts||!opts.url) return;
		opts.url = A.util.script(opts.url);	
		var ajaxData = {};
		ajaxData.url = opts.url;
		ajaxData.method = opts.type = opts.type&&opts.type.toLowerCase()=='post'?'post':'get';
		if(opts.data) ajaxData.data = opts.data;
		ajaxData.successFunction = '$util._ajax_successFunction';
		ajaxData.failFunction = '$util._ajax_errorFunction';
		if(opts.headers) ajaxData.requestHeader = opts.headers;
		ajaxData.isBlock = opts.isBlock = opts.isBlock==true?true:false;
		ajaxData.timeout = opts.timeout?(opts.timeout/1000):20;
		ajaxData.reqCharset = opts.reqCharset||'utf-8';
		var ajax = new handler(ajaxData);
		
		var index = $util._cacheMap.index++;
	
		$util._cacheMap['_ajax_opts_key_'+index] = opts;
		
		ajax.setStringData('_ajax_opts_key_', index);
		
		$util._showPageLoading = A.options.showPageLoading||$util._showPageLoading;
		
		if(!ajaxData.isBlock&&$util._showPageLoading) A.showMask();
		
		ajax.send();
		
	};
	
	_util._ajax_successFunction = function(data){
		var opts = $util._ajax_getFunction(data);
		if(typeof opts.result=='undefined'){
			opts.error&&opts.error(data, '500');
		}else{
			opts.success&&opts.success(opts.result);
		}
	};
	
	_util._ajax_errorFunction = function(data){
		var opts = $util._ajax_getFunction(data);
		opts.error&&opts.error(data, data.status);
	};
	
	_util._ajax_getFunction = function(ajax){
		var result = ajax.responseText;	
		var index = ajax.getStringData('_ajax_opts_key_');
		
		var opts = $util._cacheMap['_ajax_opts_key_'+index];
		opts.dataType = opts.dataType&&opts.dataType.toLowerCase()=='json'?'json':'text';
		if(opts.dataType=='json'){
			try{
				opts.result = eval('('+result+')');
			}catch(e){
				delete opts.result;
			}
		}else{
			opts.result = ajax.responseText;
		}
		delete $util._cacheMap['_ajax_opts_key_'+index];
	
		if(!opts.isBlock&&$util._showPageLoading) A.hideMask();
		
		return opts;
	};
	
	//对应ExMobi的Ajax
	_util.server = function(opts){
		$util.go(opts, Ajax);
	};
	//对应ExMobi的DirectAjax
	_util.ajax = function(opts){
		$util.go(opts, DirectAjax);
	};
	//对应ExMobi的DirectFormSubmit 
	_util.form = function(opts, handler){
		if(!opts||!opts.url) return;
		opts.type = opts.type||'post';
		if(opts.data){
			var dataArr = $util.paramsToJSON(opts.data);
			var fileElementId = typeof opts.fileElementId=='object'?opts.fileElementId.join():fileElementId;
			fileElementId = ','+(fileElementId?fileElementId:'')+',';
			
			for(var i=0;i<dataArr.length;i++){
				var obj = {};
				for(var k in dataArr[i]){
					obj.name = k;
					obj.value = dataArr[i][k];
				}
				obj.type = fileElementId.indexOf(','+obj.name+',')==-1?0:1;
				dataArr[i] = obj;
			};
			opts.data = dataArr;
		}
		
		$util.go(opts, handler||DirectFormSubmit);
	};
	
	//对应ExMobi的FormSubmit 
	_util.serverForm = function(opts){
		_util.form(opts, FormSubmit);
	};
	
	_util.paramsToJSON = function(data){
		var arr = [];
		if(!data||(typeof data!='string')) return [];
	
		data.replace(/(\w+)=(\w+)/ig, function(a, b, c){ 
			var obj = {};
		    obj[b] = unescape(c); 
		    arr.push(obj);
		});  
			
		return arr;
	};

	return A.util.readyAlarm(_util, '$util', 'plusready');
	
})();