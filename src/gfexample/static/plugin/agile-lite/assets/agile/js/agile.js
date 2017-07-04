/*
*	Agile Lite 移动前端框架
*	Version	:	3.0.0 beta
*	Author	:	nandy007
*   License MIT @ https://git.oschina.net/nandy007/agile-lite
*/
window.PointerEvent = undefined;//解决iscroll在chrome55.X版本的bug
var agilelite = (function($){
	var Agile = function(){
		this.$ = $;
		this.options = {
			version : '3.0.0',
			clickEvent : ('ontouchstart' in window)?'tap':'click',
			agileReadyEvent : 'agileready',
			agileStartEvent : 'agilestart', //agile生命周期事件之start，需要宿主容器触发
			readyEvent : 'ready', //宿主容器的准备事件，默认是document的ready事件
			backEvent : 'backmenu', //宿主容器的返回按钮
			complete : false, //是否启动完成
			crossDomainHandler : null, //跨域请求的处理类
			showPageLoading : false, //ajax默认是否有loading界面
			viewSuffix : '.html', //加载静态文件的默认后缀
			lazyloadPlaceholder : '', //懒人加载默认图片
			usePageParam : true,
			autoInitCom : true,
			iScrollOptions : {},
			classPre : 'agile-'
		};
		
		this.pop = {
			hasAside : false,
			hasPop : false
		};
	};
	
	var _launchMap = {};
	
	/**
     * 注册Agile框架的各个部分，可扩充
     * @param {String} 控制器的唯一标识
     * @param {Object} 任意可操作的对象，建议面向对象方式返回的对象
     */
	Agile.prototype.register = function(key, obj){
		this[key] = obj;
		if(obj.launch){
			_launchMap[key] = obj.launch;
		}
	};
	var _doLaunch = function(){
		for(var k in _launchMap){ try{ _launchMap[k](); }catch(e){ console.log(e);} }
		agilelite.options.complete = true;
		$(document).trigger(agilelite.options.agileReadyEvent);
	};
	
	/**
     * 启动Agile
     * @param {Object} 要初始化的一些参数
     */
	Agile.prototype.launch = function(opts){
		if(agilelite.options.complete==true) return;
		_initPageInfo.apply(this);
		$.extend(this.options, opts);
		var _this = this;
		if(!this.options.readyEvent){
			_doLaunch();
		}else if($(document)[this.options.readyEvent]){
			$(document)[this.options.readyEvent](_doLaunch);
		}else{
			$(document).on(this.options.readyEvent, _doLaunch);
		}
	};
	
	var _initPageInfo = function(){
		var urlObj = agilelite.util.parseURL(location.href);
		var params = agilelite.JSON.stringify(urlObj.getQueryobject());
		var hash = location.hash.replace('#','');
		this.pageInfo = agilelite.util.decodeHash(hash);
	};
	
	Agile.prototype.noConflict = function(){ return this; };
	
	return window.A ? new Agile() : window.A = new Agile();
})(window.Zepto||window.jQuery||window.$);

//控制器
(function($){
	/**
     * 控制器的基本结构，形如{key:{selector:'控制器选择器'，默认为body',handler:function($trigger){}}}
     * selector为选择器，handler为处理函数
     * @private
     */
	var _controllers = {
		default : {//默认控制器
			selector : '[data-toggle="view"]',
			handler : function(hash, el){			
				$el = $(el);				
				var toggleType = $el.data('toggle');
				var urlObj = agilelite.util.parseURL(hash);
				var hashObj = urlObj.getHashobject();				
				var controllerObj = _controllers[toggleType]||{};				
				var $target = $(hashObj.tag);
				var $container = $(controllerObj.container);				
				function _event($target, $current){									
					var targetRole = $target.data('role')||'';					
					var show = function($el){
						if(!$el.hasClass('active')) $el.addClass('active');					
						if(!$el.attr('__init_controller__')){												
							$el.attr('__init_controller__', true);
							$el.trigger(targetRole+'load');
						}			
						$el.trigger(targetRole+'show');
					};
					
					var hide = function($el){
						$el.removeClass('active').trigger(targetRole+'hide');
					};
					if(controllerObj.isToggle||$el.data('istoggle')){
						if($target.hasClass('active')){
							hide($target);
						}else{
							show($target);
						}
						return;
					}

					show($target);
					if($current&&$current.length>0) {
						hide($current);
					}
				}
				
				function _setDefaultTransition($el){
					var targetTransition = $el.data('transition')||controllerObj.transition;
					if(targetTransition) $el.data('transition', targetTransition);
				};
				
				function _next(){														
					var targetRole = $target.data('role');
					var toggleSelector = targetRole?'[data-role="'+targetRole+'"].active':'.active';
					$target.data('url', urlObj.getURL());
					if(urlObj.getQueryobject()) $target.data('params', agilelite.JSON.stringify(urlObj.getQueryobject()));
					if($target.hasClass('active')){
						_event($target);
						controllerObj.complete&&controllerObj.complete($target, {result:'thesame'});
					}else{
						var $current = $target.siblings(toggleSelector);
						_setDefaultTransition($current);_setDefaultTransition($target);
						agilelite.anim.change($current, $target, false, function(){				
							_event($target, $current);
						});	
						controllerObj.complete&&controllerObj.complete($target, {result:'success'});
					}
				}
				
				if($target.length==0){
					if($container.length==0){
						controllerObj.complete&&controllerObj.complete($target, {result:'nocontainer'});
						return;
					};
					if($el.attr('data-useTemplate')){
						$target = $(controllerObj.template.replace('{id}', hashObj.tag.replace('#', '')));
						$container.append($target);
						_next();
						return;
					}
					if(agilelite.options.showPageLoading) agilelite.showMask();				
					agilelite.util.getHTML(urlObj.getURL(), function(html){
						if(agilelite.options.showPageLoading) agilelite.hideMask();
						if(!html){
							controllerObj.complete&&controllerObj.complete($target, {result:'requesterror'});
							return;
						}
						$target = $(html);
						$container.append($target);
						$target = $(hashObj.tag);
						_next();
					});
				}else{
					_next();
				}

			}
		},
		html : {//多页模式请复写此控制器
			selector : '[data-toggle="html"]',
			handler : function(hash){
				var urlObj = agilelite.util.parseURL(hash);
				location.href = urlObj.getURL();
			}
		},
		section : {
			selector : '[data-toggle="section"]',
			container : '#section_container',
			transition : 'slide',
			template : '<section id="{id}" data-role="section"></section>',
			history : [],
			complete : function($target, msg){
				var _add2History = function(hash,noState){
			    	var hashObj = agilelite.util.parseURL(hash).getHashobject();
			    	var _history = _controllers.section.history;
			    	if(_history.length==0||(_history.length>1&&$(_history[0].tag).data('cache')==false||$(_history[0].tag).data('cache')=='false')){
			    		noState = true;
			    	}
			    	var encodeHash = hash;
			    	if(agilelite.options.usePageParam){
			    		var encodeHashObj = {
				    		params : agilelite.Component.params(hash),
				    		hash : hash.replace('#', ''),
				    		url : $(hash).data('url')
				    	};
				    	encodeHash = '#'+agilelite.util.encodeHash(agilelite.JSON.stringify(encodeHashObj));
			    	}

			        if(noState){//不添加浏览器历史记录
			            _history.shift(hashObj);
			            window.history.replaceState(hashObj,'',encodeHash);
			        }else{
			            window.history.pushState(hashObj,'',encodeHash);
			        }
			        _history.unshift(hashObj);
			    };
				var hash = '#'+$target.attr('id');
				msg = msg||{};
				if(msg.result=='thesame'){
					_add2History(hash ,true);
				}else if(msg.result=='success'){
					_add2History(hash ,false);
				}
			}
		},
		article : {
			selector : '[data-toggle="article"]',
			container : '[data-role="section"].active',
			template : '<article id="{id}" data-role="article"></article>'
		},
		modal : {
			selector : '[data-toggle="modal"]',
			container : 'body',
			template : '<div id="{id}" data-role="modal" class="modal"></div>',
			isToggle : true
		},
		aside : {
			selector : '[data-toggle="aside"]',
			container : '#aside_container',
			handler : function(hash){
				if(hash){
					agilelite.Aside.show(hash);
				}else{
					agilelite.Aside.hide();
				}
				
			}
		},
		back : {
			selector : '[data-toggle="back"]',
			handler : function(hash){
				if(arguments.length==0||typeof hash=='string'||typeof hash=='undefined'||typeof hash=='object'){
					$(document).trigger(agilelite.options.backEvent, typeof hash=='object'?[hash]:null);
					return;
				}
				var _history = _controllers.section.history;				
				var codeHashObject = agilelite.options.usePageParam?agilelite.JSON.parse(agilelite.util.decodeHash(location.hash.replace('#', ''))):{hash:location.hash};
				if(!codeHashObject||('#'+codeHashObject.hash==_history[0].tag)){
					return;
				}

				if(_history.length<2) return;
				var $current = $(_history.shift().tag);
		    	var $target = $(_history[0].tag);
		    	var backParams = $current.data('back-params');
		    	if(backParams) $target.data('params', backParams);
		    	$(document).trigger('sectiontransition', [$current, $target]);
			}
		},
		page : {
			selector : '[data-toggle="page"]',
			handler : function(el){
				var $el = $(el);
				var $scroll = $el.parent();
				var $parent = $scroll.parent();
				var _index = $scroll.children('[data-role]').index($el);
				var slider = agilelite.Slider($parent);
				slider.index(_index);
			}
		},
		scrollTop : {
			selector : '[data-toggle="scrollTop"]',
			handler : function(el){
				var $el = $(el);
				if($el.data('scroll')){
					var scroll = agilelite.Scroll(el);
					scroll.scrollTo(0, 0, 1000, IScroll.utils.ease.circular);
				}else{
					$el.animate({ scrollTop: 0 }, 200);
				}
				
			}
		}
	};	
		
	var controller = {};//定义控制器
	/**
     * 添加控制器
     * @param {Object} 控制器对象，形如{key:{selector:'',handler:function(){}}}
     */
	controller.add = function(c){	
		$.extend(_controllers, c);
	};
	
	/**
     * 获取全部控制器
     * @param {String} 控制器的key，如果有key则返回当前key的控制器，没有key则返回全部控制器
     */
	controller.get = function(key){
		return key?_controllers[key]:_controllers;
	};	
	
	/**
     * 为所有控制器创建调用方法
     * @private 只能初始化一次，启动后添加的不生效
     */
	var _makeHandler = function(){
		for(var k in _controllers){			
			(function(k){
				//定义JS调用函数
				controller[k] = function(hash, el){				
					var toggleType = k;
					var $el = el?$(el):$('<a data-toggle="'+k+'" href="'+hash+'"></a>');
					var curr = _controllers[toggleType]?(_controllers[toggleType].handler?toggleType:'default'):'default';
					_controllers[curr].handler(hash, $el);
				};
				//定义点击触发事件
				$(document).on(agilelite.options.clickEvent, _controllers[k]['selector'], function(){
					$('input:focus,textarea:focus').blur();
					var k = $(this).data('toggle');
					var hash = $(this).attr('href')||'#';
					(controller[k]||controller['default'])(hash, $(this));
					return false;
				});
			})(k);
		}
	};
	
	/**
     * 启动控制器，如果需要在启动agile的时候启动，则函数名必须叫launch
     */
	controller.launch = function(){
		_makeHandler();
	};

	agilelite.register('Controller', controller);
})(agilelite.$);

