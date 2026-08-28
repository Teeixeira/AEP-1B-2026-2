import L from 'leaflet';
import markerIconUrl from "leaflet/dist/images/marker-icon.png";
import markerIconRetinaUrl from "leaflet/dist/images/marker-icon-2x.png";
import markerShadowUrl from "leaflet/dist/images/marker-shadow.png";

export async function setupMap(element) {
    L.Icon.Default.prototype.options.iconUrl = markerIconUrl;
    L.Icon.Default.prototype.options.iconRetinaUrl = markerIconRetinaUrl;
    L.Icon.Default.prototype.options.shadowUrl = markerShadowUrl;
    L.Icon.Default.imagePath = ""; 

    const map = L.map('map').setView([-23.4213, -51.9331], 13);

    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 19,
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(
            async (position) => {
                const location = {
                    lat: position.coords.latitude,
                    lng: position.coords.longitude
                };

                await getNear(map, location, 20000);
            },
            (error) => {
                console.error(`Error getting location: ${error.message}`);
            }
        );
    } else {
        console.error("Geolocation is not supported by this browser.");
    }
}

export async function getNear(map, location, radiusInMeters) {
    const response = await fetch(`/api/crimes/proximos?lat=${location.lat}&lng=${location.lng}&raio=${radiusInMeters}`)

    let crimes = await response.json()

    crimes.forEach(crime => {
        console.log(crime);
        
        L.marker([crime.localizacao.coordinates[1], crime.localizacao.coordinates[0]]).addTo(map)
            .bindPopup(
                `<div>
                    <h1>${crime.tipo}</h1>
                    <p>${crime.descricao}</p>
                    <p>${(new Date(crime.data_hora)).toLocaleString()}</p>
                </div>`
            )
    });
}