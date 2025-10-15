<script lang="ts">
	import { colorScheme } from '$lib/state/dark.svelte';
	import { DropdownMenu, type WithoutChild } from 'bits-ui';
	import type { Snippet } from 'svelte';

	type Props = DropdownMenu.RootProps & {
		trigger: Snippet;
		content: Snippet;
		contentProps?: WithoutChild<DropdownMenu.ContentProps>;
		// other component props if needed
	};

	let { open = $bindable(false), content, trigger, contentProps, ...restProps }: Props = $props();
</script>

<DropdownMenu.Root bind:open {...restProps}>
	<DropdownMenu.Trigger>
		{@render trigger()}
	</DropdownMenu.Trigger>
	<DropdownMenu.Portal>
		<DropdownMenu.Content
			class="mt-4 flex flex-col shadow-sm {colorScheme.dark
				? 'dark'
				: ''} items-center rounded-lg bg-custom-off-white p-6 dark:bg-custom-mid-grey dark:text-white"
			{...contentProps}
		>
			{@render content()}
		</DropdownMenu.Content>
	</DropdownMenu.Portal>
</DropdownMenu.Root>
