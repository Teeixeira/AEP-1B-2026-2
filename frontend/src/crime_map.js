import L from 'leaflet';
import { HeatLayer } from 'leaflet.heat';
import markerIconUrl from 'leaflet/dist/images/marker-icon.png';
import markerIconRetinaUrl from 'leaflet/dist/images/marker-icon-2x.png';
import markerShadowUrl from 'leaflet/dist/images/marker-shadow.png';

export function setupMap(element) {
    L.Icon.Default.prototype.options.iconUrl = markerIconUrl;
    L.Icon.Default.prototype.options.iconRetinaUrl = markerIconRetinaUrl;
    L.Icon.Default.prototype.options.shadowUrl = markerShadowUrl;
    L.Icon.Default.imagePath = '';

    const map = L.map(element).setView([-23.4213, -51.9331], 13);
    const points = L.featureGroup().addTo(map);
    const heat = new HeatLayer([], {
        radius: 45,
        blur: 30,
        max: 4,
        maxZoom: 13,
        minOpacity: 0.2,
        gradient: { 0.2: '#ffd166', 0.5: '#f59e36', 0.75: '#ed633c', 1: '#c62828' }
    }).addTo(map);
    const heatToggle = document.querySelector('#show-heat');
    const heatLegend = document.querySelector('#heat-legend');
    const form = document.querySelector('#search-form');
    const radiusInput = document.querySelector('#search-radius');
    const searchButton = document.querySelector('#search-button');
    const locationButton = document.querySelector('#use-location');
    const testButton = document.querySelector('#add-test-points');
    const retryButton = document.querySelector('#retry-search');
    const status = document.querySelector('#search-status');
    let userMarker;
    let location;

    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 19,
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    function showStatus(message, type = 'info', retry = false) {
        status.textContent = message;
        status.dataset.type = type;
        retryButton.hidden = !retry;
    }

    function setBusy(busy) {
        locationButton.disabled = busy;
        radiusInput.disabled = busy;
        searchButton.disabled = busy;
        testButton.disabled = busy || !location;
        retryButton.disabled = busy;
        form.setAttribute('aria-busy', String(busy));
    }

    async function searchCrimes() {
        const radiusKm = Number(radiusInput.value);
        const radiusInMeters = radiusKm * 1000;

        if (!Number.isFinite(radiusInMeters) || radiusInMeters <= 0) {
            showStatus('Informe um raio maior que zero, em quilômetros.', 'error');
            radiusInput.focus();
            return;
        }

        if (!location) {
            useLocation();
            return;
        }

        setBusy(true);
        showStatus('Carregando ocorrências...', 'loading');

        try {
            const count = await getNear(points, location, radiusInMeters, heat);
            const distance = radiusKm.toLocaleString('pt-BR', { maximumSignificantDigits: 15 });
            if (count === 0) {
                showStatus(`Nenhuma ocorrência encontrada em um raio de ${distance} km. Tente aumentar a distância.`, 'empty');
            } else {
                map.fitBounds(points.getBounds().extend([location.lat, location.lng]), {
                    padding: [32, 32],
                    maxZoom: 15
                });
                const result = count === 1 ? '1 ocorrência encontrada' : `${count} ocorrências encontradas`;
                showStatus(`${result} em um raio de ${distance} km.`, 'success');
            }
        } catch (error) {
            showStatus('Não foi possível buscar as ocorrências. Tente novamente.', 'error', true);
        } finally {
            setBusy(false);
        }
    }

    function useLocation() {
        if (!navigator.geolocation) {
            showStatus('Seu navegador não suporta localização.', 'error');
            return;
        }

        setBusy(true);
        showStatus('Buscando sua localização...', 'loading');

        navigator.geolocation.getCurrentPosition(
            async (position) => {
                location = { lat: position.coords.latitude, lng: position.coords.longitude };
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

                setBusy(false);
                await searchCrimes();
            },
            (error) => {
                const message = error.code === 1
                    ? 'Permita o acesso à localização no navegador para buscar ocorrências perto de você.'
                    : 'Não foi possível obter sua localização. Tente novamente.';
                showStatus(message, 'error', true);
                setBusy(false);
            },
            { enableHighAccuracy: true, timeout: 15000, maximumAge: 0 }
        );
    }

    form.addEventListener('submit', async (event) => {
        event.preventDefault();
        await searchCrimes();
    });
    locationButton.addEventListener('click', useLocation);
    retryButton.addEventListener('click', () => location ? searchCrimes() : useLocation());
    heatToggle.addEventListener('change', () => {
        if (heatToggle.checked) {
            heat.addTo(map);
        } else {
            map.removeLayer(heat);
        }
        heatLegend.hidden = !heatToggle.checked;
    });

    testButton.addEventListener('click', async () => {
        setBusy(true);
        showStatus('Adicionando pontos de teste...', 'loading');
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
                if (!response.ok) throw new Error('Erro ao criar ponto de teste');
            }
        } catch (error) {
            showStatus('Não foi possível adicionar todos os pontos de teste. Atualize a busca antes de tentar novamente.', 'error', true);
            return;
        } finally {
            setBusy(false);
        }

        await searchCrimes();
    });

    useLocation();
}

function updateHeat(points, heat) {
    heat.setLatLngs(points.getLayers().map(marker => marker.getLatLng()));
}

export async function getNear(points, location, radiusInMeters, heat) {
    points.clearLayers();
    updateHeat(points, heat);
    const response = await fetch(`/api/crimes/proximos?lat=${location.lat}&lng=${location.lng}&raio=${radiusInMeters}`);
    if (!response.ok) throw new Error('Erro ao buscar ocorrências');
    const crimes = (await response.json()) || [];

    crimes.forEach(crime => {
        const marker = L.marker([crime.localizacao.coordinates[1], crime.localizacao.coordinates[0]], {
            draggable: true
        }).addTo(points);

        const popup = document.createElement('div');
        const title = document.createElement('h3');
        const description = document.createElement('p');
        const date = document.createElement('p');
        title.textContent = crime.tipo;
        description.textContent = crime.descricao;
        date.textContent = new Date(crime.data_hora).toLocaleString('pt-BR');
        popup.append(title, description, date);
        marker.bindPopup(popup);

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
                if (!response.ok) throw new Error('Erro ao atualizar localização');
            } catch (error) {
                crime.localizacao.coordinates = coordinates;
                marker.setLatLng([coordinates[1], coordinates[0]]);
                alert('Não foi possível salvar a nova localização. Tente novamente.');
            } finally {
                marker.dragging.enable();
                updateHeat(points, heat);
            }
        });
    });

    updateHeat(points, heat);
    return crimes.length;
}
