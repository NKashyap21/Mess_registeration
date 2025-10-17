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
</script>

<form
	use:enhance
	class="flex flex-col gap-y-16 px-16 pt-16 **:text-center dark:text-white"
	method="post"
>
	<h1 class="absolute top-12 right-0 left-0 mx-auto w-fit text-[2.5rem] font-semibold">
		{page.url.searchParams.get('veg') == 'true' ? 'Veg' : 'Regular'} Mess Registration
	</h1>
	<input hidden value={page.url.searchParams.get('veg')} name="veg" />
	{#each Object.keys(data.messStats.stats) as messName}
		<section class="flex flex-row">
			<h2 class="mr-12 w-[8rem] text-3xl font-bold">{messName}</h2>
			<div class="flex flex-col gap-y-12">
				{#each Object.keys(data.messStats.stats[messName]) as hall, idx}
					<div class="flex items-center">
						<h3 class="w-[8rem] text-2xl">{hall}</h3>
						<Progress
							outerClass="w-[40rem]"
							innerClass={idx % 2 == 0
								? 'bg-[#919191] dark:bg-[#B4B4B4]'
								: 'bg-custom-light-orange'}
							value={1}
						/>
						<h3 class="px-6 text-2xl font-medium">
							{data.messStats.stats[messName][hall].count}/{data.messStats.stats[messName][hall]
								.capacity}
						</h3>
						<input
							class=""
							type="radio"
							name="mess"
							value={data.messStats.stats[messName][hall].id}
						/>
					</div>
				{/each}
			</div>
		</section>
	{/each}

	<Button type="submit" class="right-0  left-0 mx-auto">Register</Button>
</form>
