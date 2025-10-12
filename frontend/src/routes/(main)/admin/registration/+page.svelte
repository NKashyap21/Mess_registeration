<script lang="ts">
	import Button from '$lib/components/common/Button.svelte';
</script>

<section class="flex flex-col">
	<div>Registration Stats</div>
	<div class="flex">
		<div class="flex flex-col">
			<div></div>
			<div class="mt-20"></div>
		</div>
		<div class="flex flex-col gap-y-6 *:px-8 **:w-full **:text-lg">
			<Button>Import New List</Button>
			<Button>Start Regular Registration</Button>
			<Button>Stop Veg Registration</Button>
			<Button>Enable Latest Registration</Button>
		</div>
	</div>
</section>
