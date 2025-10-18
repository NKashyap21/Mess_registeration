<script>
	import { enhance } from '$app/forms';
	import Button from '$lib/components/common/Button.svelte';

	let { form } = $props();
	$effect(() => {
		if (form != undefined) {
			if (form.message != undefined) {
				alert(form.message);
			}
		}
	});
</script>

<div class="flex flex-col dark:text-white">
	<form
		id="add-form"
		class=" grid w-fit grid-rows-5 gap-y-6 text-lg *:grid *:grid-cols-2 *:grid-rows-subgrid **:[input]:w-[20rem] **:[input]:rounded-sm **:[input]:bg-white **:[input]:px-2.5 **:[input]:py-2 **:[input]:outline **:[input]:outline-black **:[input]:focus:ring-1 **:[input]:focus:ring-custom-light-orange **:[input]:focus:outline-none **:[input]:dark:bg-custom-mid-grey **:[input]:dark:outline-custom-dark-white/50"
		method="POST"
		use:enhance
	>
		<label>
			Name *
			<input required name="name" />
		</label>
		<label>
			Email *
			<input required name="email" type="email" />
		</label>
		<label>
			Roll Number *
			<input required name="roll" />
		</label>
		<label>
			Phone Number
			<input name="phone" />
		</label>
		<label>
			Designation *
			<select
				class="rounded-md bg-white px-2.5 py-1.5 outline-1 outline-black focus:ring-1 focus:ring-custom-light-orange focus:outline-none dark:bg-custom-mid-grey dark:outline-custom-dark-white/50"
				required
				name="type"
			>
				<option selected value={0}>Student</option>
				<option value={1}>Staff</option>
				<option value={2}>Admin</option>
			</select>
		</label>
		<label>
			Register to *
			<select
				class="rounded-md bg-white px-2.5 py-2 outline-1 outline-black focus:ring-1 focus:ring-custom-light-orange focus:outline-none dark:bg-custom-mid-grey dark:outline-custom-dark-white/50"
				required
				name="mess"
			>
				<option value={1}>Mess A - LDH</option>
				<option value={2}>Mess A - UDH</option>
				<option value={3}>Mess B - LDH</option>
				<option value={4}>Mess B - UDH</option>
				<option value={5}>Veg Mess</option>
				<option selected value={0}>None</option>
			</select>
		</label>
	</form>

	<Button
		type="submit"
		form="add-form"
		class="mt-12 self-center !bg-white !py-1.5 !text-lg dark:!bg-custom-mid-grey"
		>Confirm Registration</Button
	>
</div>
