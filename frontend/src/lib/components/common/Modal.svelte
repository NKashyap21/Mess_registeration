<script lang="ts">
	import { Dialog, type WithoutChild } from 'bits-ui';
	import Button from './Button.svelte';
	import type { ClassValue } from 'svelte/elements';
	import xSvg from '$lib/assets/cross-red.svg';
	import { colorScheme } from '$lib/state/dark.svelte';

	type Props = Dialog.RootProps & {
		buttonText: string;
		contentProps?: WithoutChild<Dialog.ContentProps>;
		class?: ClassValue;
		// ...other component props if you wish to pass them
	};

	let {
		open = $bindable(false),
		children,
		buttonText,
		contentProps,
		class: className,
		...restProps
	}: Props = $props();
</script>

<Dialog.Root bind:open {...restProps}>
	<Dialog.Trigger class="{className} {colorScheme.dark ? 'dark' : ''}">
		<Button>
			{buttonText}
		</Button>
	</Dialog.Trigger>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 h-screen w-screen bg-black/80" />
		<Dialog.Content
			class="fixed top-1/2 right-1/2 bottom-1/2 left-1/2 dark:text-white {colorScheme.dark
				? 'dark'
				: ''} z-50 flex h-fit w-fit -translate-x-1/2 -translate-y-1/2 flex-col rounded-lg bg-[#e1e1e1] p-5 dark:bg-custom-black"
		>
			<Dialog.Close class="self-end">
				<img src={xSvg} alt="Close" />
			</Dialog.Close>
			{@render children?.()}
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
