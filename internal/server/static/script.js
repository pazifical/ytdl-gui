const baseUrl = "/download?url="

const downloadButton = document.getElementById("btn-download");
downloadButton.addEventListener("click", () => {
    const videoUrl = document.getElementById("txt-url").value;
    console.log(videoUrl);

    fetch(baseUrl + videoUrl).then(d => console.log(d));
    window.location.reload();
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