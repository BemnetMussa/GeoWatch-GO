<script lang="ts">
	import { Settings, Wifi, WifiOff, Bell, BellOff, Clock, Activity } from 'lucide-svelte';
	
	interface Props {
		autoRefresh: boolean;
		refreshInterval: number;
		notifications: boolean;
		connectionStatus: 'connected' | 'disconnected' | 'syncing';
		lastUpdate: Date | null;
		updateCount: number;
		onToggleAutoRefresh: () => void;
		onChangeInterval: (interval: number) => void;
		onToggleNotifications: () => void;
	}
	
	let { 
		autoRefresh,
		refreshInterval,
		notifications,
		connectionStatus,
		lastUpdate,
		updateCount,
		onToggleAutoRefresh,
		onChangeInterval,
		onToggleNotifications
	}: Props = $props();
	
	let isOpen = $state(false);
	
	function formatLastUpdate(date: Date | null) {
		if (!date) return 'Never';
		const now = new Date();
		const diff = now.getTime() - date.getTime();
		const minutes = Math.floor(diff / 60000);
		const seconds = Math.floor((diff % 60000) / 1000);
		
		if (minutes > 0) return `${minutes}m ${seconds}s ago`;
		return `${seconds}s ago`;
	}
	
	function getConnectionIcon() {
		switch (connectionStatus) {
			case 'connected': return Wifi;
			case 'disconnected': return WifiOff;
			case 'syncing': return Activity;
		}
	}
	
	function getConnectionColor() {
		switch (connectionStatus) {
			case 'connected': return 'text-green-600';
			case 'disconnected': return 'text-red-600';
			case 'syncing': return 'text-orange-600';
		}
	}
</script>

<div class="absolute bottom-4 right-4 z-10">
	<div class="bg-white rounded-lg shadow-lg border border-slate-200 overflow-hidden">
		<button
			onclick={() => isOpen = !isOpen}
			class="flex items-center space-x-2 px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 transition-colors"
		>
			<Settings class="w-4 h-4" />
			<span>Real-time</span>
			<div class="flex items-center space-x-1 ml-2">
				{#if connectionStatus === 'syncing'}
					<Activity class="w-3 h-3 {getConnectionColor()} animate-pulse" />
				{:else}
					{@const SvelteComponent = getConnectionIcon()}
					<SvelteComponent class="w-3 h-3 {getConnectionColor()}" />
				{/if}
			</div>
		</button>
		
		{#if isOpen}
			{@const SvelteComponent_1 = getConnectionIcon()}
			<div class="border-t border-slate-200 p-4 w-72">
				<!-- Auto Refresh Toggle -->
				<div class="flex items-center justify-between mb-4">
					<div class="flex items-center space-x-2">
						<Clock class="w-4 h-4 text-slate-400" />
						<span class="text-sm font-medium text-slate-700">Auto Refresh</span>
					</div>
					<button
						onclick={onToggleAutoRefresh}
						class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors {autoRefresh ? 'bg-orange-600' : 'bg-slate-200'}"
					>
						<span class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform {autoRefresh ? 'translate-x-6' : 'translate-x-1'}"></span>
					</button>
				</div>
				
				<!-- Refresh Interval -->
				{#if autoRefresh}
					<div class="mb-4">
						<label for="refreshInterval" class="block text-sm font-medium text-slate-700 mb-2">
							Update Interval
						</label>
						<select 
							id="refreshInterval" 
							value={refreshInterval} 
							onchange={(e: any) => onChangeInterval(Number(e.target.value))}
							class="w-full px-3 py-2 border border-slate-300 rounded-md text-sm"
						>
							<option value={30000}>30 seconds</option>
							<option value={60000}>1 minute</option>
							<option value={300000}>5 minutes</option>
							<option value={600000}>10 minutes</option>
							<option value={1800000}>30 minutes</option>
						</select>
					</div>
				{/if}
				
				<!-- Notifications Toggle -->
				<div class="flex items-center justify-between mb-4">
					<div class="flex items-center space-x-2">
						{#if notifications}
							<Bell class="w-4 h-4 text-slate-400" />
						{:else}
							<BellOff class="w-4 h-4 text-slate-400" />
						{/if}
						<span class="text-sm font-medium text-slate-700">Notifications</span>
					</div>
					<button
						onclick={onToggleNotifications}
						class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors {notifications ? 'bg-orange-600' : 'bg-slate-200'}"
					>
						<span class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform {notifications ? 'translate-x-6' : 'translate-x-1'}"></span>
					</button>
				</div>
				
				<!-- Status Information -->
				<div class="pt-4 border-t border-slate-100 space-y-2">
					<div class="flex items-center justify-between text-xs">
						<span class="text-slate-500">Connection</span>
						<div class="flex items-center space-x-1">
							<SvelteComponent_1 class="w-3 h-3 {getConnectionColor()}" />
							<span class="capitalize {getConnectionColor()}">{connectionStatus}</span>
						</div>
					</div>
					
					<div class="flex items-center justify-between text-xs">
						<span class="text-slate-500">Last Update</span>
						<span class="text-slate-700">{formatLastUpdate(lastUpdate)}</span>
					</div>
					
					<div class="flex items-center justify-between text-xs">
						<span class="text-slate-500">Updates Today</span>
						<span class="text-slate-700">{updateCount}</span>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
