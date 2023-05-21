const baseUrl = "/download?url="

const downloadButton = document.getElementById("btn-download");
downloadButton.addEventListener("click", () => {
    const element = document.getElementById("txt-url");
    const videoUrl = element.value;
    console.log(videoUrl);

    fetch(baseUrl + videoUrl).then(d => console.log(d));
    element.value = "";
});

function setStatus() {
    fetch("/status")
        .then(response => response.json())
        .then(response => {
            document.getElementById("status").innerText = response.message;
        });
}

var intervalId = window.setInterval(function () {
    setStatus();
}, 2000);