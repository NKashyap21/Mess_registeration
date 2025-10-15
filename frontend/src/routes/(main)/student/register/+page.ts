import { PUBLIC_API_URL } from '$env/static/public';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	return {
		messStats: await (
			await fetch(PUBLIC_API_URL + '/students/messStatsGrouped', { method: 'GET' })
		).json()
	};
};
