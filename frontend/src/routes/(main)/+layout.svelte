<script lang="ts">
	let { children } = $props();
</script>

<div class="rounded-xl bg-[#f8f8f8] p-12 dark:bg-custom-dark-grey">
	{@render children?.()}
</div>
