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
    fetchDownloadItems();
}, 2000);


function fetchDownloadItems() {
    fetch("/items")
        .then(response => response.json())
        .then(response => {
            setDownloadItems(response.items);
        });
}

function setDownloadItems(downloadItems) {
    const itemsList = document.getElementById("download-items-list");
    itemsList.innerHTML = "";

    Object.keys(downloadItems).forEach(function(key, index) {
        const item = downloadItems[key];
        const li = document.createElement('li'); 
        const text =  item.Title + " (Status: " + item.Status + ")";
        li.appendChild(document.createTextNode(text)); 
        itemsList.appendChild(li); 
    });
}
