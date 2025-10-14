<script lang="ts">
	import Button from '$lib/components/common/Button.svelte';
	import FilterSelector from '$lib/components/common/FilterSelector.svelte';

	const filterItems = ['All', 'Unread', 'Mess A', 'Mess B'];
</script>

<section>
	<FilterSelector items={filterItems} />
	<div class="mx-auto mt-12 flex w-fit gap-x-6 self-center *:px-6">
		<Button>Allow Selected ID Card</Button>
		<Button>Block Selected ID Card</Button>
		<Button>Send Email to Selected</Button>
	</div>
</section>
