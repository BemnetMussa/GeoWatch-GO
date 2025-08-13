<script lang="ts">
	import { AlertTriangle, Shield, TrendingUp, Zap } from 'lucide-svelte';
	import type { FireData } from '$lib/types';
	
	interface Props {
		fires: FireData[];
	}
	
	let { fires }: Props = $props();
	
	let riskAssessment = $derived.by(() => {
		if (!fires.length) return null;
		
		// Calculate risk metrics
		const highConfidenceFires = fires.filter(f => f.confidence >= 80);
		const highIntensityFires = fires.filter(f => f.brightness > 400);
		const highFRPFires = fires.filter(f => (f.frp || 0) > 100);
		
		// Risk categories
		const criticalRisk = fires.filter(f => f.confidence >= 80 && f.brightness > 400 && (f.frp || 0) > 100);
		const highRisk = fires.filter(f => f.confidence >= 70 && (f.brightness > 350 || (f.frp || 0) > 50));
		const mediumRisk = fires.filter(f => f.confidence >= 50 && f.brightness > 300);
		const lowRisk = fires.filter(f => f.confidence < 50 || f.brightness <= 300);
		
		// Calculate overall risk score (0-100)
		const riskScore = Math.min(100, Math.round(
			(criticalRisk.length * 4 + highRisk.length * 3 + mediumRisk.length * 2 + lowRisk.length * 1) / fires.length * 25
		));
		
		// Risk level determination
		let riskLevel: 'Low' | 'Medium' | 'High' | 'Critical';
		let riskColor: string;
		
		if (riskScore >= 80) {
			riskLevel = 'Critical';
			riskColor = 'text-red-600 bg-red-50 border-red-200';
		} else if (riskScore >= 60) {
			riskLevel = 'High';
			riskColor = 'text-orange-600 bg-orange-50 border-orange-200';
		} else if (riskScore >= 40) {
			riskLevel = 'Medium';
			riskColor = 'text-yellow-600 bg-yellow-50 border-yellow-200';
		} else {
			riskLevel = 'Low';
			riskColor = 'text-green-600 bg-green-50 border-green-200';
		}
		
		// Recommendations based on risk level
		const recommendations = [];
		if (criticalRisk.length > 0) {
			recommendations.push('Immediate evacuation may be necessary for areas near critical risk fires');
			recommendations.push('Deploy emergency response teams to critical fire locations');
		}
		if (highRisk.length > 5) {
			recommendations.push('Increase monitoring frequency for high-risk areas');
			recommendations.push('Prepare firefighting resources for rapid deployment');
		}
		if (riskScore > 60) {
			recommendations.push('Issue fire weather warnings to local communities');
			recommendations.push('Restrict outdoor activities and fire-prone operations');
		}
		if (recommendations.length === 0) {
			recommendations.push('Continue routine monitoring and fire prevention measures');
			recommendations.push('Maintain readiness for potential fire activity increases');
		}
		
		return {
			riskScore,
			riskLevel,
			riskColor,
			criticalRisk: criticalRisk.length,
			highRisk: highRisk.length,
			mediumRisk: mediumRisk.length,
			lowRisk: lowRisk.length,
			highConfidenceFires: highConfidenceFires.length,
			highIntensityFires: highIntensityFires.length,
			highFRPFires: highFRPFires.length,
			recommendations
		};
	});
</script>

