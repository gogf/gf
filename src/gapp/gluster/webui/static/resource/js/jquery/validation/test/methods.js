(function($) {

function methodTest( methodName ) {
	var v = jQuery("#form").validate(),
		method = $.validator.methods[methodName],
		element = $("#firstname")[0];

	return function(value, param) {
		element.value = value;
		return method.call( v, value, element, param );
	};
}

module("methods");

test("default messages", function() {
	var m = $.validator.methods;
	$.each(m, function(key) {
		ok( jQuery.validator.messages[key], key + " has a default message." );
	});
});

test("digit", function() {
	var method = methodTest("digits");
	ok( method( "123" ), "Valid digits" );
	ok(!method( "123.000" ), "Invalid digits" );
	ok(!method( "123.000,00" ), "Invalid digits" );
	ok(!method( "123.0.0,0" ), "Invalid digits" );
	ok(!method( "x123" ), "Invalid digits" );
	ok(!method( "100.100,0,0" ), "Invalid digits" );
});

test("url", function() {
	var method = methodTest("url");
	ok( method( "http://bassistance.de/jquery/plugin.php?bla=blu" ), "Valid url" );
	ok( method( "https://bassistance.de/jquery/plugin.php?bla=blu" ), "Valid url" );
	ok( method( "ftp://bassistance.de/jquery/plugin.php?bla=blu" ), "Valid url" );
	ok( method( "http://www.føtex.dk/" ), "Valid url, danish unicode characters" );
	ok( method( "http://bösendorfer.de/" ), "Valid url, german unicode characters" );
	ok( method( "http://142.42.1.1" ), "Valid IP Address" );
	ok( method( "http://pro.photography" ), "Valid long TLD" );
	ok( method( "//code.jquery.com/jquery-1.11.3.min.js" ), "Valid protocol-relative url" );
	ok( method( "//142.42.1.1" ), "Valid protocol-relative IP Address" );
	ok(!method( "htp://code.jquery.com/jquery-1.11.3.min.js" ), "Invalid protocol" );
	ok(!method( "http://192.168.8." ), "Invalid IP Address" );
	ok(!method( "http://bassistance" ), "Invalid url" ); // valid
	ok(!method( "http://bassistance." ), "Invalid url" ); // valid
	ok(!method( "http://bassistance,de" ), "Invalid url" );
	ok(!method( "http://bassistance;de" ), "Invalid url" );
	ok(!method( "http://.bassistancede" ), "Invalid url" );
	ok(!method( "bassistance.de" ), "Invalid url" );
});

test("url2 (tld optional)", function() {
	var method = methodTest("url2");
	ok( method( "http://bassistance.de/jquery/plugin.php?bla=blu" ), "Valid url" );
	ok( method( "https://bassistance.de/jquery/plugin.php?bla=blu" ), "Valid url" );
	ok( method( "ftp://bassistance.de/jquery/plugin.php?bla=blu" ), "Valid url" );
	ok( method( "http://www.føtex.dk/" ), "Valid url, danish unicode characters" );
	ok( method( "http://bösendorfer.de/" ), "Valid url, german unicode characters" );
	ok( method( "http://192.168.8.5" ), "Valid IP Address" );
	ok(!method( "http://192.168.8." ), "Invalid IP Address" );
	ok( method( "http://bassistance" ), "Invalid url" );
	ok( method( "http://bassistance." ), "Invalid url" );
	ok(!method( "http://bassistance,de" ), "Invalid url" );
	ok(!method( "http://bassistance;de" ), "Invalid url" );
	ok(!method( "http://.bassistancede" ), "Invalid url" );
	ok(!method( "bassistance.de" ), "Invalid url" );
});

test("email", function() {
	var method = methodTest("email");
	ok( method( "name@domain.tld" ), "Valid email" );
	ok( method( "name@domain.tl" ), "Valid email" );
	ok( method( "bart+bart@tokbox.com" ), "Valid email" );
	ok( method( "bart+bart@tokbox.travel" ), "Valid email" );
	ok( method( "n@d.tld" ), "Valid email" );
	ok( method( "bla.blu@g.mail.com"), "Valid email" );
	ok( method( "name@domain" ), "Valid email" );
	ok( method( "name.@domain.tld" ), "Valid email" );
	ok( method( "name@website.a" ), "Valid email" );
	ok( method( "name@pro.photography" ), "Valid email" );
	ok(!method( "ole@føtex.dk"), "Invalid email" );
	ok(!method( "jörn@bassistance.de"), "Invalid email" );
	ok(!method( "name" ), "Invalid email" );
	ok(!method( "test@test-.com" ), "Invalid email" );
	ok(!method( "name@" ), "Invalid email" );
	ok(!method( "name,@domain.tld" ), "Invalid email" );
	ok(!method( "name;@domain.tld" ), "Invalid email" );
	ok(!method( "name;@domain.tld." ), "Invalid email" );
});

test("number", function() {
	var method = methodTest("number");
	ok( method( "123" ), "Valid number" );
	ok( method( "-123" ), "Valid number" );
	ok( method( "123,000" ), "Valid number" );
	ok( method( "-123,000" ), "Valid number" );
	ok( method( "123,000.00" ), "Valid number" );
	ok( method( "-123,000.00" ), "Valid number" );
	ok(!method( "-" ), "Invalid number" );
	ok(!method( "123.000,00" ), "Invalid number" );
	ok(!method( "123.0.0,0" ), "Invalid number" );
	ok(!method( "x123" ), "Invalid number" );
	ok(!method( "100.100,0,0" ), "Invalid number" );

	ok( method( "" ), "Blank is valid" );
	ok( method( "123" ), "Valid decimal" );
	ok( method( "123000" ), "Valid decimal" );
	ok( method( "123000.12" ), "Valid decimal" );
	ok( method( "-123000.12" ), "Valid decimal" );
	ok( method( "123.000" ), "Valid decimal" );
	ok( method( "123,000.00" ), "Valid decimal" );
	ok( method( "-123,000.00" ), "Valid decimal" );
	ok( method( ".100" ), "Valid decimal" );
	ok(!method( "1230,000.00" ), "Invalid decimal" );
	ok(!method( "123.0.0,0" ), "Invalid decimal" );
	ok(!method( "x123" ), "Invalid decimal" );
	ok(!method( "100.100,0,0" ), "Invalid decimal" );
});

/* disabled for now, need to figure out how to test localized methods
test("numberDE", function() {
	var method = methodTest("numberDE");
	ok( method( "123" ), "Valid numberDE" );
	ok( method( "-123" ), "Valid numberDE" );
	ok( method( "123.000" ), "Valid numberDE" );
	ok( method( "-123.000" ), "Valid numberDE" );
	ok( method( "123.000,00" ), "Valid numberDE" );
	ok( method( "-123.000,00" ), "Valid numberDE" );
	ok(!method( "123,000.00" ), "Invalid numberDE" );
	ok(!method( "123,0,0.0" ), "Invalid numberDE" );
	ok(!method( "x123" ), "Invalid numberDE" );
	ok(!method( "100,100.0.0" ), "Invalid numberDE" );

	ok( method( "" ), "Blank is valid" );
	ok( method( "123" ), "Valid decimalDE" );
	ok( method( "123000" ), "Valid decimalDE" );
	ok( method( "123000,12" ), "Valid decimalDE" );
	ok( method( "-123000,12" ), "Valid decimalDE" );
	ok( method( "123.000" ), "Valid decimalDE" );
	ok( method( "123.000,00" ), "Valid decimalDE" );
	ok( method( "-123.000,00" ), "Valid decimalDE" )
	ok(!method( "123.0.0,0" ), "Invalid decimalDE" );
	ok(!method( "x123" ), "Invalid decimalDE" );
	ok(!method( "100,100.0.0" ), "Invalid decimalDE" );
});
*/

test("date", function() {
	var method = methodTest("date");
	ok( method( "06/06/1990" ), "Valid date" );
	ok( method( "6/6/06" ), "Valid date" );
	ok(!method( "1990x-06-06" ), "Invalid date" );
});

test("dateISO", function() {
	var method = methodTest("dateISO");
	ok( method( "1990-06-06" ), "Valid date" );
	ok( method( "1990-01-01" ), "Valid date" );
	ok( method( "1990-01-31" ), "Valid date" );
	ok( method( "1990-12-01" ), "Valid date" );
	ok( method( "1990-12-31" ), "Valid date" );
	ok( method( "1990/06/06" ), "Valid date" );
	ok( method( "1990-6-6" ), "Valid date" );
	ok( method( "1990/6/6" ), "Valid date" );
	ok(!method( "1990-106-06" ), "Invalid date" );
	ok(!method( "190-06-06" ), "Invalid date" );
	ok(!method( "1990-00-06" ), "Invalid date" );
	ok(!method( "1990-13-01" ), "Invalid date" );
	ok(!method( "1990-01-00" ), "Invalid date" );
	ok(!method( "1990-01-32" ), "Invalid date" );
	ok(!method( "1990-13-32" ), "Invalid date" );
});

/* disabled for now, need to figure out how to test localized methods
test("dateDE", function() {
	var method = methodTest("dateDE");
	ok( method( "03.06.1984" ), "Valid dateDE" );
	ok( method( "3.6.84" ), "Valid dateDE" );
	ok(!method( "6-6-06" ), "Invalid dateDE" );
	ok(!method( "1990-06-06" ), "Invalid dateDE" );
	ok(!method( "06/06/1990" ), "Invalid dateDE" );
	ok(!method( "6/6/06" ), "Invalid dateDE" );
});
*/

test("required", function() {
	var v = jQuery("#form").validate(),
		method = $.validator.methods.required,
		e = $("#text1, #text1b, #hidden2, #select1, #select2");
	ok( method.call( v, e[0].value, e[0]), "Valid text input" );
	ok(!method.call( v, e[1].value, e[1]), "Invalid text input" );
	ok(!method.call( v, e[2].value, e[2]), "Invalid text input" );

	ok(!method.call( v, e[3].value, e[3]), "Invalid select" );
	ok( method.call( v, e[4].value, e[4]), "Valid select" );

	e = $("#area1, #area2, #pw1, #pw2");
	ok( method.call( v, e[0].value, e[0]), "Valid textarea" );
	ok(!method.call( v, e[1].value, e[1]), "Invalid textarea" );
	ok( method.call( v, e[2].value, e[2]), "Valid password input" );
	ok(!method.call( v, e[3].value, e[3]), "Invalid password input" );

	e = $("#radio1, #radio2, #radio3");
	ok(!method.call( v, e[0].value, e[0]), "Invalid radio" );
	ok( method.call( v, e[1].value, e[1]), "Valid radio" );
	ok( method.call( v, e[2].value, e[2]), "Valid radio" );

	e = $("#check1, #check2");
	ok( method.call( v, e[0].value, e[0]), "Valid checkbox" );
	ok(!method.call( v, e[1].value, e[1]), "Invalid checkbox" );

	e = $("#select1, #select2, #select3, #select4");
	ok(!method.call( v, e[0].value, e[0]), "Invalid select" );
	ok( method.call( v, e[1].value, e[1]), "Valid select" );
	ok( method.call( v, e[2].value, e[2]), "Valid select" );
	ok( method.call( v, e[3].value, e[3]), "Valid select" );
});

test("required with dependencies", function() {
	var v = jQuery("#form").validate(),
		method = $.validator.methods.required,
		e = $("#hidden2, #select1, #area2, #radio1, #check2");
	ok( method.call( v, e[0].value, e[0], "asffsaa" ), "Valid text input due to dependency not met" );
	ok(!method.call( v, e[0].value, e[0], "input" ), "Invalid text input" );
	ok( method.call( v, e[0].value, e[0], function() { return false; }), "Valid text input due to dependency not met" );
	ok(!method.call( v, e[0].value, e[0], function() { return true; }), "Invalid text input" );
	ok( method.call( v, e[1].value, e[1], "asfsfa" ), "Valid select due to dependency not met" );
	ok(!method.call( v, e[1].value, e[1], "input" ), "Invalid select" );
	ok( method.call( v, e[2].value, e[2], "asfsafsfa" ), "Valid textarea due to dependency not met" );
	ok(!method.call( v, e[2].value, e[2], "input" ), "Invalid textarea" );
	ok( method.call( v, e[3].value, e[3], "asfsafsfa" ), "Valid radio due to dependency not met" );
	ok(!method.call( v, e[3].value, e[3], "input" ), "Invalid radio" );
	ok( method.call( v, e[4].value, e[4], "asfsafsfa" ), "Valid checkbox due to dependency not met" );
	ok(!method.call( v, e[4].value, e[4], "input" ), "Invalid checkbox" );
});

test("minlength", function() {
	var v = jQuery("#form").validate(),
		method = $.validator.methods.minlength,
		param = 2,
		e = $("#text1, #text1c, #text2, #text3");
	ok( method.call( v, e[0].value, e[0], param), "Valid text input" );
	ok( method.call( v, e[1].value, e[1], param), "Valid text input" );
	ok(!method.call( v, e[2].value, e[2], param), "Invalid text input" );
	ok( method.call( v, e[3].value, e[3], param), "Valid text input" );

	e = $("#check1, #check2, #check3");
	ok(!method.call( v, e[0].value, e[0], param), "Valid checkbox" );
	ok( method.call( v, e[1].value, e[1], param), "Valid checkbox" );
	ok( method.call( v, e[2].value, e[2], param), "Invalid checkbox" );

	e = $("#select1, #select2, #select3, #select4, #select5");
	ok(method.call( v, e[0].value, e[0], param), "Valid select " + e[0].id );
	ok(!method.call( v, e[1].value, e[1], param), "Invalid select " + e[1].id );
	ok( method.call( v, e[2].value, e[2], param), "Valid select " + e[2].id );
	ok( method.call( v, e[3].value, e[3], param), "Valid select " + e[3].id );
	ok( method.call( v, e[4].value, e[4], param), "Valid select " + e[4].id );
});

test("maxlength", function() {
	var v = jQuery("#form").validate(),
		method = $.validator.methods.maxlength,
		param = 4,
		e = $("#text1, #text2, #text3");

	ok( method.call( v, e[0].value, e[0], param), "Valid text input" );
	ok( method.call( v, e[1].value, e[1], param), "Valid text input" );
	ok(!method.call( v, e[2].value, e[2], param), "Invalid text input" );

	e = $("#check1, #check2, #check3");
	ok( method.call( v, e[0].value, e[0], param), "Valid checkbox" );
	ok( method.call( v, e[1].value, e[1], param), "Invalid checkbox" );
	ok(!method.call( v, e[2].value, e[2], param), "Invalid checkbox" );

	e = $("#select1, #select2, #select3, #select4");
	ok( method.call( v, e[0].value, e[0], param), "Valid select" );
	ok( method.call( v, e[1].value, e[1], param), "Valid select" );
	ok( method.call( v, e[2].value, e[2], param), "Valid select" );
	ok(!method.call( v, e[3].value, e[3], param), "Invalid select" );
});

test("rangelength", function() {
	var v = jQuery("#form").validate(),
		method = $.validator.methods.rangelength,
		param = [ 2, 4 ],
		e = $("#text1, #text2, #text3");

	ok( method.call( v, e[0].value, e[0], param), "Valid text input" );
	ok(!method.call( v, e[1].value, e[1], param), "Invalid text input" );
	ok(!method.call( v, e[2].value, e[2], param), "Invalid text input" );
});

test("min", function() {
	var v = jQuery("#form").validate(),
		method = $.validator.methods.min,
		param = 8,
		e = $("#value1, #value2, #value3");

	ok(!method.call( v, e[0].value, e[0], param), "Invalid text input" );
	ok( method.call( v, e[1].value, e[1], param), "Valid text input" );
	ok( method.call( v, e[2].value, e[2], param), "Valid text input" );
});

test("max", function() {
	var v = jQuery("#form").validate(),
		method = $.validator.methods.max,
		param = 12,
		e = $("#value1, #value2, #value3");

	ok( method.call( v, e[0].value, e[0], param), "Valid text input" );
	ok( method.call( v, e[1].value, e[1], param), "Valid text input" );
	ok(!method.call( v, e[2].value, e[2], param), "Invalid text input" );
});

test("range", function() {
	var v = jQuery("#form").validate(),
		method = $.validator.methods.range,
		param = [ 4, 12 ],
		e = $("#value1, #value2, #value3");

	ok(!method.call( v, e[0].value, e[0], param), "Invalid text input" );
	ok( method.call( v, e[1].value, e[1], param), "Valid text input" );
	ok(!method.call( v, e[2].value, e[2], param), "Invalid text input" );
});

test("equalTo", function() {
	var v = jQuery("#form").validate(),
		method = $.validator.methods.equalTo,
		e = $("#text1, #text2");

	ok( method.call( v, "Test", e[0], "#text1" ), "Text input" );
	ok( method.call( v, "T", e[1], "#text2" ), "Another one" );
});

test("creditcard", function() {
	var method = methodTest("creditcard");
	ok( method( "4111-1111-1111-1111" ), "Valid creditcard number" );
	ok( method( "4111 1111 1111 1111" ), "Valid creditcard number" );
	ok(!method( "41111" ), "Invalid creditcard number" );
	ok(!method( "asdf" ), "Invalid creditcard number" );
});

test("extension", function() {
	var method = methodTest("extension"),
		v;
	ok( method( "picture.gif" ), "Valid default accept type" );
	ok( method( "picture.jpg" ), "Valid default accept type" );
	ok( method( "picture.jpeg" ), "Valid default accept type" );
	ok( method( "picture.png" ), "Valid default accept type" );
	ok(!method( "picture.pgn" ), "Invalid default accept type" );

	v = jQuery("#form").validate();
	method = function(value, param) {
		return $.validator.methods.extension.call(v, value, $("#text1")[0], param);
	};
	ok( method( "picture.doc", "doc" ), "Valid custom accept type" );
	ok( method( "picture.pdf", "doc|pdf" ), "Valid custom accept type" );
	ok( method( "picture.pdf", "pdf|doc" ), "Valid custom accept type" );
	ok(!method( "picture.pdf", "doc" ), "Invalid custom accept type" );
	ok(!method( "picture.doc", "pdf" ), "Invalid custom accept type" );

	ok( method( "picture.pdf", "doc,pdf" ), "Valid custom accept type, comma separated" );
	ok( method( "picture.pdf", "pdf,doc" ), "Valid custom accept type, comma separated" );
	ok(!method( "picture.pdf", "gop,top" ), "Invalid custom accept type, comma separated" );
});

asyncTest("remote", function() {
	expect(7);
	var e = $("#username"),
		v = $("#userForm").validate({
			rules: {
				username: {
					required: true,
					remote: "users.php"
				}
			},
			messages: {
				username: {
					required: "Please",
					remote: jQuery.validator.format("{0} in use")
				}
			},
			submitHandler: function() {
				ok( false, "submitHandler may never be called when validating only elements");
			}
		});

	$(document).ajaxStop(function() {
		$(document).unbind("ajaxStop");
		equal( 1, v.size(), "There must be one error" );
		equal( "Peter in use", v.errorList[0].message );

		$(document).ajaxStop(function() {
			$(document).unbind("ajaxStop");
			equal( 1, v.size(), "There must be one error" );
			equal( "Peter2 in use", v.errorList[0].message );
			start();
		});
		e.val("Peter2");
		strictEqual( v.element(e), true, "new value, new request; dependency-mismatch considered as valid though" );
	});
	strictEqual( v.element(e), false, "invalid element, nothing entered yet" );
	e.val("Peter");
	strictEqual( v.element(e), true, "still invalid, because remote validation must block until it returns; dependency-mismatch considered as valid though" );
});

asyncTest("remote, customized ajax options", function() {
	expect(2);
	$("#userForm").validate({
		rules: {
			username: {
				required: true,
				remote: {
					url: "users.php",
					type: "POST",
					beforeSend: function(request, settings) {
						deepEqual(settings.type, "POST");
						deepEqual(settings.data, "username=asdf&email=email.com");
					},
					data: {
						email: function() {
							return "email.com";
						}
					},
					complete: function() {
						start();
					}
				}
			}
		}
	});
	$("#username").val("asdf");
	$("#userForm").valid();
});

asyncTest("remote extensions", function() {
	expect(5);
	var e = $("#username"),
		v = $("#userForm").validate({
			rules: {
				username: {
					required: true,
					remote: "users2.php"
				}
			},
			messages: {
				username: {
					required: "Please"
				}
			},
			submitHandler: function() {
				ok( false, "submitHandler may never be called when validating only elements");
			}
		});

	$(document).ajaxStop(function() {
		$(document).unbind("ajaxStop");
		if ( v.size() !== 0 ) {
			ok( "There must be one error" );
			equal( v.errorList[0].message, "asdf is already taken, please try something else" );
			v.element(e);
			equal( v.errorList[0].message, "asdf is already taken, please try something else", "message doesn't change on revalidation" );
		}
		start();
	});
	strictEqual( v.element(e), false, "invalid element, nothing entered yet" );
	e.val("asdf");
	strictEqual( v.element(e), true, "still invalid, because remote validation must block until it returns; dependency-mismatch considered as valid though" );
});

module("additional methods");

test("phone (us)", function() {
	var method = methodTest("phoneUS");
	ok( method( "1(212)-999-2345" ), "Valid US phone number" );
	ok( method( "212 999 2344" ), "Valid US phone number" );
	ok( method( "212-999-0983" ), "Valid US phone number" );
	ok(!method( "111-123-5434" ), "Invalid US phone number. Area Code cannot start with 1" );
	ok(!method( "212 123 4567" ), "Invalid US phone number. NXX cannot start with 1" );
	ok(!method( "234-911-5678" ), "Invalid US phone number, because the exchange code cannot be in the form N11" );
	ok(!method( "911-333-5678" ), "Invalid US phone number, because the area code cannot be in the form N11" );
	ok(method( "234-912-5678" ), "Valid US phone number" );
});

test("phoneUK", function() {
	var method = methodTest("phoneUK");
	ok( method( "0117 333 5555" ), "Valid UK Phone Number" );
	ok( method( "0121 555 5555" ), "Valid UK Phone Number" );
	ok( method( "01633 555555" ), "Valid UK Phone Number" );
	ok( method( "01298 28555" ), "Valid UK Phone Number" );
	ok( method( "015395 55555" ), "Valid UK Phone Number" );
	ok( method( "016977 3999" ), "Valid UK Phone Number" );
	ok( method( "020 3000 5555" ), "Valid UK Phone Number" );
	ok( method( "024 7500 5555" ), "Valid UK Phone Number" );
	ok( method( "0333 555 5555" ), "Valid UK Phone Number" );
	ok( method( "0500 555555" ), "Valid UK Phone Number" );
	ok( method( "055 3555 5555" ), "Valid UK Phone Number" );
	ok( method( "07122 555555" ), "Valid UK Phone Number" );
	ok( method( "07222 555555" ), "Valid UK Phone Number" );
	ok( method( "07322 555555" ), "Valid UK Phone Number" );
	ok( method( "0800 555 5555" ), "Valid UK Phone Number" );
	ok( method( "0800 355555" ), "Valid UK Phone Number" );
	ok( method( "0843 555 5555" ), "Valid UK Phone Number" );
	ok( method( "0872 555 5555" ), "Valid UK Phone Number" );
	ok( method( "0903 555 5555" ), "Valid UK Phone Number" );
	ok( method( "0983 555 5555" ), "Valid UK Phone Number" );
	ok( method( "(07122) 555555" ), "Valid UK Phone Number" );
	ok( method( "(07222) 555555" ), "Valid UK Phone Number" );
	ok( method( "(07322) 555555" ), "Valid UK Phone Number" );
	ok( method( "+44 7122 555 555" ), "Valid UK Phone Number" );
	ok( method( "+44 7222 555 555" ), "Valid UK Phone Number" );
	ok( method( "+44 7322 555 555" ), "Valid UK Phone Number" );
	ok(!method( "7222 555555" ), "Invalid UK Phone Number" );
	ok(!method( "+44 07222 555555" ), "Invalid UK Phone Number" );
});

test("mobileUK", function() {
	var method = methodTest("mobileUK");
	ok( method( "07134234323" ), "Valid UK Mobile Number" );
	ok( method( "07334234323" ), "Valid UK Mobile Number" );
	ok( method( "07624234323" ), "Valid UK Mobile Number" );
	ok( method( "07734234323" ), "Valid UK Mobile Number" );
	ok( method( "+447134234323" ), "Valid UK Mobile Number" );
	ok( method( "+447334234323" ), "Valid UK Mobile Number" );
	ok( method( "+447624234323" ), "Valid UK Mobile Number" );
	ok( method( "+447734234323" ), "Valid UK Mobile Number" );
	ok(!method( "07034234323" ), "Invalid UK Mobile Number" );
	ok(!method( "0753423432" ), "Invalid UK Mobile Number" );
	ok(!method( "07604234323" ), "Invalid UK Mobile Number" );
	ok(!method( "077342343234" ), "Invalid UK Mobile Number" );
	ok(!method( "044342343234" ), "Invalid UK Mobile Number" );
	ok(!method( "+44753423432" ), "Invalid UK Mobile Number" );
	ok(!method( "+447604234323" ), "Invalid UK Mobile Number" );
	ok(!method( "+4477342343234" ), "Invalid UK Mobile Number" );
	ok(!method( "+4444342343234" ), "Invalid UK Mobile Number" );
});

test("dateITA", function() {
	var method = methodTest("dateITA");
	ok( method( "01/01/1900" ), "Valid date ITA" );
	ok( method( "17/10/2010" ), "Valid date ITA" );
	ok(!method( "01/13/1990" ), "Invalid date ITA" );
	ok(!method( "01.01.1900" ), "Invalid date ITA" );
	ok(!method( "01/01/199" ), "Invalid date ITA" );
});

test("dateFA", function() {
	var method = methodTest("dateFA");

	ok( method( "1342/12/29" ), "Valid date FA" );
	ok( method( "1342/12/30" ), "Valid date FA" );
	ok( method( "1361/6/31" ), "Valid date FA" );
	ok( method( "1321/11/30" ), "Valid date FA" );
	ok( method( "1361/1/1" ), "Valid date FA" );
	ok( method( "1020/3/3" ), "Valid date FA" );
	ok( method( "1020/03/3" ), "Valid date FA" );
	ok( method( "1020/3/03" ), "Valid date FA" );
	ok( method( "1020/03/03" ), "Valid date FA" );
	ok( method( "1001/7/30" ), "Valid date FA" );

	ok(!method( "1000/1/32" ), "Invalid date FA" );
	ok(!method( "1323/12/31" ), "Invalid date FA" );
	ok(!method( "1361/0/11" ), "Invalid date FA" );
	ok(!method( "63/4/4" ), "Invalid date FA" );
	ok(!method( "15/6/1361" ), "Invalid date FA" );
});

test("iban", function() {
	var method = methodTest("iban");
	ok( method( "NL20INGB0001234567"), "Valid IBAN");
	ok( method( "DE68 2105 0170 0012 3456 78"), "Valid IBAN");
	ok( method( "NL20 INGB0001234567"), "Valid IBAN: invalid spacing");
	ok( method( "NL20 INGB 00 0123 4567"), "Valid IBAN: invalid spacing");
	ok( method( "XX40INGB000123456712341234"), "Valid (more or less) IBAN: unknown country, but checksum OK");

	ok(!method( "NL20INGB000123456"), "Invalid IBAN: too short");
	ok(!method( "NL20INGB00012345678"), "Invalid IBAN: too long");
	ok(!method( "NL20INGB0001234566"), "Invalid IBAN: checksum incorrect");
	ok(!method( "DE68 2105 0170 0012 3456 7"), "Invalid IBAN: too short");
	ok(!method( "DE68 2105 0170 0012 3456 789"), "Invalid IBAN: too long");
	ok(!method( "DE68 2105 0170 0012 3456 79"), "Invalid IBAN: checksum incorrect");

	ok(!method( "NL54INGB00012345671234"), "Invalid IBAN too long, BUT CORRECT CHECKSUM");
	ok(!method( "XX00INGB000123456712341234"), "Invalid IBAN: unknown country and checksum incorrect");

	// sample IBANs for different countries
	ok( method( "AL47 2121 1009 0000 0002 3569 8741"), "Valid IBAN - AL");
	ok( method( "AD12 0001 2030 2003 5910 0100"), "Valid IBAN - AD");
	ok( method( "AT61 1904 3002 3457 3201"), "Valid IBAN - AT");
	ok( method( "AZ21 NABZ 0000 0000 1370 1000 1944"), "Valid IBAN - AZ");
	ok( method( "BH67 BMAG 0000 1299 1234 56"), "Valid IBAN - BH");
	ok( method( "BE62 5100 0754 7061"), "Valid IBAN - BE");
	ok( method( "BA39 1290 0794 0102 8494"), "Valid IBAN - BA");
	ok( method( "BG80 BNBG 9661 1020 3456 78"), "Valid IBAN - BG");
	ok( method( "HR12 1001 0051 8630 0016 0"), "Valid IBAN - HR");
	ok( method( "CH93 0076 2011 6238 5295 7"), "Valid IBAN - CH");
	ok( method( "CY17 0020 0128 0000 0012 0052 7600"), "Valid IBAN - CY");
	ok( method( "CZ65 0800 0000 1920 0014 5399"), "Valid IBAN - CZ");
	ok( method( "DK50 0040 0440 1162 43"), "Valid IBAN - DK");
	ok( method( "EE38 2200 2210 2014 5685"), "Valid IBAN - EE");
	ok( method( "FO97 5432 0388 8999 44"), "Valid IBAN - FO");
	ok( method( "FI21 1234 5600 0007 85"), "Valid IBAN - FI");
	ok( method( "FR14 2004 1010 0505 0001 3M02 606"), "Valid IBAN - FR");
	ok( method( "GE29 NB00 0000 0101 9049 17"), "Valid IBAN - GE");
	ok( method( "DE89 3704 0044 0532 0130 00"), "Valid IBAN - DE");
	ok( method( "GI75 NWBK 0000 0000 7099 453"), "Valid IBAN - GI");
	ok( method( "GR16 0110 1250 0000 0001 2300 695"), "Valid IBAN - GR");
	ok( method( "GL56 0444 9876 5432 10"), "Valid IBAN - GL");
	ok( method( "HU42 1177 3016 1111 1018 0000 0000"), "Valid IBAN - HU");
	ok( method( "IS14 0159 2600 7654 5510 7303 39"), "Valid IBAN - IS");
	ok( method( "IE29 AIBK 9311 5212 3456 78"), "Valid IBAN - IE");
	ok( method( "IL62 0108 0000 0009 9999 999"), "Valid IBAN - IL");
	ok( method( "IT40 S054 2811 1010 0000 0123 456"), "Valid IBAN - IT");
	ok( method( "LV80 BANK 0000 4351 9500 1"), "Valid IBAN - LV");
	ok( method( "LB62 0999 0000 0001 0019 0122 9114"), "Valid IBAN - LB");
	ok( method( "LI21 0881 0000 2324 013A A"), "Valid IBAN - LI");
	ok( method( "LT12 1000 0111 0100 1000"), "Valid IBAN - LT");
	ok( method( "LU28 0019 4006 4475 0000"), "Valid IBAN - LU");
	ok( method( "MK07 2501 2000 0058 984"), "Valid IBAN - MK");
	ok( method( "MT84 MALT 0110 0001 2345 MTLC AST0 01S"), "Valid IBAN - MT");
	ok( method( "MU17 BOMM 0101 1010 3030 0200 000M UR"), "Valid IBAN - MU");
	ok( method( "MD24 AG00 0225 1000 1310 4168"), "Valid IBAN - MD");
	ok( method( "MC93 2005 2222 1001 1223 3M44 555"), "Valid IBAN - MC");
	ok( method( "ME25 5050 0001 2345 6789 51"), "Valid IBAN - ME");
	ok( method( "NL39 RABO 0300 0652 64"), "Valid IBAN - NL");
	ok( method( "NO93 8601 1117 947"), "Valid IBAN - NO");
	ok( method( "PK36 SCBL 0000 0011 2345 6702"), "Valid IBAN - PK");
	ok( method( "PL60 1020 1026 0000 0422 7020 1111"), "Valid IBAN - PL");
	ok( method( "PT50 0002 0123 1234 5678 9015 4"), "Valid IBAN - PT");
	ok( method( "RO49 AAAA 1B31 0075 9384 0000"), "Valid IBAN - RO");
	ok( method( "SM86 U032 2509 8000 0000 0270 100"), "Valid IBAN - SM");
	ok( method( "SA03 8000 0000 6080 1016 7519"), "Valid IBAN - SA");
	ok( method( "RS35 2600 0560 1001 6113 79"), "Valid IBAN - RS");
	ok( method( "SK31 1200 0000 1987 4263 7541"), "Valid IBAN - SK");
	ok( method( "SI56 1910 0000 0123 438"), "Valid IBAN - SI");
	ok( method( "ES80 2310 0001 1800 0001 2345"), "Valid IBAN - ES");
	ok( method( "SE35 5000 0000 0549 1000 0003"), "Valid IBAN - SE");
	ok( method( "CH93 0076 2011 6238 5295 7"), "Valid IBAN - CH");
	ok( method( "TN59 1000 6035 1835 9847 8831"), "Valid IBAN - TN");
	ok( method( "TR33 0006 1005 1978 6457 8413 26"), "Valid IBAN - TR");
	ok( method( "AE07 0331 2345 6789 0123 456"), "Valid IBAN - AE");
	ok( method( "GB29 NWBK 6016 1331 9268 19"), "Valid IBAN - GB");
});

/**
 * BIC tests (For BIC definition take a look on the implementation itself)
 */
test("bic", function() {
	var method = methodTest( "bic" );

	ok( !method( "PBNKDEF" ), "Invalid BIC: too short" );
	ok( !method( "DEUTDEFFA1" ), "Invalid BIC: disallowed length" );
	ok( !method( "PBNKDEFFXXX1" ), "Invalid BIC: too long" );
	ok( !method( "1BNKDEFF" ), "Invalid BIC: invalid digit" );
	ok( !method( "PBNKDE1F" ), "Invalid BIC: invalid digit" );
	ok( !method( "PBNKDEF3" ), "Invalid BIC: invalid digit" );
	ok( !method( "PBNKDEFO" ), "Invalid BIC: invalid char" );
	ok( !method( "INGDDEFFXAA" ), "Invalid BIC: invalid char" );
	ok( !method( "DEUTDEF0" ), "Invalid BIC: invalid digit" );

	ok( method( "DEUTDEFF" ), "Valid BIC" );
	ok( method( "DEUTDEFFXXX" ), "Valid BIC" );
	ok( method( "PBNKDE2F" ), "Valid BIC" );
	ok( method( "INGDDEFF101" ), "Valid BIC" );
	ok( method( "INGDDEF2134" ), "Valid BIC" );
	ok( method( "INGDDE91XXX" ), "Valid BIC" );
	ok( method( "INGDDEF2" ), "Valid BIC" );
	ok( method( "AAFFFRP1" ), "Valid BIC" );
	ok( method( "DEUTDEFFAB1" ), "Valid BIC" );
	ok( method( "DEUTDEFFAXX" ), "Valid BIC" );
});

test("postcodeUK", function() {
	var method = methodTest("postcodeUK");
	ok( method( "AA9A 9AA" ), "Valid postcode" );
	ok( method( "A9A 9AA" ), "Valid postcode" );
	ok( method( "A9 9AA" ), "Valid postcode" );
	ok( method( "A99 9AA" ), "Valid postcode" );
	ok( method( "AA9 9AA" ), "Valid postcode" );
	ok( method( "AA99 9AA" ), "Valid postcode" );

	// Channel Island
	ok(!method( "AAAA 9AA" ), "Invalid postcode" );
	ok(!method( "AA-2640" ), "Invalid postcode" );

	ok(!method( "AAA AAA" ), "Invalid postcode" );
	ok(!method( "AA AAAA" ), "Invalid postcode" );
	ok(!method( "A AAAA" ), "Invalid postcode" );
	ok(!method( "AAAAA" ), "Invalid postcode" );
	ok(!method( "999 999" ), "Invalid postcode" );
	ok(!method( "99 9999" ), "Invalid postcode" );
	ok(!method( "9 9999" ), "Invalid postcode" );
	ok(!method( "99999" ), "Invalid postcode" );
});

test("dateNL", function() {
	var method = methodTest("dateNL");
	ok( method( "01-01-1900" ), "Valid date NL" );
	ok( method( "01.01.1900" ), "Valid date NL" );
	ok( method( "01/01/1900" ), "Valid date NL" );
	ok( method( "01-01-00" ), "Valid date NL" );
	ok( method( "1-01-1900" ), "Valid date NL" );
	ok( method( "10-10-1900" ), "Valid date NL" );
	ok(!method( "0-01-1900" ), "Invalid date NL" );
	ok(!method( "00-01-1900" ), "Invalid date NL" );
	ok(!method( "35-01-1990" ), "Invalid date NL" );
	ok(!method( "01.01.190" ), "Invalid date NL" );
});

test("phoneNL", function() {
	var method = methodTest("phoneNL");
	ok( method( "0701234567"), "Valid phone NL");
	ok( method( "0687654321"), "Valid phone NL");
	ok( method( "020-1234567"), "Valid phone NL");
	ok( method( "020 - 12 34 567"), "Valid phone NL");
	ok( method( "010-2345678"), "Valid phone NL");
	ok( method( "+3120-1234567"), "Valid phone NL");
	ok( method( "+31(0)10-2345678"), "Valid phone NL");
	ok(!method( "020-123456"), "Invalid phone NL: too short");
	ok(!method( "020-12345678"), "Invalid phone NL: too long");
	ok(!method( "-0201234567"), "Invalid phone NL");
	ok(!method( "+310201234567"), "Invalid phone NL: no 0 after +31 allowed");
});

test("mobileNL", function() {
	var method = methodTest("mobileNL");
	ok( method( "0612345678"), "Valid NL Mobile Number");
	ok( method( "06-12345678"), "Valid NL Mobile Number");
	ok( method( "06-12 345 678"), "Valid NL Mobile Number");
	ok( method( "+316-12345678"), "Valid NL Mobile Number");
	ok( method( "+31(0)6-12345678"), "Valid NL Mobile Number");
	ok(!method( "abcdefghij"), "Invalid NL Mobile Number: text");
	ok(!method( "0123456789"), "Invalid NL Mobile Number: should start with 06");
	ok(!method( "0823456789"), "Invalid NL Mobile Number: should start with 06");
	ok(!method( "06-1234567"), "Invalid NL Mobile Number: too short");
	ok(!method( "06-123456789"), "Invalid NL Mobile Number: too long");
	ok(!method( "-0612345678"), "Invalid NL Mobile Number");
	ok(!method( "+310612345678"), "Invalid NL Mobile Number: no 0 after +31 allowed");
});

test("postalcodeNL", function() {
	var method = methodTest("postalcodeNL");
	ok( method( "1234AB"), "Valid NL Postal Code");
	ok( method( "1234ab"), "Valid NL Postal Code");
	ok( method( "1234 AB"), "Valid NL Postal Code");
	ok( method( "6789YZ"), "Valid NL Postal Code");
	ok(!method( "123AA"), "Invalid NL Postal Code: not enough digits");
	ok(!method( "12345ZZ"), "Invalid NL Postal Code: too many digits");
	ok(!method( "1234  AA"), "Invalid NL Postal Code: too many spaces");
	ok(!method( "AA1234"), "Invalid NL Postal Code");
	ok(!method( "1234-AA"), "Invalid NL Postal Code");
});

test("bankaccountNL", function() {
	var method = methodTest("bankaccountNL");
	ok( method( "755490975"), "Valid NL bank account");
	ok( method( "75 54 90 975"), "Valid NL bank account");
	ok( method( "123456789"), "Valid NL bank account");
	ok( method( "12 34 56 789"), "Valid NL bank account");
	ok(!method( "12 3456789"), "Valid NL bank account: inconsistent spaces");
	ok(!method( "123 45 67 89"), "Valid NL bank account: incorrect spaces");
	ok(!method( "755490971"), "Invalid NL bank account");
	ok(!method( "755490973"), "Invalid NL bank account");
	ok(!method( "755490979"), "Invalid NL bank account");
	ok(!method( "123456781"), "Invalid NL bank account");
	ok(!method( "123456784"), "Invalid NL bank account");
	ok(!method( "123456788"), "Invalid NL bank account");
});

test("giroaccountNL", function() {
	var method = methodTest("giroaccountNL");
	ok( method( "123"), "Valid NL giro  account");
	ok( method( "1234567"), "Valid NL giro account");
	ok(!method( "123456788"), "Invalid NL giro account");
});

test("bankorgiroaccountNL", function() {
	var method = methodTest("bankorgiroaccountNL");
	ok( method( "123"), "Valid NL giro account");
	ok( method( "1234567"), "Valid NL giro account");
	ok( method( "123456789"), "Valid NL bank account");
	ok(!method( "12345678"), "Invalid NL bank or giro account");
	ok(!method( "123456788"), "Invalid NL bank or giro account");
});

test("time", function() {
	var method = methodTest("time");
	ok( method( "00:00" ), "Valid time, lower bound" );
	ok( method( "23:59" ), "Valid time, upper bound" );
	ok( method( "3:59" ), "Valid time, single digit hour" );
	ok(!method( "12" ), "Invalid time" );
	ok(!method( "29:59" ), "Invalid time" );
	ok(!method( "00:60" ), "Invalid time" );
	ok(!method( "24:60" ), "Invalid time" );
	ok(!method( "24:00" ), "Invalid time" );
	ok(!method( "30:00" ), "Invalid time" );
	ok(!method( "29:59" ), "Invalid time" );
	ok(!method( "120:00" ), "Invalid time" );
	ok(!method( "12:001" ), "Invalid time" );
	ok(!method( "12:00a" ), "Invalid time" );
});

test("time12h", function() {
	var method = methodTest("time12h");
	ok( method( "12:00 AM" ), "Valid time, lower bound, am" );
	ok( method( "11:59 AM" ), "Valid time, upper bound, am" );
	ok( method( "12:00AM" ), "Valid time, no space, am" );
	ok( method( "12:00PM" ), "Valid time, no space, pm" );
	ok( method( "12:00 PM" ), "Valid time, lower bound, pm" );
	ok( method( "11:59 PM" ), "Valid time, upper bound, pm" );
	ok( method( "11:59 am" ), "Valid time, also accept lowercase" );
	ok( method( "11:59 pm" ), "Valid time, also accept lowercase" );
	ok( method( "1:59 pm" ), "Valid time, single hour, no leading 0" );
	ok( method( "01:59 pm" ), "Valid time, single hour, leading 0" );
	ok(!method( "12:00" ), "Invalid time" );
	ok(!method( "9" ), "Invalid time" );
	ok(!method( "9 am"), "Invalid time" );
	ok(!method( "12:61 am" ), "Invalid time" );
	ok(!method( "13:00 am" ), "Invalid time" );
	ok(!method( "00:00 am" ), "Invalid time" );
});

test("minWords", function() {
	var method = methodTest("minWords");
	ok( method( "hello worlds", 2 ), "plain text, valid" );
	ok( method( "<b>hello</b> world", 2 ), "html, valid" );
	ok(!method( "hello", 2 ), "plain text, invalid" );
	ok(!method( "<b>world</b>", 2 ), "html, invalid" );
	ok(!method( "world <br/>", 2 ), "html, invalid" );
});

test("maxWords", function() {
	var method = methodTest("maxWords");
	ok( method( "hello", 2 ), "plain text, valid" );
	ok( method( "<b>world</b>", 2 ), "html, valid" );
	ok( method( "world <br/>", 2 ), "html, valid" );
	ok( method( "hello worlds", 2 ), "plain text, valid" );
	ok( method( "<b>hello</b> world", 2 ), "html, valid" );
	ok(!method( "hello 123 world", 2 ), "plain text, invalid" );
	ok(!method( "<b>hello</b> 123 world", 2 ), "html, invalid" );
});

test("rangeWords", function() {
	var method = methodTest("rangeWords");
	ok( method( "hello", [ 0, 2 ] ), "plain text, valid" );
	ok( method( "hello worlds", [ 0, 2 ] ), "plain text, valid" );
	ok( method( "<b>hello</b> world", [ 0, 2 ] ), "html, valid" );
	ok(!method( "hello worlds what is up", [ 0, 2 ] ), "plain text, invalid" );
	ok(!method( "<b>Hello</b> <b>world</b> <b>hello</b>", [ 0, 2 ] ), "html, invalid" );
});

test("pattern", function() {
	var method = methodTest("pattern");
	ok( method( "AR1004", "AR\\d{4}" ), "Correct format for the given RegExp" );
	ok( method( "AR1004", /^AR\d{4}$/ ), "Correct format for the given RegExp" );
	ok(!method( "BR1004", /^AR\d{4}$/ ), "Invalid format for the given RegExp" );
	ok( method( "1ABC", "[0-9][A-Z]{3}" ), "Correct format for the given RegExp" );
	ok(!method( "ABC", "[0-9][A-Z]{3}" ), "Invalid format for the given RegExp" );
	ok(!method( "1ABC DEF", "[0-9][A-Z]{3}" ), "Invalid format for the given RegExp" );
	ok( method( "1ABCdef", "[a-zA-Z0-9]+" ), "Correct format for the given RegExp" );
	ok(!method( "1ABC def", "[a-zA-Z0-9]+" ), "Invalid format for the given RegExp" );
	ok( method( "2014-10-02", "[0-9]{4}-(0[1-9]|1[012])-(0[1-9]|1[0-9]|2[0-9]|3[01])" ), "Correct format for the given RegExp" );
	ok(!method( "02-10-2014", "[0-9]{4}-(0[1-9]|1[012])-(0[1-9]|1[0-9]|2[0-9]|3[01])" ), "Invalid format for the given RegExp" );
});

function testCardTypeByNumber(number, cardname, expected) {
	$("#cardnumber").val(number);
	var actual = $("#ccform").valid();
	equal(actual, expected, $.validator.format("Expect card number {0} to validate to {1}, actually validated to ", number, expected));
}

test("creditcardtypes, all", function() {
	$("#ccform").validate({
		rules: {
			cardnumber: {
				creditcard: true,
				creditcardtypes: {
					all: true
				}
			}
		}
	});

	testCardTypeByNumber( "4111-1111-1111-1111", "VISA", true );
	testCardTypeByNumber( "5111-1111-1111-1118", "MasterCard", true );
	testCardTypeByNumber( "6111-1111-1111-1116", "Discover", true );
	testCardTypeByNumber( "3400-0000-0000-009", "AMEX", true );

	testCardTypeByNumber( "4111-1111-1111-1110", "VISA", false );
	testCardTypeByNumber( "5432-1111-1111-1111", "MasterCard", false );
	testCardTypeByNumber( "6611-6611-6611-6611", "Discover", false );
	testCardTypeByNumber( "3777-7777-7777-7777", "AMEX", false );
});

test("creditcardtypes, visa", function() {
	$("#ccform").validate({
		rules: {
			cardnumber: {
				creditcard: true,
				creditcardtypes: {
					visa: true
				}
			}
		}
	});

	testCardTypeByNumber( "4111-1111-1111-1111", "VISA", true );
	testCardTypeByNumber( "5111-1111-1111-1118", "MasterCard", false );
	testCardTypeByNumber( "6111-1111-1111-1116", "Discover", false );
	testCardTypeByNumber( "3400-0000-0000-009", "AMEX", false );
});

test("creditcardtypes, mastercard", function() {
	$("#ccform").validate({
		rules: {
			cardnumber: {
				creditcard: true,
				creditcardtypes: {
					mastercard: true
				}
			}
		}
	});

	testCardTypeByNumber( "5111-1111-1111-1118", "MasterCard", true );
	testCardTypeByNumber( "6111-1111-1111-1116", "Discover", false );
	testCardTypeByNumber( "3400-0000-0000-009", "AMEX", false );
	testCardTypeByNumber( "4111-1111-1111-1111", "VISA", false );
});

/*
function fillFormWithValuesAndExpect(formSelector, inputValues, expected) {
	var i, actual;

	for (i = 0; i < inputValues.length; i++) {
		$(formSelector + " input:eq(" + i + ")").val(inputValues[i]);
	}
	actual = $(formSelector).valid();
	equal(actual, expected, $.validator.format("Filled inputs of form '{0}' with {1} values ({2})", formSelector, inputValues.length, inputValues.toString()));

}

test("require_from_group", function() {
	$("#productInfo").validate({
		rules: {
			partnumber: { require_from_group: [ 2, ".productInfo" ] },
			description: { require_from_group: [ 2, ".productInfo" ] },
			discount: { require_from_group: [ 2, ".productInfo" ] }
		}
	});

	fillFormWithValuesAndExpect("#productInfo", [], false);
	fillFormWithValuesAndExpect("#productInfo", [ 123 ], false);
	$("#productInfo input[type='checkbox']").attr("checked", "checked");
	fillFormWithValuesAndExpect("#productInfo", [ 123 ], true);
	$("#productInfo input[type='checkbox']").removeAttr("checked");
	fillFormWithValuesAndExpect("#productInfo", [ 123, "widget" ], true);
	fillFormWithValuesAndExpect("#productInfo", [ 123, "widget", "red" ], true);
	fillFormWithValuesAndExpect("#productInfo", [ 123, "widget", "red" ], true);
});

test("require_from_group preserve other rules", function() {
	$("#productInfo").validate({
		rules: {
			partnumber: { require_from_group: [ 2, ".productInfo" ] },
			description: { require_from_group: [ 2, ".productInfo" ] },
			color: { require_from_group: [ 2, ".productInfo" ] },
			supplier: { required: true }
		}
	});

	fillFormWithValuesAndExpect("#productInfo", [], false);
	fillFormWithValuesAndExpect("#productInfo", [ 123 ], false);
	fillFormWithValuesAndExpect("#productInfo", [ 123, "widget" ], false);
	fillFormWithValuesAndExpect("#productInfo", [ "", "", "", "Acme" ], false);
	fillFormWithValuesAndExpect("#productInfo", [ 123, "", "", "Acme" ], false);
	fillFormWithValuesAndExpect("#productInfo", [ 123, "widget", "", "Acme" ], true);
	fillFormWithValuesAndExpect("#productInfo", [ 123, "widget", "red", "Acme" ], true);
});

test("skip_or_fill_minimum", function() {
	$("#productInfo").validate({
		rules: {
			partnumber:  { skip_or_fill_minimum: [ 2, ".productInfo" ] },
			description: { skip_or_fill_minimum: [ 2, ".productInfo" ] },
			color:       { skip_or_fill_minimum: [ 2, ".productInfo" ] }
		}
	});

	fillFormWithValuesAndExpect("#productInfo", [], true);
	fillFormWithValuesAndExpect("#productInfo", [ 123 ], false);
	fillFormWithValuesAndExpect("#productInfo", [ 123, "widget" ], true);
	fillFormWithValuesAndExpect("#productInfo", [ 123, "widget", "red" ], true);
});

test("skip_or_fill_minimum preserve other rules", function() {
	$("#productInfo").validate({
		rules: {
			partnumber:  { skip_or_fill_minimum: [ 2, ".productInfo" ] },
			description: { skip_or_fill_minimum: [ 2, ".productInfo" ] },
			color:       { skip_or_fill_minimum: [ 2, ".productInfo" ] },
			supplier: { required: true }
		}
	});

	fillFormWithValuesAndExpect("#productInfo", [], false);
	fillFormWithValuesAndExpect("#productInfo", [ "", "", "", "Acme" ], true);
	fillFormWithValuesAndExpect("#productInfo", [ 123, "", "", "Acme" ], false);
	fillFormWithValuesAndExpect("#productInfo", [ 123, "widget", "", "Acme" ], true);
	fillFormWithValuesAndExpect("#productInfo", [ 123, "widget", "red", "Acme" ], true);
});
*/

test("zipcodeUS", function() {
	var method = methodTest("zipcodeUS");
	ok( method( "12345" ), "Valid zip" );
	ok( method( "12345-2345" ), "Valid zip" );
	ok( method( "90210-4567" ), "Valid zip" );
	ok(!method( "1" ), "Invalid zip" );
	ok(!method( "1234" ), "Invalid zip" );
	ok(!method( "123-23" ), "Invalid zip" );
	ok(!method( "12345-43" ), "Invalid zip" );
	ok(!method( "123456-7890" ), "Invalid zip" );
});

test("nifES", function() {
	var method = methodTest("nifES");
	ok( method( "11441059P" ), "NIF valid" );
	ok( method( "80054306T" ), "NIF valid" );
	ok( method( "76048581R" ), "NIF valid" );
	ok( method( "28950849J" ), "NIF valid" );
	ok( method( "34048598L" ), "NIF valid" );
	ok( method( "28311529R" ), "NIF valid" );
	ok( method( "34673804Q" ), "NIF valid" );
	ok( method( "92133247P" ), "NIF valid" );
	ok( method( "77149717N" ), "NIF valid" );
	ok( method( "15762034L" ), "NIF valid" );
	ok( method( "05122654W" ), "NIF valid" );
	ok( method( "05122654w" ), "NIF valid: lower case" );
	ok(!method( "1144105R" ), "NIF invalid: less than 8 digits without zero" );
	ok(!method( "11441059 R" ), "NIF invalid: white space" );
	ok(!method( "11441059" ), "NIF invalid: no letter" );
	ok(!method( "11441059PR" ), "NIF invalid: two letters" );
	ok(!method( "11440059R" ), "NIF invalid: wrong number" );
	ok(!method( "11441059S" ), "NIF invalid: wrong letter" );
	ok(!method( "114410598R" ), "NIF invalid: > 8 digits" );
	ok(!method( "11441059-R" ), "NIF invalid: dash" );
	ok(!method( "asdasdasd" ), "NIF invalid: all letters" );
	ok(!method( "11.144.059R" ), "NIF invalid: two dots" );
	ok(!method( "05.122.654R" ), "NIF invalid: starts with 0 and dots" );
	ok(!method( "5.122.654-R" ), "NIF invalid:  dots and dash" );
	ok(!method( "05.122.654-R" ), "NIF invalid: starts with zero and dot and dash" );
});

test("nieES", function() {
	var method = methodTest("nieES");
	ok( method( "X0093999K" ), "NIE valid" );
	ok( method( "X1923000Q" ), "NIE valid" );
	ok( method( "Z9669587R" ), "NIE valid" );
	ok( method( "Z8945005B" ), "NIE valid" );
	ok( method( "Z6663465W" ), "NIE valid" );
	ok( method( "Y7875935J" ), "NIE valid" );
	ok( method( "X3390130E" ), "NIE valid" );
	ok( method( "Y7699182S" ), "NIE valid" );
	ok( method( "Y1524243R" ), "NIE valid" );
	ok( method( "X3744072V" ), "NIE valid" );
	ok( method( "X7436800A" ), "NIE valid" );
	ok( method( "y7875935j" ), "NIE valid: lower case" );
	ok(!method( "X0093999 K" ), "NIE inválido: white space" );
	ok(!method( "X 0093999 K" ), "NIE inválido:  white space" );
	ok(!method( "11441059" ), "NIE inválido: no letter" );
	ok(!method( "11441059PR" ), "NIE inválido: two letters" );
	ok(!method( "11440059R" ), "NIE inválido: wrong number" );
	ok(!method( "11441059S" ), "NIE inválido: wrong letter" );
	ok(!method( "114410598R" ), "NIE inválido: > 8 digits" );
	ok(!method( "11441059-R" ), "NIE inválido: dash" );
	ok(!method( "asdasdasd" ), "NIE inválido: all letters" );
	ok(!method( "11.144.059R" ), "NIE inválido: two dots" );
	ok(!method( "05.122.654R" ), "NIE inválido: starts with 0 and dots" );
	ok(!method( "5.122.654-R" ), "NIE inválido: dots and dash" );
	ok(!method( "05.122.654-R" ), "NIE inválido: starts with zero and dot and dash" );
});

test("cifES", function() {
	var method = methodTest("cifES");
	ok( method( "A79082244" ), "CIF valid" );
	ok( method( "A60917978" ), "CIF valid" );
	ok( method( "A39000013" ), "CIF valid" );
	ok( method( "B43522192" ), "CIF valid" );
	ok( method( "B38624334" ), "CIF valid" );
	ok( method( "G72102064" ), "CIF valid" );
	ok( method( "F41190612" ), "CIF valid" );
	ok( method( "J85081081" ), "CIF valid" );
	ok( method( "S98038813" ), "CIF valid" );
	ok( method( "G32937757" ), "CIF valid" );
	ok( method( "B46125746" ), "CIF valid" );
	ok( method( "C27827559" ), "CIF valid" );
	ok( method( "E48911572" ), "CIF valid" );
	ok( method( "s98038813" ), "CIF valid: lower case" );
	ok(!method( "K48911572" ), "CIF invalid: starts with K" );
	ok(!method( "L48911572" ), "CIF invalid: starts with L" );
	ok(!method( "M48911572" ), "CIF invalid: starts with M" );
	ok(!method( "X48911572" ), "CIF invalid: starts with X" );
	ok(!method( "Y48911572" ), "CIF invalid: starts with Y" );
	ok(!method( "Z48911572" ), "CIF invalid: starts with Z" );
	ok(!method( "M15661515" ), "CIF invalid" );
	ok(!method( "Z98038813" ), "CIF invalid: wrong letter" );
	ok(!method( "B 43522192" ), "CIF invalid: white spaces" );
	ok(!method( "43522192" ), "CIF invalid: missing letter" );
	ok(!method( "BB43522192" ), "CIF invalid: two letters" );
	ok(!method( "B53522192" ), "CIF invalid: wrong number" );
	ok(!method( "B433522192" ), "CIF invalid: > 8 digits" );
	ok(!method( "B3522192" ), "CIF invalid: < 8 digits" );
	ok(!method( "B-43522192" ), "CIF invalid: dash" );
	ok(!method( "Basdasdas" ), "CIF invalid: all letters" );
	ok(!method( "B43.522.192" ), "CIF invalid: dots" );
	ok(!method( "B-43.522.192" ), "CIF invalid: dots and dash" );
});

test("maxWords", function() {
	var method = methodTest("maxWords"),
		maxWords = 6;

	ok( method( "I am a sentence", maxWords), "Max Words");
	ok(!method( "I'm way too long for this sentence!", maxWords), "Too many words");
	ok(method( "Don’t “count” me as too long", maxWords), "Right amount of words with smartquotes");
	ok(!method( "But you can “count” me as too long", maxWords), "Too many words with smartquotes");
	ok(method( "<div>Don’t “count” me as too long</div>", maxWords), "Right amount of words with smartquotes w/ HTML");
	ok(!method( "<div>But you can “count” me as too long</div>", maxWords), "Too many words with smartquotes w/ HTML");
});

test("minWords", function() {
	var method = methodTest("minWords"),
		minWords = 6;

	ok(!method( "I am a short sentence", minWords), "Max Words");
	ok( method( "I'm way too long for this sentence!", minWords), "Too many words");
	ok(!method( "Don’t “count” me as short.", minWords), "Right amount of words with smartquotes");
	ok( method( "But you can “count” me as too short", minWords), "Too many words with smartquotes");
	ok(!method( "<div>“Count” me as too short.</div>", minWords), "Right amount of words with smartquotes w/ HTML");
	ok( method( "<div>But you can “count” me as too long</div>", minWords), "Too many words with smartquotes w/ HTML");
});

test("rangeWords", function() {
	var method = methodTest("rangeWords"),
		rangeWords = [ 3, 6 ];

	ok(!method( "I'm going to be longer than “six words!”", rangeWords), "Longer than 6 with smartquotes");
	ok( method( "I'm just the right amount!", rangeWords), "In between");
	ok( method( "Super short sentence’s.", rangeWords), "Low end");
	ok(!method( "I", rangeWords), "Too short");
	ok( method( "<div>“Count” me as perfect.</div>", rangeWords), "Right amount of words with smartquotes w/ HTML");
	ok(!method( "<div>But you can “count” me as too long</div>", rangeWords), "Too many words with smartquotes w/ HTML");
});

test("currency", function() { // Works with any symbol
	var method = methodTest( "currency" );
	ok( method( "£9", "£"), "Symbol no decimal" );
	ok( method( "£9.9", "£"), "£, one decimal" );
	ok( method( "£9.99", "£"), "£, two decimal" );
	ok( method( "£9.90", "£"), "Valid currency" );
	ok( method( "£9,999.9", "£"), "£, thousand, comma separator, one decimal" );
	ok( method( "£9,999.99", "£"), "£, thousand, comma separator, two decimal" );
	ok( method( "£9,999,999.9", "£"), "£, million, comma separators, one decimal" );
	ok( method( "9", [ "£", false ]), "Valid currency" );
	ok( method( "9.9", [ "£", false ]), "Valid currency" );
	ok( method( "9.99", [ "£", false ]), "Valid currency" );
	ok( method( "9.90", [ "£", false ]), "Valid currency" );
	ok( method( "9,999.9", [ "£", false ]), "Valid currency" );
	ok( method( "9,999.99", [ "£", false ]), "Valid currency" );
	ok( method( "9,999,999.9", [ "£", false ]), "Valid currency" );
	ok(!method( "9,", "£"), "Invalid currency" );
	ok(!method( "9,99.99", "£"), "Invalid currency" );
	ok(!method( "9,", "£"), "Invalid currency" );
	ok(!method( "9.999", "£"), "Invalid currency" );
	ok(!method( "9.999", "£"), "Invalid currency" );
	ok(!method( "9.99,9", "£"), "Invalid currency" );
});

test("postalCodeCA", function() {
	var method = methodTest("postalCodeCA");
	ok( method( "H0H 0H0"), "Valid CA Postal Code; Single space" );
	ok( !method( "H0H0H0"), "Inalid CA Postal Code; No space" );
	ok( !method( "H0H-0H0"), "Invalid CA Postal Code; Single dash" );
	ok( !method( "H0H 0H"), "Invalid CA Postal Code; Too Short" );
	ok( !method( "Z0H 0H"), "Invalid CA Postal Code; Only 'ABCEGHJKLMNPRSTVXY' are valid starting characters" );
	ok( !method( "h0h 0h0"), "Invalid CA Postal Code; Only upper case characters" );
});

test("stateUS", function() {
	var method = methodTest("stateUS");
	ok( method( "AZ" ), "Valid US state" );
	ok( method( "OH" ), "Valid US state" );
	ok( method( "DC" ), "Valid US state" );
	ok( method( "PR", { includeTerritories: true } ), "Valid US territory" );
	ok( method( "AA", { includeMilitary: true } ), "Valid US military zone" );
	ok( method( "me", { caseSensitive: false } ), "Valid US state" );
	ok(!method( "az", { caseSensitive: true } ), "Must be capital letters" );
	ok(!method( "mp", { caseSensitive: false, includeTerritories: false } ), "US territories not allowed" );
});

test("postalcodeBR", function() {
	var method = methodTest("postalcodeBR");
	ok( method( "99999-999"), "Valid BR Postal Code");
	ok( method( "99999999"), "Valid BR Postal Code");
	ok( method( "99.999-999"), "Valid BR Postal Code");
	ok( !method( "99.999999"), "Invalid BR Postal Code");
});

test("cpfBR", function() {
	var method = methodTest("cpfBR");
	ok( method( "11144477735"), "Valid CPF Number");
	ok( method( "263.946.533-30"), "Valid CPF Number");
	ok( method( "325 861 044 47"), "Valid CPF Number");
	ok( method( "859-684-732-40"), "Valid CPF Number");
	ok( !method( "99999999999"), "Invalid CPF Number: dump data");
	ok( !method( "1114447773"), "Invalid CPF Number: < 11 digits");
	ok( !method( "111444777355"), "Invalid CPF Number: > 11 digits");
	ok( !method( "11144477715"), "Invalid CPF Number: 1st check number failed");
	ok( !method( "11144477737"), "Invalid CPF Number: 2nd check number failed");
});

})(jQuery);
