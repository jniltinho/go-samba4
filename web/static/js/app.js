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

    // ── Generic modal triggers ────────────────────────────────────────────────
    $(document).on('click', '.js-modal-trigger', function(e) {
        e.preventDefault();
        var target = $(this).attr('data-target');
        $(target).removeClass('hidden');
    });

    $(document).on('click', '.js-modal-close', function() {
        $(this).closest('.modal-overlay').addClass('hidden');
    });

    // Close on overlay background click
    $(document).on('click', '.modal-overlay', function(e) {
        if ($(e.target).hasClass('modal-overlay')) {
            $(this).addClass('hidden');
        }
    });

    // ── Delete trigger from the user list ─────────────────────────────────────
    // Populates the shared #modal-delete-list with the correct username/action
    $(document).on('click', '.js-delete-trigger', function(e) {
        e.preventDefault();
        // Use .attr() for reliable DOM attribute reading in jQuery 4
        var sam = $(this).attr('data-sam');

        $('#modal-delete-sam-label').text(sam);
        $('#form-delete-list').attr('action', '/users/' + sam + '/delete');
        $('#confirm_name_list')
            .val('')
            .attr('data-confirm-value', sam)
            .attr('placeholder', 'Type ' + sam + ' to confirm');
        $('#btn-confirm-delete-list').prop('disabled', true);

        $('#modal-delete-list').removeClass('hidden');
    });

    // ── Confirm-name inputs ───────────────────────────────────────────────────
    // List page modal input
    $(document).on('input', '#confirm_name_list', function() {
        var expected = $(this).attr('data-confirm-value');
        $('#btn-confirm-delete-list').prop('disabled', $(this).val() !== expected);
    });

    // Detail page modal input
    $(document).on('input', '#confirm_name', function() {
        var expected = $(this).attr('data-confirm-value');
        $('#btn-confirm-delete').prop('disabled', $(this).val() !== expected);
    });

    // ── Auto-fill UPN from username (new user form) ───────────────────────────
    $(document).on('input', '#SAMAccountName', function() {
        var domain = $(this).attr('data-upn-domain');
        if (!domain) return;
        var sam = $(this).val();
        $('#UserPrincipalName').val(sam ? sam + '@' + domain : '');
    });

    // ── General form confirmation (browser native) ─────────────────────────────
    $(document).on('submit', '.js-confirm-form', function(e) {
        var msg = $(this).attr('data-confirm') || "Are you sure?";
        if (!confirm(msg)) {
            e.preventDefault();
        }
    });
});
