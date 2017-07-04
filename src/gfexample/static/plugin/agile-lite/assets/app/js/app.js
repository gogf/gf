//启动agile
var $config = {
	exmobiSevice: '${$native.getAppSetting().domain}/process/service/${ClientUtil.getAppId()}'
};
(function(){
	A.event.add({
		beforeunload : function(){
			$(document).on('beforeunload', function(){
				A.Controller.close();
			});
		}
	});
	
	A.Controller.add({
		close : {
			selector : '[data-toggle="close"]',
			handler : function(hash){
				try{
					ExMobiWindow;
					$native.close();
				}catch(e){
					history.go(-1);
				}
				
			}
		},
		html : {
			selector : '[data-toggle="html"]',
			handler : function(hash, el){
				try{
					ExMobiWindow;
					var $el = $(el);
					var isBlank = $el.attr('target')=='_self'?false:true;
					var transition = $el.data('transition')||'';
					$native.openWebview(hash, isBlank, transition);
				}catch(e){
					location.href = A.util.parseURL(hash).getURL();
				}
				
			}
		}
	});
})();
A.launch({
	readyEvent: 'ready', //触发ready的事件，在ExMobi中为plusready
	backEvent: 'backmonitor',
	crossDomainHandler: function(opts) {
		$util.server(opts);
	}
});
$(document).on(A.options.clickEvent, '#ratchet_form_article span', function() {
	A.alert('提示', $(this).attr('class').replace(/.* /, ''));
	return false;
});