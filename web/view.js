$(document).ready(function () {
    loadUser();
});

function showReminderToast() {
    $('.toast').toast('show')
}

function showSpinner() {
    let spinner = document.getElementById("loading-spinner");
    spinner.style.display = "block";
}

function hideSpinner() {
    let spinner = document.getElementById("loading-spinner");
    spinner.style.display = "none";
}