//列表扩展
(function($) {

	ExList = {};

	function _controller(getli, option) {

		var $selector_li = null;

		function _init() {
			$selector_li = getli();

			_init_dom();
		}

		var swipe_option_width = null;

		var exlist_left_div_width = null;

		var status_code = {
			'left': 0,
			'normal': 1,
			'right': 2
		};

		var status = status_code.normal;

		function _init_dom() {
			$selector_li.addClass('swipe_block');

			$selector_li.each(function() {
				if ($(this).find('.swipe_option').length === 0) {
					$(this).append('<div class="swipe_option"></div>');
				}
				if ($(this).find('.exlist_left_div').length === 0) {
					$(this).append('<div class="exlist_left_div"></div>');
				}
			});

			$selector_li.find('.swipe_option').html(option.rightContent);

			$selector_li.find('.exlist_left_div').html(option.leftContent);

			swipe_option_width = $selector_li.children('.swipe_option').css('width');

			exlist_left_div_width = $selector_li.children('.exlist_left_div').css('width');

			$selector_li.children('.swipe_option').css('right', '-' + swipe_option_width);

			$selector_li.children('.exlist_left_div').css('left', '-' + exlist_left_div_width);

			$selector_li.off('swipeleft').on('swipeleft', swipeleft);
		}

		function fallback_swipeleft(e) {
			$selector_li.parent().find('li').animate({
				left: '0px'
			}, 100);
			status = status_code.normal;
			$selector_li.off('swipeleft').on('swipeleft', swipeleft);
		}

		function swipeleft(e) {
			var _li_element = $(e.currentTarget);

			fallback_swipeleft();

			_li_element.animate({
				left: '-' + swipe_option_width
			}, 100);

			$selector_li.parent().find('li').children(':not(.swipe_option)').off('tap', fallback_swipeleft).on('tap', fallback_swipeleft);

			_li_element.off('swiperight', fallback_swipeleft).on('swiperight', fallback_swipeleft);

			_li_element.children(':not(.swipe_option)').off('tap').on('tap', function() {
				_li_element.animate({
					left: '0px'
				}, 100);
				_li_element.children(':not(.swipe_option)').off('tap');
				return false;
			});

			_li_element.children('.swipe_option').off('tap').on('tap', function(e) {
				option.swipeOptionOnTap && option.swipeOptionOnTap(_li_element, $(e.target));
				return false;
			});

			status = status_code.right;
		}

		function _showLeft() {
			$selector_li.animate({
				left: exlist_left_div_width
			}, 100);
			$selector_li.off('swipeleft');
			$selector_li.parent().find('li').children(':not(.swipe_option)').off('tap', fallback_swipeleft);
			status = status_code.left;
		}

		function _hideLeft() {
			$selector_li.animate({
				left: '0px'
			}, 100);
			$selector_li.off('swipeleft').on('swipeleft', swipeleft);
			status = status_code.normal;
		}

		function _refresh() {
			_init();
		}

		_init();

		return {
			showLeft: _showLeft,
			hideLeft: _hideLeft,
			refresh: _refresh,
			hideOne: function(el) {
				var _el = $(el);
				_el.animate({
					left: '0px'
				}, 100);
			}
		};
	}

	ExList.liController = function($selector, opt) {
		return (function($selector, opt) {
			var option = {
				leftContent: '<div style="width:100%;height:100%;padding: 12px;" data-role="checkbox"><input class="exlist_checkbox" type="checkbox"/></div>',
				rightContent: '<div style="width:100%;height:100%;background-color: #FF2D2D;text-align: center;padding: 12px;color: #FFFFFF;">删除</div>'
			};

			$.extend(option, opt);

			return _controller(function() {
				var _$selector = $($selector);
				return _$selector;
			}, option);

		})($selector, opt);
	};

	ExList.ulController = function($selector, opt) {

		return (function($selector, opt) {
			var option = {
				leftContent: '<div style="width:100%;height:100%;padding: 12px;" data-role="checkbox"><input class="exlist_checkbox" type="checkbox"/></div>',
				rightContent: '<div style="width:100%;height:100%;background-color: #FF2D2D;text-align: center;padding: 12px;color: #FFFFFF;">删除</div>'
			};

			$.extend(option, opt);

			return _controller(function() {
				var _$selector = $($selector);
				return _$selector.find('li');
			}, option);

		})($selector, opt);
	};

	A.register('ExList', ExList);
})(A.$);