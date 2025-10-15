<script>
	import Button from '$lib/components/common/Button.svelte';
	import CustomSelect from '$lib/components/common/CustomSelect.svelte';
</script>

<div class="flex flex-col">
	<form
		class=" grid w-fit grid-rows-5 gap-y-6 *:grid *:grid-cols-2 *:grid-rows-subgrid **:[input]:rounded-md **:[input]:bg-white **:[input]:px-2.5 **:[input]:py-1.5 **:[input]:outline **:[input]:outline-black **:[input]:focus:ring-1 **:[input]:focus:ring-custom-light-orange **:[input]:focus:outline-none"
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
			<select
				class="rounded-md bg-white px-2.5 py-1.5 outline-1 outline-black focus:ring-1 focus:ring-custom-light-orange focus:outline-none"
			>
				<option>Mess A</option>
			</select>
		</label>
		<label>
			Duration
			<input type="number" placeholder="Months" />
		</label>
	</form>

	<Button class="mt-12 self-center !py-1.5 !text-lg">Confirm Registration</Button>
</div>
