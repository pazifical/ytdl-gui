const baseUrl = "/download?url="

const downloadButton = document.getElementById("btn-download");
downloadButton.addEventListener("click", () => {
    const videoUrl = document.getElementById("txt-url").value;
    console.log(videoUrl);

    fetch(baseUrl+videoUrl).then(d => console.log(d));
});