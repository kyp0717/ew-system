function openTab(event, tabId) {
    // Hide all tab content
    var tabContents = document.getElementsByClassName('tab-content');
    for (var i = 0; i < tabContents.length; i++) {
        tabContents[i].classList.remove('active');
    }

    // Remove active class from all tabs
    var tabs = document.getElementsByClassName('tab');
    for (var i = 0; i < tabs.length; i++) {
        tabs[i].classList.remove('active');
    }

    // Show current tab content
    document.getElementById(tabId).classList.add('active');

    // Set clicked tab as active
    event.currentTarget.classList.add('active');
}

// Set default active tab
document.addEventListener('DOMContentLoaded', function() {
    var firstTab = document.getElementsByClassName('tab')[0];
    if (firstTab) {
        firstTab.click();
    }
});
