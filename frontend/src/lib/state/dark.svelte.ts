import { browser } from '$app/environment';

if (browser) {
	if (localStorage.getItem('colorDark') == null) {
		localStorage.setItem('colorDark', '1');
	}
}
export let colorScheme = $state({
	dark: browser ? parseInt(localStorage.getItem('colorDark')!) == 1 : false
});
