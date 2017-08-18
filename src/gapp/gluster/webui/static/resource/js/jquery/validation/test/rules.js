module("rules");

test("rules() - internal - input", function() {
	var element = $("#firstname");

	$("#testForm1").validate();

	deepEqual( element.rules(), { required: true, minlength: 2 } );
});

test("rules(), ignore method:false", function() {
	var element = $("#firstnamec");

	$("#testForm1clean").validate({
		rules: {
			firstnamec: { required: false, minlength: 2 }
		}
	});

	deepEqual( element.rules(), { minlength: 2 } );
});

test("rules() HTML5 required (no value)", function() {
	var element = $("#testForm11text1");

	$("#testForm11").validate();

	deepEqual( element.rules(), { required: true } );
});

test("rules() - internal - select", function() {
	var element = $("#meal");

	$("#testForm3").validate();

	deepEqual( element.rules(), { required: true } );
});

test("rules() - external", function() {
	var element = $("#text1");

	$("#form").validate({
		rules: {
			action: { date: true, min: 5 }
		}
	});

	deepEqual( element.rules(), { date: true, min: 5 } );
});

test("rules() - external - complete form", function() {
	expect(1);

	var methods = $.extend({}, $.validator.methods),
		messages = $.extend({}, $.validator.messages),
		v;

	$.validator.addMethod("verifyTest", function() {
		ok( true, "method executed" );
		return true;
	});
	v = $("#form").validate({
		rules: {
			action: { verifyTest: true }
		}
	});
	v.form();

	$.validator.methods = methods;
	$.validator.messages = messages;
});

test("rules() - internal - input", function() {
	var element = $("#form8input");

	$("#testForm8").validate();

	deepEqual( element.rules(), { required: true, number: true, rangelength: [ 2, 8 ] } );
});

test("rules(), merge min/max to range, minlength/maxlength to rangelength", function() {
	jQuery.validator.autoCreateRanges = true;

	$("#testForm1clean").validate({
		rules: {
			firstnamec: {
				min: -15,
				max: 0
			},
			lastname: {
				minlength: 0,
				maxlength: 10
			}
		}
	});

	deepEqual( $("#firstnamec").rules(), { range: [ -15, 0 ] } );
	deepEqual( $("#lastnamec").rules(), { rangelength: [ 0, 10 ] } );

	jQuery.validator.autoCreateRanges = false;
});

test("rules(), guarantee that required is at front", function() {
	$("#testForm1").validate();
	var v = $("#v2").validate();
	$("#subformRequired").validate();
	function flatRules(element) {
		var result = [];
		jQuery.each($(element).rules(), function(key) { result.push(key); });
		return result.join(" ");
	}
	equal( "required minlength", flatRules("#firstname") );
	equal( "required minlength maxlength", flatRules("#v2-i6") );
	equal( "required maxlength", flatRules("#co_name") );

	QUnit.reset();
	jQuery.validator.autoCreateRanges = true;
	v = $("#v2").validate();
	equal( "required rangelength", flatRules("#v2-i6") );

	$("#subformRequired").validate({
		rules: {
			co_name: "required"
		}
	});
	$("#co_name").removeClass();
	equal( "required maxlength", flatRules("#co_name") );
	jQuery.validator.autoCreateRanges = false;
});

test("rules(), evaluate dynamic parameters", function() {
	expect(2);

	$("#testForm1clean").validate({
		rules: {
			firstnamec: {
				min: function(element) {
					equal( $("#firstnamec")[0], element );
					return 12;
				}
			}
		}
	});

	deepEqual( $("#firstnamec").rules(), { min: 12 });
});

test("rules(), class and attribute combinations", function() {

	$.validator.addMethod("customMethod1", function() {
		return false;
	}, "");
	$.validator.addMethod("customMethod2", function() {
		return false;
	}, "");

	$("#v2").validate({
		rules: {
			"v2-i7": {
				required: true,
				minlength: 2,
				customMethod: true
			}
		}
	});

	deepEqual( $("#v2-i1").rules(), { required: true });
	deepEqual( $("#v2-i2").rules(), { required: true, email: true });
	deepEqual( $("#v2-i3").rules(), { url: true });
	deepEqual( $("#v2-i4").rules(), { required: true, minlength: 2 });
	deepEqual( $("#v2-i5").rules(), { required: true, minlength: 2, maxlength: 5, customMethod1: "123" });
	jQuery.validator.autoCreateRanges = true;
	deepEqual( $("#v2-i5").rules(), { required: true, customMethod1: "123", rangelength: [ 2, 5 ] });
	deepEqual( $("#v2-i6").rules(), { required: true, customMethod2: true, rangelength: [ 2, 5 ] });
	jQuery.validator.autoCreateRanges = false;
	deepEqual( $("#v2-i7").rules(), { required: true, minlength: 2, customMethod: true });

	delete $.validator.methods.customMethod1;
	delete $.validator.messages.customMethod1;
	delete $.validator.methods.customMethod2;
	delete $.validator.messages.customMethod2;
});

