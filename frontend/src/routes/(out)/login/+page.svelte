<script lang="ts">
	import diningSvg from '$lib/assets/dining.svg';
	import diningDarkSvg from '$lib/assets/dining-dark.svg';
	import googleIcon from '$lib/assets/google.webp';
</script>

<img src={diningSvg} alt="dining icon" class="dark:hidden" />
<img src={diningDarkSvg} alt="dining icon" class="hidden dark:block" />

<p class="mt-8 mb-2 text-2xl text-custom-black dark:text-custom-light-grey">Welcome to</p>
<h2 class="text-4xl font-bold text-custom-black dark:text-white">IITH Mess Portal</h2>

<form class="auth-form" method="post" action="?/OAuth2">
	<button
		class="mt-16 flex flex-row items-center gap-x-4 rounded-md bg-white px-8 py-5 text-xl font-semibold text-custom-black shadow-xs shadow-custom-dark-grey hover:brightness-95 active:scale-[99%] active:brightness-90 dark:bg-custom-dark-grey dark:text-white"
	>
		<img src={googleIcon} class="size-8" alt="google icon" />
		<p class="leading-none">Sign in with Google</p>
	</button>
</form>
