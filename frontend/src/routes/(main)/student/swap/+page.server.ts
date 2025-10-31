import { PRIVATE_API_URL } from '$env/static/private';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	return {
		swapData: (
			await (
				await fetch(PRIVATE_API_URL + '/students/getSwaps', {
					method: 'GET',
					credentials: 'include'
				})
			).json()
		).data as Array<any>
	};
};
