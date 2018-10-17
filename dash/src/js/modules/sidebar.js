import SimpleBar from 'simplebar'

$(function(){
    /* Sidebar toggle behaviour */
    $('.sidebar-toggle').on('click', function () {
        $('.sidebar').toggleClass('toggled');
    });

    const active = $('.sidebar .active');

    if (active.length && active.parent('.collapse').length) {
        const parent = active.parent('.collapse');

        parent.prev('a').attr('aria-expanded', true);
        parent.addClass('show');
    }

    /* Scrollable sidebar */
    if($('.sidebar-sticky').length > 0) {
        new SimpleBar($('.sidebar-sticky .sidebar-content')[0])
    }
});
