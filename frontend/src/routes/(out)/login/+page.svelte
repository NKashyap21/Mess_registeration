<script lang="ts">
	import snackSvg from '$lib/assets/dining.svg';
	import googleIcon from '$lib/assets/google.webp';
</script>

<img src={snackSvg} alt="dining icon" class="" />

<p class="mt-8 mb-2 text-2xl text-black dark:text-custom-light-grey">Welcome to</p>
<h2 class="text-4xl font-bold text-black dark:text-white">IITH Mess Portal</h2>

<button
	class="mt-16 flex flex-row items-center gap-x-4 rounded-md bg-white px-8 py-5 text-xl font-semibold text-black shadow-xs shadow-custom-dark-grey hover:brightness-95 active:scale-[99%] active:brightness-90 dark:bg-custom-dark-grey dark:text-white"
>
	<img src={googleIcon} class="size-8" alt="google icon" />
	<p class="leading-none">Sign in with Google</p>
</button>
