<script>
	import Button from '$lib/components/common/Button.svelte';
</script>

<div class="flex flex-col">
	<form
		class="mt-12 grid w-fit grid-rows-5 gap-y-4 *:grid *:grid-cols-2 *:grid-rows-subgrid **:[input]:border"
		method="POST"
	>
		<label>
			Name
			<input />
		</label>
		<label>
			Roll Number
			<input />
		</label>
		<label>
			Designation
			<input />
		</label>
		<label>
			Register to
			<input />
		</label>
		<label>
			Duration
			<input type="number" placeholder="months" />
		</label>
	</form>

	<Button class="mt-12 self-center !text-lg">Confirm Registration</Button>
</div>
