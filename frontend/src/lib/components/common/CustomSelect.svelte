<script lang="ts">
	import { Select } from 'bits-ui';
	import chevronDown from '$lib/assets/chevron-down.svg';
	import type { ClassValue } from 'svelte/elements';
	import { colorScheme } from '$lib/state/dark.svelte';
	let {
		value = $bindable(''),
		name,
		items,
		widthClass
	}: {
		value: string;
		name?: string;
		items: { label: string; value: string }[];
		widthClass?: ClassValue;
	} = $props();
</script>

<Select.Root {name} type="single" bind:value {items}>
	<Select.Trigger
		class="flex {widthClass} {colorScheme.dark
			? 'dark'
			: ''} w-[14rem] items-center justify-center gap-x-3 rounded-md border border-custom-red bg-white pt-4 pb-4 text-xl leading-0 font-medium shadow-xs dark:border-custom-light-grey dark:bg-custom-mid-grey dark:text-white"
	>
		{value}
		<img
			src={chevronDown}
			alt="chevron down"
			class="-mb-1 size-4 brightness-0 dark:brightness-100"
		/>
	</Select.Trigger>
	<Select.Portal>
		<Select.Content
			class="z-[9999] mt-2 w-[14rem] {colorScheme.dark
				? 'dark'
				: ''} {widthClass} rounded-md border border-custom-red bg-custom-off-white px-2 py-3 dark:border-custom-light-grey dark:bg-custom-mid-grey dark:text-white "
		>
			<Select.ScrollUpButton />
			<Select.Viewport class="dark:bg-custom-mid-grey">
				{#each items as item (item.value)}
					<Select.Item
						value={item.value}
						class="flex w-full flex-col items-start rounded-sm bg-custom-off-white px-2 py-4 text-center leading-0 select-none not-last:mb-2 data-highlighted:brightness-95 dark:bg-custom-mid-grey"
						label={item.label}>{item.label}</Select.Item
					>
				{/each}
				<Select.ScrollDownButton />
			</Select.Viewport>
		</Select.Content>
	</Select.Portal>
</Select.Root>
