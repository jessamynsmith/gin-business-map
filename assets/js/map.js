var markers = [];

function clearMarkers(map) {
    markers.forEach(function(marker) {
        map.removeLayer(marker);
    });

    markers = [];
}

function updateBusinesses(map, term) {
    clearMarkers(map);

    var url = "/api/v1/businesses/search/";
    var location = map.getCenter();
    fetch(`${url}?term=${term}&latitude=${location.lat}&longitude=${location.lng}`)
        .then(response => {
            response.json().then(responseJson => {
                if (responseJson.businesses) {

                    for (let i = 0; i < responseJson.businesses.length; i++) {
                        let business = responseJson.businesses[i];
                        let markerCoords = [business.coordinates.latitude, business.coordinates.longitude];
                        var marker = L.marker(markerCoords).addTo(map);
                        markers.push(marker);
                        var businessDesc = `${business.rating}`;
                        if (business.price) {
                            businessDesc = `${businessDesc} - ${business.price}`;
                        }
                        marker.bindPopup(`<b>${business.name}</b><br>${businessDesc}`);
                    }
                }
            });
        });
}


addEventListener("DOMContentLoaded",  (event) => {
    let coords = [51.0271596,-114.4174673];

    var map = L.map('map').setView(coords, 13);

    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 19,
        attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
    }).addTo(map);

    if ("geolocation" in navigator) {
        navigator.geolocation.getCurrentPosition((position) => {
            coords = [position.coords.latitude, position.coords.longitude];
            map = map.setView(coords, 13);
        });
    }

    document.getElementById("searchForm").addEventListener('submit', function(evt) {
        evt.preventDefault();

        let term = document.getElementById("search").value;
        updateBusinesses(map, term);
    });

});
