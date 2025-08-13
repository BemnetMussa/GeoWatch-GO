<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';
	import type { FireData, ChangeDetectionResult, NotificationData } from '$lib/types';
	import AnalyticsDashboard from './AnalyticsDashboard.svelte';
	import MapControls from './MapControls.svelte';
	import ChangeDetectionPanel from './ChangeDetectionPanel.svelte';
	import RealTimeSettings from './RealTimeSettings.svelte';
	import NotificationSystem from './NotificationSystem.svelte';
	import FireInfoPanel from './FireInfoPanel.svelte';

	type LeafletMap = any;
	type LeafletLayerGroup = any;
	type LeafletTileLayer = any;

	// Map state
	let mapContainer: HTMLDivElement;
	let map: LeafletMap;
	let fireLayer: LeafletLayerGroup;
	let changeLayer: LeafletLayerGroup;
	let L: any; // Leaflet will be dynamically imported

	// UI state
	let showAnalytics = $state(false);
	let activeLayer = $state('satellite');
	let showFires = $state(true);
	let showChanges = $state(false);
	let loading = $state(false);
	let changeLoading = $state(false);
	let selectedFire = $state<FireData | null>(null);

	// Data state
	let fires = $state<FireData[]>([]);
	let fireCount = $state(0);
	let changeResult = $state<ChangeDetectionResult | null>(null);
	let newFires = $state<FireData[]>([]);

	// Real-time state
	let autoRefresh = $state(false);
	let refreshIntervalMs = $state(300000); // 5 minutes
	let notifications = $state(false);
	let connectionStatus = $state<'connected' | 'disconnected' | 'error'>('disconnected');
	let lastUpdate = $state<Date | null>(null);
	let updateCount = $state(0);
	let refreshInterval: number | null = null;

	// Layer definitions - will be initialized after Leaflet loads
	let layers: Record<string, LeafletTileLayer> = {};

	onMount(async () => {
		if (browser && mapContainer) {
			try {
				// Dynamic import of Leaflet to avoid SSR issues
				const leafletModule = await import('leaflet');
				L = leafletModule.default;
				
				// Import Leaflet CSS
				await import('leaflet/dist/leaflet.css');

				// Initialize layers after Leaflet is loaded
				layers = {
					satellite: L.tileLayer('https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}', {
						attribution: '&copy; Esri &mdash; Source: Esri, i-cubed, USDA, USGS, AEX, GeoEye, Getmapping, Aerogrid, IGN, IGP, UPR-EGP, and the GIS User Community'
					}),
					terrain: L.tileLayer('https://{s}.tile.opentopomap.org/{z}/{x}/{y}.png', {
						attribution: 'Map data: &copy; OpenStreetMap contributors, SRTM | Map style: &copy; OpenTopoMap (CC-BY-SA)'
					}),
					street: L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
						attribution: '&copy; OpenStreetMap contributors'
					})
				};

				// Initialize map
				map = L.map(mapContainer, {
					center: [39.8283, -98.5795], // Center of US
					zoom: 5
				});

				// Add default layer
				layers.satellite.addTo(map);

				// Initialize fire and change layers
				fireLayer = L.layerGroup().addTo(map);
				changeLayer = L.layerGroup();

				// Load initial data
				loadFireData();

				// Set up auto-refresh if enabled
				if (autoRefresh) {
					startAutoRefresh();
				}
			} catch (error) {
				console.error('Error loading Leaflet:', error);
			}
		}
	});

	onDestroy(() => {
		if (refreshInterval) {
			clearInterval(refreshInterval);
		}
	});

	// Functions
	function toggleAnalytics() {
		showAnalytics = !showAnalytics;
	}

	function switchLayer(layerName: string) {
		if (!L || !map) return; // Guard against SSR
		
		// Remove current layer
		Object.values(layers).forEach(layer => map.removeLayer(layer));
		
		// Add new layer
		layers[layerName as keyof typeof layers].addTo(map);
		activeLayer = layerName;
	}

	function toggleFires() {
		if (!L || !map || !fireLayer) return; // Guard against SSR
		
		showFires = !showFires;
		if (showFires) {
			map.addLayer(fireLayer);
		} else {
			map.removeLayer(fireLayer);
		}
	}

	function toggleChanges() {
		if (!L || !map || !changeLayer) return; // Guard against SSR
		
		showChanges = !showChanges;
		if (showChanges) {
			map.addLayer(changeLayer);
		} else {
			map.removeLayer(changeLayer);
		}
	}

	async function loadFireData() {
		if (!L || !map) return; // Guard against SSR
		
		loading = true;
		connectionStatus = 'error';
		
		try {
			const bounds = map.getBounds();
			const response = await fetch(`/api/fires?west=${bounds.getWest()}&south=${bounds.getSouth()}&east=${bounds.getEast()}&north=${bounds.getNorth()}`);
			
			if (response.ok) {
				const data = await response.json();
				fires = data.fires;
				fireCount = fires.length;
				
				// Update fire markers
				updateFireMarkers();
				
				connectionStatus = 'connected';
				lastUpdate = new Date();
				updateCount++;
			} else {
				connectionStatus = 'disconnected';
			}
		} catch (error) {
			console.error('Error loading fire data:', error);
			connectionStatus = 'disconnected';
		} finally {
			loading = false;
		}
	}

	function updateFireMarkers() {
		if (!L || !fireLayer) return; // Guard against SSR
		
		fireLayer.clearLayers();
		
		fires.forEach(fire => {
			const color = getFireColor(fire.confidence);
			const marker = L.circleMarker([fire.latitude, fire.longitude], {
				radius: Math.max(4, Math.min(12, fire.brightness / 50)),
				fillColor: color,
				color: color,
				weight: 2,
				opacity: 0.8,
				fillOpacity: 0.6
			});

			marker.bindPopup(`
				<div class="p-2">
					<h3 class="font-semibold">Fire Hotspot</h3>
					<p><strong>Confidence:</strong> ${fire.confidence}%</p>
					<p><strong>Brightness:</strong> ${fire.brightness}K</p>
					<p><strong>Power:</strong> ${fire.frp} MW</p>
					<p><strong>Satellite:</strong> ${fire.satellite}</p>
					<p><strong>Date:</strong> ${new Date(fire.acq_date).toLocaleDateString()}</p>
				</div>
			`);

			marker.on('click', () => {
				selectedFire = fire;
			});

			fireLayer.addLayer(marker);
		});
	}

	function getFireColor(confidence: number): string {
		if (confidence >= 80) return '#dc2626'; // High confidence - red
		if (confidence >= 50) return '#ea580c'; // Medium confidence - orange
		return '#facc15'; // Low confidence - yellow
	}

	async function refreshData() {
		await loadFireData();
	}

	async function analyzeChanges(basePeriod: number, comparePeriod: number) {
		if (!L || !map) return; // Guard against SSR
		
		changeLoading = true;
		
		try {
			const bounds = map.getBounds();
			const response = await fetch('/api/change-detection', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					bounds: {
						west: bounds.getWest(),
						south: bounds.getSouth(),
						east: bounds.getEast(),
						north: bounds.getNorth()
					},
					period1: `${basePeriod}d`,
					period2: `${comparePeriod}d`
				})
			});

			if (response.ok) {
				changeResult = await response.json();
				updateChangeMarkers();
			}
		} catch (error) {
			console.error('Error analyzing changes:', error);
		} finally {
			changeLoading = false;
		}
	}

	function updateChangeMarkers() {
		if (!L || !changeLayer || !changeResult) return; // Guard against SSR
		
		changeLayer.clearLayers();
		
		const allChanges = [
			...changeResult.newFires.map(fire => ({ fire, type: 'new_fire' })),
			...changeResult.extinguishedFires.map(fire => ({ fire, type: 'extinguished_fire' })),
			...changeResult.growingFires.map(fire => ({ fire, type: 'growing_fire' })),
			...changeResult.diminishingFires.map(fire => ({ fire, type: 'diminishing_fire' }))
		];

		allChanges.forEach(change => {
			const color = getChangeColor(change.type);
			const marker = L.circleMarker([change.fire.latitude, change.fire.longitude], {
				radius: 6,
				fillColor: color,
				color: color,
				weight: 2,
				opacity: 0.9,
				fillOpacity: 0.7
			});

			marker.bindPopup(`
				<div class="p-2">
					<h3 class="font-semibold">${change.type.replace('_', ' ').toUpperCase()}</h3>
					<p><strong>Confidence:</strong> ${change.fire.confidence}%</p>
					<p><strong>Brightness:</strong> ${change.fire.brightness}K</p>
				</div>
			`);

			changeLayer.addLayer(marker);
		});
	}

	function getChangeColor(type: string): string {
		switch (type) {
			case 'new_fire': return '#10b981'; // Green
			case 'extinguished_fire': return '#6b7280'; // Gray
			case 'growing_fire': return '#f59e0b'; // Amber
			case 'diminishing_fire': return '#3b82f6'; // Blue
			default: return '#8b5cf6'; // Purple
		}
	}

	function toggleAutoRefresh() {
		autoRefresh = !autoRefresh;
		if (autoRefresh) {
			startAutoRefresh();
		} else {
			stopAutoRefresh();
		}
	}

	function startAutoRefresh() {
		if (refreshInterval) clearInterval(refreshInterval);
		refreshInterval = setInterval(loadFireData, refreshIntervalMs);
	}

	function stopAutoRefresh() {
		if (refreshInterval) {
			clearInterval(refreshInterval);
			refreshInterval = null;
		}
	}

	function changeRefreshInterval(interval: number) {
		refreshIntervalMs = interval;
		if (autoRefresh) {
			startAutoRefresh();
		}
	}

	function toggleNotifications() {
		notifications = !notifications;
		if (notifications && browser && 'Notification' in window) { // Added browser check
			Notification.requestPermission();
		}
	}

	function handleNotificationClick(notification: NotificationData) {
		if (!L || !map) return; // Guard against SSR
		
		if (notification.location) {
			map.setView([notification.location.lat, notification.location.lng], 10);
		}
	}
