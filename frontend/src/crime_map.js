import L from 'leaflet';
import markerIconUrl from 'leaflet/dist/images/marker-icon.png';
import markerIconRetinaUrl from 'leaflet/dist/images/marker-icon-2x.png';
import markerShadowUrl from 'leaflet/dist/images/marker-shadow.png';
import { APIError, parseAPIError } from './api_error.js';

export function setupMap(element) {
    L.Icon.Default.prototype.options.iconUrl = markerIconUrl;
    L.Icon.Default.prototype.options.iconRetinaUrl = markerIconRetinaUrl;
    L.Icon.Default.prototype.options.shadowUrl = markerShadowUrl;
    L.Icon.Default.imagePath = '';

    const map = L.map(element).setView([-23.4213, -51.9331], 13);
    const points = L.featureGroup().addTo(map);
    const form = document.querySelector('#search-form');
    const radiusInput = document.querySelector('#search-radius');
    const searchButton = document.querySelector('#search-button');
    const locationButton = document.querySelector('#use-location');
    const testButton = document.querySelector('#add-test-points');
    const retryButton = document.querySelector('#retry-search');
    const status = document.querySelector('#search-status');
    let userMarker;
    let location;
    let radiusCircle;
    let radiusTimeout;

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

    function showRadius(radiusInMeters) {
        clearTimeout(radiusTimeout);
        if (radiusCircle) radiusCircle.remove();

        radiusCircle = L.circle([location.lat, location.lng], {
            radius: radiusInMeters,
            color: '#0071e3',
            weight: 1.5,
            fillOpacity: 0.12,
            interactive: false
        }).addTo(map);
        radiusCircle.bringToBack();
        map.fitBounds(radiusCircle.getBounds(), { padding: [32, 32], maxZoom: 16 });
        radiusTimeout = setTimeout(() => radiusCircle.remove(), 2500);
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
        showRadius(radiusInMeters);
        showStatus('Carregando ocorrências...', 'loading');

        try {
            const count = await getNear(points, location, radiusInMeters);
            const distance = radiusKm.toLocaleString('pt-BR', { maximumSignificantDigits: 15 });
            if (count === 0) {
                showStatus(`Nenhuma ocorrência encontrada em um raio de ${distance} km. Tente aumentar a distância.`, 'empty');
            } else {
                const result = count === 1 ? '1 ocorrência encontrada' : `${count} ocorrências encontradas`;
                showStatus(`${result} em um raio de ${distance} km.`, 'success');
            }
        } catch (error) {
            const userMsg = error instanceof APIError
                ? error.toUserMessage()
                : 'Não foi possível buscar as ocorrências. Tente novamente.';
            const retryable = error instanceof APIError ? error.isRetryable() : true;
            showStatus(userMsg, 'error', retryable);
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

    testButton.addEventListener('click', async () => {
        setBusy(true);
        showStatus('Adicionando pontos de teste...', 'loading');
        const offsets = [[0.002, 0.002], [-0.002, 0.001], [0.001, -0.002]];

        try {
            for (let i = 0; i < offsets.length; i++) {
                let response;
                try {
                    response = await fetch('/api/crimes', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            tipo: `Ponto de teste ${i + 1}`,
                            descricao: 'Registro fictício para testar o arraste no mapa.',
                            data_hora: new Date().toISOString(),
                            localizacao: {
                                latitude: location.lat + offsets[i][1],
                                longitude: location.lng + offsets[i][0]
                            }
                        })
                    });
                } catch (netErr) {
                    throw new APIError('ERRO_CONEXAO', 'Falha ao conectar com o servidor.', 0, netErr.message);
                }

                if (!response.ok) {
                    throw await parseAPIError(response, `Erro ao criar ponto de teste ${i + 1}`);
                }
            }
        } catch (error) {
            const userMsg = error instanceof APIError
                ? error.toUserMessage()
                : 'Não foi possível adicionar todos os pontos de teste.';
            showStatus(userMsg, 'error', error instanceof APIError ? error.isRetryable() : true);
            return;
        } finally {
            setBusy(false);
        }

        await searchCrimes();
    });

    useLocation();
}

export async function getNear(points, location, radiusInMeters) {
    points.clearLayers();
    let response;
    try {
        response = await fetch(`/api/crimes/proximos?lat=${location.lat}&lng=${location.lng}&raio=${radiusInMeters}`);
    } catch (networkError) {
        throw new APIError('ERRO_CONEXAO', 'Falha ao conectar com o servidor.', 0, networkError.message);
    }

    if (!response.ok) {
        throw await parseAPIError(response, 'Erro ao buscar ocorrências próximas.');
    }

    const crimes = (await response.json()) || [];

    crimes.forEach(crime => {
        const marker = L.marker([crime.localizacao.latitude, crime.localizacao.longitude], {
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
            const originalLat = crime.localizacao.latitude;
            const originalLng = crime.localizacao.longitude;
            const position = marker.getLatLng();
            marker.dragging.disable();
            crime.localizacao.latitude = position.lat;
            crime.localizacao.longitude = position.lng;

            try {
                let response;
                try {
                    response = await fetch(`/api/crimes/${crime.id}`, {
                        method: 'PUT',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(crime)
                    });
                } catch (netErr) {
                    throw new APIError('ERRO_CONEXAO', 'Falha ao conectar com o servidor.', 0, netErr.message);
                }

                if (!response.ok) {
                    throw await parseAPIError(response, 'Erro ao atualizar localização da ocorrência.');
                }
            } catch (error) {
                crime.localizacao.latitude = originalLat;
                crime.localizacao.longitude = originalLng;
                marker.setLatLng([originalLat, originalLng]);
                const msg = error instanceof APIError
                    ? error.toUserMessage()
                    : 'Não foi possível salvar a nova localização. Tente novamente.';
                alert(msg);
            } finally {
                marker.dragging.enable();
            }
        });
    });

    return crimes.length;
}
