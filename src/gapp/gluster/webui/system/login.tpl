<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <title>{$config['Sites']['admin']['name']}</title>
    <meta name="description" content="{$config['Site']['admin']['name']}" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <link rel="shortcut icon" href="/static/resource/images/favicon.ico" type="image/x-icon" />

    <!-- basic styles -->
    <link href="/static/plugin/ace-admin/assets/css/bootstrap.min.css" rel="stylesheet" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/font-awesome.min.css" />
    <!--[if IE 7]>
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/font-awesome-ie7.min.css" />
    <![endif]-->
    <!-- page specific plugin styles -->
    <!-- ace styles -->
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/ace.min.css" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/ace-rtl.min.css" />
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/login.css" />
    <!--[if lte IE 8]>
    <link rel="stylesheet" href="/static/plugin/ace-admin/assets/css/ace-ie.min.css" />
    <![endif]-->
    <!-- inline styles related to this page -->
    <!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!--[if lt IE 9]>
    <script src="/static/plugin/ace-admin/assets/js/html5shiv.js"></script>
    <script src="/static/plugin/ace-admin/assets/js/respond.min.js"></script>
    <![endif]-->
    <style type="text/css">
        .help-block {color:#DD5A43;}
        #LoginInfo {margin:0;}
    </style>
</head>

<body class="login-layout theme-clean page-signin" id="animate-area">

<div class="main-container">
    <div class="main-content">
        <div class="row">
            <div class="col-sm-10 col-sm-offset-1">
                <div class="login-container">


                    <div class="position-relative" style="margin-top:20%;">

                        <div class="signin-form">
                            <!-- Form -->
                            <form action="/login/dologin" method="post" id="validation-form">
                                <div class="signin-text">
                                    <span class="ng-binding">请输入您的登录信息</span>
                                </div> <!-- / .signin-text -->

                                {if $_SESSION['message']}
                                    {foreach from=$_SESSION['message'] index=$index key=$key item=$item}
                                        <div class="alert alert-danger" id="LoginInfo" style="margin-bottom:10px;">
                                            <button data-dismiss="alert" class="close" type="button" onclick="$('#LoginInfo').hide()">
                                                <i class="icon-remove"></i>
                                            </button>
                                            <strong>
                                                <i class="icon-remove"></i>
                                                登录失败
                                            </strong>
                                            {$item['message']}
                                            <br>
                                        </div>
                                    {/foreach}
                                {/if}


                                <div class="form-group w-icon">
                                    <input type="text" name="passport" class="form-control input-lg ng-pristine ng-valid ng-touched" placeholder="账号" autocomplete="off" autofocus="">
                                    <span class="ace-icon fa fa-user signin-form-icon"></span>
                                </div> <!-- / Username -->

                                <div class="form-group w-icon">
                                    <input type="password" name="password" class="form-control input-lg ng-pristine ng-valid ng-touched" placeholder="密码" autocomplete="off">
                                    <span class="ace-icon fa fa-lock signin-form-icon"></span>
                                </div> <!-- / Password -->

                                <div class="form-actions">
                                    <button type="submit" class="width-35  btn btn-sm btn-primary" style="height:45px;">
                                        <i class="ace-icon fa fa-key"></i>
                                        登录
                                    </button>
                                </div> <!-- / .form-actions -->
                            </form>
                            <!-- / Form -->
                        </div>



                    </div><!-- /position-relative -->
                </div>
            </div><!-- /.col -->
        </div><!-- /.row -->
    </div>
</div><!-- /.main-container -->

<!-- basic scripts -->

<!--[if !IE]> -->
<script type="text/javascript">
    window.jQuery || document.write("<script src='/static/plugin/ace-admin/assets/js/jquery-2.0.3.min.js'>"+"<"+"/script>");
</script>
<!-- <![endif]-->

<!--[if IE]>
<script type="text/javascript">
    window.jQuery || document.write("<script src='/static/plugin/ace-admin/assets/js/jquery-1.10.2.min.js'>"+"<"+"/script>");
</script>
<![endif]-->

<script type="text/javascript">
    if("ontouchend" in document) document.write("<script src='/static/plugin/ace-admin/assets/js/jquery.mobile.custom.min.js'>"+"<"+"/script>");
</script>

<script src="/static/plugin/ace-admin/assets/js/jquery.validate.min.js"></script>
<script src="/static/resource/js/md5-min.js"></script>

<!-- inline scripts related to this page -->

<script type="text/javascript">
    function show_box(id) {
        jQuery('.widget-box.visible').removeClass('visible');
        jQuery('#'+id).addClass('visible');
    }

    jQuery(function($) {
        $('#validation-form').validate({
            errorElement: 'div',
            errorClass: 'help-block',
            focusInvalid: true,
            rules: {
                passport: {
                    required: true
                },
                password: {
                    required: true
                },
            },
            messages: {
                passport: {
                    required: "请输入账号"
                },
                password: {
                    required: "请输入密码"
                },

            },
            submitHandler: function(form) {
                $("input[name='password']").val(hex_md5($("input[name='password']").val()));
                $(form).ajaxSubmit();
            },
            highlight: function (e) {
                $(e).closest('.form-group').removeClass('has-info').addClass('has-error');
            },

            success: function (e) {
                $(e).closest('.form-group').removeClass('has-error').addClass('has-info');
                $(e).remove();
            },

            errorPlacement: function (error, element) {
                if(element.is(':checkbox') || element.is(':radio')) {
                    var controls = element.closest('div[class*="col-"]');
                    if(controls.find(':checkbox,:radio').length > 1) {
                        controls.append(error);
                    } else {
                        error.insertAfter(element.nextAll('.lbl:eq(0)').eq(0));
                    }
                } else if(element.is('.select2')) {
                    error.insertAfter(element.siblings('[class*="select2-container"]:eq(0)'));
                } else if(element.is('.chosen-select')) {
                    error.insertAfter(element.siblings('[class*="chosen-container"]:eq(0)'));
                } else {
                    error.insertAfter(element.parent());
                }
            }
        });
    });
</script>
</body>
</html>
