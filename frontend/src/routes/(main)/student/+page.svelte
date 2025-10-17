<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$lib/components/common/Button.svelte';
	import MainPanelHeader from '$lib/components/common/MainPanelHeader.svelte';

	let { data } = $props();
	let registered = $state(data.user.mess_id != 'No mess assigned');
	console.log(data);
</script>

<MainPanelHeader>Mess Portal</MainPanelHeader>

<section
	class="grid grid-cols-2 gap-x-8 gap-y-14 px-24 pt-24 pb-6 text-3xl leading-none font-medium text-custom-black dark:text-custom-off-white"
>
	<p>Current Registered Mess :</p>
	<p>{data.user.mess_id}</p>
	<p>Next registration date :</p>
	<p>Unknown</p>
	<p>Registration Status :</p>
	<p>
		{data.regData.regular && data.regData.veg
			? 'Regular and Veg Registration Active Now'
			: data.regData.veg
				? 'Veg Registration Active Now'
				: data.regData.regular
					? 'Regular Registration Active Now'
					: 'Inactive'}
	</p>

	<div class="col-span-2 mx-auto mt-16 flex gap-x-8">
		{#if !registered}
			{#if data.regData.regular}
				<Button
					disabled={data.userMessData.status == 'pending_sync' ||
						(data.userMessData.status == 'confirmed' && !registered)}
					onclick={() => {
						goto('student/register?veg=false');
					}}
					class="">Go for Regular Registration</Button
				>
			{/if}
			{#if data.regData.veg}
				<Button
					disabled={data.userMessData.status == 'pending_sync' ||
						(data.userMessData.status == 'confirmed' && !registered)}
					onclick={() => {
						goto('student/register?veg=true');
					}}
					class="">Go for Veg Registration</Button
				>
			{/if}
		{:else}
			<Button
				onclick={() => {
					goto('student/swap');
				}}
				class="">Swap Mess</Button
			>
		{/if}
	</div>
</section>
