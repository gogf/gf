module( "placement" );

test( "elements() order", function() {
	var container = $( "#orderContainer" ),
		v = $( "#elementsOrder" ).validate({
			errorLabelContainer: container,
			wrap: "li"
		});

	deepEqual(
		v.elements().map( function() {
			return $( this ).attr( "id" );
		}).get(),
		[
			"order1",
			"order2",
			"order3",
			"order4",
			"order5",
			"order6"
		],
		"elements must be in document order"
	);

	v.form();
	deepEqual(
		container.children().map( function() {
			return $( this ).attr( "id" );
		}).get(),
		[
			"order1-error",
			"order2-error",
			"order3-error",
			"order4-error",
			"order5-error",
			"order6-error"
		],
		"labels in error container must be in document order"
	);
});

test( "error containers, simple", function() {
	expect( 14 );
	var container = $( "#simplecontainer" ),
		v = $( "#form" ).validate({
			errorLabelContainer: container,
			showErrors: function() {
				container.find( "h3" ).html( jQuery.validator.format( "There are {0} errors in your form.", this.size()) );
				this.defaultShowErrors();
			}
		});

	v.prepareForm();
	ok( v.valid(), "form is valid" );
	equal( 0, container.find( ".error:not(input)" ).length, "There should be no error labels" );
	equal( "", container.find( "h3" ).html() );

	v.prepareForm();
	v.errorList = [
		{
			message: "bar",
			element: {
				name: "foo"
			}
		},
		{
			message: "necessary",
			element: {
				name: "required"
			}
		}
	];

	ok( !v.valid(), "form is not valid after adding errors manually" );
	v.showErrors();
	equal( container.find( ".error:not(input)" ).length, 2, "There should be two error labels" );
	ok( container.is( ":visible" ), "Check that the container is visible" );
	container.find( ".error:not(input)" ).each(function() {
		ok( $( this ).is( ":visible" ), "Check that each label is visible" );
	});
	equal( "There are 2 errors in your form.", container.find( "h3" ).html() );

	v.prepareForm();
	ok( v.valid(), "form is valid after a reset" );
	v.showErrors();
	equal( container.find( ".error:not(input)" ).length, 2, "There should still be two error labels" );
	ok( container.is( ":hidden" ), "Check that the container is hidden" );
	container.find( ".error:not(input)" ).each(function() {
		ok( $( this ).is( ":hidden" ), "Check that each label is hidden" );
	});
});

test( "error containers, with labelcontainer I", function() {
	expect( 16 );
	var container = $( "#container" ),
		labelcontainer = $( "#labelcontainer" ),
		v = $( "#form" ).validate({
			errorContainer: container,
			errorLabelContainer: labelcontainer,
			wrapper: "li"
		});

	ok( v.valid(), "form is valid" );
	equal( 0, container.find( ".error:not(input)" ).length, "There should be no error labels in the container" );
	equal( 0, labelcontainer.find( ".error:not(input)" ).length, "There should be no error labels in the labelcontainer" );
	equal( 0, labelcontainer.find( "li" ).length, "There should be no lis labels in the labelcontainer" );

	v.errorList = [
		{
			message: "bar",
			element: {
				name: "foo"
			}
		},
		{
			name: "required",
			message: "necessary",
			element: {
				name: "required"
			}
		}
	];

	ok( !v.valid(), "form is not valid after adding errors manually" );
	v.showErrors();
	equal( 0, container.find( ".error:not(input)" ).length, "There should be no error label in the container" );
	equal( 2, labelcontainer.find( ".error:not(input)" ).length, "There should be two error labels in the labelcontainer" );
	equal( 2, labelcontainer.find( "li" ).length, "There should be two error lis in the labelcontainer" );
	ok( container.is( ":visible" ), "Check that the container is visible" );
	ok( labelcontainer.is( ":visible" ), "Check that the labelcontainer is visible" );
	labelcontainer.find( ".error:not(input)" ).each(function() {
		ok( $( this ).is( ":visible" ), "Check that each label is visible1" );
		equal( "li", $( this ).parent()[0].tagName.toLowerCase(), "Check that each label is wrapped in an li" );
		ok( $( this ).parent( "li" ).is( ":visible" ), "Check that each parent li is visible" );
	});
});

test( "errorcontainer, show/hide only on submit", function() {
	expect( 14 );
	var container = $( "#container" ),
		labelContainer = $( "#labelcontainer" ),
		v = $( "#testForm1" ).bind( "invalid-form.validate", function() {
			ok( true, "invalid-form event triggered called" );
		}).validate({
			errorContainer: container,
			errorLabelContainer: labelContainer,
			showErrors: function() {
				container.html( jQuery.validator.format( "There are {0} errors in your form.", this.numberOfInvalids()) );
				ok( true, "showErrors called" );
				this.defaultShowErrors();
			}
		});

	equal( "", container.html(), "must be empty" );
	equal( "", labelContainer.html(), "must be empty" );
	// validate whole form, both showErrors and invalidHandler must be called once
	// preferably invalidHandler first, showErrors second
	ok( !v.form(), "invalid form" );
	equal( 2, labelContainer.find( ".error:not(input)" ).length );
	equal( "There are 2 errors in your form.", container.html() );
	ok( labelContainer.is( ":visible" ), "must be visible" );
	ok( container.is( ":visible" ), "must be visible" );

	$( "#firstname" ).val( "hix" ).keyup();
	$( "#testForm1" ).triggerHandler( "keyup", [
			jQuery.event.fix({
				type: "keyup",
				target: $( "#firstname" )[ 0 ]
			})
		]);
	equal( 1, labelContainer.find( ".error:visible" ).length );
	equal( "There are 1 errors in your form.", container.html() );

	$( "#lastname" ).val( "abc" );
	ok( v.form(), "Form now valid, trigger showErrors but not invalid-form" );
});