test("rules(), dependency checks", function() {
	var v = $("#testForm1clean").validate({
			rules: {
				firstnamec: {
					min: {
						param: 5,
						depends: function(el) {
							return (/^a/).test($(el).val());
						}
					}
				},
				lastname: {
					max: {
						param: 12
					},
					email: {
						depends: function() { return true; }
					}
				}
			}
		}),
		rules = $("#firstnamec").rules();

	equal( 0, v.objectLength(rules) );

	$("#firstnamec").val("ab");
	deepEqual( $("#firstnamec").rules(), { min: 5 });

	deepEqual( $("#lastnamec").rules(), { max: 12, email: true });
});

test("rules(), add and remove", function() {
	$.validator.addMethod("customMethod1", function() {
		return false;
	}, "");
	$("#v2").validate();
	$("#v2-i5").removeClass("required").removeAttr("minlength maxlength");
	deepEqual( $("#v2-i5").rules(), { customMethod1: "123" });

	$("#v2-i5").addClass("required").attr({
		minlength: 2,
		maxlength: 5
	});
	deepEqual( $("#v2-i5").rules(), { required: true, minlength: 2, maxlength: 5, customMethod1: "123" });

	$("#v2-i5").addClass("email").attr({ min: 5 });
	deepEqual( $("#v2-i5").rules(), { required: true, email: true, minlength: 2, maxlength: 5, min: 5, customMethod1: "123" });

	$("#v2-i5").removeClass("required email").removeAttr("minlength maxlength customMethod1 min");
	deepEqual( $("#v2-i5").rules(), {});

	delete $.validator.methods.customMethod1;
	delete $.validator.messages.customMethod1;
});

test("rules(), add and remove static rules", function() {

	$("#testForm1clean").validate({
		rules: {
			firstnamec: "required date"
		}
	});

	deepEqual( $("#firstnamec").rules(), { required: true, date: true } );

	$("#firstnamec").rules("remove", "date");
	deepEqual( $("#firstnamec").rules(), { required: true } );
	$("#firstnamec").rules("add", "email");
	deepEqual( $("#firstnamec").rules(), { required: true, email: true } );

	$("#firstnamec").rules("remove", "required");
	deepEqual( $("#firstnamec").rules(), { email: true } );

	deepEqual( $("#firstnamec").rules("remove"), { email: true } );
	deepEqual( $("#firstnamec").rules(), { } );

	$("#firstnamec").rules("add", "required email");
	deepEqual( $("#firstnamec").rules(), { required: true, email: true } );

	deepEqual( $("#lastnamec").rules(), {} );
	$("#lastnamec").rules("add", "required");
	$("#lastnamec").rules("add", {
		minlength: 2
	});
	deepEqual( $("#lastnamec").rules(), { required: true, minlength: 2 } );

	var removedRules = $("#lastnamec").rules("remove", "required email");
	deepEqual( $("#lastnamec").rules(), { minlength: 2 } );
	$("#lastnamec").rules("add", removedRules);
	deepEqual( $("#lastnamec").rules(), { required: true, minlength: 2 } );
});

test("rules(), add messages", function() {
	$("#firstnamec").attr("title", null);
	var v = $("#testForm1clean").validate({
		rules: {
			firstnamec: "required"
		}
	});
	$("#testForm1clean").valid();
	$("#firstnamec").valid();
	deepEqual( v.settings.messages.firstname, undefined );

	$("#firstnamec").rules("add", {
		messages: {
			required: "required"
		}
	});

	$("#firstnamec").valid();
	deepEqual( v.errorList[0] && v.errorList[0].message, "required" );

	$("#firstnamec").val("test");
	$("#firstnamec").valid();
	equal(v.errorList.length, 0);
});

test( "rules(), rangelength attribute as array", function() {
	$("#testForm13").validate();
	deepEqual( $("#cars-select").rules(), {
		required: true,
		rangelength: [ 2, 3 ]
	});
});
