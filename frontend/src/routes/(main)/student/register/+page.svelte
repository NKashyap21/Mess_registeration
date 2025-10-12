<script lang="ts">
	import Button from '$lib/components/common/Button.svelte';

	const messData: Record<string, string[]> = {
		'Mess-A': ['LDH', 'UDH'],
		'Mess-B': ['LDH', 'UDH']
	};
</script>

<form
	class="flex flex-col gap-y-16 px-16 py-16 **:text-center dark:text-white"
	method="post"
	action="?/asd"
>
	<h1 class="absolute top-12 right-0 left-0 mx-auto w-fit text-[2.5rem] font-semibold">
		Mess Registration
	</h1>
	{#each Object.keys(messData) as messName}
		<section class="flex flex-row">
			<h2 class="mr-12 w-[8rem] text-3xl font-bold">{messName}</h2>
			<div class="flex flex-col gap-y-12">
				{#each messData[messName] as hall, idx}
					<div class="flex items-center">
						<h3 class="w-[8rem] text-2xl">{hall}</h3>
						<div class="h-5 w-[40rem] overflow-hidden rounded-full bg-[#E0E0E0] p-1 dark:bg-black">
							<div
								class="h-full rounded-full {idx % 2 == 0
									? 'bg-[#919191] dark:bg-[#B4B4B4]'
									: 'bg-custom-light-orange'}"
								style="width: {40}%;"
							></div>
						</div>
						<h3 class="px-6 text-2xl font-medium">1000/2500</h3>
						<input class="" type="radio" name="mess" value={hall} />
					</div>
				{/each}
			</div>
		</section>
	{/each}

	<Button type="submit" class="absolute right-0 bottom-32 left-0 mx-auto">Register</Button>
</form>