test( "test label used as error container", function(assert) {
	expect( 8 );
	var form = $( "#testForm16" ),
		field = $( "#testForm16text" );

	form.validate({
		errorPlacement: function( error, element ) {
			// Append error within linked label
			$( "label[for='" + element.attr( "id" ) + "']" ).append( error );
		},
		errorElement: "span"
	});

	ok( !field.valid() );
	equal( "Field Label", field.next( "label" ).contents().first().text(), "container label isn't disrupted" );
	assert.hasError(field, "missing");
	ok( !field.attr( "aria-describedby" ), "field does not require aria-describedby attribute" );

	field.val( "foo" );
	ok( field.valid() );
	equal( "Field Label", field.next( "label" ).contents().first().text(), "container label isn't disrupted" );
	ok( !field.attr( "aria-describedby" ), "field does not require aria-describedby attribute" );
	assert.noErrorFor(field);
});

test( "test error placed adjacent to descriptive label", function(assert) {
	expect( 8 );
	var form = $( "#testForm16" ),
		field = $( "#testForm16text" );

	form.validate({
		errorElement: "span"
	});

	ok( !field.valid() );
	equal( 1, form.find( "label" ).length );
	equal( "Field Label", form.find( "label" ).text(), "container label isn't disrupted" );
	assert.hasError( field, "missing" );

	field.val( "foo" );
	ok( field.valid() );
	equal( 1, form.find( "label" ).length );
	equal( "Field Label", form.find( "label" ).text(), "container label isn't disrupted" );
	assert.noErrorFor( field );
});

test( "test descriptive label used alongside error label", function(assert) {
	expect( 8 );
	var form = $( "#testForm16" ),
		field = $( "#testForm16text" );

	form.validate({
		errorElement: "label"
	});

	ok( !field.valid() );
	equal( 1, form.find( "label.title" ).length );
	equal( "Field Label", form.find( "label.title" ).text(), "container label isn't disrupted" );
	assert.hasError( field, "missing" );

	field.val( "foo" );
	ok( field.valid() );
	equal( 1, form.find( "label.title" ).length );
	equal( "Field Label", form.find( "label.title" ).text(), "container label isn't disrupted" );
	assert.noErrorFor( field );
});

test( "test custom errorElement", function(assert) {
	expect( 4 );
	var form = $( "#userForm" ),
		field = $( "#username" );

	form.validate({
		messages: {
			username: "missing"
		},
		errorElement: "label"
	});

	ok( !field.valid() );
	assert.hasError( field, "missing", "Field should have error 'missing'" );
	field.val( "foo" );
	ok( field.valid() );
	assert.noErrorFor( field, "Field should not have a visible error" );
});

test( "test existing label used as error element", function(assert) {
	expect( 4 );
	var form = $( "#testForm14" ),
		field = $( "#testForm14text" );

	form.validate({ errorElement: "label" });

	ok( !field.valid() );
	assert.hasError( field, "required" );

	field.val( "foo" );
	ok( field.valid() );
	assert.noErrorFor( field );
});

test( "test existing non-label used as error element", function(assert) {
	expect( 4 );
	var form = $( "#testForm15" ),
		field = $( "#testForm15text" );

	form.validate({ errorElement: "span" });

	ok( !field.valid() );
	assert.hasError( field, "required" );

	field.val( "foo" );
	ok( field.valid() );
	assert.noErrorFor( field );
});

test( "test existing non-error aria-describedby", function( assert ) {
	expect( 8 );
	var form = $( "#testForm17" ),
		field = $( "#testForm17text" );

	equal( field.attr( "aria-describedby" ), "testForm17text-description" );
	form.validate({ errorElement: "span" });

	ok( !field.valid() );
	equal( field.attr( "aria-describedby" ), "testForm17text-description testForm17text-error" );
	assert.hasError( field, "required" );

	field.val( "foo" );
	ok( field.valid() );
	assert.noErrorFor( field );

	strictEqual( "This is where you enter your data", $("#testForm17text-description").text() );
	strictEqual( "", $("#testForm17text-error").text(), "Error label is empty for valid field" );
});

test( "test pre-assigned non-error aria-describedby", function( assert ) {
	expect( 7 );
	var form = $( "#testForm17" ),
		field = $( "#testForm17text" );

	// Pre-assign error identifier
	field.attr( "aria-describedby", "testForm17text-description testForm17text-error" );
	form.validate({ errorElement: "span" });

	ok( !field.valid() );
	equal( field.attr( "aria-describedby" ), "testForm17text-description testForm17text-error" );
	assert.hasError( field, "required" );

	field.val( "foo" );
	ok( field.valid() );
	assert.noErrorFor( field );

	strictEqual( "This is where you enter your data", $("#testForm17text-description").text() );
	strictEqual( "", $("#testForm17text-error").text(), "Error label is empty for valid field" );
});

test( "test id/name containing brackets", function( assert ) {
	var form = $( "#testForm18" ),
		field = $( "#testForm18\\[text\\]" );

	form.validate({
		errorElement: "span"
	});

	form.valid();
	field.valid();
	assert.hasError( field, "required" );
});

test( "test id/name containing $", function( assert ) {
	var form = $( "#testForm19" ),
		field = $( "#testForm19\\$text" );

	form.validate({
		errorElement: "span"
	});

	field.valid();
	assert.hasError( field, "required" );
});
