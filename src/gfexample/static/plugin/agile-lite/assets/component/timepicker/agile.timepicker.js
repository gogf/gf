//时间选择
(function($) {

	function addZero(number) {
		return (number < 10 ? '0' + number : number);
	}

	function getPage($el) {
		$el = $el.first().prev();
		var height = $el.get(0).offsetHeight;//$el.height();
		var eTop = $el.offset().top;
		var pTop = $el.closest('.agile-popup').offset().top;
		return Math.round((pTop - eTop) / height);
	}

	function solarDays(y, m) {
		var solarMonth = new Array(31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31);
		if (m == 1)
			return (((y % 4 === 0) && (y % 100 !== 0) || (y % 400 === 0)) ? 29 : 28);
		else
			return (solarMonth[m]);
	}

	var _Timepicker = function(opts) {
		var _this = this;

		var _onStart = null,
			_onSelected = null;

		var option = {
			hasSecond: true,
			isCustomLeftButton: false
		};

		$.extend(option, opts);

		_this.picker_selector_html = '' + '<div class="datepickerscroll" id="hourpickerscroll"><div class="scroller"><ul id="hourpicker" class="listitem"></ul></div></div>';

		if (option.hasSecond) {
			_this.picker_selector_html += '<div style="position: absolute;left: 72px;top: 52px;">:</div>';
			_this.picker_selector_html += '<div style="position: absolute;left: 148px;top: 52px;">:</div>';
		} else {
			_this.picker_selector_html += '<div style="position: absolute;left: 108px;top: 52px;">:</div>';
		}

		_this.picker_selector_html += '<div class="datepickerscroll" id="minutepickerscroll"><div class="scroller"><ul id="minutepicker" class="listitem"></ul></div></div>';

		if (option.hasSecond) {
			_this.picker_selector_html += '<div class="datepickerscroll" id="secondpickerscroll"><div class="scroller"><ul id="secondpicker" class="listitem"></ul></div></div>';
		}

		if (option.isCustomLeftButton) {
			_this.picker_selector_html += '<button id="clear_date_picker" style="width: 50%;background-color: #E8E8E8;color: #000000;" class="width-full">' + option.customLeftButtonName + '</button>';
		} else {
			_this.picker_selector_html += '<button id="cancel_date_picker" style="width: 50%;background-color: #E8E8E8;color: #000000;" class="width-full">取消</button>';
		}

		_this.picker_selector_html += '<button id="confirm_date_picker" style="width: 50%;background-color: #E8E8E8;color: #3779d0;" class="width-full">确定</button>';



		this.onStart = function(callback) {
			if (typeof callback == 'function') {
				_onStart = callback;
			}
		};

		this.onSelected = function(callback) {
			if (typeof callback == 'function') {
				_onSelected = callback;
			}
		};

		this.customLeftButtonCallback = function(callback) {
			if (typeof callback == 'function') {
				_this._customLeftButtonCallback = callback;
			}
		};

		this.destroy = function() {
			_this.hourScroll.destroy();
			_this.minuteScroll.destroy();
			if (option.hasSecond) {
				_this.secondScroll.destroy();
			}
		};

		if (typeof option.onSelected == 'function') {
			_this.onSelected(option.onSelected);
		}

		if (typeof option.onStart == 'function') {
			_this.onStart(option.onStart);
		}

		if (typeof option.customLeftButtonCallback == 'function') {
			_this.customLeftButtonCallback(option.customLeftButtonCallback);
		}


		this.select = function(selectedDate, callback) {
			if (typeof _onStart == 'function') {
				_onStart();
			}

			var _date_picker_popup = A.popup({
				html: _this.picker_selector_html,
				css: {
					width: '230px'
				},
				pos: 'center',
				id: 'timepickerpopup'
			});

			if (!option.hasSecond) {
				$('#hourpickerscroll').css('width', '50%').css('padding-left', '50px').css('padding-right', '30px');
				$('#minutepickerscroll').css('width', '50%').css('padding-right', '55px').css('padding-left', '25px');
			}

			var nbsp_html = '<li style="background-color: #E8E8E8;">&nbsp;</li>',
				html_hour = '<li style="background-color: #E8E8E8;" id="hourli-1">&nbsp;</li>',
				html_minute = '<li style="background-color: #E8E8E8;" id="minuteli-1">&nbsp;</li>',
				html_second = '<li style="background-color: #E8E8E8;" id="secondli-1">&nbsp;</li>';

			var $hourpicker = $('#hourpicker'),
				$minutepicker = $('#minutepicker'),
				$secondpicker = null;
			if (option.hasSecond) {
				$secondpicker = $('#secondpicker');
			}

			var _selected_date = (typeof selectedDate == 'object') ? selectedDate : new Date();

			_this.hourSelected = _selected_date.getHours();
			_this.minuteSelected = _selected_date.getMinutes();
			_this.secondSelected = _selected_date.getSeconds();

			var i = 0;

			for (i = 0; i < 24; i++) {
				html_hour += '<li style="background-color: #E8E8E8;" class="hourli" id="hourli' + i + '">' + addZero(i) + '</li>';
			}

			for (i = 0; i < 60; i++) {
				html_minute += '<li style="background-color: #E8E8E8;" class="minuteli" id="minuteli' + i + '">' + addZero(i) + '</li>';
			}

			for (i = 0; i < 60; i++) {
				html_second += '<li style="background-color: #E8E8E8;" class="secondli" id="secondli' + i + '">' + addZero(i) + '</li>';
			}

			html_hour += nbsp_html;
			html_minute += nbsp_html;
			html_second += nbsp_html;

			$hourpicker.html(html_hour);
			$minutepicker.html(html_minute);
			if (option.hasSecond) {
				$secondpicker.html(html_second);
			}

			_this.hourScroll = A.Scroll('#hourpickerscroll', {
				scrollbars: false,
				snap: 'li'
			});
			_this.minuteScroll = A.Scroll('#minutepickerscroll', {
				scrollbars: false,
				snap: 'li'
			});
			if (option.hasSecond) {
				_this.secondScroll = A.Scroll('#secondpickerscroll', {
					scrollbars: false,
					snap: 'li'
				});
			}

			//_this.hourScroll.scrollToElement('#hourli' + (_selected_date.getHours() - 1), 0);
			_this.hourScroll.goToPage(0, _selected_date.getHours(), 0);
			$('#hourli' + _selected_date.getHours()).addClass('selectedli');
			//_this.minuteScroll.scrollToElement('#minuteli' + (_selected_date.getMinutes() - 1), 0);
			_this.minuteScroll.goToPage(0, _selected_date.getMinutes(), 0);
			$('#minuteli' + _selected_date.getMinutes()).addClass('selectedli');
			if (option.hasSecond) {
				//_this.secondScroll.scrollToElement('#secondli' + (_selected_date.getSeconds() - 1), 0);
				_this.secondScroll.goToPage(0, _selected_date.getSeconds(), 0);
				$('#secondli' + _selected_date.getSeconds()).addClass('selectedli');
			}

			_this.hourScroll.on('scrollEnd', function(e) {
				$('.hourli').removeClass('selectedli');
				//var _hour = Math.round(_this.hourScroll.y / (-41));
				var _hour = getPage($('.hourli'));
				_this.hourSelected = _hour;
				$('#hourli' + _hour).addClass('selectedli');
			});

			_this.minuteScroll.on('scrollEnd', function(e) {
				$('.minuteli').removeClass('selectedli');
				//var _minute = Math.round(_this.minuteScroll.y / (-41));
				var _minute = getPage($('.minuteli'));
				_this.minuteSelected = _minute;
				$('#minuteli' + _minute).addClass('selectedli');
			});
			if (option.hasSecond) {
				_this.secondScroll.on('scrollEnd', function(e) {
					$('.secondli').removeClass('selectedli');
					//var _second = Math.round(_this.secondScroll.y / (-41));
					var _second = getPage($('.secondli'));
					_this.secondSelected = _second;
					$('#secondli' + _second).addClass('selectedli');
				});
			}
			$('#confirm_date_picker').off(A.options.clickEvent);

			$('#confirm_date_picker').on(A.options.clickEvent, function() {
				var _full_date = new Date(_selected_date.getFullYear(), _selected_date.getMonth(), _selected_date.getDate(), _this.hourSelected, _this.minuteSelected, option.hasSecond ? _this.secondSelected : 0),
					_full_string = _full_date.getFullYear() + '-' + addZero(_full_date.getMonth() + 1) + '-' + addZero(_full_date.getDate()) + ' ' + addZero(_full_date.getHours()) + ':' + addZero(_full_date.getMinutes()) + ':' + addZero(_full_date.getSeconds()),
					_date_string = _full_date.getFullYear() + '-' + addZero(_full_date.getMonth() + 1) + '-' + addZero(_full_date.getDate()),
					_time_string = '';
				if (option.hasSecond) {
					_time_string = addZero(_full_date.getHours()) + ':' + addZero(_full_date.getMinutes()) + ':' + addZero(_full_date.getSeconds());
				} else {
					_time_string = addZero(_full_date.getHours()) + ':' + addZero(_full_date.getMinutes());
				}
				if (typeof callback == 'function') {
					callback({
						date: _full_date,
						dateString: _date_string,
						timeString: _time_string,
						fullString: _full_string
					});
				}
				if (typeof _onSelected == 'function') {
					callback({
						date: _full_date,
						dateString: _date_string,
						timeString: _time_string,
						fullString: _full_string
					});
				}
				_date_picker_popup.close();
				return false;
			});

			$('#cancel_date_picker').on(A.options.clickEvent, function() {
				_date_picker_popup.close();
				return false;
			});


			$('#clear_date_picker').on(A.options.clickEvent, function() {
				if (typeof _this._customLeftButtonCallback == 'function') {
					_this._customLeftButtonCallback();
				}
				_date_picker_popup.close();
				return false;
			});

			_date_picker_popup.on('popupclose', function() {
				_this.destroy();
			});
		};
	};

	var _Datepicker = function(opts) {
		var _this = this;

		var _onStart = null,
			_onSelected = null;

		var option = {
			isCustomLeftButton: false
		};

		$.extend(option, opts);

		_this.picker_selector_html = '' + '<div class="datepickerscroll" id="yearpickerscroll" style="padding-right: 0px;"><div class="scroller"><ul id="yearpicker" class="listitem"></ul></div></div>' + '<div class="datepickerscroll" id="monthpickerscroll" style="padding-right: 0px;padding-left: 18px;"><div class="scroller"><ul id="monthpicker" class="listitem"></ul></div></div>' + '<div class="datepickerscroll" id="daypickerscroll"><div class="scroller"><ul id="daypicker" class="listitem"></ul></div></div>';

		if (option.isCustomLeftButton) {
			_this.picker_selector_html += '<button id="clear_date_picker" style="width: 50%;background-color: #E8E8E8;color: #000000;" class="width-full">' + option.customLeftButtonName + '</button>';
		} else {
			_this.picker_selector_html += '<button id="cancel_date_picker" style="width: 50%;background-color: #E8E8E8;color: #000000;" class="width-full">取消</button>';
		}

		_this.picker_selector_html += '<button id="confirm_date_picker" style="width: 50%;background-color: #E8E8E8;color: #3779d0;" class="width-full">确定</button>';


		this.onStart = function(callback) {
			if (typeof callback == 'function') {
				_onStart = callback;
			}
		};

		this.onSelected = function(callback) {
			if (typeof callback == 'function') {
				_onSelected = callback;
			}
		};

		this.customLeftButtonCallback = function(callback) {
			if (typeof callback == 'function') {
				_this._customLeftButtonCallback = callback;
			}
		};

		this.reCountDay = function() {
			var number_of_day = parseInt(solarDays(_this.yearSelected, (parseInt(_this.monthSelected, 10) - 1)), 10);
			var html_day = '<li style="background-color: #E8E8E8;" id="dayli0">&nbsp;</li>';
			for (var i = 1; i <= number_of_day; i++) {
				html_day += '<li style="background-color: #E8E8E8;" class="dayli" id="dayli' + i + '">' + addZero(i) + '</li>';
			}
			html_day += '<li style="background-color: #E8E8E8;">&nbsp;</li>';
			$('#daypicker').html(html_day);

			_this.dayScroll.destroy();
			_this.dayScroll = A.Scroll('#daypickerscroll', {
				scrollbars: false,
				snap: 'li'
			});
			if (_this.daySelected > number_of_day) {
				//_this.dayScroll.scrollToElement('#dayli' + (number_of_day - 1), 0);
				_this.dayScroll.goToPage(0, number_of_day - 1, 0);
				$('#dayli' + number_of_day).addClass('selectedli');
				_this.daySelected = number_of_day;
			} else {
				//_this.dayScroll.scrollToElement('#dayli' + (parseInt(_this.daySelected) - 1), 0);
				_this.dayScroll.goToPage(0, (parseInt(_this.daySelected, 10) - 1), 0);
				$('#dayli' + parseInt(_this.daySelected, 10)).addClass('selectedli');
			}

			_this.dayScroll.off('scrollEnd');

			_this.dayScroll.on('scrollEnd', function(e) {
				$('.dayli').removeClass('selectedli');
				var _day = Math.round(_this.dayScroll.y / (-41) + 1);
				_this.daySelected = _day;
				$('#dayli' + _day).addClass('selectedli');
			});
		};

		this.destroy = function() {
			_this.yearScroll.destroy();
			_this.monthScroll.destroy();
			_this.dayScroll.destroy();
		};

		if (typeof option.onSelected == 'function') {
			_this.onSelected(option.onSelected);
		}

		if (typeof option.onStart == 'function') {
			_this.onStart(option.onStart);
		}

		if (typeof option.customLeftButtonCallback == 'function') {
			_this.customLeftButtonCallback(option.customLeftButtonCallback);
		}

		this.select = function(selectedDate, callback) {
			if (typeof _onStart == 'function') {
				_onStart();
			}

			var _date_picker_popup = A.popup({
				html: _this.picker_selector_html,
				css: {
					width: '230px'
				},
				pos: 'center',
				id: 'timepickerpopup'
			});

			var nbsp_html = '<li style="background-color: #E8E8E8;">&nbsp;</li>',
				html_year = '<li style="background-color: #E8E8E8;" id="yearli1899">&nbsp;</li>',
				html_month = '<li style="background-color: #E8E8E8;" id="monthli0">&nbsp;</li>',
				html_day = '<li style="background-color: #E8E8E8;" id="dayli0">&nbsp;</li>';

			var $yearpicker = $('#yearpicker'),
				$monthpicker = $('#monthpicker'),
				$daypicker = $('#daypicker');

			var _selected_date = (typeof selectedDate == 'object') ? selectedDate : new Date();

			_this.yearSelected = _selected_date.getFullYear();
			_this.monthSelected = _selected_date.getMonth() + 1;
			_this.daySelected = _selected_date.getDate();

			var i = 0;

			for (i = 1900; i < 2050; i++) {
				html_year += '<li style="background-color: #E8E8E8;" class="yearli" id="yearli' + i + '">' + i + '</li>';
			}

			for (i = 1; i < 13; i++) {
				html_month += '<li style="background-color: #E8E8E8;" class="monthli" id="monthli' + i + '">' + i + '月</li>';
			}

			var number_of_day = solarDays(_this.yearSelected, (parseInt(_this.monthSelected, 10) - 1));

			for (i = 1; i <= number_of_day; i++) {
				html_day += '<li style="background-color: #E8E8E8;" class="dayli" id="dayli' + i + '">' + addZero(i) + '</li>';
			}

			html_year += nbsp_html;
			html_month += nbsp_html;
			html_day += nbsp_html;

			$yearpicker.html(html_year);
			$monthpicker.html(html_month);
			$daypicker.html(html_day);

			_this.yearScroll = A.Scroll('#yearpickerscroll', {
				mouseWheel: false,
				scrollbars: false,
				snap: 'li'
			});
			_this.monthScroll = A.Scroll('#monthpickerscroll', {
				mouseWheel: false,
				scrollbars: false,
				snap: 'li'
			});
			_this.dayScroll = A.Scroll('#daypickerscroll', {
				mouseWheel: false,
				scrollbars: false,
				snap: 'li'
			});
			var _yGap = 1900,
				_mGap = 1,
				_dGap = 1;
			//_this.yearScroll.scrollToElement('#yearli' + (_selected_date.getFullYear() - 1), 0);
			_this.yearScroll.goToPage(0, _selected_date.getFullYear() - _yGap, 0);
			$('#yearli' + _selected_date.getFullYear()).addClass('selectedli');
			//_this.monthScroll.scrollToElement('#monthli' + _selected_date.getMonth(), 0);
			_this.monthScroll.goToPage(0, _selected_date.getMonth() + 1 - _mGap, 0);
			$('#monthli' + (_selected_date.getMonth() + 1)).addClass('selectedli');
			//_this.dayScroll.scrollToElement('#dayli' + (_selected_date.getDate() - 1), 0);
			_this.dayScroll.goToPage(0, _selected_date.getDate() - _dGap, 0);
			$('#dayli' + _selected_date.getDate()).addClass('selectedli');

			_this.yearScroll.on('scrollEnd', function(e) {
				$('.yearli').removeClass('selectedli');
				//var _year = Math.round(_this.yearScroll.y / (-41) + 1900);
				var _year = getPage($('.yearli')) + _yGap;
				_this.yearSelected = _year;
				$('#yearli' + _year).addClass('selectedli');
				_this.reCountDay();
			});

			_this.monthScroll.on('scrollEnd', function(e) {
				$('.monthli').removeClass('selectedli');
				//var _month = Math.round(_this.monthScroll.y / (-41) + 1);
				var _month = getPage($('.monthli')) + _mGap;
				_this.monthSelected = _month;
				$('#monthli' + _month).addClass('selectedli');
				_this.reCountDay();
			});

			_this.dayScroll.on('scrollEnd', function(e) {
				$('.dayli').removeClass('selectedli');
				//var _day = Math.round(_this.dayScroll.y / (-41) + 1);
				var _day = getPage($('.dayli')) + _dGap;
				_this.daySelected = _day;
				$('#dayli' + _day).addClass('selectedli');
			});

			$('#confirm_date_picker').off(A.options.clickEvent);

			$('#confirm_date_picker').on(A.options.clickEvent, function() {
				var _full_date = new Date(_this.yearSelected, _this.monthSelected - 1, _this.daySelected, _selected_date.getHours(), _selected_date.getMinutes(), _selected_date.getSeconds()),
					_full_string = _full_date.getFullYear() + '-' + addZero(_full_date.getMonth() + 1) + '-' + addZero(_full_date.getDate()) + ' ' + addZero(_full_date.getHours()) + ':' + addZero(_full_date.getMinutes()) + ':' + addZero(_full_date.getSeconds()),
					_date_string = _full_date.getFullYear() + '-' + addZero(_full_date.getMonth() + 1) + '-' + addZero(_full_date.getDate()),
					_time_string = addZero(_full_date.getHours()) + ':' + addZero(_full_date.getMinutes()) + ':' + addZero(_full_date.getSeconds());

				if (typeof callback == 'function') {
					callback({
						date: _full_date,
						dateString: _date_string,
						timeString: _time_string,
						fullString: _full_string
					});
				}
				if (typeof _onSelected == 'function') {
					_onSelected({
						date: _full_date,
						dateString: _date_string,
						timeString: _time_string,
						fullString: _full_string
					});
				}
				_date_picker_popup.close();

				return false;
			});

			$('#cancel_date_picker').on(A.options.clickEvent, function() {
				_date_picker_popup.close();
				return false;
			});

			$('#clear_date_picker').on(A.options.clickEvent, function() {
				if (typeof _this._customLeftButtonCallback == 'function') {
					_this._customLeftButtonCallback();
				}
				_date_picker_popup.close();
				return false;
			});

			_date_picker_popup.on('popupclose', function() {
				_this.destroy();
			});
		};

	};

	function Datepicker(option) {
		return new _Datepicker(option);
	}

	function Timepicker(option) {
		return new _Timepicker(option);
	}

	A.register('Datepicker', Datepicker);

	A.register('Timepicker', Timepicker);
})(A.$);