{#if riskAssessment}
	<div class="space-y-8">
		<!-- Overall Risk Score -->
		<div class="bg-white rounded-lg border border-slate-200 p-6">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-lg font-semibold text-slate-900">Overall Risk Assessment</h3>
				<div class="flex items-center space-x-2 px-4 py-2 rounded-lg border {riskAssessment.riskColor}">
					<AlertTriangle class="w-4 h-4" />
					<span class="font-semibold">{riskAssessment.riskLevel} Risk</span>
				</div>
			</div>
			
			<div class="flex items-center space-x-4 mb-4">
				<div class="flex-1">
					<div class="flex items-center justify-between mb-2">
						<span class="text-sm font-medium text-slate-600">Risk Score</span>
						<span class="text-2xl font-bold text-slate-900">{riskAssessment.riskScore}/100</span>
					</div>
					<div class="w-full bg-slate-200 rounded-full h-4">
						<div 
							class="h-4 rounded-full transition-all duration-500 {riskAssessment.riskScore >= 80 ? 'bg-red-500' : riskAssessment.riskScore >= 60 ? 'bg-orange-500' : riskAssessment.riskScore >= 40 ? 'bg-yellow-500' : 'bg-green-500'}"
							style="width: {riskAssessment.riskScore}%"
						></div>
					</div>
				</div>
			</div>
		</div>
		
		<!-- Risk Categories -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
			<div class="bg-gradient-to-br from-red-50 to-red-100 rounded-lg p-4 border border-red-200">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-red-600">Critical Risk</p>
						<p class="text-xl font-bold text-red-900">{riskAssessment.criticalRisk}</p>
					</div>
					<AlertTriangle class="w-6 h-6 text-red-500" />
				</div>
			</div>
			
			<div class="bg-gradient-to-br from-orange-50 to-orange-100 rounded-lg p-4 border border-orange-200">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-orange-600">High Risk</p>
						<p class="text-xl font-bold text-orange-900">{riskAssessment.highRisk}</p>
					</div>
					<TrendingUp class="w-6 h-6 text-orange-500" />
				</div>
			</div>
			
			<div class="bg-gradient-to-br from-yellow-50 to-yellow-100 rounded-lg p-4 border border-yellow-200">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-yellow-600">Medium Risk</p>
						<p class="text-xl font-bold text-yellow-900">{riskAssessment.mediumRisk}</p>
					</div>
					<Zap class="w-6 h-6 text-yellow-500" />
				</div>
			</div>
			
			<div class="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-4 border border-green-200">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-green-600">Low Risk</p>
						<p class="text-xl font-bold text-green-900">{riskAssessment.lowRisk}</p>
					</div>
					<Shield class="w-6 h-6 text-green-500" />
				</div>
			</div>
		</div>
		
		<!-- Risk Factors -->
		<div class="bg-white rounded-lg border border-slate-200 p-6">
			<h3 class="text-lg font-semibold text-slate-900 mb-6">Risk Factors Analysis</h3>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
				<div class="text-center">
					<div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-3">
						<span class="text-2xl font-bold text-red-600">{riskAssessment.highConfidenceFires}</span>
					</div>
					<p class="text-sm font-medium text-slate-900">High Confidence</p>
					<p class="text-xs text-slate-500">Fires with ≥80% confidence</p>
				</div>
				
				<div class="text-center">
					<div class="w-16 h-16 bg-orange-100 rounded-full flex items-center justify-center mx-auto mb-3">
						<span class="text-2xl font-bold text-orange-600">{riskAssessment.highIntensityFires}</span>
					</div>
					<p class="text-sm font-medium text-slate-900">High Intensity</p>
					<p class="text-xs text-slate-500">Brightness >400K</p>
				</div>
				
				<div class="text-center">
					<div class="w-16 h-16 bg-yellow-100 rounded-full flex items-center justify-center mx-auto mb-3">
						<span class="text-2xl font-bold text-yellow-600">{riskAssessment.highFRPFires}</span>
					</div>
					<p class="text-sm font-medium text-slate-900">High Power</p>
					<p class="text-xs text-slate-500">FRP >100 MW</p>
				</div>
			</div>
		</div>
		
		<!-- Recommendations -->
		<div class="bg-white rounded-lg border border-slate-200 p-6">
			<h3 class="text-lg font-semibold text-slate-900 mb-4">Recommendations</h3>
			<div class="space-y-3">
				{#each riskAssessment.recommendations as recommendation, index}
					<div class="flex items-start space-x-3">
						<span class="inline-flex items-center justify-center w-6 h-6 bg-orange-100 text-orange-800 text-xs font-bold rounded-full mt-0.5">
							{index + 1}
						</span>
						<p class="text-sm text-slate-700 flex-1">{recommendation}</p>
					</div>
				{/each}
			</div>
		</div>
	</div>
{:else}
	<div class="text-center py-12">
		<AlertTriangle class="w-12 h-12 text-slate-300 mx-auto mb-4" />
		<p class="text-slate-500">No data available for risk assessment</p>
	</div>
{/if}
