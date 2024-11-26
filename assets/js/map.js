
addEventListener("DOMContentLoaded",  (event) => {
    let coords = [];

    var map = L.map('map').setView([51.0271596,-114.4174673], 13);

    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 19,
        attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
    }).addTo(map);

    if ("geolocation" in navigator) {
        navigator.geolocation.getCurrentPosition((position) => {
            map = map.setView([position.coords.latitude, position.coords.longitude], 13);
        });
    }

});
