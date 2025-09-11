<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$lib/components/common/Button.svelte';
	import MainPanelHeader from '$lib/components/common/MainPanelHeader.svelte';

	let registered = $state(false);
</script>

<MainPanelHeader>Mess Portal</MainPanelHeader>

<section
	class="grid grid-cols-2 gap-x-8 gap-y-14 px-24 pt-24 pb-6 text-3xl leading-none font-medium text-custom-black dark:text-custom-off-white"
>
	<p>Current Registered Mess :</p>
	<p>Unregistered</p>
	<p>Next registration date :</p>
	<p>Today</p>
	<p>Time left for next registration :</p>
	<p>Registration Active Now</p>

	{#if registered}
		<Button
			onclick={() => {
				goto('student/register');
			}}
			class="col-span-2 mx-auto mt-16">Go for Registration</Button
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
