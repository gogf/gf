<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <title>{$config['Sites']['admin']['name']}</title>
    <meta name="description" content="" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <link rel="shortcut icon" href="/static/resource/images/favicon16.ico" type="image/x-icon" />
    <link rel="shortcut icon" href="/static/resource/images/favicon32.ico" type="image/x-icon" />

    <!-- basic styles -->
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/bootstrap.min.css"/>
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/bootstrap-editable.css" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/font-awesome.min.css" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/ace-fonts.css" />

    <!-- page specific plugin styles -->

    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/select2.css" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/jquery.gritter.css" />


    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/jquery-ui-1.10.3.custom.min.css" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/chosen.css" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/datepicker.css" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/bootstrap-timepicker.css" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/daterangepicker.css" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/colorpicker.css" />
    <!-- ace styles -->
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/ace.min.css" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/ace-rtl.min.css" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/ace-skins.min.css" />



    <script type="text/javascript">
         window.jQuery || document.write("<script src='/static/plugin/ace-admin/assets/js/jquery-2.0.3.min.js'>"+"<"+"/script>");
        //window.jQuery || document.write("<script src='//cdn.bootcss.com/jquery/3.1.1/jquery.min.js'>"+"<"+"/script>");
    </script>


    <!-- inline styles related to this page -->
    <!-- ace settings handler -->
    <script src="/static/plugin/ace-admin/assets/js/ace-extra.min.js"></script>
    <!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media queries -->

    <!-- ace scripts -->
    <script src="/static/plugin/ace-admin/assets/js/ace-elements.min.js"></script>
    <script src="/static/plugin/ace-admin/assets/js/ace.min.js"></script>
</head>

<body class="no-skin">

<div id="navbar" class="navbar navbar-default">
    <script type="text/javascript">
        try{ace.settings.check('navbar' , 'fixed')}catch(e){}
    </script>

    <div class="navbar-container" id="navbar-container">
        <div class="navbar-header pull-left">
            <!--<a href="/">
                <div class="navbar-brand-text">{$config['Sites']['admin']['name']}</div>
            </a>-->
            <a href="#" class="navbar-brand" style="padding-top: 6px;">
				<small style="height: 30px;line-height: 30px;font-size: 20px;">
					<!--<i class="fa fa-leaf"></i>-->
					<img src="/static/plugin/ace-admin/assets/img/agentlogo.png" style="width: 30px;height: 30px;margin-top: -5px;"/>
                    {$config['Sites']['admin']['name']}
				</small>
			</a>
        </div><!-- /.navbar-header -->

        <div class="navbar-buttons navbar-header pull-right" role="navigation">
            <ul class=   "nav ace-nav">
            	<!--消息中心-->
    			<li class="green">
					<a  class="dropdown-toggle" href="/message/index">
						<i class="ace-icon fa fa-envelope icon-animated-vertical"></i>
						<span class="badge badge-success">{$_SESSION['user']['unreadMessageCount']}</span>
					</a>
				</li>
				<!--个人中心-->
                <li class="light-blue">
                    <a data-toggle="dropdown" href="#" class="dropdown-toggle">
                        欢迎您！ <b>{$_SESSION['user']['nickname']}</b>
                        <i class="ace-icon fa fa-caret-down"></i>
                    </a>

                    <ul class="user-menu dropdown-menu-right dropdown-menu dropdown-yellow dropdown-caret dropdown-close">
                        {*<li>*}
                            {*<a href="##">*}
                                {*<i class="ace-icon fa fa-user"></i>*}
                              		{*个人信息*}
                            {*</a>*}
                        {*</li>*}
                        <li>
                            <a href="javascript:confirmLogout();"><i class="ace-icon fa fa-power-off"></i>注销登录</a>
                        </li>
                    </ul>
                </li>
            </ul><!-- /.ace-nav -->
        </div><!-- /.navbar-header -->
    </div><!-- /.container -->
</div>

