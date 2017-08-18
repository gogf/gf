<?php

// Make the page validate
ini_set('session.use_trans_sid', '0');

// Include the random string file
require 'rand.php';

// Begin the session
session_start();

// Set the session contents
$_SESSION['captcha_id'] = $str;

?>
<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>AJAX CAPTCHA</title>
	<script src="../../lib/jquery.js"></script>
	<script src="../../dist/jquery.validate.js"></script>
	<script src="captcha.js"></script>
	<link rel="stylesheet" href="style.css">
	<style>
	img {
		border: 1px solid #eee;
	}
	p#statusgreen {
		font-size: 1.2em;
		background-color: #fff;
		color: #0a0;
	}
	p#statusred {
		font-size: 1.2em;
		background-color: #fff;
		color: #a00;
	}
	fieldset label {
		display: block;
	}
	fieldset div#captchaimage {
		float: left;
		margin-right: 15px;
	}
	fieldset input#captcha {
		width: 25%;
		border: 1px solid #ddd;
		padding: 2px;
	}
	fieldset input#submit {
		display: block;
		margin: 2% 0% 0% 0%;
	}
	#captcha.success {
		border: 1px solid #49c24f;
		background: #bcffbf;
	}
	#captcha.error {
		border: 1px solid #c24949;
		background: #ffbcbc;
	}
	</style>
</head>
<body>
<h1><acronym title="Asynchronous JavaScript And XML">AJAX</acronym> <acronym title="Completely Automated Public Turing test to tell Computers and Humans Apart">CAPTCHA</acronym></h1>
<form id="captchaform" action="">
	<fieldset>
		<div id="captchaimage"><a href="<?php echo htmlEntities($_SERVER['PHP_SELF'], ENT_QUOTES); ?>" id="refreshimg" title="Click to refresh image"><img src="images/image.php?<?php echo time(); ?>" width="132" height="46" alt="Captcha image"></a></div>
		<label for="captcha">Enter the characters as seen on the image above (case insensitive):</label>
		<input type="text" maxlength="6" name="captcha" id="captcha">
		<input type="submit" name="submit" id="submit" value="Check">
	</fieldset>
</form>
<p>If you can&#39;t decipher the text on the image, click it to dynamically generate a new one.</p>
</body>
</html>
