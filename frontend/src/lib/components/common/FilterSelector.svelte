<script lang="ts">
	let { items, active = $bindable(items[0]) }: { items: string[]; active?: string } = $props();
</script>

<div class="flex gap-x-3">
	{#each items as item}
		<button
			class="{active == item
				? 'bg-custom-orange text-white'
				: 'bg-custom-dark-white text-custom-mid-grey outline outline-custom-mid-grey'} rounded-lg px-3 py-3.5 leading-0 font-semibold"
			onclick={() => {
				active = item;
			}}>{item}</button
		>
	{/each}
</div>
