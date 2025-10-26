<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$lib/components/common/Button.svelte';
	import MainPanelHeader from '$lib/components/common/MainPanelHeader.svelte';

	let { data } = $props();
	let registered = $state(data.userMessData.next_mess != 0);
	console.log(data);
</script>

<MainPanelHeader>Mess Portal</MainPanelHeader>

<section
	class="grid grid-cols-1 gap-y-2 px-6 pt-8 pb-10 text-lg font-medium text-custom-black *:even:mb-4 sm:text-xl md:px-10 lg:grid-cols-2 lg:gap-x-4 lg:gap-y-10 lg:px-12 lg:pt-12 lg:pb-6 lg:text-2xl 2xl:gap-x-8 2xl:gap-y-14 2xl:px-24 2xl:pt-24 2xl:text-3xl dark:text-custom-off-white"
>
	<!-- Current Mess -->
	<p class="text-gray-700 dark:text-custom-off-white">Current Registered Mess :</p>
	<p class="text-gray-900 dark:text-custom-off-white">
		{data.user.mess_id ?? 'Unknown'}
	</p>

	<!-- Next Registration -->
	<p class="text-gray-700 dark:text-custom-off-white">Next Registration :</p>
	<p>
		{data.userMessData?.next_mess_name ?? 'Unknown'}
		{#if data.userMessData.status == 'pending_sync'}
			<span
				class="mt-2 block text-base text-gray-500 sm:text-lg md:text-xl 2xl:text-2xl dark:text-gray-300"
			>
				Pending assignment to {data.userMessData.mess == 1
					? 'Mess A (LDH)'
					: data.userMessData.mess == 2
						? 'Mess A (UDH)'
						: data.userMessData.mess == 3
							? 'Mess B (LDH)'
							: data.userMessData.mess == 4
								? 'Mess B (UDH)'
								: data.userMessData.mess == 5
									? 'Veg Mess'
									: 'Unknown'}
			</span>
		{/if}
	</p>

	<!-- Registration Status -->
	<p class="text-gray-700 dark:text-custom-off-white">Registration Status :</p>
	<p>
		{data.regData.regular && data.regData.veg
			? 'Regular and Veg Registration Active Now'
			: data.regData.veg
				? 'Veg Registration Active Now'
				: data.regData.regular
					? 'Regular Registration Active Now'
					: 'Inactive'}
	</p>

	<!-- Buttons -->
	<div
		class="col-span-1 mt-6 flex flex-col gap-y-4 sm:flex-row sm:justify-center sm:gap-x-4 md:col-span-2 md:mt-8 lg:mx-auto lg:mt-6 lg:flex-row lg:gap-x-6 2xl:mt-16 2xl:gap-x-8"
	>
		{#if !registered}
			{#if data.regData.regular}
				<Button
					disabled={data.userMessData.status == 'pending_sync' ||
						(data.userMessData.status == 'confirmed' && registered)}
					onclick={() => goto('student/register?veg=false')}
					class="w-full sm:w-auto"
				>
					Go for Regular Registration
				</Button>
			{/if}
			{#if data.regData.veg}
				<Button
					disabled={data.userMessData.status == 'pending_sync' ||
						(data.userMessData.status == 'confirmed' && registered)}
					onclick={() => goto('student/register?veg=true')}
					class="w-full sm:w-auto"
				>
					Go for Veg Registration
				</Button>
			{/if}
			{#if data.userMessData.current_mess != 0 && (data.userMessData?.current_mess ?? 5) != 5}
				<Button onclick={() => goto('student/swap')} class="w-full sm:w-auto">Swap Mess</Button>
			{/if}
		{/if}
	</div>
</section>
