import { PUBLIC_API_URL } from '$env/static/public';
import type { PageLoad } from './$types';

export const ssr = false;

export const load: PageLoad = async ({ fetch }) => {
	return {
		registrationState: await (
			await fetch(PUBLIC_API_URL + '/students/', { method: 'get', credentials: 'include' })
		).json(),
		messStats: await (
			await fetch(PUBLIC_API_URL + '/office/messStatsGrouped', {
				method: 'GET',
				credentials: 'include'
			})
		).json()
	};
};
