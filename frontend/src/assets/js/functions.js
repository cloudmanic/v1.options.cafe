;(function($, window, document, undefined) {
	var $win = $(window);
	var $doc = $(document);

	// Set Padding Top of Content Block
	function setContentPadding() {
		var $content = $('.content');
		var padding;

		if ($win.width() < 1320) {
			padding = $('.banner--fixed').length ? $('.header-wrapper').outerHeight() + 50 : $('.header-wrapper').outerHeight() + 20;
		} else {
			padding = $('.header-wrapper').outerHeight() + 20;
		}

		if ($content.length) {
			$content.css('padding-top', padding);
		}
	}

	$doc.ready(function() {
		// Mobile Menu
		$('.btn-menu').on('click', function(event) {
			$(this).toggleClass('active');
			$('.sidebar, .nav').toggleClass('open');
			$('.content, .header').toggleClass('moved');

			if ($win.width() < 768) {
				$('.sidebar .sidebar__content').scrollLock();
			}

			event.preventDefault();
		});

		// Popup
		$('.popup--trigger').magnificPopup({
			type: 'ajax'
		});

		$('body').on('click', '.popup--close', function(event) {
			$.magnificPopup.close();

			event.preventDefault();
		});

		// Set Padding Top of Content Block
		setContentPadding();
	});

	$win.on('resize', function() {
		setContentPadding();
	});

})(jQuery, window, document);