//组件
(function($){
	var component = {};
	
	var _components = {
		default : {
			selector : '[data-role="view"]',
			handler : function(el, roleType){				
				var $el = $(el);
				roleType = roleType=='default'?$el.data('role'):roleType;
				var componentObj = _components[roleType]||{};
				if(componentObj.isToggle||$el.data('istoggle')){
					$el[$el.hasClass('active')?'removeClass':'addClass']('active');
					return;
				}else if($el.hasClass('active')){
					return;
				}
				
				var $current;
				var curSelector = '[data-role="'+roleType+'"].active';
				if(componentObj.container){
					$current = $el.parents(componentObj.container).first().find(curSelector);
				}else{
					$current = $el.siblings(curSelector+'.active');
				}
				$el.addClass('active');
				$current.removeClass('active');
			}
		},
		tab : {
			selector : '[data-role="tab"]',
			handler : function(el, roleType){
				var $el = $(el);
				var toggleType = $el.data('toggle');
				if((!toggleType)||toggleType=='section'||toggleType=='modal'||toggleType=='html'){
					return;
				}
				_components['default'].handler(el, roleType);
			}
		},
		scroll : {
			selector : '[data-role="article"].active',
			event : 'articleload',
			handler : function(el, roleType){
				var options = {
					verticle : { },
					horizontal : {
						scrollbars : false, scrollX : true, scrollY : false, bindToWrapper : true
					},
					scroll : {
						scrollX : true, scrollY : true
					},
					pulldown : {},
					pullup : {},
					pull : {}
				};

				var _doScroll = function($scroll){
					var opts = options[$scroll.data('scroll')];					
					agilelite.Scroll('#'+$scroll.attr('id'), opts).refresh();
				};
				
				var _doRefresh = function($el){
					agilelite.Refresh('#'+$el.attr('id'), options[$el.data('scroll')]).refresh();
				};
				
				var $el = $(el||this);
				var scrollType = $el.data('scroll');
				
				if(scrollType=='verticle'||scrollType=='horizontal'||scrollType=='scroll') _doScroll($el);	
							
				if(scrollType=='pulldown'||scrollType=='pullup'||scrollType=='pull') _doRefresh($el);
				
				var scrolls = $el.find('[data-scroll="verticle"],[data-scroll="horizontal"],[data-scroll="scroll"]');
				for(var i=0;i<scrolls.length;i++){
					(function($this){
						setTimeout(function(){
							_doScroll($this);
						}, 100);
					})($(scrolls[i]));
					
				}
				
				scrolls = $el.find('[data-scroll="pulldown"],[data-scroll="pullup"],[data-scroll="pull"]');
				for(var i=0;i<scrolls.length;i++){
					(function($this){
						setTimeout(function(){
							_doRefresh($this);
						}, 100);
					})($(scrolls[i]));
				}
			}
		},
		formcheck : {
			selector : '[data-role="article"].active',
			event : 'articleload',
			handler : function(el, roleType){
				
				var $el = $(el);	
				var _doInit = function($el){	
					$el.removeAttr('href');					
					$el.on(agilelite.options.clickEvent, function(e){
						var isInScroll = $el.data('__isinscroll__');
						if(typeof isInScroll=='undefined'){
							isInScroll = $el.closest('[data-scroll]').length==0?false:true;
							$el.data('__isinscroll__', isInScroll);
						}
						var tag = e.target.tagName.toLowerCase();
			    		if((isInScroll&&tag!='input')||(!isInScroll&&tag!='input'&&tag!='label')){
			    			var checkObj = $el.find('input')[0];
			    			if(checkObj.type=='radio'){
			    				checkObj.checked = true;
			    			}else{
			    				checkObj.checked = !checkObj.checked;
			    			}
			    		}
			    	});
				};
				
				if($el.data('role')=='checkbox'||$el.data('role')=='radio'){
					_doInit($el);
				}else{
					var els = $el.find('[data-role="checkbox"],[data-role="radio"]');
					for(var i=0;i<els.length;i++){
						_doInit($(els[i]));
					}
				}
			}
		},
		toggle : {
			selector : '[data-role="article"].active',
			event : 'articleload',
			handler : function(el, roleType){
				var $el = $(el);			
			    var _getVal = function(isActive, $el){
			    	var onValue = $el.data('on-value');
				    var offValue = $el.data('off-value');
				    var val = isActive?onValue:offValue;
				    return typeof val=='undefined'?'':val;
			    };
				var _doToggle = function($el){	
					var isActive = $el.hasClass('active');	
					if($el.find('div.toggle-handle').length>0){//已经初始化	
						$el[isActive?'removeClass':'addClass']('active');
				    	$el.find('input[data-com-refer="toggle"]').val(_getVal(isActive?false:true, $el));
			            return;
			        }
			        var name = $el.attr('name');
			        //添加隐藏域，方便获取值
			        $el.append('<input type="hidden" data-com-refer="toggle" '+(name?'name="'+name+'" ':'')+' value="'+_getVal(isActive, $el)+'"/>');
			        $el.append('<div class="toggle-handle"></div>');
			        $el.on(agilelite.options.clickEvent, function(){
			            var $t = $(this),v = _getVal($el.hasClass('active')?false:true, $el);
			            $t.toggleClass('active').trigger('datachange', [v]);//定义toggle事件
			            $t.find('input').val(v);
			            return false;
			        });
				};
				
				if($el.data('role')=='toggle'){
					_doToggle($el);
				}else{
					var toggles = $el.find('[data-role="toggle"]');
					for(var i=0;i<toggles.length;i++){
						_doToggle($(toggles[i]));
					}
				}
			}
		},
		'lazyload' : {
			selector : '[data-role="article"].active',
			event : 'articleload',
			handler : function(el, roleType){
				var $el = $(el);
		    	var _doLazyload = function($el){
		    		var placeholder = $el.attr('placeholder')||agilelite.options.lazyloadPlaceholder;
			    	if(!$el.attr('src')&&placeholder) $el.attr('src', placeholder);
					var source = agilelite.util.script($el.data('source'));
			    	var eTop = $el.offset().top;//元素的位置
			    	var validateDistance = 100;
			    	var winHeight = $(window).height()+validateDistance;
			    	if(eTop<0||eTop>winHeight) return;
			    	
			    	var type = $el.data('type');
			    	if(type=='base64'){
			    		agilelite.ajax({
			    			url : source,
			    			success : function(data){
			    				_injectImg($el, data);
			    			},
			    			isBlock : false
			    		});
			    	}else{
			    		var img = new Image();
			    		img.src = source;
			    		if(img.complete) {
			    			_injectImg($el, source);
			    			img = null;
			    	    }else{
			    	    	img.onload = function(){
			    				_injectImg($el, source);
			    				img = null;
			        		};
			        		img.onerror = function(){
			    				_injectImg($el, placeholder||'');
			    				img = null;
			        		};
			    	    }
			    	}
		    	};
		    	
		    	var _injectImg = function($el, data){
		    		if(!$el.data('source')) return;
		    		agilelite.anim.run($el,'fadeOut', function(){
		    			$el.css('opacity', '0');		    				    			
		    			$el[0].onload = function(){			
		    				agilelite.anim.run($el,'fadeIn', function(){
		    					$el.css('opacity', '1');
		    					var $sd = $el.closest('[data-role="slider"]');
		    					if($sd.hasClass('active')) agilelite.Slider($sd).refresh();
		    					agilelite.Component.scroll($el.closest('[data-scroll]'));
		    				});
		        		};
		        		$el.attr('src', data);	
		    		});
		    		$el.removeAttr('data-source');
		    	};

		    	if($el.data('source')){
		    		_doLazyload($el);
		    	}else{
		    		var lazyloads = $el.find('img[data-source]');
		    		for(var i=0;i<lazyloads.length;i++){
		    			_doLazyload($(lazyloads[i]));
		    		}
		    	}		    	
			}
		},
		scrollTop : {
			selector : '[data-role="article"].active',
			event : 'articleload',
			handler : function(el, roleType){				
				var _work = function($el){
					$el.on(agilelite.options.clickEvent, function(){ $el.removeClass('active');});
					$target = $($el.attr('href'));
					if($target.data('scroll')){
						var scroll = agilelite.Scroll($el.attr('href'));
						scroll.on('scrollEnd', function(){
							$el[this.y<-120?'addClass':'removeClass']('active');
						});
					}else{
						var scroll = $target;
						scroll.on('scrollEnd', function(){
							$el[scroll.scrollTop()>120?'addClass':'removeClass']('active');
						});
						agilelite.util.jquery.scrollEnd(scroll);
					}
				};
				var $el = $(el);
				if($el.data('role')=='scrollTop'){
					_work($el);
				}else{
					var comps = $el.find('[data-role="scrollTop"]');
					for(var i=0;i<comps.length;i++){
						_work($(comps[i]));
					}
				}
			}
		},
		pictureShow : {
			type : 'function',
			handler : function(opts){
				var id = opts.id;
				var index = opts.index||0;
				var list = opts.list;
				var htmlArr = [];
				htmlArr.push('<section id="'+id+'" data-role="section" data-cache="false" style="background:#000;">');
				htmlArr.push('<header class="picture_show_header">');
				htmlArr.push(opts.header||'<a data-toggle="back" href="#" style="float:right;color:#fff;text-align:center; padding:8px; font-size:18px;">X</a>');
				htmlArr.push('</header>');
				htmlArr.push('<article id="'+id+'_article" data-role="article" class="active" style="overflow:hidden;background:rgba(255, 255, 255, 0.1) none repeat scroll 0 0 !important;">');
				htmlArr.push('<div id="'+id+'_slide" data-role="slider" class="full" style="overflow: hidden;"><div class="scroller">');
				for(var i=0;i<list.length;i++){
					htmlArr.push('<div class="slide" style="padding-top:100px;">');
					htmlArr.push('<img data-source="'+list[i].imgURL+'" class="full-width"/>');
					htmlArr.push('</div>');
				}
				htmlArr.push('</div></div>');
				
				htmlArr.push('</article>');
				
				htmlArr.push('<div class="picture_show_footer" style="position:relative;bottom:0px;left:0px;right:0px;padding:8px;color:#fff;min-height:140px;background:#111;">');
				htmlArr.push('<div><span style="float:left;font-size:20px;">'+opts.title+'</span><p class="picture_show_num" style="float:right;"></p></div>');
				htmlArr.push('<div style="clear:both;line-height:20px;"><span class="picture_show_content"></span></div>');
				htmlArr.push('</div>');

				htmlArr.push('</section>');
				$('#section_container').append(htmlArr.join(''));
				agilelite.Controller.section('#'+id);
				var $num = $('#'+id+' .picture_show_num');
				var $content = $('#'+id+' .picture_show_content');
				var $header = $('#'+id+' .picture_show_header');
				var $footer = $('#'+id+' .picture_show_footer');
				var $slider = agilelite.Slider('#'+id+'_slide', {
					dots : 'hide',
					change : function(i){
						$num.html('<span style="font-size:18px;">'+(i+1)+'</span>/'+list.length);
						$content.html(list[i].content||opts.title);
					}
				});
				$slider.index(index);
				$('#'+id+'_article').on(agilelite.options.clickEvent, function(){
					$header.toggle();
					$footer.toggle();
					return false;
				});
			}
		}
	};
	
	/**
     * 添加组件
     * @param {Object} 控制器对象，形如{key:{selector:'',handler:function(){}}}
     */
	component.add = function(c){
		$.extend(_components, c);
	};
	
	/**
     * 获取全部组件
     * @param {String} 组件的key，如果有key则返回当前key的组件，没有key则返回全部组件
     */
	component.get = function(key){
		return key?_components[key]:_components;
	};
	
	/**
     * 获取全部组件
     * @param 要初始化的组件选择器，可以是$对象
     */
	component.initComponents = function(selector){
		var $el = $(selector);
		for(var k in _components){
			if(_components[k].type=='function'||k=='default') continue;
			agilelite.Component[k]($el);
		}
	};
	
	
	/**
     * 初始化组件
     * @private 仅能调用一次
     */
	_makeHandler = function(){
		for(var k in _components){	
			if(_components[k].type=='function'){
				component[k] = _components[k].handler;
				continue;
			}	
			(function(k){	
				//定义JS调用函数
				component[k] = function(hash){
					var roleType = k;
					var curr = _components[roleType]?(_components[roleType].handler?roleType:'default'):'default';
					return _components[curr].handler(hash, roleType);
				};
				//定义触发事件	
				if(!_components[k]['selector']) return;
				$(document).on(_components[k].event||agilelite.options.clickEvent, _components[k]['selector'], function(){					
					var curr = _components[k].handler?k:'default';				
					_components[curr].handler(this, k);
					return false;
				});
				
			})(k);
		}

	};
	/**
     * 获取控制器传给组件的参数,hash为被控制的组件的#id
     */
	component.params = function(hash){
		return agilelite.JSON.parse($(hash).data('params'))||{};
	};
	
	/**
     * 设置并返回原始初始化状态
     */
	component.isInit = function($el){
		if($el.data('com-init')){
			return true;
		}else{
			$el.data('com-init','init');
			return false;
		}
	};
	
	/**
     * type组件的名称，hash组件的#id
     */
	component.getObject = function(type, hash){
		var $el = $(hash);
		var returnObj = {};
		var comObj = component.get(type);
		if(!comObj||!comObj.extend) return returnObj;
		var _extend = comObj.extend;
		for(var k in _extend){
			(function(k){
				returnObj[k] = function(){
					return _extend[k].apply($el, arguments);
				};
			})(k);
		}
		return returnObj;
	};
	
	//启动组件
	component.launch = function(){
		_makeHandler();
	};
	
	agilelite.register('Component', component);
})(agilelite.$);

