<script>
	import CustomScrollbar from '$lib/components/common/CustomScrollbar.svelte';
</script>

<div class="flex px-16 py-6 dark:text-white">
	<h1 class="absolute top-12 right-0 left-0 mx-auto w-fit text-[2.5rem] font-semibold">
		Mess Swapping
	</h1>

	<div class="flex gap-x-10">
		<div class="flex flex-col items-center gap-y-4">
			<h1 class="text-4xl font-medium">Mess A to Mess B</h1>
			<CustomScrollbar>
				<div class="relative grid h-[16rem] w-[26rem] grid-cols-3 text-center *:font-medium">
					<div>Sl no.</div>
					<div>Name</div>
					<div>Roll No.</div>
				</div>
			</CustomScrollbar>
		</div>
		<div class="flex flex-col items-center gap-y-4">
			<h1 class="text-4xl font-medium">Mess B to Mess A</h1>
			<CustomScrollbar>
				<div class="relative grid h-[16rem] w-[26rem] grid-cols-3 *:font-medium **:text-center">
					<div>Sl no.</div>
					<div>Name</div>
					<div>Roll No.</div>
				</div>
			</CustomScrollbar>
		</div>
	</div>
</div>
