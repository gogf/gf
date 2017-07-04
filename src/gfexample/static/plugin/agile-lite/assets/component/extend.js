//扩展日期选择器组件
A.Component.add({
	datetime: {
		selector: '[data-role="article"].active',
		event: 'articleload',
		handler: function(el, roleType) {
			var $el = $(el);
			var _work = function($el) {
				var returnObj = {
					open : function(){
						$el.trigger(A.options.clickEvent);
					},
					clear : function(){
						$el.find('label').html($el.data('placeholder')||'');
						$el.find('input').val('');
					}
				};
				if(A.Component.isInit($el)){
					return returnObj;
				}
				var $label = $el.find('label'),$input,placeholder;
				if($label.length==1){
					$input = $el.find('input');
					placeholder = $label.html();
				}else{
					$input = $el;
					placeholder = '';
				}
				$el.data('placeholder', placeholder);
				$el.on(A.options.clickEvent, function(e) {
					function _clear() {
						triggerChange({
							date : '',
							str : ''
						});
					}
					
					function triggerChange(data){
						var _date = data.date,
						_date_string = data.str;
						if($input.val()!=_date_string){
							$label.html(_date_string?_date_string:placeholder);
							$input.val(_date_string||'');
							var _changeFunc = $el.trigger('datachange', [_date_string||'']).data('change');
							if(!_changeFunc) return;
							var _replace = function(){
								try{ eval(_changeFunc);}catch(e){ console.log(e); };
							};
							_replace.apply($input[0]);
						}
					}

					function timepicker_callback(e) {
						var _date = null,
							_date_string_array = $input.val().split(':'),
							today = new Date();

						if (_date_string_array.length == 2) {
							_date = new Date(today.getFullYear(), today.getMonth(), today.getDate(), parseInt(_date_string_array[0], 10), parseInt(_date_string_array[1], 10), 0);
						} else {
							_date = today;
						}

						picker.select(_date, function(data) {
							triggerChange({
								date : data.date,
								str : data.timeString
							});
						});

						return false;
					}

					function datepicker_callback() {
						var _date = null,
							_date_string_array = $input.val().split('-');

						if (_date_string_array.length == 3) {
							_date = new Date(_date_string_array[0], parseInt(_date_string_array[1]) - 1, _date_string_array[2]);
						} else {
							_date = new Date();
						}

						return picker.select(_date, function(data) {
							triggerChange({
								date : data.date,
								str : data.dateString
							});
						});
					}

					var type = $el.data('role')||$el.attr('type');

					var picker;

					if (type == 'time') {
						picker = A.Timepicker({
							hasSecond: false,
							isCustomLeftButton: true,
							customLeftButtonName:'清除',
							customLeftButtonCallback: _clear
						});
						timepicker_callback(e);
					} else if (type == 'date') {
						picker = A.Datepicker({
							hasClear: true,
							isCustomLeftButton: true,
							customLeftButtonName:'清除',
							customLeftButtonCallback: _clear
						});
						datepicker_callback(e);
					}

					return false;
				});
				$label.html($input.val() || placeholder);
				return returnObj;
			};

			if ($el.data('role') == 'date' || $el.data('role') == 'time' || $el.attr('type') == 'date' || $el.attr('type') == 'time') {
				return _work($el);
			} else {
				var components = $el.find('[data-role="date"],[data-role="time"],input[type="date"],input[type="time"]');
				for (var i = 0; i < components.length; i++) {
					_work($(components[i]));
				}
			}

		},
		extend : {
			open : function(){
				this.trigger(A.options.clickEvent);
			},
			clear : function(){
				var $el = this;
				$el.find('label').html($el.data('placeholder')||'');
				$el.find('input').val('');
			}
		}
	}
});