//动画封装
(function($){
	
	var anim = {};
	
	anim.classes = {//形如：[[currentOut,targetIn],[currentOut,targetIn]]
		empty : [['empty','empty'],['empty','empty']],
		slide : [['slideLeftOut','slideLeftIn'],['slideRightOut','slideRightIn']],
		cover : [['','slideLeftIn'],['slideRightOut','']],
		slideUp : [['','slideUpIn'],['slideDownOut','']],
		slideDown : [['','slideDownIn'],['slideUpOut','']],
		popup : [['','scaleIn'],['scaleOut','']],
		flip : [['slideLeftOut','flipOut'],['slideRightOut','flipIn']]
	};
	
	anim.addClass = function(cssObj){
		$.extend(anim.classes, cssObj||{});
	};
	
	/**
     * 添加控制器
     * @param {Object} 要切换动画的jQuery对象
     * @param {String} 动画样式
     * @param {Number} 动画时间
     * @param {Function} 动画结束的回调函数
     */
    
	anim.run = function($el, cls, cb){
		var $el = $($el);
		var _eName = 'webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend';
		if($el.length==0){
			cb&&cb();
			return;
		}
		if(typeof cls=='object'){
			$el.animate(cls,250,'linear',function(){ cb&&cb();});
			return;
		}
		cls = (cls||'empty')+' anim';
		$el.addClass(cls).one(_eName,function(){			
			$el.removeClass(cls);
			cb&&cb();
			$el.unbind(_eName);
		});
	};
	
	anim.css3 = function($el, cssObj, cb){
		$el = $($el);
		var _eName = 'webkitTransitionEnd otransitionend transitionend';
		$el.addClass('comm_anim').css(cssObj).one(_eName, function(){		
			$el.removeClass('comm_anim');
			cb&&cb();
			$el.unbind(_eName);
		});
	};
	
	/**
     * 添加控制器
     * @param {String} 要切换的当前对象selector
     * @param {String} 要切换的目标对象selector
     * @param {Boolean} 是否返回
     * @param {Function} 动画结束的回调函数
     */
	anim.change = function(current, target, isBack ,callback){
		
		var $current = $(current);
		var $target = $(target);
		
		if($current.length+$target.length==0){
			callback&&callback();
			return;
		}
		
		var targetTransition = $target.data('transition');
		if(!anim.classes[targetTransition]){
			callback&&callback();
			return;
		}
		var transitionType = anim.classes[targetTransition][isBack?1:0];
		anim.run($current, transitionType[0]);							
		anim.run($target, transitionType[1], function(){	
			callback&&callback();
		});
	};
	
	agilelite.register('anim', anim);

})(agilelite.$);

//ajax封装
(function($){
	
	//用法继承jQuery的$.ajax
	var ajax = function(opts){
		
		var success = opts.success;
		var error = opts.error;
		
		var random = '__ajax_random__='+Math.random();
    	opts.url += (opts.url.split(/\?|&/i).length==1?'?':'&')+random;

		var _success = function(data){
			success&&success(data);
		};
		var _error = function(data){
			error&&error();
		};

		opts.success = _success;
		opts.error = _error;
		
		$.ajax(opts);
	};
	
	agilelite.register('ajax', ajax);
})(agilelite.$);

