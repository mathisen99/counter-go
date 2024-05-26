document.addEventListener("DOMContentLoaded", function() {
    // Set initial unix timestamp
    let unixTimestamp = Math.floor(Date.now());
    let countElement = document.getElementById('count');
    let logoElement = document.getElementById('logo');
    let speed = parseInt(countElement.getAttribute('data-speed'), 10);
    let viewSpeed = parseInt(countElement.getAttribute('data-viewspeed'), 10);
    let fields = countElement.getAttribute('data-fields').split('|').filter(field => field.trim() !== '');
    let counterValue = parseInt(countElement.textContent, 10);
    let currentIndex = -1;
    let isCounterVisible = true;

    function updateCounter() {
        counterValue += 1;
        if (isCounterVisible) {
            countElement.textContent = counterValue;
        }
    }

    function switchDisplay() {
        currentIndex = (currentIndex + 1) % (fields.length + 1);

        if (currentIndex === fields.length) {
            // Show the counter and logo
            isCounterVisible = true;
            countElement.textContent = counterValue;
            logoElement.style.display = 'block';
        } else {
            // Show the field text and hide the logo
            isCounterVisible = false;
            countElement.textContent = fields[currentIndex];
            logoElement.style.display = 'none';
        }
    }

    // Update the counter value every 'speed' milliseconds
    setInterval(updateCounter, speed);

    // Switch display between counter and fields every 'viewSpeed' milliseconds
    setInterval(switchDisplay, viewSpeed);

    // Additional script to update time every second
    setInterval(function() {
        let timeDifference = Math.floor(Date.now()) - unixTimestamp;

        if (Math.floor(timeDifference / speed)) {
            // Calculate new value based on time difference and speed
            counterValue += Math.floor(timeDifference / speed);

            if (isCounterVisible) {
                countElement.textContent = counterValue;
            }
            // Update unix timestamp variable
            unixTimestamp = math.floor(Date.now());
        }
    }, 1000);
});