<div class="main-container" id="main-container">
    <div class="main-container-inner">
        <a class="menu-toggler" id="menu-toggler" href="#">
            <span class="menu-text"></span>
        </a>
        <div class="sidebar responsive-min" id="sidebar">
            <ul class="nav nav-list">
                {* 一级菜单遍历 *}
                {foreach from=$menus index=$i key=$menuKey item=$menu }
                    {if $menu['subs']}
                        <li {if $menu['active']}class="active open"{/if}>
                            <a href="{$menu['url']}" class="dropdown-toggle">
                                <i class="{$menu['icon']}"></i>
                                <span class="menu-text"> {$menu['name']} </span>
                                <b class="arrow fa fa-angle-down"></b>
                            </a>
                            <ul class="submenu">
                                {* 二级菜单遍历 *}
                                {foreach from=$menu['subs'] index=$index key=$subKey item=$subMenu}
                                    {if $subMenu['subs']}
                                        <li class="hsub {if $subMenu['active']}active open{/if}">
                                            <a href="#" class="dropdown-toggle">
                                                <i class="menu-icon fa fa-caret-right"></i>
                                                {$subMenu['name']}
                                                <b class="arrow fa fa-angle-down"></b>
                                            </a>
                                            <ul class="submenu">
                                                {foreach from=$subMenu['subs'] index=$index2 key=$subKey2 item=$subMenu2}
                                                    <li class="{if $subMenu2['active']}active{/if}">
                                                        <a href="{$subMenu2['url']}">
                                                            {if $subMenu2['icon']}
                                                                <i class="{$subMenu2['icon']}"></i>
                                                            {else}
                                                                <i class="menu-icon fa fa-caret-right"></i>
                                                            {/if}
                                                            &nbsp;{$subMenu2['name']}
                                                        </a>
                                                    </li>
                                                {/foreach}
                                            </ul>
                                        </li>
                                    {else}
                                        <li {if $subMenu['active']}class="active"{/if}>
                                            <a href="{$subMenu['url']}">
                                                <i class="menu-icon fa fa-caret-right"></i>
                                                {$subMenu['name']}
                                            </a>
                                        </li>
                                    {/if}
                                {/foreach}
                            </ul>
                        </li>
                    {else}
                        {* 只有一级菜单，没有下级 *}
                        <li {if $menu['active']}class="active"{/if}>
                            <a href="{$menu['url']}">
                                <i class="menu-icon fa {$menu['icon']}"></i>
                                <span class="menu-text"> {$menu['name']} </span>
                            </a>
                        </li>
                    {/if}
                {/foreach}
            </ul><!-- /.nav-list -->

            <div class="sidebar-toggle sidebar-collapse" id="sidebar-collapse">
                <i class="ace-icon fa fa-angle-double-left" data-icon1="ace-icon fa fa-angle-double-left" data-icon2="ace-icon fa fa-angle-double-right"></i>
            </div>

            <script type="text/javascript">
                try{ace.settings.check('sidebar' , 'collapsed')}catch(e){}
            </script>
        </div>

        <div class="main-content">
        	
           {if $breadCrumbs}
                <div class="breadcrumbs" id="breadcrumbs">
                    <script type="text/javascript">
                        try{ace.settings.check('breadcrumbs' , 'fixed')}catch(e){ }
                    </script>

                    <ul class="breadcrumb">
                        {foreach from=$breadCrumbs index=$index key=$key item=$item}
                            {if $item['url']}
                                <li>
                                    {if $item['icon']}<i class="{$item['icon']} home-icon"></i>{/if}
                                    <a href="{$item['url']}">{$item['name']}</a>
                                </li>
                            {else}
                                <li class="active">{$item['name']}</li>
                            {/if}
                        {/foreach}
                    </ul><!-- .breadcrumb -->
                </div>
            {/if}
            <div class="page-content">
                {if $mainTpl}
                    {include $mainTpl}
                {/if}
            </div><!-- /.page-content -->
        </div><!-- /.main-content -->
    </div><!-- /.main-container-inner -->

    <a href="#" id="btn-scroll-up" class="btn-scroll-up btn btn-sm btn-inverse">
        <i class="ace-icon fa fa-angle-double-up icon-only bigger-110"></i>
    </a>
</div><!-- /.main-container -->

<script type="text/javascript">
    //if("ontouchend" in document) document.write("<script src='/static/plugin/ace-admin/assets/js/jquery.mobile.custom.min.js'>"+"<"+"/script>");
</script>
<script src="/static/plugin/ace-admin/assets/js/bootstrap.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/typeahead-bs2.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/jquery-ui-1.10.3.custom.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/jquery.ui.touch-punch.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/chosen.jquery.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/fuelux/fuelux.spinner.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/date-time/bootstrap-datepicker.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/date-time/bootstrap-timepicker.min.js"></script>

<script src="/static/plugin/ace-admin/assets/js/date-time/moment.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/date-time/daterangepicker.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/date-time/bootstrap-datetimepicker.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/bootstrap-colorpicker.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/jquery.knob.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/jquery.autosize.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/jquery.inputlimiter.1.3.1.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/jquery.maskedinput.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/bootstrap-tag.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/dropzone.min.js"></script>

<script src="/static/plugin/ace-admin/assets/js/jquery.slimscroll.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/jquery.easy-pie-chart.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/jquery.sparkline.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/flot/jquery.flot.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/flot/jquery.flot.pie.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/flot/jquery.flot.resize.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/jquery.gritter.min.js"></script>

<script src="/static/plugin/ace-admin/assets/js/fuelux/fuelux.wizard.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/jquery.form.js"></script>
<script src="/static/plugin/ace-admin/assets/js/jquery.validate.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/additional-methods.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/bootbox.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/jquery.maskedinput.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/select2.min.js"></script>

<script src="/static/plugin/ace-admin/assets/js/x-editable/bootstrap-editable.min.js"></script>
<script src="/static/plugin/ace-admin/assets/js/x-editable/ace-editable.min.js"></script>

<script src="{$sysurl}/assets/js/common.js"></script>


<script type="text/javascript">
    jQuery(function($) {
        {* 提示插件 *}
        {if $_SESSION['message']}
            {foreach from=$_SESSION['message'] index=$index key=$key item=$item}
                {if $item['type'] == 'success'}
                    $.gritter.add({
                        title: '成功提示',
                        text: '{$item['message']}',
                        class_name: 'gritter-success {if $item['align'] == 'center'}gritter-center{/if}'
                    });
                {elseif $item['type'] == 'error'}
                    $.gritter.add({
                        title: '错误提示',
                        text: '{$item['message']}',
                        class_name: 'gritter-error {if $item['align'] == 'center'}gritter-center{/if}'
                    });
                {elseif $item['type'] == 'info'}
                    $.gritter.add({
                        title: '消息提示',
                        text: '{$item['message']}',
                        class_name: 'gritter-info {if $item['align'] == 'center'}gritter-center{/if}'
                    });
                {/if}
            {/foreach}
        {/if}
    });
</script>
</body>
</html>