</script>

<div class="relative h-full w-full">
	<div bind:this={mapContainer} class="h-full w-full z-0"></div>
	
	<MapControls 
		{activeLayer}
		{showFires}
		{loading}
		{fireCount}
		{showChanges}
		changeCount={changeResult?.summary.totalNew || 0}
		onLayerChange={switchLayer}
		onToggleFires={toggleFires}
		onToggleChanges={toggleChanges}
		onRefresh={refreshData}
		onShowAnalytics={toggleAnalytics}
	/>
	
	<ChangeDetectionPanel 
		onAnalyze={analyzeChanges}
		result={changeResult}
		loading={changeLoading}
	/>
	
	<RealTimeSettings 
		{autoRefresh}
		refreshInterval={refreshIntervalMs}
		{notifications}
		{connectionStatus}
		{lastUpdate}
		{updateCount}
		onToggleAutoRefresh={toggleAutoRefresh}
		onChangeInterval={changeRefreshInterval}
		onToggleNotifications={toggleNotifications}
	/>
	
	<NotificationSystem 
		enabled={notifications}
		{newFires}
		onNotificationClick={handleNotificationClick}
	/>
	
	<AnalyticsDashboard 
		{fires}
		isOpen={showAnalytics}
		onClose={() => showAnalytics = false}
	/>
	
	{#if selectedFire}
		<FireInfoPanel 
			fire={selectedFire} 
			onClose={() => selectedFire = null}
		/>
	{/if}
</div>
