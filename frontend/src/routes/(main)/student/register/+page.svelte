<script lang="ts">
	import Button from '$lib/components/common/Button.svelte';
	import Progress from '$lib/components/common/Progress.svelte';

	const messData: Record<string, string[]> = {
		'Mess-A': ['LDH', 'UDH'],
		'Mess-B': ['LDH', 'UDH']
	};
</script>

<form class="flex flex-col gap-y-16 px-16 pt-16 **:text-center dark:text-white" method="post">
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
						<Progress
							outerClass="w-[40rem]"
							innerClass={idx % 2 == 0
								? 'bg-[#919191] dark:bg-[#B4B4B4]'
								: 'bg-custom-light-orange'}
							value={1}
						/>
						<h3 class="px-6 text-2xl font-medium">1000/2500</h3>
						<input class="" type="radio" name="mess" value={hall} />
					</div>
				{/each}
			</div>
		</section>
	{/each}

	<Button type="submit" class="right-0  left-0 mx-auto">Register</Button>
</form>
