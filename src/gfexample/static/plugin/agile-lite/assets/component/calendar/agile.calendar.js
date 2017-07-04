//日历控件
(function($) {
	var _millisecondsPerDay = 86400000;

	function calendar(y, m, d) {
		function solarDays(y, m) {
			var solarMonth = new Array(31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31);
			if (m == 1)
				return (((y % 4 == 0) && (y % 100 != 0) || (y % 400 == 0)) ? 29 : 28);
			else
				return (solarMonth[m]);
		}

		var sDObj = new Date(y, m, 1);
		this.firstDayDataObj = sDObj; //当月第一天的日期
		if (d) {
			this.dayDataObj = new Date(y, m, d);
		} else {
			this.dayDataObj = sDObj;
		}
		this.length = solarDays(y, m); //公历当月天数
		this.firstWeek = sDObj.getDay(); //公历当月1日星期几
		this.solarDays = solarDays;
		this.previousMonthLength = m == 0 ? solarDays(y - 1, 11) : solarDays(y, m - 1);
	}

	function _initMonthDom(selector, _calendarObj) {
		var $selector = $(selector);

		var selector_id = $selector.attr('id');

		var html = "";

		html += '<table class="calendar_table" style="width: 100%;" cellspacing="0px">' +
			'<thead>' +
			'<tr style="text-align: left;">' +
			'<td colspan="7">' +

			'<div style="float: left;">' +
			'<span id="' + selector_id + '_year_decrease" class="iconfont iconline-arrow-left"></span>' +

			'<select id="' + selector_id + '_year_selector" style="padding-right: 0px;padding-left: 0px;font-size: 11pt;padding-bottom: 1px;border-width: 0px;background-color: transparent;">' +
			'</select>' +

			'<span id="' + selector_id + '_year_increase" class="iconfont iconline-arrow-right"></span>' +
			'</div>' +

			'<div style="float: left;">' +
			'&nbsp;&nbsp;&nbsp;&nbsp;' +
			'</div>' +

			'<div style="float: left;">' +
			'<span id="' + selector_id + '_month_decrease" class="iconfont iconline-arrow-left"></span>' +

			'<select id="' + selector_id + '_month_selector" style="padding-right: 0px;padding-left: 0px;font-size: 11pt;padding-bottom: 1px;border-width: 0px;background-color: transparent;">' +
			'</select>' +

			'<span id="' + selector_id + '_month_increase" class="iconfont iconline-arrow-right"></span>' +
			'</div>' +

			'</td>' +
			'</tr>' +
			'<tr class="calendar_header">' +
			'<td width="54" style="color: #FF0000;padding-bottom: 6px;">日</td>' +
			'<td width="54" style="padding-bottom: 6px;">一</td>' +
			'<td width="54" style="padding-bottom: 6px;">二</td>' +
			'<td width="54" style="padding-bottom: 6px;">三</td>' +
			'<td width="54" style="padding-bottom: 6px;">四</td>' +
			'<td width="54" style="padding-bottom: 6px;">五</td>' +
			'<td width="54" style="color: #FF0000;padding-bottom: 6px;">六</td>' +
			'</tr>' +
			'</thead>' +
			'<tbody class="calendar_tbody">';
		var gNum;
		for (var i = 0; i < _calendarObj.opts.row; i++) {
			html += '<tr class="calendar_tr">';
			for (var j = 0; j < 7; j++) {
				gNum = i * 7 + j;
				html += '<td id="' + selector_id + '_TD_' + gNum + '" class="calendar_td">';
				/*
				html+='<font id="'+selector_id+'_DF_' + gNum +'"></font>';
				html+='<br/><font id="'+selector_id+'_MF_' + gNum + '"></font>';
				*/
				html += '</td>';
			}
			html += '</tr>';
		}
		html += '</tbody>' +
			'</table>';

		$selector.html(html);

	};

	function _isSameDate(dateA, dateB) {
		if (dateA.getFullYear() == dateB.getFullYear() && dateA.getMonth() == dateB.getMonth() && dateA.getDate() == dateB.getDate()) {
			return true;
		}

		return false;
	}

	function _drawMonthDate(selector, _calendar, _calendarObj) {
		var $selector = $(selector);

		var selector_id = $selector.attr('id');

		var year = _calendar.firstDayDataObj.getFullYear(),
			month = _calendar.firstDayDataObj.getMonth();

		var numOfTd = _calendarObj.opts.row * 7;

		for (var i = 0; i < numOfTd; i++) {

			var thisDate = _MonthTableTodate(_calendarObj, i);

			var td_html = '';

			td_html += '<font id="' + selector_id + '_MK_' + i + '" ' + 'class="calendar_MK_font">&nbsp;</font><br/>';

			td_html += '<font id="' + selector_id + '_DF_' + i + '" ' + 'class="calendar_DF_font';

			td_html += '" ';

			if (year != thisDate.getFullYear() || month != thisDate.getMonth()) {
				td_html += 'style="opacity: 0.3;"';
			}

			td_html += ' data-day="' + thisDate.getFullYear() + '-' + (thisDate.getMonth() + 1) + '-' + thisDate.getDate() + '">' + thisDate.getDate() + '</font>' + '<br/><font class="calendar_MF_font" id="' + selector_id + '_MF_' + i + '">&nbsp;</font>';

			$('#' + selector_id + '_TD_' + i).removeClass().addClass('calendar_td');

			$('#' + selector_id + '_TD_' + i).html(td_html);


		}

		var html = "";

		for (var i = 1900; i < 2050; i++) {
			html += '<option value="' + i + '" ';
			if (i == year) {
				html += 'selected';
			}
			html += '>' + i + '</option>';
		};

		$('#' + selector_id + '_year_selector').html(html);

		html = "";

		for (var i = 1; i < 13; i++) {
			html += '<option value="' + (i - 1) + '" ';
			if (i == (month + 1)) {
				html += 'selected';
			}
			html += '>' + (i > 9 ? i : '&nbsp;' + i) + '</option>';
		};

		$('#' + selector_id + '_month_selector').html(html);
	}

	function _drawMonthMark(_calendarObj, marks) {
		var selector_id = _calendarObj.$selector.attr('id');

		var color_plan = ['plan_one', 'plan_two'];

		$('#' + selector_id + ' .calendar_MF_font').html('&nbsp;').removeClass().addClass('calendar_MF_font');

		$('#' + selector_id + ' .calendar_MK_font').html('&nbsp;').removeClass().addClass('calendar_MK_font');

		$('#' + selector_id + ' .calendar_td').removeClass().addClass('calendar_td');

		var today_date = new Date(),
			todayTdNum = _dateToMonthTable(_calendarObj, today_date);

		if (todayTdNum != null) {
			$('#' + selector_id + '_TD_' + todayTdNum).addClass('today_highlight');
		}

		if (_calendarObj.opts.selectedDate == null && todayTdNum != null) {
			_calendarObj.opts.selectedDate = today_date;
			_calendarObj.opts.selectedDateString = today_date.getFullYear() + '-' + (today_date.getMonth() + 1) + '-' + today_date.getDate();
			$('#' + selector_id + '_TD_' + todayTdNum).addClass('tap_highlight');
		} else {
			var selectedTdNum = _dateToMonthTable(_calendarObj, _calendarObj.opts.selectedDate);

			if (selectedTdNum != null) {
				$('#' + selector_id + '_TD_' + selectedTdNum).addClass('tap_highlight');
			} else {
				_calendarObj.opts.selectedDate = null;
				_calendarObj.opts.selectedDateString = null;
			}
		}

		for (var key in marks) {
			var _date_string = key.split('-');

			var _date = new Date(_date_string[0], parseInt(_date_string[1], 10) - 1, _date_string[2]);

			var tdNum = _dateToMonthTable(_calendarObj, _date);

			if (tdNum == null) {
				continue;
			}

			var $td = $('#' + selector_id + '_TD_' + tdNum);

			if (marks[key].formClass) {
				$td.addClass(marks[key].formClass);
			}

			if (marks[key].bottom) {
				if (marks[key].bottom.class) {
					$('#' + selector_id + '_MF_' + tdNum).addClass(marks[key].bottom.class);
				}
				if (marks[key].bottom.content) {
					$('#' + selector_id + '_MF_' + tdNum).html(marks[key].bottom.content);
				}
			}

			if (marks[key].top) {
				if (marks[key].top.class) {
					$('#' + selector_id + '_MK_' + tdNum).addClass(marks[key].top.class);
				}
				if (marks[key].top.content) {
					$('#' + selector_id + '_MK_' + tdNum).html(marks[key].top.content);
				}
			}

			if (!marks[key].hideCount) {
				if (marks[key].data && marks[key].data.length > color_plan.length) {
					$td.addClass('plan_many');
				} else if (marks[key].data && marks[key].data.length != 0) {
					$td.addClass(color_plan[marks[key].data.length - 1]);
				}
			}
		}
	}

	function _dateToMonthTable(_calendarObj, date) {
		if (!date) {
			return null;
		}
		var _calendar = new calendar(_calendarObj.selectedYear, _calendarObj.selectedMonth);

		var startDate = new Date(_calendarObj.selectedYear, _calendarObj.selectedMonth, 1),
			endDate = new Date(_calendarObj.selectedYear, _calendarObj.selectedMonth, _calendar.length, 23, 59, 59),
			numOfBeforeDay = _calendarObj.opts.startRow * 7 + _calendar.firstWeek,
			numOfAfterDay = 7 * _calendarObj.opts.row - (numOfBeforeDay + _calendar.length);

		if (date.getTime() < startDate.getTime()) {
			if (date.getTime() < (startDate.getTime() - _millisecondsPerDay * (numOfBeforeDay))) {
				return null;
			}
			return numOfBeforeDay - Math.ceil((startDate.getTime() - date.getTime()) / _millisecondsPerDay);
		} else if (date.getTime() > endDate.getTime()) {
			if (date.getTime() > (endDate.getTime() + _millisecondsPerDay * (numOfAfterDay))) {
				return null;
			}
			return numOfBeforeDay + Math.floor((date.getTime() - startDate.getTime()) / _millisecondsPerDay);
		} else {
			return numOfBeforeDay + Math.floor((date.getTime() - startDate.getTime()) / _millisecondsPerDay);
		}
	}

	function _MonthTableTodate(_calendarObj, tdNum) {
		var _calendar = new calendar(_calendarObj.selectedYear, _calendarObj.selectedMonth);

		var startDate = new Date(_calendarObj.selectedYear, _calendarObj.selectedMonth, 1),
			numOfBeforeDay = _calendarObj.opts.startRow * 7 + _calendar.firstWeek,
			firstTdDate = new Date(startDate.getTime() - numOfBeforeDay * _millisecondsPerDay);

		return new Date(firstTdDate.getTime() + tdNum * _millisecondsPerDay);
	}

	var _Calendar = function(selector, opts) {

		var defaultOpts = {
			show_day: null,
			marks: {},
			row: 6,
			startRow: null,
			autoJump: true
		};

		$.extend(defaultOpts, opts);

		if (!defaultOpts.startRow) {
			if (defaultOpts.row > 6) {
				defaultOpts.startRow = Math.round((defaultOpts.row - 7) / 2);
			} else {
				defaultOpts.startRow = 0;
			}
		}

		var _this = this;

		_this.$selector = $(selector);
		_this.selectedYear = 1900;
		_this.selectedMonth = 0;
		_this.opts = defaultOpts;
		_this.opts.selectedDate = null;
		_this.opts.selectedDateString = '';
		_this.opts.refreshShot = {};

		var selector_id = _this.$selector.attr('id');

		_initMonthDom(_this.$selector, _this);

		var today = new Date();

		if (_this.opts.show_day) {
			_this.selectedYear = _this.opts.show_day.getFullYear();
			_this.selectedMonth = _this.opts.show_day.getMonth();
		} else {
			_this.selectedYear = today.getFullYear();
			_this.selectedMonth = today.getMonth();
		}

		_this.opts.refreshShot.selectedYear = _this.selectedYear;
		_this.opts.refreshShot.selectedMonth = _this.selectedMonth;

		var _calendar = new calendar(_this.selectedYear, _this.selectedMonth);

		_drawMonthDate(_this.$selector, _calendar, _this);

		$('#' + selector_id + '_year_selector').change(function(e) {
			_this.selectedYear = parseInt($(e.target).val());
			_this.refresh(1);
		});

		$('#' + selector_id + '_month_selector').change(function(e) {
			_this.selectedMonth = parseInt($(e.target).val());
			_this.refresh(1);
		});

		_drawMonthMark(_this, _this.opts.marks);

		_this.onTap();

		_this.$selector.on('swipeleft', function(e) {
			A.anim.run(_this.$selector.find('tbody'), 'slideLeftOut');
			_this.monthJump(1);
		});

		_this.$selector.on('swiperight', function(e) {
			A.anim.run(_this.$selector.find('tbody'), 'slideRightOut');
			_this.monthJump(-1);
		});

		$('#' + selector_id + '_year_decrease').on('tap', function() {
			if (_this.selectedYear < 1900)
				return;
			_this.selectedYear -= 1;
			A.anim.run(_this.$selector, 'fadeIn');
			_this.refresh(1);
		});

		$('#' + selector_id + '_year_increase').on('tap', function() {
			if (_this.selectedYear > 2050)
				return;
			_this.selectedYear += 1;
			A.anim.run(_this.$selector, 'fadeIn');
			_this.refresh(1);
		});

		$('#' + selector_id + '_month_decrease').on('tap', function() {
			A.anim.run(_this.$selector.find('tbody'), 'slideRightOut');
			_this.monthJump(-1);
		});

		$('#' + selector_id + '_month_increase').on('tap', function() {
			A.anim.run(_this.$selector.find('tbody'), 'slideLeftOut');
			_this.monthJump(1);
		});
	};

	_Calendar.prototype.onTap = function(callback) {

		var _this = this;

		var selector_id = _this.$selector.attr('id');

		function _oneClick(e) {
			var $target = $(e.target);

			if ($target.get(0).tagName == 'FONT') {
				$target = $target.parent('td');
			}

			var $calendar_DF_font = $target.children(".calendar_DF_font");

			if ($calendar_DF_font.length != 0) {
				var day = $calendar_DF_font.data('day');

				if (day) {

					var _day = day.split('-');

					var date_data = new Date(_day[0], parseInt(_day[1]) - 1, _day[2]);

					_this.opts.selectedDate = date_data;

					_this.opts.selectedDateString = day;

					if (_this.opts.autoJump &&
						(date_data.getFullYear() != _this.selectedYear || date_data.getMonth() != _this.selectedMonth)) {
						_this.selectedYear = date_data.getFullYear();
						_this.selectedMonth = date_data.getMonth();
						A.anim.run(_this.$selector, 'fadeIn');
						_this.refresh(1);
						$target = $('#' + selector_id + '_TD_' + _dateToMonthTable(_this, date_data));
					} else {
						_this.refresh(0);
					}

					var marks_array = _this.opts.marks[day] || {};

					if (marks_array.data && marks_array.data.length == 0) {
						marks_array = {};
					}

					callback && callback({
						date: date_data,
						mark: marks_array,
						element: $target
					});
				}

			}

		}
		$('#' + selector_id + ' .calendar_td').off('tap');

		$('#' + selector_id + ' .calendar_td').on('tap', _oneClick);

		//$('.calendar_td').on('click',_oneClick);

	};

	_Calendar.prototype.refresh = function(level, callback) {
		var _this = this;

		var _refresh_shot = _this.opts.refreshShot;

		//$('.calendar_table').animateObj('slideLeftOut');

		//A.anim.run(_this.$selector,'fadeIn');

		if (!level) {
			level = 0;
		}

		if (level >= 2) {
			_initMonthDom(_this.$selector, _this)
		}

		var _calendar = new calendar(_this.selectedYear, _this.selectedMonth);

		_this._calendar = _calendar;

		if (level >= 1) {
			_drawMonthDate(_this.$selector, _calendar, _this);
		}

		_drawMonthMark(_this, _this.opts.marks);

		if (_this.selectedYear != _refresh_shot.selectedYear || _this.selectedMonth != _refresh_shot.selectedMonth) {
			_this.onChangeCallback && _this.onChangeCallback({
				selectedYear: _this.selectedYear,
				selectedMonth: parseInt(_this.selectedMonth) + 1
			});
		}

		_refresh_shot.selectedYear = _this.selectedYear;

		_refresh_shot.selectedMonth = _this.selectedMonth;

		_this.opts.refreshShot = _refresh_shot;

		callback && callback();
	};

	_Calendar.prototype.onChange = function(callback) {
		this.onChangeCallback = callback;
	}

	_Calendar.prototype.goto = function(data, callback) {
		if (!data) {
			return;
		}

		var _this = this;

		_this.selectedYear = data.getFullYear();
		_this.selectedMonth = data.getMonth();

		_this.opts.selectedDate = data;
		_this.opts.selectedDateString = data.getFullYear() + '-' + (data.getMonth() + 1) + '-' + data.getDate();


		A.anim.run(_this.$selector, 'fadeIn');

		_this.refresh(1, callback);
	};

	_Calendar.prototype.monthJump = function(step, callback) {
		if (isNaN(step)) {
			return;
		}

		var _this = this;

		var month = _this.selectedMonth,
			difference = parseInt(month) + step;

		if (difference < 0) {
			_this.selectedYear = parseInt(_this.selectedYear) + (difference - difference % 12 - 12) / 12;
			_this.selectedMonth = 11 + (difference + 1) % 12;
		} else if (difference > 11) {
			_this.selectedYear = parseInt(_this.selectedYear) + (difference - difference % 12) / 12;
			_this.selectedMonth = difference % 12;
		} else {
			_this.selectedMonth = parseInt(_this.selectedMonth) + step;
		}

		//A.anim.run(_this.$selector,'fadeIn');

		_this.refresh(1, callback);
	};

	_Calendar.prototype.addMarkData = function(date_string, data) {
		if (this.opts.marks[date_string] && data) {
			if (this.opts.marks[date_string].data) {
				this.opts.marks[date_string].data.push(data);
			} else {
				this.opts.marks[date_string].data = [data];
			}
		} else if (data) {
			this.opts.marks[date_string] = {
				data: [data]
			};
		}
	};

	_Calendar.prototype.setMarkAllData = function(date_string, data) {
		if (this.opts.marks[date_string] && data) {
			this.opts.marks[date_string].data = data;
		}
	};

	_Calendar.prototype.getMarkAllData = function(date_string) {
		return this.opts.marks[date_string] ? this.opts.marks[date_string].data : null;
	};

	_Calendar.prototype.getMark = function(date_string) {
		return this.opts.marks[date_string] ? this.opts.marks[date_string] : null;
	};

	_Calendar.prototype.setMark = function(date_string, data) {
		this.opts.marks[date_string] = data;
	};

	_Calendar.prototype.removeMarkAllData = function(date_string) {
		if (this.opts.marks[date_string]) {
			this.opts.marks[date_string].data = [];
		}
	};

	_Calendar.prototype.removeMark = function(date_string) {
		if (this.opts.marks[date_string]) {
			delete this.opts.marks[date_string];
		}
	};

	var Calendar = function(selector, opts) {
		return new _Calendar(selector, opts);
	};

	A.register('Calendar', Calendar);
})(A.$);