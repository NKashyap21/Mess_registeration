<script lang="ts">
	let { children } = $props();
</script>

<div class="size-[40rem] rounded-xl bg-[#f8f8f8] dark:bg-custom-dark-grey">
	{@render children?.()}
</div>
