// Initialize jQuery Document Ready
$(document).ready(function() {
    console.log("Samba4 Admin Panel Initialized");

    // CSRF Setup for AJAX calls
    var csrfToken = $('meta[name="csrf-token"]').attr('content');
    $.ajaxSetup({
        headers: {
            'X-CSRF-Token': csrfToken
        }
    });

    // Simple Modal triggers
    $(document).on('click', '.js-modal-trigger', function(e) {
        e.preventDefault();
        var target = $(this).data('target');
        $(target).removeClass('hidden');
    });

    $(document).on('click', '.js-modal-close, .modal-overlay', function(e) {
        if ($(e.target).closest('.modal-content').length === 0 || $(e.target).hasClass('js-modal-close')) {
            $('.modal-overlay').addClass('hidden');
        }
    });

    // Form confirmation
    $(document).on('submit', '.js-confirm-form', function(e) {
        var msg = $(this).data('confirm') || "Are you sure?";
        if (!confirm(msg)) {
            e.preventDefault();
        }
    });
});
