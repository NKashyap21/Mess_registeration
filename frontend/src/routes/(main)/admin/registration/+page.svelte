<script lang="ts">
	import Button from '$lib/components/common/Button.svelte';
</script>

<section class="flex flex-col">
	<div class="text-2xl font-bold">Registration Stats</div>
	<div class="flex">
		<div class="flex flex-col">
			<div></div>
			<div class="mt-12 grid grid-cols-2 gap-x-4 *:odd:font-medium">
				<p>Next Veg Registration Date:</p>
				<p>26 Aug, 5pm</p>
				<p>Next Regular Registration Date:</p>
				<p>26 Aug, 5pm</p>
			</div>
		</div>
		<div class="ml-auto flex flex-col gap-y-6 *:px-8 **:w-full **:text-lg">
			<Button>Import New List</Button>
			<Button>Start Regular Registration</Button>
			<Button>Stop Veg Registration</Button>
			<Button>Enable Latest Registration</Button>
		</div>
	</div>
</section>
