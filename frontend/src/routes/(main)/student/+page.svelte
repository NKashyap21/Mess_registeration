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
	<p>Today</p>
	<p>Registration Status :</p>
	<p>
		{data.regData.regular
			? 'Registration Active Now'
			: data.regData.veg
				? 'Veg Registration Active Now'
				: 'Inactive'}
	</p>

	{#if !registered}
		<Button
			disabled={!(data.regData.regular || data.regData.veg)}
			onclick={() => {
				goto('student/register');
			}}
			class="col-span-2 mx-auto mt-16">Go for {data.regData.veg ? 'Veg' : ''} Registration</Button
		>
	{:else}
		<Button
			onclick={() => {
				goto('student/swap');
			}}
			class="col-span-2 mx-auto mt-16">Swap Mess</Button
		>
	{/if}
</section>
