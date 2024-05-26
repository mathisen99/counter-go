document.addEventListener('DOMContentLoaded', function () {
    const countElement = document.getElementById('count');
    const speed = parseInt(countElement.getAttribute('data-speed'), 10);
    const viewSpeed = parseInt(countElement.getAttribute('data-viewspeed'), 10);
    const fields = countElement.getAttribute('data-fields').split('|').filter(field => field.trim() !== '');

    let currentIndex = 0;
    let counterValue = parseInt(countElement.textContent, 10);

    function updateCounter() {
        counterValue += 1;
        countElement.textContent = counterValue;
    }

    function switchDisplay() {
        if (fields.length === 0) return;

        // Show the counter
        countElement.textContent = counterValue;
        setTimeout(() => {
            // Cycle through the fields
            for (let i = 0; i < fields.length; i++) {
                setTimeout(() => {
                    countElement.textContent = fields[i];
                }, i * viewSpeed);
            }
            // Set the counter to be shown again after all fields
            setTimeout(() => {
                countElement.textContent = counterValue;
            }, fields.length * viewSpeed);
        }, viewSpeed);
    }

    // Update the counter every 'speed' milliseconds
    setInterval(updateCounter, speed);

    // Switch display between counter and fields
    setInterval(switchDisplay, (fields.length + 1) * viewSpeed);
});