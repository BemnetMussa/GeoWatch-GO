<script lang="ts">
	import { onMount } from 'svelte';
	import { Bell, X, Flame, TrendingUp } from 'lucide-svelte';
	import type { FireData } from '$lib/types';
	
	interface Notification {
		id: string;
		type: 'new-fire' | 'fire-growth' | 'system';
		title: string;
		message: string;
		timestamp: Date;
		data?: any;
	}
	
	interface Props {
		enabled: boolean;
		newFires: FireData[];
		onNotificationClick?: (notification: Notification) => void;
	}
	
	let { enabled, newFires, onNotificationClick }: Props = $props();
	
	let notifications: Notification[] = $state([]);
	let permissionGranted = $state(false);
	
	onMount(() => {
		// Request notification permission
		if ('Notification' in window) {
			if (Notification.permission === 'granted') {
				permissionGranted = true;
			} else if (Notification.permission !== 'denied') {
				Notification.requestPermission().then(permission => {
					permissionGranted = permission === 'granted';
				});
			}
		}
	});
	
	// Watch for new fires and create notifications
	$effect(() => {
		if (!enabled || !permissionGranted || newFires.length === 0) return;
		
		const notification: Notification = {
			id: Date.now().toString(),
			type: 'new-fire',
			title: 'New Fire Detected',
			message: `${newFires.length} new fire${newFires.length > 1 ? 's' : ''} detected in your area`,
			timestamp: new Date(),
			data: newFires
		};
		
		notifications = [notification, ...notifications.slice(0, 4)]; // Keep only 5 most recent
		
		// Show browser notification
		if ('Notification' in window && permissionGranted) {
			const browserNotification = new Notification(notification.title, {
				body: notification.message,
				icon: '/favicon.png',
				badge: '/favicon.png',
				tag: 'fire-alert'
			});
			
			browserNotification.onclick = () => {
				window.focus();
				onNotificationClick?.(notification);
				browserNotification.close();
			};
			
			// Auto close after 5 seconds
			setTimeout(() => browserNotification.close(), 5000);
		}
	});
	
	function dismissNotification(id: string) {
		notifications = notifications.filter(n => n.id !== id);
	}
	
	function getNotificationIcon(type: string) {
		switch (type) {
			case 'new-fire': return Flame;
			case 'fire-growth': return TrendingUp;
			default: return Bell;
		}
	}
	
	function getNotificationColor(type: string) {
		switch (type) {
			case 'new-fire': return 'border-red-200 bg-red-50';
			case 'fire-growth': return 'border-orange-200 bg-orange-50';
			default: return 'border-blue-200 bg-blue-50';
		}
	}
</script>

{#if notifications.length > 0}
	<div class="fixed top-4 left-1/2 transform -translate-x-1/2 z-50 space-y-2">
		{#each notifications as notification}
			{@const SvelteComponent = getNotificationIcon(notification.type)}
			<div class="bg-white rounded-lg shadow-lg border {getNotificationColor(notification.type)} p-4 max-w-sm animate-in slide-in-from-top-2">
				<div class="flex items-start space-x-3">
					<SvelteComponent class="w-5 h-5 text-orange-600 mt-0.5" />
					
					<div class="flex-1 min-w-0">
						<h4 class="text-sm font-semibold text-slate-900">{notification.title}</h4>
						<p class="text-sm text-slate-600 mt-1">{notification.message}</p>
						<p class="text-xs text-slate-500 mt-2">
							{notification.timestamp.toLocaleTimeString()}
						</p>
					</div>
					
					<button
						onclick={() => dismissNotification(notification.id)}
						class="text-slate-400 hover:text-slate-600 transition-colors"
					>
						<X class="w-4 h-4" />
					</button>
				</div>
			</div>
		{/each}
	</div>
{/if}
