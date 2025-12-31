import { PRIVATE_API_URL } from '$env/static/private';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const [swapsRes, mySwapRes] = await Promise.all([
		fetch(PRIVATE_API_URL + '/students/getSwaps', {
			method: 'GET',
			credentials: 'include'
		}),
		fetch(PRIVATE_API_URL + '/students/getSwapByID', {
			method: 'GET',
			credentials: 'include'
		})
	]);

	const swapsData = await swapsRes.json();
	const mySwapData = await mySwapRes.json();

	console.log(await mySwapRes);

	return {
		swapData: (swapsData.data ?? []) as Array<any>,
		mySwap: mySwapRes.ok ? mySwapData.data : null
	};
};