(function($){
	var util = {};
	
	util.jquery = {
		scrollEnd : function($el){
			var _step;
			var _handler = function(){
				if(_step) clearTimeout(_step); _step = null;
				_step = setTimeout(function(){
					clearTimeout(_step); _step = null;
					$el.trigger('scrollEnd');
				}, 100);
			};
			$el.off('scroll', _handler).on('scroll', _handler);
		}
	};
	
	util.getCSSTransform = function($el){
		var transform = $el.css('transform');
		if(!transform||transform=='none') return [];
		var matrix = new Function('return ['+transform.replace(/matrix.([^\\)]+)./, '$1')+'];');
		return matrix();
	};
	
	util.script = function(str){	
		str = (str||'').trim();
		var tag = false;
		
    	str = str.replace(/\$\{([^\}]*)\}/g, function(s, s1){
    		try{
    			tag = true;
    			return eval(s1.trim());
    		}catch(e){
    			return '';
    		}

    	});

    	return tag?util.script(str):str;
    };	
	
	util.provider = function(str, data){
    	str = (str||'').trim();
    	return str.replace(/\$\{([^\}]*)\}/g, function(s, s1){   		
    		return data[s1.trim()]||'';
    	});
    };
    
	var URLParser = function(url) {

	    this._fields = {
	        'Username' : 4,   //用户
	        'Password' : 5,   //密码
	        'Port' : 7,       //端口
	        'Protocol' : 2,   //协议
	        'Host' : 6,       //主机
	        'Pathname' : 8,   //路径
	        'URL' : 0,
	        'Querystring' : 9,//查询字符串
	        'Fragment' : 10,   //锚点
	        'Filename' : -1,//文件名,不含后缀
	        'Queryobject' : -1, //查询字串对象
	        'Hashobject' : -1 //hash对象
	    };

	    this._values = {};
	    this._regex = null;
	    this.version = 0.1;
	    this._regex = /^((\w+):\/\/)?((\w+):?(\w+)?@)?([^\/\?:]+):?(\d+)?(\/?[^\?#]+)?\??([^#]+)?#?(\w*)/;
	    for(var f in this._fields)
	    {
	        this['get' + f] = this._makeGetter(f);
	    }

	    if (typeof url != 'undefined')
	    {
	        this._parse(url);
	    }
	};
	URLParser.prototype.setURL = function(url) {
	    this._parse(url);
	};

	URLParser.prototype._initValues = function() {
	    for(var f in this._fields)
	    {
	        this._values[f] = '';
	    }
	};
	
	var removeURLPre = function(url){
		if(url.indexOf('file:')==0){
			url = (url.replace(/file:\/*/,''));
		}
		return url;
	};
	
	URLParser.prototype._path = function(url) {
	    var up = new URLParser(removeURLPre(location.href.split(/\?\#/)[0]));
	    var baseUrl = up.getProtocol()+'://'+up.getHost()+(up.getPort()?(':'+up.getPort()):'');
		url = removeURLPre(url);
		if(url.indexOf('/')==0){
			return baseUrl + url;
		}else{
			var curPaths = up.getPathname().split('/');
			curPaths.length -= url.split('../').length;
			return baseUrl + curPaths.join('/') + '/' + url.replace(/[\.]{1,2}\//g,'');
		}
	};

	URLParser.prototype._parse = function(url) {
	    this._initValues();
	    url = removeURLPre(url);
	    var r = this._regex.exec(url);
	    if (!r){	    	
	    	url = this._path(url);
	    	r = this._regex.exec(url);
	    	if(!r){
	    		console.log('url invalid:'+url);
	    		return;
	    	}
	    }

	    for(var f in this._fields) if (typeof r[this._fields[f]] != 'undefined')
	    {
	        this._values[f] = r[this._fields[f]];
	    }
	};
	URLParser.prototype._makeGetter = function(field) {
	    return function() {
	    	if(field=='Filename'){
	    		var path = (this._values['URL'].split('?')[0]).split('/');
	    		return path[path.length-1].split('.')[0];
	    	}else if(field=='Queryobject'){
	    		var query = this._values['Querystring']||'';
	            var seg = query.split('&');
	            var obj = {};
	            for(var i=0;i<seg.length;i++){
	                var s = seg[i].split('=');
	                obj[s[0]] = s[1];
	            }
	            return obj;
	    	}else if(field=='Hashobject'){
	    		return {
		        	url : this.getURL(),
		            hash : '#'+(this.getFragment()||this.getFilename()),
		            tag : '#'+(this.getFragment()||this.getFilename()),
		            query : this.getQuerystring(),
		            param : this.getQueryobject()
		       };
	    	}
	    	
	        return this._values[field];
	    };
	};
	
	//url解析类
	util.parseURL = function(url){
		url = util.script(url||'');
		if(url.indexOf('#')==0){
			url = url.replace('#','')+(agilelite.options.viewSuffix?agilelite.options.viewSuffix:'');
		}	

		return new URLParser(url);
   	};
   	
   	util.isCrossDomain = function(url){
    	if(!url||url.indexOf(':')<0) return false;
    	var urlOpts = util.parseURL(url);
    	if(!urlOpts.getProtocol()) return false;
    	return !((location.protocol.replace(':','')+location.host)==(urlOpts.getProtocol()+urlOpts.getHost()+':'+urlOpts.getPort()));
    };

    util.checkBoolean = function(){
    	var result = false;
    	for(var i=0;i<arguments.length;i++){
    		if(typeof arguments[i]=='boolean'){
    			return arguments[i];
    		}
    	}
    	return result;
    };
	
	var _isCrossDomain = function(url){
    	if(!url||url.indexOf(':')<0) return false;

    	var urlOpts = agilelite.util.parseURL(url);

    	if(!urlOpts.getProtocol()) return false;

    	return !((location.protocol.replace(':','')+location.host)==(urlOpts.getProtocol()+urlOpts.getHost()+':'+urlOpts.getPort()));
    };
    
    /*
     * $.ajax函数封装，判断是否有跨域，并且设置跨域处理函数
     * @param ajax json对象，进行简单统一处理，如果需要完整功能请使用$.ajax
     * */
    util.ajax = function(opts){
    	if(!opts||!opts.url) return;
    	opts.url = util.script(opts.url);
    	//var _isBlock = util.checkBoolean(opts.isBlock,agilelite.options.showPageLoading);
		var _isBlock = opts.isBlock;
    	opts.dataType = (opts.dataType||'text').toLowerCase();
    	var ajaxData = {
        	url : opts.url,
            timeout : 20000,
            type : opts.type||'get',
            dataType : opts.dataType,
            success : function(html){
            	if(_isBlock) agilelite.hideMask();
                opts.success && opts.success(opts.dataType=='text'?util.script(html):html);
            },
            error : function(html){
            	if(_isBlock) agilelite.hideMask();
               	opts.error && opts.error(null);
            },
            reqCharset : opts.reqCharset||'utf-8'
        };
        if(opts.data) ajaxData.data = opts.data;
        if(opts.headers) ajaxData.headers = opts.headers;
    	var isCross = _isCrossDomain(opts.url);
    	var handler = agilelite.ajax;		
    	if(isCross){
    		ajaxData.dataType = agilelite.options.crossDomainHandler?ajaxData.dataType:'jsonp';
    		handler = agilelite.options.crossDomainHandler||handler;
    	}
    	if(ajaxData.dataType.toLowerCase()=='jsonp'){
    		ajaxData.jsonp = opts.jsonp||'agilecallback';
    	}
    	if(_isBlock) agilelite.showMask();
    	handler(ajaxData);
    };
    
    /*
     * $.ajax函数封装，不允许跨域
     * @param url地址，必选
     * @param callback回调函数，可选
     * 
     * */
    util.getHTML = function(url, callback){
    	agilelite.ajax({
    		url : url,
    		type : 'get',
    		dataType : 'text',
    		success : function(html){
    			callback&&callback(agilelite.util.script(html||''));
    		},
    		error : callback
    	});
    };
    
    util.readyAlarm = function($inner,targetName,eventName){
		var _return = {};		
		for(var k in $inner){
			if(typeof $inner[k]!='function'){
				_return[k] = $inner[k];
				continue;
			}
			_return[k] = (function(k){
				return function(){
							try{
								return $inner[k].apply(this, arguments);
							}catch(e){
								console.log('提示', '请在'+(eventName||agilelite.options.readyEvent)+'之后调用'+targetName+'.'+k+'桥接函数');
							}						
						};
			})(k);
		}
		return _return;
	};
	
	util.encodeHash = function(hash){
		return encodeURIComponent(agilelite.Base64.encode(hash));
	};
	
	util.decodeHash = function(enHash){
		return agilelite.Base64.decode(decodeURIComponent(enHash));
	};
	
    agilelite.register('util', util);

})(agilelite.$);

(function($){
	var _index_key_ = {};
	
	var scrollObj = {};
	
	//重写IScroll的click方法
	var iscrollRewrite = {};
	iscrollRewrite.click = function (e) {
		var target = e.target,
			ev;
		var isSpecialInputType = (/INPUT/i).test(target.tagName)&&(/DATE|COLOR/i).test(target.getAttribute('type'));
		if (isSpecialInputType || !(/(SELECT|INPUT|TEXTAREA)/i).test(target.tagName) ) {
			ev = document.createEvent('MouseEvents');
			ev.initMouseEvent('click', true, true, e.view, 1,
				target.screenX, target.screenY, target.clientX, target.clientY,
				e.ctrlKey, e.altKey, e.shiftKey, e.metaKey,
				0, null);

			ev._constructed = true;
			target.dispatchEvent(ev);
		}
		e.stopPropagation();//强制阻止事件冒泡
	};
	
	var scroll = function(selector, opts){	
		if(iscrollRewrite){
			for(var k in iscrollRewrite){
				IScroll.utils[k] = iscrollRewrite[k];
			}
			iscrollRewrite = null;
		}	
		var $el = $(selector);
		var eId = $el.attr('id');
		if($el.length==0||!eId){
			return null;
		}else if(_index_key_[eId]){
			return _index_key_[eId];
		}
		
		var $scroll;
		
		var costomOpts = {
			scrollTop : 'scrollTop',
			scrollBottom : 'scrollBottom'
		};
		
		var options = {
			mouseWheel: true,
			scrollbars : 'custom',
			fadeScrollbars : true,
			//click : true,
			preventDefaultException: { tagName: /^(INPUT|TEXTAREA|BUTTON|SELECT|IMG)$/ }
		};
		options = $.extend(options, agilelite.options.iScrollOptions||{});
		$.extend(options, opts||{});
		var _attr_options = $el.attr('data-scroll-options');
		_attr_options = typeof _attr_options!='object'?agilelite.JSON.parse(_attr_options):_attr_options;
		$.extend(options, _attr_options||{});
		IScroll.utils.isBadAndroid = false;//处理页面抖动
		$scroll = new IScroll(selector, options);
		$el.attr('__com_iscroll__', true);
		$el.off('touchmove.iscroll').on('touchmove.iscroll', 'textarea[data-scroll-diabled="true"]', function(){
			$scroll._execEvent('__setdiabled');
		}).off('touchend.iscroll blur.iscroll').on('touchend.iscroll blur.iscroll', 'textarea[data-scroll-diabled="true"]', function(){
			$scroll._execEvent('__setenabled');
		});

		$scroll.on('__setdiabled', function(){
			this.enabled = false;
		});
		$scroll.on('__setenabled', function(){
			this.enabled = true;
		});
		//解决iscroll在苹果设备超出边界不反弹的问题
		if(/AppleWebKit/i.test(window.navigator.appVersion)){
			var scrollTimer,$container = $el.closest('[data-role="section"], [data-role="modal"]');
			$el.off('touchmove.bounce').on('touchmove.bounce', function(e){
				if(scrollTimer){clearTimeout(scrollTimer);scrollTimer=null;}
				if($scroll.y>0||$scroll.y<$scroll.maxScrollY||$scroll.x>0||$scroll.x<$scroll.maxScrollX){
					scrollTimer = setTimeout(function(){
						var screenHeight = $container.height(),screenWidth = $container.width(),_touch = e.originalEvent.targetTouches[0];
						var point = {x : _touch.pageX, y : _touch.pageY};
						if(point.y<=0||point.y>=screenHeight||point.x<=0||point.x>=screenWidth){ $scroll._end(e);$scroll.refresh(); }
						clearTimeout(scrollTimer);scrollTimer = null;
					},200);
				}
			}).off('touchend.bounce').on('touchend.bounce', function(e){
				if(scrollTimer){ clearTimeout(scrollTimer);scrollTimer=null;}
			});
		}
		
		/*//处理软键盘被盖住
		var _focusProcess;
		$el.on('focus', 'input[type="url"], input[type="search"], input[type="email"], input[type="password"], input[type="number"], input[type="text"], textarea', function(){		
			var _this = this;
			if($el.height() - $(this).offset().top < 100){
				_focusProcess = setTimeout(function(){
					$scroll.scrollToElement(_this);
					clearTimeout(_focusProcess);
					_focusProcess = null;
				}, 500);
			}
			
		});
		$scroll.on('beforeScrollStart', function(){
			$el.find('input:focus, textarea:focus').blur();
			if(_focusProcess){
				clearTimeout(_focusProcess);
				_focusProcess = null;
			}			
		});*/
		//事件结束一律终止
		$el.off('touchend.reset').on('touchend.reset', function(e){ $scroll._execEvent('__setenabled');$scroll._end(e);});
		$scroll.on('scrollEnd' , function(){
			if(this.y==0){
            	this._execEvent(costomOpts.scrollTop);//自定义事件滑动到顶部
            }else if(this.y==this.maxScrollY){
            	this._execEvent(costomOpts.scrollBottom);//自定义事件滑动到底部
            }
			agilelite.Component.lazyload($el);//初始化懒人加载
		});		
		_index_key_[eId] = $scroll;
		$scroll.on('destroy', function(){ delete _index_key_[eId];});
		$scroll.refresh();
		$el.trigger('scrollInit');//自定义scroll初始化事件
		return _index_key_[eId];
	};
    agilelite.register('Scroll', scroll);
})(agilelite.$);

(function($){	
	var _index_key_ = {};	
	/*
     * refresh上拉下拉刷新
     * @param selector，scroll容器的id选择器
     * @param 选项，目前仅支持onPullDown事件和onPullUp事件，{onPullDown:function(){},onPullUp:function(){}}
     * 
     * */
    function refresh(selector, opts) {
    	var $el = $(selector);
    	var eId = $el.attr('id'); 	
    	if(_index_key_[eId]) return _index_key_[eId];
    	
		var $scroller = $el.children().first();
		var myScroll, opts=opts||{}, scrollType = $el.data('scroll'); 
    	var $pullDownEl, $pullDownL;  
    	var $pullUpEl, $pullUpL; 
    	var pullDownHeight, pullUpHeight; 
    	var refreshStep = -1;//加载状态-1默认，0显示提示下拉信息，1显示释放刷新信息，2执行加载数据，只有当为-1时才能再次加载
    	
    	var pullDownOpts = {
    		id : 'agile-pulldown',
    		iconStyle : 'agile-pulldown-icon',
    		labelStyle : 'agile-pulldown-label',
    		normalLabel : '下拉刷新',
    		releaseLabel : '释放加载',
    		refreshLabel : '加载中，请稍后',
    		normalStyle : '',
    		releaseStyle : 'release',
    		refreshStyle : 'refresh',
    		callback : null,
    		callbackEvent : 'pulldown'
    	};
    	var pullUpOpts = {
    		id : 'agile-pullup',
    		iconStyle : 'agile-pullup-icon',
    		labelStyle : 'agile-pullup-label',
    		normalLabel : '上拉刷新',
    		releaseLabel : '释放加载',
    		refreshLabel : '加载中，请稍后',
    		normalStyle : '',
    		releaseStyle : 'release',
    		refreshStyle : 'refresh',
    		callback : null,
    		callbackEvent : 'pullup'
    	};
    	
    	if(opts.onPullDown||scrollType=='pull'||scrollType=='pulldown'){
    		$pullDownEl = $('<div id="'+pullDownOpts.id+'"><div class="'+pullDownOpts.iconStyle+'"></div><div class="'+pullDownOpts.labelStyle+'">'+pullDownOpts.normalLabel+'</div></div>').prependTo($scroller);
	        $pullDownL = $pullDownEl.find('.'+pullDownOpts.labelStyle);   
	        pullDownHeight = $pullDownEl.height();
	        $pullDownEl.attr('class','').hide();
	        pullDownOpts.callback = opts.onPullDown||function(){};
    	}
    	
    	if(opts.onPullUp||scrollType=='pull'||scrollType=='pullup'){
    		$pullUpEl = $('<div id="'+pullUpOpts.id+'"><div class="'+pullUpOpts.iconStyle+'"></div><div class="'+pullUpOpts.labelStyle+'">'+pullUpOpts.normalLabel+'</div></div>').appendTo($scroller);
	        $pullUpL = $pullUpEl.find('.'+pullUpOpts.labelStyle);   
	        pullUpHeight = $pullUpEl.height();
	        $pullUpEl.attr('class','').hide();
	        pullUpOpts.callback = opts.onPullUp||function(){};
    	}

        myScroll = agilelite.Scroll('#'+eId, {  
			probeType: 2,//probeType：1对性能没有影响。在滚动事件被触发时，滚动轴是不是忙着做它的东西。probeType：2总执行滚动，除了势头，反弹过程中的事件。这类似于原生的onscroll事件。probeType：3发出的滚动事件与到的像素精度。注意，滚动被迫requestAnimationFrame（即：useTransition：假）。  
        });

        if(!myScroll) return null;
        //滚动时  
        myScroll.on('scroll', function(){
        	if(refreshStep==2){
        		//do nothing
        	}else if($pullDownEl && this.y > 0 && this.y < pullDownHeight){
           		$pullDownEl.show().removeClass(pullDownOpts.releaseStyle);
           		$pullDownL.html(pullDownOpts.normalLabel);  
           		refreshStep = 0;
           	}else if($pullDownEl && this.y >= pullDownHeight){
           		$pullDownEl.addClass(pullDownOpts.releaseStyle);
           		$pullDownL.html(pullDownOpts.releaseLabel);
           		refreshStep = 1;
           	}else if($pullUpEl && this.y < 0 && -this.y > -(this.maxScrollY - pullUpHeight)){
	            $pullUpEl.addClass(pullUpOpts.releaseStyle);
           		$pullUpL.html(pullUpOpts.releaseLabel); 
           		refreshStep = 1;
           	}else if($pullUpEl && this.y < 0 && -this.y > -this.maxScrollY && -this.y <= -(this.maxScrollY - pullUpHeight) ){
           		$pullUpEl.show().removeClass(pullUpOpts.releaseStyle);  
	            $pullUpL.html(pullUpOpts.normalLabel);  
	           	refreshStep = 0;
           	}
            
            if(this.y >0 && refreshStep > -1 && refreshStep < 2){
            	$pullDownEl.css('margin-top', - Math.max(pullDownHeight - this.y, 0));
            }
        }); 
        //滚动完毕  
        myScroll.on('scrollEnd', function(){
            if(refreshStep == 1){
	            if ($pullUpEl && $pullUpEl.hasClass(pullUpOpts.releaseStyle)) {  
	                $pullUpEl.removeClass(pullUpOpts.releaseStyle).addClass(pullUpOpts.refreshStyle);  
	                $pullUpL.html(pullUpOpts.refreshLabel);  
	                refreshStep = 2;  
	                pullUpOpts.callback.call(this);
	                this._execEvent(pullUpOpts.callbackEvent);
	            }else if($pullDownEl && $pullDownEl.hasClass(pullDownOpts.releaseStyle)){  	            	
	                $pullDownEl.removeClass(pullDownOpts.releaseStyle).addClass(pullDownOpts.refreshStyle);  
	                $pullDownL.html(pullDownOpts.refreshLabel);  
	                refreshStep = 2;  
	                pullDownOpts.callback.call(this);
	                this._execEvent(pullDownOpts.callbackEvent);
	            }
	            window.setTimeout(function(){
	            	if(refreshStep==2) myScroll.refresh();
	            },5000);
            }else{
            	myScroll.refresh();
            }
        });  
        
        myScroll.on('refresh', function(){
        	if(refreshStep!=-1){
	        	refreshStep = -1;	        	 
	            if($pullDownEl){
	            	$pullDownEl.attr('class','').hide();  
	            	$pullDownEl.css('margin-top', 0);
	            }	            
	            if($pullUpEl){
	            	$pullUpEl.attr('class','').hide();
	            }
        	}
        	
        });

        myScroll.setConfig = function(opts){
        	opts = opts||{};
        	if(opts.pullDownOpts){
        		pullDownOpts = $.extend(pullDownOpts, opts.pullDownOpts);
        	}
        	if(opts.pullUpOpts){
        		pullUpOpts = $.extend(pullUpOpts, opts.pullUpOpts);
        	}
        };

        _index_key_[eId] = myScroll;
        
        myScroll.on('destroy', function(){ delete _index_key_[eId]; });
        
        $el.trigger('refreshInit');
        
        return _index_key_[eId];
    }
    
    agilelite.register('Refresh', refresh);
    
})(agilelite.$);

(function($){
	
	var _index_key_ = {};
	/*
     * slider滑动容器
     * @param selector，slider容器的id选择器
     * @param 选项，auto：是否自动播放|自动播放的毫秒数；dots：指示点显示的位置，center|right|left|hide；change：每次切换slider后的事件，{auto:false,dots:'right',change:function(){}}
     * */
	var slider = function(selector, opts){
		var $el = $(selector);
		var eId = $el.attr('id');
		
		if(_index_key_[eId]) return _index_key_[eId];
		
		var sliderOpts = {
			auto : false,
			change : function(){
				
			},
			dots : '',
			loop : false
		};
		$.extend(sliderOpts, opts);
		if($el.children('scroller').children('[data-role="page"]').length>0){
			sliderOpts.loop = false;
			sliderOpts.dots = 'hide';
		}
		
		var $scroller,$slide,outWidth,slideNum,$dots,index = 0,vIndex = 1;
		var _initHeight = $el.height();
		var init = function(){
			if(!$el.hasClass('active')) $el.addClass('active');
			$scroller = $el.children('.scroller');
			$slide = $scroller.children('.slide').remove('.clone');
			outWidth = $el.parent().width();
			slideNum = $slide.length;
			$scroller.css({height:'auto'});
			$slide.width(outWidth);
			if(!_initHeight){
				if($slide.height()){
					$el.height($slide.height());
				}else{
					$el.height('1px');
				}	
			}
			var trueNum = slideNum;
			if(sliderOpts.loop){
				$scroller.prepend($slide.last().clone().addClass('clone').removeClass('active')).append($slide.first().clone().addClass('clone').removeClass('active'));
				$slide = $scroller.children('.slide');
				trueNum += 2;
			}
			$scroller.css({height:$el.height(), width:outWidth*trueNum});
		};
		
		var createDots = function(){
			index = $scroller.index($slide.filter('.active'));index = index<0?0:index;
			vIndex = _loopUtil.getVIndex(index);
			$el.children('.dots').remove();			
			var arr = [];
			arr.push('<div class="dots">');
			if(sliderOpts.loop) arr.push('<div class="dotty hidden"></div>');
			for(var i=0;i<slideNum;i++){
				arr.push('<div class="dotty"></div>');
			}
			if(sliderOpts.loop) arr.push('<div class="dotty hidden"></div>');
			arr.push('</div>');
			$dots = $(arr.join('')).appendTo($el).addClass(sliderOpts.dots).find('.dotty');
			$($dots.get(vIndex)).addClass('active').siblings('.active').removeClass('active');
			$($slide.get(vIndex)).addClass('active').siblings('.active').removeClass('active');
		};
		var _loopUtil = {
			getVIndex : function(ti){
				if(sliderOpts.loop) return ti<1?1:(ti>slideNum?slideNum:ti);
				return ti;
			},
			getTIndex : function(vi){
				if(sliderOpts.loop) return vi+1;
				return vi;
			},
			setTIndex : function(vi){
				if(sliderOpts.loop) return vi-1;
				return vi;
			}
		};
		
		init();//第一步初始化
		createDots();//第二步创建dots

		var options = {
			scrollbars : false,
			scrollX: true,
			scrollY: false,
			momentum: false,
			snap: true,
			keyBindings: true,
			bounceEasing : 'circular',
			probeType: 2
		};		
		var myScroll = agilelite.Scroll('#'+eId, options);		
		var outerSlider;
		$el.off('touchstart.slider').on('touchstart.slider', function(e){
			var $outerSlider = $(this).parent().closest('[data-role="slider"]');
			if($outerSlider.length==0) return;
			outerSlider = agilelite.Slider($outerSlider[0]);
			outerSlider._execEvent('__setdiabled');
		});
		$el.off('touchmove.slider').on('touchmove.slider', function(e){
			var $outerSlider = $(this).parent().closest('[data-role="slider"]');
			if($outerSlider.length==0) return;
			outerSlider = agilelite.Slider($outerSlider[0]);
			outerSlider._execEvent('__setdiabled');
		});
		$el.off('touchend.slider').on('touchend.slider', function(e){
			outerSlider&&outerSlider._execEvent('__setenabled');
		});
		var transformX = 'ready';
		myScroll.on('scrollStart', function(){
			transformX = agilelite.util.getCSSTransform($scroller)[4];
		});
		myScroll.on('scroll', function(){
			if(transformX=='ready') return;
			if(Math.abs(agilelite.util.getCSSTransform($scroller)[4]-transformX)>=40){
				$el.closest('[data-scroll]').add($el.find('[data-scroll]:first')).each(function(){
					A.Scroll('#'+this.id)._execEvent('__setdiabled');
				});
				transformX = 'ready';
			}
		});

		myScroll.on('scrollEnd', function(){
			index = this.currentPage.pageX;
			var curSlide = $($slide.get(index));
			var curDots = $($dots.get(index));
			if(!curSlide.hasClass('active')){
				if(sliderOpts.loop){
					if(index<1){
						myScroll.goToPage(slideNum, 0, 0);myScroll._execEvent('scrollEnd');
						return;
					}else if(index>slideNum){
						myScroll.goToPage(1, 0, 0);myScroll._execEvent('scrollEnd');
						return;
					}
				}
				sliderOpts.change&&sliderOpts.change.call(this, index-(sliderOpts.loop?1:0));
				curSlide.addClass('active').trigger('slidershow').siblings('.active').removeClass('active').trigger('sliderhide');
				curDots.addClass('active').siblings('.active').removeClass('active');
			}
		});		
		$el.off(agilelite.options.clickEvent+'.slider').on(agilelite.options.clickEvent+'.slider', function(){
			myScroll.goToPage(index, 0, options.snapSpeed);
		});
		myScroll.on('init', function(type){
			init();
			createDots();
			myScroll._execEvent('scrollEnd');
		});
		myScroll.trigger = function(eName, params){
			myScroll._execEvent(eName, params);
		};
		myScroll.index = function(){
			if(arguments.length==0) return _loopUtil.setTIndex(index);
			myScroll.goToPage(_loopUtil.getTIndex(arguments[0]), 0, arguments[1]||options.snapSpeed);
			myScroll._execEvent('scrollEnd');
		};

		myScroll.goToPage(vIndex, 0, 0);
		myScroll.refresh();
		var _auto;
		if(sliderOpts.auto){
			_auto = setInterval(function(){
				index = index==slideNum?0:index+1;
				myScroll.goToPage(_loopUtil.getVIndex(index), 0, options.snapSpeed);
			}, typeof sliderOpts.auto=='boolean'?5000:sliderOpts.auto);
		}
		
		_index_key_[eId] = myScroll;
		
		myScroll.on('destroy', function(){ $dots.remove(); delete _index_key_[eId]; clearInterval(_auto); });

		return _index_key_[eId];
	};
	
	agilelite.register('Slider', slider);
})(agilelite.$);

//侧边栏
(function($){
	var $asideContainer=$('#aside_container'), $sectionContainer=$('#section_container'), $sectionMask=$('#section_container_mask');
	var _this;
	var Aside = function(){	
		_this = this;
	};
	
	Aside.prototype.launch = function(){
		if($asideContainer.length==0) $asideContainer = $('<div id="aside_container"></div>').appendTo('body');
        if($sectionContainer.length==0) $sectionContainer = $('<div id="section_container"></div>').appendTo('body');
        if($sectionMask.length==0) $sectionMask = $('<div id="section_container_mask"></div>').appendTo('#section_container');
  		$sectionMask.on(agilelite.options.clickEvent, function(){
        	_this.hide();
        	return false;
        });
		$sectionMask.on('swipeleft', function(e){
			var $activeAside = $('#aside_container aside.active');
			if($activeAside.data('position') == 'left'){
				_this.hide();
			}
		});
		$sectionMask.on('swiperight', function(e){
			var $activeAside = $('#aside_container aside.active');
			if($activeAside.data('position') == 'right'){
				_this.hide();
			}
		});
	};
	
	 /**
     * 打开侧边菜单
     * @param selector css选择器或element实例
     * @param callback 动画完毕回调函数
     */
	Aside.prototype.show = function(selector, callback){
		var $aside = $(selector).addClass('active'),
            transition = $aside.data('transition'),// push overlay  reveal
            position = $aside.data('position') || 'left',
            showClose = $aside.data('show-close'),
            width = $aside.width(),
            translateX = (position=='left'?'':'-')+width+'px';
        var cssName = {
        	aside : {},
        	section : {}
        };
        var funcName = 'run';
        
        var _finish = function(){
        	callback&&callback();$aside.trigger('asideshow');
        };

        if($aside.data('role')=='aside'){
        	funcName = 'css3';
        	cssName.aside['-webkit-transform'] = 'translateX(0%)';
        	cssName.section['-webkit-transform'] = 'translateX('+translateX+')';
        	$aside.css({'transition':'all 250ms'});$sectionContainer.css({'transition':'all 250ms'});
        }else{
        	cssName.aside[position] = width+'px';
        	cssName.section['left'] = translateX;
        	$aside.css({'transition':'none'});$sectionContainer.css({'transition':'none'});
        }
        
        if(transition == 'overlay'){
            agilelite.anim[funcName]($aside, cssName.aside, _finish);
        }else if(transition == 'reveal'){
        	agilelite.anim[funcName]($sectionContainer, cssName.section, _finish);
        }else{//默认为push
        	agilelite.anim[funcName]($aside, cssName.aside);
        	agilelite.anim[funcName]($sectionContainer, cssName.section, _finish);
        }
        $sectionMask.show();
        agilelite.pop.hasAside = true;
	};
	/**
     * 关闭侧边菜单
     * @param callback 动画完毕回调函数
     */
   Aside.prototype.hide = function(callback){
        var $aside = $('#aside_container aside.active'),
            transition = $aside.data('transition'),// push overlay  reveal
            position = $aside.data('position') || 'left',
            translateX = position == 'left'?'-100%':'100%';

        var _finish = function(){
            $aside.removeClass('active');
            agilelite.pop.hasAside = false;
            callback&&callback.call(this);$aside.trigger('asidehide');
        };
        
        var cssName = {
        	aside : {},
        	section : {}
        };
        var funcName = 'run';
        if($aside.data('role')=='aside'){
        	funcName = 'css3';
        	cssName.aside['-webkit-transform'] = 'translateX('+translateX+')';
        	cssName.section['-webkit-transform'] = 'translateX(0%)';
        }else{
        	cssName.aside[position] = 0;
        	cssName.section['left'] = 0;
        }
        
        if(transition == 'overlay'){
            agilelite.anim[funcName]($aside, cssName.aside, _finish);
        }else if(transition == 'reveal'){
            agilelite.anim[funcName]($sectionContainer, cssName.section, _finish);
        }else{//默认为push
            agilelite.anim[funcName]($aside, cssName.aside);
            agilelite.anim[funcName]($sectionContainer, cssName.section, _finish);
        }
        $sectionMask.hide();
    };
    agilelite.register('Aside', new Aside());
})(agilelite.$);

/**
 * 弹出框组件
 */
(function($){
	var _popMap= {index:0};
    var transitionMap = {
    	default : ['fadeIn','fadeOut'],
        top : ['slideDownIn','slideUpOut'],
        bottom : ['slideUpIn','slideDownOut'],
        center : ['bounceIn','bounceOut'],
        left : ['slideRightIn','slideLeftOut'],
        right : ['slideLeftIn','slideRightOut']           
	};
 
    var Popup = function(opts){
    	var _index = _popMap.index++;  
    	var options = {
    		id : '__popup__'+_index,//popup对象的id
    		html : '',//位于pop中的内容
    		pos : 'default',//pop显示的位置和样式,default|top|bottom|center|left|right|custom
    		css : {},//自定义的样式
    		isBlock : false//是否禁止关闭，false为不禁止，true为禁止
    	};
    	$.extend(options,opts);
    	if(_popMap[options.id]) return _popMap[options.id];
    	var _this = _popMap[options.id] = this;
    	$('<div data-refer="'+options.id+'" class="'+agilelite.options.classPre+'popup-mask"></div><div id="'+options.id+'" class="'+agilelite.options.classPre+'popup"></div>').appendTo('body');
    	var $popup = $('#'+options.id), $mask = $('[data-refer="'+options.id+'"]');
    	$popup.data('block', options.isBlock);
    	$popup.css(options.css);
    	$popup.addClass(agilelite.options.classPre+options.pos);
		$popup.html(options.html);
		$mask.on(agilelite.options.clickEvent, function(){
        	if(options.isBlock) return false;
        	_this.close();
        	return false;
        });
        $popup.on(agilelite.options.clickEvent,'[data-toggle="popup"]',function(){
        	_this.close();
			return false;
		});
		this.options = options;
		this.popup = $popup;
		this.mask = $mask;
		this._event = {};
    };    
	Popup.prototype.on = function(eType, callback){
    	if(this._event[eType]){
    		this._event[eType].push(callback);
    	}else{
    		this._event[eType] = [callback];
    	}
    	return this;
    };
    Popup.prototype.trigger = function(eType, callback){
    	var arr = this._event[eType];
    	for(var i=0;i<(arr||[]).length;i++){
    		arr[i].apply(this);
    	};
    	callback&&callback.apply(this);
    	return this;
    };
    
    Popup.prototype.open = function(callback){
    	var _this=this, $popup=this.popup, $mask=this.mask, options=this.options;
    	if($mask.hasClass('active')) return _popMap[options.id];
    	$('body').children('.'+agilelite.options.classPre+'popup-mask.active').removeClass('active');
    	$mask.addClass('active').show();
        $popup.show();
        var popHeight = $popup.height();
        if(options.pos == 'center') $popup.css('margin-top','-'+popHeight/2+'px');
        var transition = transitionMap[options.pos];
        if(transition) {
        	_this._status = 'opening';
        	agilelite.anim.run($popup,transition[0], function(){
        		_this._status = 'opened';
        	});
        }else{
        	_this._status = 'opened';
        }
		this.trigger('popupopen');callback&&callback.call(this);
		agilelite.pop.hasPop = $popup;
		if(agilelite.options.autoInitCom) agilelite.Component.initComponents($popup);
		return this;
    };
    
    var _finish = function(callback){   	
    	var $popup=this.popup, $mask=this.mask;
    	$popup.remove();$mask.remove();this.trigger('popupclose');setTimeout(function(){ callback&&callback(); }, 200);
    	$last = $('body').children('.'+agilelite.options.classPre+'popup-mask').last().addClass('active');
    	agilelite.pop.hasPop = $last.length==0?false:$('body').children('.agile-popup').last();
    };
    
    Popup.prototype.close = function(callback){
    	var _this = this; setTimeout(function(){ _this._close(callback); }, 100);
    };
    
    Popup.prototype._close = function(callback){
    	if(this._status=='opening') return;
    	this.trigger('popupbeforeclose');
    	var _this=this,$popup=this.popup, $mask=this.mask, options=this.options;
        var transition = transitionMap[options.pos];
        if(transition){
            agilelite.anim.run($popup,transition[1],function(){ _finish.call(_this, callback); });
        }else{
            _finish.call(_this, callback);
        }
        delete _popMap[options.id];
        //return this;
    };
    
    agilelite.register('Popup' , Popup);
    
    var _ext = {};
    
    _ext.popup = function(opts){
    	if(_popMap[opts.id]) return _popMap[opts.id];
    	return new Popup(opts).open();
    };
    
    _ext.closePopup = function(callback){
    	var _id = agilelite.pop.hasPop.attr('id');
    	if(_popMap[_id]) _popMap[_id].close(callback);
    };
    
    /**
     * alert组件
     * @param title 标题
     * @param content 内容
     * @param callback 点击确定按钮后的回调
     */
    _ext.alert = function(title,content,callback){
    	if(!content){
    		content = title;
    		title = '提示';
    	}else if(typeof content=='function'){
    		callback = content;
    		content = title;
    		title = '提示';
    	}
        return new Popup({
            html : agilelite.util.provider('<div class="'+agilelite.options.classPre+'popup-title">${title}</div><div class="'+agilelite.options.classPre+'popup-content">${content}</div><div class="'+agilelite.options.classPre+'popup-handler"><a data-toggle="popup" class="'+agilelite.options.classPre+'popup-handler-ok">${ok}</a></div>', {title : title, content:content, ok:'确定'}),
            pos : 'center',
            isBlock : true
        }).on('popupclose', function(){
        	callback&&callback();
        }).open();
    };
    
    /**
     * confirm 组件
     * @param title 标题
     * @param content 内容
     * @param okCallback 确定按钮handler
     * @param cancelCallback 取消按钮handler
     */
    _ext.confirm = function(title,content,okCallback,cancelCallback){
    	if(typeof content=='function'){
    		cancelCallback = okCallback;
    		okCallback = content;
    		content = title;
    		title = '提示';
    	}
    	var clickBtn = okCallback;
        return new Popup({
            html : agilelite.util.provider('<div class="'+agilelite.options.classPre+'popup-title">${title}</div><div class="'+agilelite.options.classPre+'popup-content">${content}</div><div class="'+agilelite.options.classPre+'popup-handler"><a data-toggle="popup" class="'+agilelite.options.classPre+'popup-handler-cancel">${cancel}</a><a data-toggle="popup" class="'+agilelite.options.classPre+'popup-handler-ok">${ok}</a></div>', {title : title, content:content, cancel:'取消', ok:'确定'}),
            pos : 'center',
            isBlock : true
        }).on('popupclose', function(){       	
        	clickBtn&&clickBtn.call(this);
        }).open(function(){    	
        	var $popup = $(this.popup), _this=this;
        	$popup.find('.'+agilelite.options.classPre+'popup-handler-ok').on(agilelite.options.clickEvent, function(){
	            clickBtn = okCallback;
	        });
	        $popup.find('.'+agilelite.options.classPre+'popup-handler-cancel').on(agilelite.options.clickEvent, function(){
	            clickBtn = cancelCallback;
	        });
        });        
    };

    /**
     * loading组件
     * @param text 文本，默认为“加载中...”
     * @param callback 函数，当loading被人为关闭的时候的触发事件
     */
    _ext.showMask = function(text,callback){
    	if(typeof text=='function'){
    		callback = text;
    		text = '加载中';
    	}
        return new Popup({
        	id : 'popup_loading',
            html : agilelite.util.provider('<div><i class="'+agilelite.options.classPre+'popup-spinner"></i><p>${title}</p></div>', {title : text||'加载中'}),
            pos : 'loading',
            isBlock : true
        }).on('popupclose', function(){
        	callback&&callback();
        }).open();
    };
    
    _ext.hideMask = function(callback){
    	if(_popMap['popup_loading']) _popMap['popup_loading'].close();
    };

    /**
     * actionsheet组件
     * @param buttons 按钮集合
     * [{css:'red',text:'btn',handler:function(){}},{css:'red',text:'btn',handler:function(){}}]
     */
    _ext.actionsheet = function(buttons,showCancel){
        var markMap = ['<div class="'+agilelite.options.classPre+'actionsheet"><div class="'+agilelite.options.classPre+'actionsheet-group">'];
        var defaultCalssName = agilelite.options.classPre+"popup-actionsheet-normal";
        var defaultCancelCalssName = agilelite.options.classPre+"popup-actionsheet-cancel";
        var showCancel = showCancel==false?false:(showCancel||true);
        $.each(buttons,function(i,n){
            markMap.push('<button data-toggle="popup" class="'+(n.css||defaultCalssName)+'">'+ n.text +'</button>');
        });
        markMap.push('</div>');
        if(showCancel) markMap.push('<button data-toggle="popup" class="'+(typeof showCancel=='string'?showCancel:defaultCancelCalssName)+'">取消</button>');
        markMap.push('</div>');
        return new Popup({
            html : markMap.join(''),
            pos : 'bottom',
            css : {'background':'transparent'},
       	}).open(function(){           	
 			$(this.popup).find('button').each(function(i,button){              	
            	$(button).on(agilelite.options.clickEvent,function(){
                	if(buttons[i] && buttons[i].handler){
                    	buttons[i].handler.call(button);
                    }
                });
            });
    	});
    };
    
    /**
     * 带箭头的弹出框
     * @param html 弹出的内容可以是html文本也可以输button数组
     * @param el 弹出位置相对的元素对象
     */
    _ext.popover = function(html,el){
    	var markMap = [];
    	markMap.push('<div class="'+agilelite.options.classPre+'popover-angle"></div>');
    	if(typeof html=='object'){
    		markMap.push('<ul class="'+agilelite.options.classPre+'popover-items">');
    		for(var i=0;i<html.length;i++){
    			markMap.push('<li data-toggle="popup" class="'+(html[i].css||'')+'">'+html[i].text+'</li>');
    		}
    		markMap.push('</ul>');
    	}else{
    		markMap.push(html);
    	}

        return new Popup({
            html : markMap.join('')
        }).on('popupopen', function(){
        	var $del = $(document);
			var dHeight = $del.height();
			var dWidth = $del.width();
			
			var $rel = $(el);
			var rHeight = $rel.height();
			var rWidth = $rel.width();
			var rTop = $rel.offset().top;
			var rLeft = $rel.offset().left;
			var rCenter = rLeft+(rWidth/2);
			
			var $el = $(this.popup).addClass(agilelite.options.classPre+'popover');
			var $angle = $($el.find('.'+agilelite.options.classPre+'popover-angle').get(0));
			var gapH = $angle.height()/2;
			var gapW = Math.ceil(($angle.width()-2)*Math.sqrt(2));			
			var height = $el.height();
			var width = $el.width();
			
			var posY = dHeight-height-rHeight<0&&rTop>height?'up':'down';			
			var elCss = {}, anCss = {};   		
    		if(posY=='up'){
    			elCss.top = rTop - (height + gapH);
    			anCss.bottom = -gapH+4;
    		}else{
    			elCss.top = rTop + (rHeight + gapH);
    			anCss.top = -gapH+4;
    		}
    		elCss.left = rCenter-width/2;
			if(elCss.left+width>=dWidth-4){
				elCss.left = dWidth - width - 4;
			}else if(elCss.left<4){
				elCss.left = 4;
			}			
			anCss.left = rCenter - elCss.left - gapW/2;
    		$el.css(elCss);
    		$angle.css(anCss);
            $el.find('.'+agilelite.options.classPre+'popover-items li').each(function(i,button){             	
                $(button).on(agilelite.options.clickEvent,function(){
                    if(html[i] && html[i].handler){
                        html[i].handler.call(button);
                    }
                });
            });
        }).open();
    };
	
	for(var k in _ext){
		agilelite.register(k, _ext[k]);
	}
    
})(agilelite.$);

/*
 * Toast提示框
 * */
(function($){
    var _toastMap = {index:0, process: [], map: {}};
    var Toast = function(opts){
    	var _index = _toastMap.index++;
    	var options = {
    		id : '__toast__'+_index,//popup对象的id
    		isBlock : false,
    		duration : 2000,
    		css : '',
    		text : ''
    	};
    	$.extend(options, opts);
    	if(_toastMap[options.id]) return this;
    	var _this = _toastMap.map[options.id] = this;
    	var $toast = this.toast = $('<div id="'+options.id+'" class="agile-toast"><a>'+options.text+'</a></div>').appendTo('body');
    	$toast.addClass(options.css);
    	this.options = options;
    	_toastMap.process.push(this);
    };
    Toast.prototype.show = function(){
        if(_toastMap.process.length>0&&_toastMap.process[0]!=this){
        	return this;
        }
        var $toast = this.toast, options = this.options;        
        var _this = this;
        $toast.css({'display':'block'});
        agilelite.anim.run($toast,'scaleIn', function(){
        	if(options.isBlock==false) setTimeout(function(){ _this.hide();}, options.duration);
        });
        return this;
    };
    Toast.prototype.hide = function(){
    	var $toast = this.toast, options = this.options;
    	agilelite.anim.run($toast,'scaleOut',function(){
        	$toast.remove();
        	delete _toastMap.map[options.id];
        	_toastMap.process.shift();
        	if(_toastMap.process.length>0) _toastMap.process[0].show();
        });
        return this;
    };
    
    agilelite.register('Toast', Toast);

    var _ext = {};
    /**
     *  显示消息
     * @param text
     * @param duration 持续时间
     */
    _ext.showToast = function(text,duration){
        new Toast({
        	text : text,
        	duration : duration
        }).show();
    };   
    _ext.alarmToast = function(text,duration){
        new Toast({
        	text : text,
        	duration : duration,
        	css : 'alarm'
        }).show();
    };
    _ext.clearToast = function(){
        $('.agile-toast').remove();
        _toastMap.map = {};
        _toastMap.process = [];
    };
    _ext.getAllToasts = function(){
        return _toastMap.process;
    };
    
    for(var k in _ext){
    	agilelite.register(k, _ext[k]);
    }
    
})(agilelite.$);

/*
 * 事件处理
 * */
(function($){
	var event = {},_events = {};
	
	/**
	 * 处理modal打开顺序
	 * @private
	 */
	_events.modalCache = function(){
		var _modalCollection = {};
		$(document).on('modalshow', '[data-role="modal"]', function(){
			_modalCollection[this.id] = this.id;
		});
		$(document).on('modalhide', '[data-role="modal"]', function(){
			delete _modalCollection[this.id];
			var k;
			for(k in _modalCollection){}
			if(k){
				$('#'+k).trigger('modalshow');
			}
		});
		$(document).on('modalback', function(){
			var k;
			for(k in _modalCollection){}
			agilelite.Controller.modal('#'+k);
		});
	};
	
	/**
	 * 处理点击事件，默认阻止
	 * @private
	 */
	_events.clickHandler = function(){
		$(document).on(agilelite.options.clickEvent, '[data-click]', function(){
			var clickFunc = $(this).data('click');
			var _this = this;
			return clickFunc?eval(clickFunc.replace('.(', '.call(_this')):true;
		});
	};
	
	 /**
     * 处理浏览器的后退事件
     * @private
     */
	_events.back = function(){
		$(window).on('popstate', function(e){  
			agilelite.Controller.back(false);
			return;
    	});
	};
	
	/**
     * 处理data-cache="false"的容器缓存
     * @private
     */
	_events.removeCache = function(){
		var handler = function(){
			var $el = $(this);
			if($el.data('scroll')) agilelite.Scroll($el).destroy();
			$el.find('[data-scroll],[data-role="slider"]').each(function(){ agilelite.Scroll(this).destroy(); });
			$el.remove();
		};
		var controller = agilelite.Controller.get();
		var eName = [];
		var $doc = $(document);
		for(var k in controller){
			k = k=='default'?'view':k;
			$doc.on(k+'hide', '[data-cache="false"][data-role="'+k+'"]', handler);
		}
	};
	
	/**
     * 处理返回事件
     * @private
     */
	_events.backEvent = function(){
		var setParams = function($el, obj){
			if(obj){
				$($el).data('params', agilelite.JSON.stringify(obj));
			}
		};
		$(document).on(agilelite.options.backEvent, function(event, func){
			if(agilelite.pop.hasPop){
				if(Boolean(agilelite.pop.hasPop.data('block'))==false){
					agilelite.closePopup();
				}
			}else if($('[data-role="modal"].active').length>0){
				setParams('[data-role="section"].active:first', func);
				$(document).trigger('modalback');
			}else if(agilelite.pop.hasAside){
				setParams('[data-role="section"].active:first', func);
				agilelite.Controller.aside();
			}else if(agilelite.Controller.get('section').history.length<2){	
	    		$(document).trigger('beforeunload');//触发关闭页面事件
	    	}else{
	    		if(func) $('[data-role="section"].active:first').data('back-params', agilelite.JSON.stringify(func));
	    		window.history.go(-1);
	    	}
		});
		var _sectionMap = [],_doTransition = function(){
			if(_sectionMap.length!=1) return;
			var $current = _sectionMap[0][0];
			var $target = _sectionMap[0][1];
			agilelite.anim.change($current, $target, true, function(){	  
				_sectionMap.shift();     			       	       	
		       	var targetRole = $target.data('role');
		        $target.addClass('active').trigger(targetRole+'show');
				$current.removeClass('active').trigger(targetRole+'hide');
				$current.find(':focus').blur();//解决编辑状态切换页面的问题
				setTimeout(_doTransition, 100);
		    });
		};
		$(document).on('sectiontransition', function(e, $current, $target){
			$('.agile-popup').each(function(){ (new agilelite.Popup({ id : this.id})).close(); });
			if(agilelite.pop.hasAside){ agilelite.Controller.aside(); }
			$('[data-role="modal"].active').each(function(){ agilelite.Controller.modal('#'+this.id); });
			_sectionMap.push([$current, $target]);
			_doTransition();
		});
	};
	
	_events.zepto = function(){
		if(agilelite.options.clickEvent!='click') $(document).on('click', 'a[data-toggle]', function(){return false; });
		if($==window.Zepto){
			$(document).on('swipeLeft','[data-aside-right],[data-role="calendar"],.swipe_block', function(){$(this).trigger('swipeleft');});
			$(document).on('swipeRight','[data-aside-left],[data-role="calendar"],.swipe_block',function(){$(this).trigger('swiperight');});
		}	
	};

	/*
	 * 初始化agilestart事件
	 * */
	var _initSection = function(){
		$(document).off('sectionload.agilestart').on('sectionload.agilestart', 'section', function(){
			var $el = $(this);
			var $childred = $el.children(':first-child, :last-child').not('[data-role="article"], [data-boundary="false"]');
			$childred.off('touchmove.agilestart').on('touchmove.agilestart', function(e){
				e.preventDefault();
			});
			//处理aside事件
			if($el.data('aside-left')){
				$el.off('swiperight.agilestart').on('swiperight.agilestart', function(){
					agilelite.Controller.aside($el.data('aside-left'));
				});
			}
			if($el.data('aside-right')){
				$el.off('swipeleft.agilestart').on('swipeleft.agilestart', function(){
					agilelite.Controller.aside($el.data('aside-right'));
				});
			}
		});
		//初始化section
		var sectionSelecor = agilelite.Controller.get()['section']['container']+' [data-role="section"]';
		var $section = $(sectionSelecor+'.active:first');
		if($section.length==0) $section = $(sectionSelecor).first();
		if(agilelite.options.usePageParam){
			$section.children('article.active').off('articleload.agilestart').on('articleload.agilestart', function(){
				var pageInfo = agilelite.JSON.parse(agilelite.pageInfo);
				if(!pageInfo||$section.attr('id')==pageInfo.hash) return;
				var url = pageInfo.url||(pageInfo.hash+'.html?'+agilelite.JSON.toParams(pageInfo.params));
				setTimeout(function(){
					agilelite.Controller.section(url);
				},1000);
				
			});
		}

		agilelite.Controller.section($section.attr('id')+'.html?'+agilelite.JSON.toParams(agilelite.JSON.parse($section.data('params'))));
	};
	_events.agileStart = function(){
		var flag = true;
		$(document).on(agilelite.options.agileReadyEvent, function(){
			_initSection();
		});
		$(document).on(agilelite.options.agileStartEvent, function(){
			if(flag) {
				flag = false;
				return;	
			}
			_initSection();
		});
	};
	
	/**
     * 初始化相关组件
     */
	_events.initComponents = function(){
		$(document).on('sectionshow', 'section', function(){
			//初始化article
			agilelite.Controller.article('#'+$(this).children('[data-role="article"].active').first().attr('id'));		
		});		
		$(document).on('modalshow', '.modal', function(){
			//初始化article
			agilelite.Controller.article('#'+$(this).children('[data-role="article"].active').first().attr('id'));		
		});
		$(document).on('slidershow', '[data-role="page"].active', function(e, el){		
			var id = $(this).attr('id');			
			agilelite.Component.default($('[href="#'+id+'"]'));//初始化slider page
		});
		$(document).on('renderEnd', 'script', function(e,h){
			if(agilelite.options.autoInitCom) agilelite.Component.initComponents(h);
			var $scroller = $(h).closest('[data-scroll]');
			if($scroller.length==1) agilelite.Scroll('#'+$scroller.attr('id')).refresh();
		});
		$(document).on('modalhide', 'div.modal', function(e,h){
			_initSection();
		});
		$(document).on('scrollInit', '[data-scroll], [data-role="slider"]', function(){
			agilelite.Component.lazyload(this);
		});
		$(document).on('touchstart', '[data-readonly="true"]', function(){
			return false;
		});
		
		$(document).on('articleload', 'article[data-role="article"]', function(){
			var $el = $(this);
			if($el.data('scroll')){
				$el.on('click', 'summary', function(){			
					setTimeout(function(){
						agilelite.Scroll($el).refresh();
					}, 200);
				});
			}else{
				$el.on('scrollEnd', function(){
					agilelite.Component.lazyload($el);
				});
				agilelite.util.jquery.scrollEnd($el);
			}
		});
	};	
	
	/**
     * @param {Object} 添加的事件对象 {key:function(){}}
     */
	event.add = function(obj){
		$.extend(_events, obj||{});
	};
	/**
     * 启动event
     */
	event.launch = function(){
		for(var k in _events){
			_events[k]();
		}
	};
	agilelite.register('event', event);
})(agilelite.$);

(function($){
	var Template = function(selecotr){
		this.$el = $(selecotr);
	};
	
	Template.prototype.getTemplate = function(cb){
		var $el = this.$el;
		var source = $el.attr('src');
		var tmpl = $el.html().trim();
		if(tmpl){
			cb&&cb(tmpl);
		}else if(source){
			agilelite.util.getHTML(agilelite.util.script(source), function(html){
				$el.text(html);
				cb&&cb(html);
			});
		}else{
			cb&&cb('');
		}
	};
	
	Template.prototype.render = function(url, cb){
		var data;
		if(typeof url=='object'){
			data = url;
			url = '';
		}
		if(!data&&!url){
			cb&&cb();
			return;
		}
		var tmplHTML = this.$el.html().trim();
		if(data&&tmplHTML){
			var html = '';
			try{
				html =  template.compile(tmplHTML)(data);
			}catch(e){}
			cb&&cb(html, tmplHTML, data);
			return html;
		}
		
		this.getTemplate(function(tmpl){
			if(url){
				agilelite.util.ajax({
					url : url,
					dataType : 'json',
					success : function(data){
						var render = template.compile(tmpl);
                    	html = render(data);
                    	cb&&cb(html, tmpl, data);
					},
					error : function(){
						cb&&cb();
					}
				});
			}else{
				var render = template.compile(tmpl);
                html = render(data);
                cb&&cb(html, tmpl, data);
			}
		});

	};
	
	var _render = function(type, url, cb){
		var $el = this.$el;	
		this.render(url, function(html, tmpl, data){
			var $referObj = $el;
			var id = $el.attr('id');
			var tag = '#'+$el.attr('id');
			var $oldObj = $el.parent().find('[__inject_dependence__="'+tag+'"]');
			if(type=='replace'){
				$oldObj.remove();
				//$referObj = $el;
			}else if(type=='after'){
				$referObj = $oldObj.length==0?$el:$oldObj.last();
			}else if(type=='before'){
				//$referObj = $el;
			}
			var $html = $(html).attr('__inject_dependence__',tag);
			$referObj.after($html);			
			cb&&cb($html, tmpl, data);
			$el.trigger('renderEnd', [$html]);
		});
	};
	
	Template.prototype.renderReplace = function(url, cb){
		_render.call(this, 'replace', url, cb);
	};
	
	Template.prototype.renderAfter = function(url, cb){
		_render.call(this, 'after', url, cb);
	};
	Template.prototype.renderBefore = function(url, cb){
		_render.call(this, 'before', url, cb);
	};
	
	agilelite.register('template', function(selecotr){
		return new Template(selecotr);
	});
})(agilelite.$);

/*
 * 扩展JSON:agilelite.JSON.stringify和agilelite.JSON.parse，用法你懂
 * */
(function(){
	var JSON={};JSON.parse=function(str){try{return eval("("+str+")")}catch(e){return null}};JSON.stringify=function(o){var r=[];if(typeof o=="string"){return'"'+o.replace(/([\'\"\\])/g,"\\$1").replace(/(\n)/g,"\\n").replace(/(\r)/g,"\\r").replace(/(\t)/g,"\\t")+'"'}if(typeof o=="undefined"){return""}if(typeof o!="object"){return o.toString()}if(o===null){return null}if(o instanceof Array){for(var i=0;i<o.length;i++){r.push(this.stringify(o[i]))}r="["+r.join()+"]"}else{for(var i in o){r.push('"'+i+'":'+this.stringify(o[i]))}r="{"+r.join()+"}"}return r};JSON.toParams=function(json){var arr=[];for(var k in json){var str=[];var v=json[k];v=v instanceof Array?v:[v];for(var i=0;i<v.length;i++){str.push(v[i])}arr.push(k+"="+str.join(","))}return arr.join("&")};
	agilelite.register('JSON', JSON);
})();

(function(){
	var Base64={table:["A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","0","1","2","3","4","5","6","7","8","9","+","/"],UTF16ToUTF8:function(str){var res=[],len=str.length;for(var i=0;i<len;i++){var code=str.charCodeAt(i);if(code>0&&code<=127){res.push(str.charAt(i))}else{if(code>=128&&code<=2047){var byte1=192|((code>>6)&31);var byte2=128|(code&63);res.push(String.fromCharCode(byte1),String.fromCharCode(byte2))}else{if(code>=2048&&code<=65535){var byte1=224|((code>>12)&15);var byte2=128|((code>>6)&63);var byte3=128|(code&63);res.push(String.fromCharCode(byte1),String.fromCharCode(byte2),String.fromCharCode(byte3))}else{if(code>=65536&&code<=2097151){}else{if(code>=2097152&&code<=67108863){}else{}}}}}}return res.join("")},UTF8ToUTF16:function(str){var res=[],len=str.length;var i=0;for(var i=0;i<len;i++){var code=str.charCodeAt(i);if(((code>>7)&255)==0){res.push(str.charAt(i))}else{if(((code>>5)&255)==6){var code2=str.charCodeAt(++i);var byte1=(code&31)<<6;var byte2=code2&63;var utf16=byte1|byte2;res.push(Sting.fromCharCode(utf16))}else{if(((code>>4)&255)==14){var code2=str.charCodeAt(++i);var code3=str.charCodeAt(++i);var byte1=(code<<4)|((code2>>2)&15);var byte2=((code2&3)<<6)|(code3&63);utf16=((byte1&255)<<8)|byte2;res.push(String.fromCharCode(utf16))}else{if(((code>>3)&255)==30){}else{if(((code>>2)&255)==62){}else{}}}}}}return res.join("")},encode:function(str){if(!str){return""}var utf8=this.UTF16ToUTF8(str);var i=0;var len=utf8.length;var res=[];while(i<len){var c1=utf8.charCodeAt(i++)&255;res.push(this.table[c1>>2]);if(i==len){res.push(this.table[(c1&3)<<4]);res.push("==");break}var c2=utf8.charCodeAt(i++);if(i==len){res.push(this.table[((c1&3)<<4)|((c2>>4)&15)]);res.push(this.table[(c2&15)<<2]);res.push("=");break}var c3=utf8.charCodeAt(i++);res.push(this.table[((c1&3)<<4)|((c2>>4)&15)]);res.push(this.table[((c2&15)<<2)|((c3&192)>>6)]);res.push(this.table[c3&63])}return res.join("")},decode:function(str){if(!str){return""}var len=str.length;var i=0;var res=[];while(i<len){code1=this.table.indexOf(str.charAt(i++));code2=this.table.indexOf(str.charAt(i++));code3=this.table.indexOf(str.charAt(i++));code4=this.table.indexOf(str.charAt(i++));c1=(code1<<2)|(code2>>4);c2=((code2&15)<<4)|(code3>>2);c3=((code3&3)<<6)|code4;res.push(String.fromCharCode(c1));if(code3!=64){res.push(String.fromCharCode(c2))}if(code4!=64){res.push(String.fromCharCode(c3))}}return this.UTF8ToUTF16(res.join(""))}};
	agilelite.register('Base64', {
		encode : function(str){
			try{ return Base64.encode(str); }catch(e){ return ''; }
		},
		decode : function(str){
			try{ return Base64.decode(str); }catch(e){ return ''; }
		}
	});
})();
