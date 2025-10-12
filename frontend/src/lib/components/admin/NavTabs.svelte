<script lang="ts">
	import type { HTMLButtonAttributes } from 'svelte/elements';
	let {
		text,
		active,
		class: className,
		...rest
	}: { active: boolean; text: string } & HTMLButtonAttributes = $props();
</script>

<button
	class="rounded-t-2xl {className} px-8 py-4 text-center leading-0 font-bold text-nowrap select-none hover:cursor-pointer hover:brightness-95 {active
		? 'bg-white'
		: 'bg-custom-light-grey'}"
	{...rest}
>
	{text}
</button>
