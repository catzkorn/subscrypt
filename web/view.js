$(document).ready(function () {
    // $("#reminder-alert").hide();
    // $("#reminder-button").click($("#reminder-alert").alert());
    loadUser();
});

function showReminderToast() {
    $('.toast').toast('show')
}