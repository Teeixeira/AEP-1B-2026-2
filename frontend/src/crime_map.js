import L from 'leaflet';
import markerIconUrl from "leaflet/dist/images/marker-icon.png";
import markerIconRetinaUrl from "leaflet/dist/images/marker-icon-2x.png";
import markerShadowUrl from "leaflet/dist/images/marker-shadow.png";

export async function setupMap(element) {
    L.Icon.Default.prototype.options.iconUrl = markerIconUrl;
    L.Icon.Default.prototype.options.iconRetinaUrl = markerIconRetinaUrl;
    L.Icon.Default.prototype.options.shadowUrl = markerShadowUrl;
    L.Icon.Default.imagePath = ""; 

    const map = L.map(element).setView([-23.4213, -51.9331], 13);
    const points = L.layerGroup().addTo(map);
    const locationButton = document.querySelector('#use-location');
    const testButton = document.querySelector('#add-test-points');
    let userMarker;
    let location;

    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 19,
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    function useLocation() {
        if (!navigator.geolocation) {
            alert('Seu navegador não suporta localização.');
            return;
        }

        locationButton.disabled = true;
        testButton.disabled = true;
        locationButton.textContent = 'Buscando localização...';

        navigator.geolocation.getCurrentPosition(
            async (position) => {
                location = {
                    lat: position.coords.latitude,
                    lng: position.coords.longitude
                };

                map.setView([location.lat, location.lng], 15);

                if (userMarker) {
                    userMarker.setLatLng([location.lat, location.lng]);
                } else {
                    userMarker = L.circleMarker([location.lat, location.lng], {
                        radius: 8,
                        color: '#ffffff',
                        fillColor: '#1976d2',
                        fillOpacity: 1,
                        weight: 3
                    }).addTo(map).bindTooltip('Você está aqui');
                }

                try {
                    await getNear(points, location, 20000);
                } catch (error) {
                    alert('Não foi possível buscar os pontos. Tente novamente.');
                } finally {
                    locationButton.disabled = false;
                    testButton.disabled = false;
                    locationButton.textContent = 'Usar minha localização';
                }
            },
            (error) => {
                alert('Não foi possível obter sua localização. Confira a permissão de localização do navegador e tente novamente.');
                locationButton.disabled = false;
                testButton.disabled = !location;
                locationButton.textContent = 'Usar minha localização';
            },
            { enableHighAccuracy: true, timeout: 15000, maximumAge: 0 }
        );
    }

    locationButton.addEventListener('click', useLocation);

    testButton.addEventListener('click', async () => {
        testButton.disabled = true;
        locationButton.disabled = true;
        testButton.textContent = 'Adicionando...';
        const offsets = [[0.002, 0.002], [-0.002, 0.001], [0.001, -0.002]];

        try {
            for (let i = 0; i < offsets.length; i++) {
                const response = await fetch('/api/crimes', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        tipo: `Ponto de teste ${i + 1}`,
                        descricao: 'Registro fictício para testar o arraste no mapa.',
                        data_hora: new Date().toISOString(),
                        localizacao: {
                            type: 'Point',
                            coordinates: [location.lng + offsets[i][0], location.lat + offsets[i][1]]
                        }
                    })
                });

                if (!response.ok) {
                    throw new Error('Erro ao criar ponto de teste');
                }
            }

            await getNear(points, location, 20000);
        } catch (error) {
            alert('Não foi possível adicionar ou carregar todos os pontos de teste.');
        } finally {
            testButton.disabled = false;
            locationButton.disabled = false;
            testButton.textContent = 'Adicionar pontos de teste';
        }
    });

    useLocation();
}

export async function getNear(points, location, radiusInMeters) {
    const response = await fetch(`/api/crimes/proximos?lat=${location.lat}&lng=${location.lng}&raio=${radiusInMeters}`)

    if (!response.ok) {
        throw new Error('Erro ao buscar pontos');
    }

    const crimes = await response.json()
    points.clearLayers();

    (crimes || []).forEach(crime => {
        
        const marker = L.marker([crime.localizacao.coordinates[1], crime.localizacao.coordinates[0]], {
            draggable: true
        }).addTo(points)
            .bindPopup(
                `<div>
                    <h1>${crime.tipo}</h1>
                    <p>${crime.descricao}</p>
                    <p>${(new Date(crime.data_hora)).toLocaleString()}</p>
                </div>`
            )

        marker.on('dragend', async () => {
            const coordinates = crime.localizacao.coordinates;
            const position = marker.getLatLng();
            marker.dragging.disable();
            crime.localizacao.coordinates = [position.lng, position.lat];

            try {
                const response = await fetch(`/api/crimes/${crime.id}`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(crime)
                });

                if (!response.ok) {
                    throw new Error('Erro ao atualizar localização');
                }
            } catch (error) {
                crime.localizacao.coordinates = coordinates;
                marker.setLatLng([coordinates[1], coordinates[0]]);
                alert('Não foi possível salvar a nova localização. Tente novamente.');
            } finally {
                marker.dragging.enable();
            }
        });
    });
}
