<script lang="ts">
	import { enhance } from '$app/forms';
	import { page } from '$app/state';
	import Button from '$lib/components/common/Button.svelte';
	import Progress from '$lib/components/common/Progress.svelte';

	let { data, form } = $props();
	$effect(() => {
		if (form != undefined) {
			if (form.message == undefined) {
				alert(form.error);
			} else {
				alert(form.message);
			}
		}
	});
	console.log(data);

	let activeIdx = $state('-1');
</script>

<form
	use:enhance
	class="flex flex-col gap-y-6 **:text-center lg:gap-y-6 lg:px-2 lg:pt-6 2xl:gap-y-12 2xl:px-8 2xl:pt-10 dark:text-white"
	method="post"
>
	<h1
		class="absolute top-24 right-0 left-0 mx-auto w-fit text-3xl font-semibold lg:top-12 lg:text-4xl 2xl:text-[2.5rem]"
	>
		{page.url.searchParams.get('veg') == 'true' ? 'Veg' : 'Regular'} Mess Registration
	</h1>
	<input hidden value={page.url.searchParams.get('veg')} name="veg" />
	{#each Object.keys(data.messStats.data.stats) as messName, idx}
		<section class="flex flex-col max-lg:gap-y-6 lg:flex-row">
			<h2
				class="w-fit text-xl font-bold lg:mr-6 lg:w-24 lg:text-2xl 2xl:mr-12 2xl:w-32 2xl:text-3xl"
			>
				{messName}
			</h2>
			<div class="flex flex-col gap-y-6 lg:gap-y-4 2xl:gap-y-6">
				{#each Object.keys(data.messStats.data.stats[messName]) as hall, idx2}
					<label
						onchange={() => {
							activeIdx = idx.toString() + idx2.toString();
						}}
						class="flex {activeIdx == idx.toString() + idx2.toString()
							? 'bg-custom-orange'
							: 'hover:bg-custom-orange/30'} flex-row items-center rounded-full py-2 hover:cursor-pointer 2xl:py-4"
					>
						<h3 class="w-16 text-lg lg:w-24 lg:text-xl 2xl:w-32 2xl:text-2xl">{hall}</h3>
						<Progress
							outerClass="lg:w-[40vw] w-[30vw] max-2xl:!h-4"
							innerClass=" {activeIdx == idx.toString() + idx2.toString()
								? 'bg-custom-orange dark:bg-custom-orange'
								: 'bg-[#919191] dark:bg-[#B4B4B4]'}"
							value={(data.messStats.data.stats[messName][hall].count * 100) /
								data.messStats.data.stats[messName][hall].capacity}
						/>
						<h3 class="px-2 text-lg font-medium lg:px-4 2xl:px-6 2xl:text-2xl">
							{data.messStats.data.stats[messName][hall].count}/{data.messStats.data.stats[
								messName
							][hall].capacity}
						</h3>
						<input
							class="mr-3 size-3 appearance-none rounded-full border border-black transition-colors checked:bg-[radial-gradient(circle,#000000_40%,#ffffff_40%)] hover:cursor-pointer lg:mr-5 lg:size-4 2xl:mr-9 2xl:size-6 dark:bg-white dark:checked:bg-[radial-gradient(circle,#000000_50%,#ffffff_50%)]"
							type="radio"
							name="mess"
							value={data.messStats.data.stats[messName][hall].id}
						/>
					</label>
				{/each}
			</div>
		</section>
	{/each}

	<Button type="submit" class="right-0 left-0 mx-auto">Register</Button>
</form>
