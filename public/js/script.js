document.addEventListener("DOMContentLoaded", function() {
    let unixTimestamp = Date.now();
    let countElement = document.getElementById('count');
    let logoElement = document.getElementById('logo');
    let speed = parseInt(countElement.getAttribute('data-speed'), 10);
    let viewSpeed = parseInt(countElement.getAttribute('data-viewspeed'), 10);
    let timeCounterDisplay = parseInt(countElement.getAttribute('data-time-counter-display'), 10);
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
    setTimeout(switchDisplay, timeCounterDisplay*2);
});
