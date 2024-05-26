document.addEventListener("DOMContentLoaded", function() {
    let unixTimestamp = Date.now();
    let countElement = document.getElementById('count');
    let logoElement = document.getElementById('logo');
    let logoH1Element = logoElement.querySelector('h1'); // Select the h1 inside the logo element
    let speed = parseInt(countElement.getAttribute('data-speed'), 10);
    let viewSpeed = parseInt(countElement.getAttribute('data-viewspeed'), 10);
    let timeCounterDisplay = parseInt(countElement.getAttribute('data-time-counter-display'), 10);
    let fontSize = parseInt(countElement.getAttribute('data-fontsize'), 10);
    let logoFontSize = parseInt(countElement.getAttribute('data-logofontsize'), 10);
    let logoText = countElement.getAttribute('data-logotext');
    let fields = countElement.getAttribute('data-fields').split('|').filter(field => field.trim() !== '');
    let counterValue = parseInt(countElement.textContent, 10);
    let currentIndex = -1;
    let isCounterVisible = true;

    // Apply the chosen font size
    countElement.style.fontSize = fontSize + 'px';
    logoH1Element.style.fontSize = logoFontSize + 'px'; // Apply font size to the h1 inside logo element
    logoH1Element.textContent = logoText; // Set the logo text

    function updateCounter() {
        counterValue += 1;
        if (isCounterVisible) {
            countElement.textContent = counterValue;
        }
    }

    function switchDisplay() {
        currentIndex = (currentIndex + 1) % (fields.length + 1);

        if (currentIndex === fields.length) {
            isCounterVisible = true;
            countElement.textContent = counterValue;
            logoElement.style.display = 'block';
            setTimeout(switchDisplay, timeCounterDisplay);
        } else {
            isCounterVisible = false;
            let fieldContent = fields[currentIndex];
            const imgTags = fieldContent.match(/IMG=([^\s]+)/g) || [];
            const sizeTags = fieldContent.match(/SIZE=([^\s]+)/g) || [];

            imgTags.forEach((imgTag, index) => {
                const imgUrl = imgTag.split('=')[1];
                const imgSize = sizeTags[index] ? sizeTags[index].split('=')[1] : '100%';
                const imgHTML = `<img src="/images/${imgUrl}" style="max-width:${imgSize}; vertical-align: middle;">`;
                fieldContent = fieldContent.replace(imgTag, imgHTML);
            });

            sizeTags.forEach(sizeTag => {
                fieldContent = fieldContent.replace(sizeTag, ''); // Remove the SIZE parameter
            });

            countElement.innerHTML = fieldContent;
            logoElement.style.display = 'none';
            setTimeout(switchDisplay, viewSpeed);
        }
    }

    setInterval(updateCounter, speed);
    setTimeout(switchDisplay, timeCounterDisplay); // Start the switch display after the counter has been displayed twice
});
