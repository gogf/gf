module("aria");

test("Invalid field adds aria-invalid=true", function() {
	var ariaInvalidFirstName = $("#ariaInvalidFirstName"),
		form = $("#ariaInvalid");

	form.validate({
		rules: {
			ariaInvalidFirstName: "required"
		}
	});
	ariaInvalidFirstName.val("");
	ariaInvalidFirstName.valid();
	equal(ariaInvalidFirstName.attr("aria-invalid"), "true");
});

test("Valid field adds aria-invalid=false", function() {
	var ariaInvalidFirstName = $("#ariaInvalidFirstName"),
		form = $("#ariaInvalid");

	form.validate({
		rules: {
			ariaInvalidFirstName: "required"
		}
	});
	ariaInvalidFirstName.val("not empty");
	ariaInvalidFirstName.valid();
	equal(ariaInvalidFirstName.attr("aria-invalid"), "false");
	equal($("#ariaInvalid [aria-invalid=false]").length, 1);
});

test("resetForm(): removes all aria-invalid attributes", function() {
	var ariaInvalidFirstName = $("#ariaInvalidFirstName"),
		form = $("#ariaInvalid"),
		validator = form.validate({
			rules: {
				ariaInvalidFirstName: "required"
			}
		});

	ariaInvalidFirstName.val("not empty");
	ariaInvalidFirstName.valid();
	validator.resetForm();
	equal($("#ariaInvalid [aria-invalid]").length, 0, "resetForm() should remove any aria-invalid attributes");
});

test("Static required field adds aria-required", function() {
	var ariaRequiredStatic = $("#ariaRequiredStatic"),
		form = $("#ariaRequired");

	form.validate();
	equal(ariaRequiredStatic.attr("aria-required"), "true");
});

test("Data required field adds aria-required", function() {
	var ariaRequiredData = $("#ariaRequiredData"),
		form = $("#ariaRequired");

	form.validate();
	equal(ariaRequiredData.attr("aria-required"), "true");
});

test("Class required field adds aria-required", function() {
	var ariaRequiredClass = $("#ariaRequiredClass"),
		form = $("#ariaRequired");

	form.validate();
	equal(ariaRequiredClass.attr("aria-required"), "true");
});

test("Dynamically required field adds aria-required after valid()", function() {
	var ariaRequiredDynamic = $("#ariaRequiredDynamic"),
		form = $("#ariaRequired");

	form.resetForm();
	form.validate({
		rules: {
			ariaRequiredDynamic: "required"
		}
	});
	ariaRequiredDynamic.valid();
	equal(ariaRequiredDynamic.attr("aria-required"), "true");
});
