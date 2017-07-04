(function(){
	A.event.add({
		exmobiOnstart : function(){
			$(document).on('onstart', function(){
				$(document).trigger(A.options.agileStartEvent);
			});
		}
	});
	
	A.Controller.add({
		nativepage : {
			selector : '[data-toggle="nativepage"]',
			handler : function(hash, el){
				var $el = $(el);
				var isBlank = $el.attr('target')=='_self'?false:true;
	
				$native.openNativePage(hash, isBlank);
			}
		},
		exit : {
			selector : '[data-toggle="exit"]',
			handler : function(){
				$native.exit('是否退出程序？');
			}
		}
	});
	
	A.Component.add({
		datetime : {
			selector : '[data-role="article"].active',
			event : 'articleload',
			handler : function(el, roleType){
				var $el = $(el);		
				var _work = function($el){
					var returnObj = {
						open : function(){
							$el.trigger(A.options.clickEvent);
						},
						clear : function(){
							$el.find('label').html($el.data('placeholder'));
							$el.find('input').val('');
						}
					};
					if(A.Component.isInit($el)){
						return returnObj;
					}
					var $label = $el.find('label');
					var $input = $el.find('input');
					var placeholder = $label.html();$el.data('placeholder', placeholder||'');
					
					$el.on(A.options.clickEvent, function(e){
						$native.openDateTimeSelector({
							mode : $el.data('role'),
							val : $input.val(),
							callback : function(str){
								if($input.val()!=str){
									$label.html(str?str:placeholder);
									$input.val(str||'');
									var _changeFunc = $el.data('change');
									if(!_changeFunc) return;
									var _replace = function(){
										try{ eval(_changeFunc);}catch(e){ console.log(e); };
									};
									_replace.apply($input[0]);
								}
								
							}
						});
						return false;
					});
					$label.html($input.val()||placeholder);
					return returnObj;
				};
	
				if($el.data('role')=='date'||$el.data('role')=='time'){
					return _work($el);
				}else{
					var components = $el.find('[data-role="date"],[data-role="time"]');
					for(var i=0;i<components.length;i++){
						_work($(components[i]));
					}
				}
				
			}
		},
		scanning : {
			selector : '[data-role="article"].active',
			event : 'articleload',
			handler : function(el, roleType){
				var $el = $(el);		
				var _work = function($el){			
					var $label = $el.find('label');
					var $input = $el.find('input');
					var placeholder = $label.html();
					$el.on(A.options.clickEvent, function(e){					
						$native.openDecodeScan(function(str){
							var _iptVal = $input.val();
							$label.html(str?str:placeholder);
							$input.val(str||'');
							if(str&&(_iptVal!=str)&&$el.data('change')){
								eval($el.data('change'));
							}
						});
						return false;
					});
					$label.html($input.val()||placeholder);
				};
	
				if($el.data('role')=='barcode'||$el.data('role')=='qrcode'){
					_work($el);
				}else{
					var components = $el.find('[data-role="barcode"],[data-role="qrcode"]');
					for(var i=0;i<components.length;i++){
						_work($(components[i]));
					}
				}
			}
		},
		file : {
			selector : '[data-role="article"].active',
			event : 'articleload',
			handler : function(el, roleType){
				var $el = $(el);		
				var _work = function($el){			
					var $label = $el.find('label');
					var $input = $el.find('input');
					var placeholder = $label.html();
					$el.on(A.options.clickEvent, function(e){					
						$native.openFileGroupSelector(function(str){
							str = str?str.join(';'):'';
							var flag = false;
							if(str&&($input.val()!=str)&&$el.data('change')){
								flag = true;
							}
							$label.html(str?str:placeholder);
							$input.val(str||'');
							if(flag) eval($el.data('change'));
						});
						return false;
					});
					$label.html($input.val()||placeholder);
				};
	
				if($el.data('role')=='file'){
					_work($el);
				}else{
					var components = $el.find('[data-role="file"]');
					for(var i=0;i<components.length;i++){
						_work($(components[i]));
					}
				}
			}
		}
	});
	
})();



