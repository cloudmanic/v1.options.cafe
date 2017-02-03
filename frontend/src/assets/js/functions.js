;(function($, window, document, undefined) {
	var $win = $(window);
	var $doc = $(document);

	function contentHeight() {
		var contentHeadH = $('.content .content-head').outerHeight();

		$('.content').css('padding-top', contentHeadH);
	}

	$doc.ready(function() {
		
	});

	$win.on('load resize', function() {
		contentHeight();
	});

})(jQuery, window, document);
