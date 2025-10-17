import { PRIVATE_API_URL } from '$env/static/private';
import type { PageServerLoad } from './$types';

export const ssr = false;

export const load: PageServerLoad = async ({ fetch }) => {
	return {
		status: await (
			await fetch(PRIVATE_API_URL + '/office/status', {
				method: 'GET',
				credentials: 'include'
			})
		).json()
	};
};
